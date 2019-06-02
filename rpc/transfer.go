package rpc

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
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

func (trans *Transfer) WriteDataByProtocol(data []byte) (int, error) {
	return trans.WriteData(enProto(data))
}

func (trans *Transfer) ReadDataByProtocol () error {
	for {
		data := make([]byte, int32(len(ProtocolHead)))
		i, err := trans.conn.Read(data)
		if err != nil {
			return err
		}

		if string(i) == string(ProtocolHead) {
			data = make([]byte, 4)
			if err != nil {
				return err
			}

			var l int
			_ := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &l)

			data = make([]byte, l)
			trans.conn.Read(data)

		} else {
			return err
		}
	}
}

func enProto(oriData []byte) []byte {
	proLen := int32(len(ProtocolHead))
	oriLen := int32(len(oriData))
	totalLen := int32(proLen + 4 + oriLen)
	newBytes := make([]byte, totalLen)
	copy(newBytes[0:proLen], ProtocolHead)
	// int32 转换为bytes
	tempByteBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(tempByteBuffer, binary.BigEndian, oriLen)
	copy(newBytes[proLen:proLen + 4], tempByteBuffer.Bytes())
	copy(newBytes[proLen + 4:], oriData)

	return newBytes
}