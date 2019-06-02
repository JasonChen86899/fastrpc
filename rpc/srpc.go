package rpc

import (
	"errors"
	"net"
	"sync"
)

// NewServer 实例化一个服务者
func NewServer() *Server {
	return &Server{
		ServiceMap:  make(map[string]map[string]*Service),
		serviceLock: sync.Mutex{}}
}

// NewClient 实例化一个客户端调用者
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
		lock: sync.Mutex{}}
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
