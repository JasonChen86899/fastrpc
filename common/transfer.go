package common

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

const (
	// EachReadBytes 每个读取的字节大小
	EachReadBytes = 500
)

var (
	ProtocolHead = []byte("srpc")
	RequestID = uint64(1)
	// ReqResMapping = make(map[uint64]chan []byte, 1024)
	ReqResMapping = new(sync.Map)
)

// Transfer 封装tcp conn
type Transfer struct {
	conn net.Conn
}

// NewTransfer 实例化新的传输管道
func NewTransfer(conn net.Conn) *Transfer {
	return &Transfer{conn: conn}
}

// ReadData 读取数据
func (trans *Transfer) ReadData() ([]byte, error) {
	finalData := make([]byte, 0)
	for {
		data := make([]byte, EachReadBytes)
		i, err := trans.conn.Read(data)
		if err != nil {
			return nil, err
		}
		finalData = append(finalData, data[:i]...)
		if i < EachReadBytes {
			break
		}
	}
	return finalData, nil
}

// WriteData 写数据
func (trans *Transfer) WriteData(data []byte) (int, error) {
	num, err := trans.conn.Write(data)
	return num, err
}

func (trans *Transfer) ClientWriteDataByProtocol(data []byte) interface{} {
	enCodeData, requestID := clientEnProto(data)
	var n , totals int
	var err error
	// TODO 是否需要重新建立连接还是根据返回 连接关闭的错误 重新建立连接然后重传
	for totals < len(data) {
		n, err = trans.WriteData(enCodeData[totals:]) // 写数据如果超时或者出现err需要进行重试
		if err == nil {
			break
		}
		log.Println(err)
		totals += n
	}
	var replyChannel = make(chan interface{}, 1) // 设置容量非0 不阻塞生产端
	ReqResMapping.Store(requestID, replyChannel)
	return <- replyChannel // 消费端阻塞
}

func (trans *Transfer) ClientReadDataByProtocol ()  (uint64, []byte){
	data := make([]byte, 8) // 8字节数据
	trans.conn.Read(data)
	// 将requestID解码出来
	var requestID uint64
	_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &requestID)
	data = make([]byte, 4)
	trans.conn.Read(data)
	// 将数据len解码出来
	var len int
	_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &len)
	data = make([]byte, len)
	// 从conn读取数据
	trans.conn.Read(data)
	return requestID, data
}

func (trans *Transfer) ServerReadDataByProtocol () (uint64, []byte, error) {
	data := make([]byte, int32(len(ProtocolHead)))
	i, err := trans.conn.Read(data)
	if err != nil {
		return 0, nil, err
	}

	// 验证头部
	if string(i) != string(ProtocolHead) {
		return 0, nil, err
	}
	data = make([]byte, 8) // 初始化一个4字节的暂时存放数据的数组
	// 从read接受数据
	_, err = trans.conn.Read(data)
	if err != nil {
		return 0, nil, err
	}
	// 将requestID解码出来
	var requestID uint64
	_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &requestID)

	data = make([]byte, 4)
	// 从conn读取字节数组
	_, err = trans.conn.Read(data)
	if err != nil {
		return 0, nil, err
	}
	// 将数据len解码出来
	var len int
	_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &len)
	data = make([]byte, len)
	// 从conn读取数据
	_, err = trans.conn.Read(data)
	if err != nil {
		return 0, nil, err
	}

	return requestID,data, err
}

func (trans *Transfer) ServerWriteDataByProtocol(requestID uint64, data []byte)  {
	tempByteBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(tempByteBuffer, binary.BigEndian, requestID)
	_ = binary.Write(tempByteBuffer, binary.BigEndian, len(data))
	tempByteBuffer.Write(data)

	_, _ = trans.conn.Write(tempByteBuffer.Bytes())
}

func clientEnProto(oriData []byte) ([]byte, uint64) {
	proLen := int32(len(ProtocolHead))
	oriLen := int32(len(oriData))
	totalLen := int32(proLen + 8 + 4 + oriLen)
	newBytes := make([]byte, totalLen)
	copy(newBytes[0:proLen], ProtocolHead)
	// 写入requestID
	requestID := atomic.AddUint64(&RequestID, 1)
	tempByteBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(tempByteBuffer, binary.BigEndian, requestID)
	copy(newBytes[proLen:proLen + 8], tempByteBuffer.Bytes())
	// int32 转换为bytes
	tempByteBuffer = bytes.NewBuffer([]byte{})
	_ = binary.Write(tempByteBuffer, binary.BigEndian, oriLen)
	copy(newBytes[proLen + 8:proLen + 12], tempByteBuffer.Bytes())
	copy(newBytes[proLen + 12:], oriData)

	return newBytes, requestID
}