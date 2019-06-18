package server

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
	//@deprecated
	// EachReadBytes 每个读取的字节大小
	EachReadBytes = 500
)

var (
	ProtocolHead = []byte("srpc")
)

// Transfer 封装tcp conn
type Transfer struct {
	conn net.Conn
}

// NewTransfer 实例化新的传输管道
func NewTransfer(conn net.Conn) *Transfer {
	return &Transfer{conn: conn}
}

//@deprecated
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

func (trans *Transfer) ServerReadDataByProtocol () (uint64, []byte, error) {
	data := make([]byte, int32(len(ProtocolHead)))
	_, err := trans.conn.Read(data)
	if err != nil {
		return 0, nil, err
	}

	// 验证头部
	if string(data) != string(ProtocolHead) {
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
	var dataLen int32
	_ = binary.Read(bytes.NewBuffer(data), binary.BigEndian, &dataLen)
	data = make([]byte, dataLen)
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
	_ = binary.Write(tempByteBuffer, binary.BigEndian, int32(len(data)))
	tempByteBuffer.Write(data)

	_, _ = trans.conn.Write(tempByteBuffer.Bytes())
}