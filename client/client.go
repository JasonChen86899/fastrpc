package client

import (
	"errors"
	"fastrpc/common"
	"log"
	"net"
	"sync"
)

// Client 封装rpc客户端
type Client struct {
	conn net.Conn
	lock sync.Mutex
}

// Close 关闭客户端连接
func (client *Client) Close() {
	client.conn.Close()
}

//@deprecated
// Call rpc客户端发起远程函数请求(rpc)的入口函数
func (client *Client) Call(methodName string, req interface{}, reply interface{}) error {
	// 开启协程锁
	client.lock.Lock()
	defer client.lock.Unlock()

	// 构造一个Request
	request := common.NewRequest(methodName, req)

	// encode
	edcode, err := common.GetEdcode()
	if err != nil {
		return err
	}
	data, err := edcode.Encode(request)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// write
	trans := NewTransfer(client.conn)
	_, err = trans.WriteData(data)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//fmt.Println("** start read return value **")
	// read
	data2, err := trans.ReadData()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//fmt.Println("** end read return value **")
	// decode and assin to reply
	edcode.Decode(data2, reply)

	// return
	return nil
}

func (client *Client) SendReqByProtocol(methodName string, req interface{}, res interface{})  {
	trans := NewTransfer(client.conn)
	request := common.NewRequest(methodName, req)
	edCode, _ := common.GetEdcode()
	reqBytesData, _ := edCode.Encode(request)
	resByteData := trans.ClientWriteDataByProtocol(reqBytesData)
	_ = edCode.Decode(resByteData.([]byte), res)
}

// 内部函数
func (client *Client) loopReceiveResponseByProtocol()  {
	trans := NewTransfer(client.conn)
	for {
		requestID, bytes := trans.ClientReadDataByProtocol()
		if channel ,ok := ReqResMapping.Load(requestID); ok {
			channel.(chan interface {}) <- bytes
		} else {
			log.Println("requestID has no channel")
		}
	}
}

// NewClient 实例化一个客户端调用者
func NewClient(conn net.Conn) *Client {
	client := &Client{
		conn: conn,
		lock: sync.Mutex{}}
	// 开启 长连接循环接受返回值
	go client.loopReceiveResponseByProtocol()
	return client
}

// Dial rpc客户端向服务者建立tcp连接
func Dial(network, address string) (*Client, error) {
	if network != "tcp" {
		return nil, errors.New("Unsupported protocol")
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return NewClient(conn), nil
}
