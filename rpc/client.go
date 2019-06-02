package rpc

import (
	"log"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
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

// Call rpc客户端发起远程函数请求(rpc)的入口函数
func (client *Client) Call(methodName string, req interface{}, reply interface{}) error {
	// 开启协程锁
	client.lock.Lock()
	defer client.lock.Unlock()

	// 构造一个Request
	request := NewRequest(methodName, req)

	// encode
	edcode, err := GetEdcode()
	if err != nil {
		return err
	}
	data, err := edcode.encode(request)
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
	edcode.decode(data2, reply)

	// return
	return nil
}

var CallBackMap = new(sync.Map)
var requestID int64= 0

func (client *Client) CallByCallBack(methodName string, req interface{}, reply interface{}) error {
	eachCallChannel := make(chan interface{})
	CallBackMap.Store(atomic.AddInt64(&requestID, 1), eachCallChannel)
	client.sendReqByProtocol(methodName, req)
	eachCallChannel <- *reply.(unsafe.Pointer)
}

func (client *Client) sendReqByProtocol(methodName string, req interface{})  {
	trans := NewTransfer(client.conn)

}
