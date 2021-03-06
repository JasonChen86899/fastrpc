package server

import (
	"errors"
	"fastrpc/common"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)


// Service 封装每个rpc注册服务的属性
type Service struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

// Server 封装服务者实例，每个服务者实例包含多个Service
type Server struct {
	ServiceMap  map[string]map[string]*Service
	serviceLock sync.Mutex
	ServerType  reflect.Type
}

// Register 注册服务者
func (server *Server) Register(obj interface{}) error {
	server.serviceLock.Lock()
	defer server.serviceLock.Unlock()

	//通过obj得到其各个方法，存储在servicesMap中
	tp := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	serviceName := reflect.Indirect(val).Type().Name()
	if _, ok := server.ServiceMap[serviceName]; ok {
		return errors.New(serviceName + " already registed.")
	}

	s := make(map[string]*Service)
	numMethod := tp.NumMethod()
	for m := 0; m < numMethod; m++ {
		service := new(Service)
		method := tp.Method(m)
		mtype := method.Type
		mname := method.Name

		service.ArgType = mtype.In(1)
		service.ReplyType = mtype.In(2)
		service.Method = method
		s[mname] = service
	}
	server.ServiceMap[serviceName] = s
	server.ServerType = reflect.TypeOf(obj)
	return nil
}

//@deprecated
// ServeConn 每个服务者处理每个rpc请求的入口函数
func (server *Server) ServeConn(conn net.Conn) {
	trans := NewTransfer(conn)
	for {
		// 从conn读数据
		data, err := trans.ReadData()
		if err != nil {
			return
		}

		// decode
		var req common.Request
		edcode, err := common.GetEdcode()
		if err != nil {
			return
		}
		err = edcode.Decode(data, &req)
		if err != nil {
			return
		}

		// 根据MethodName拿到service
		methodStr := strings.Split(req.MethodName, ".")
		if len(methodStr) != 2 {
			return
		}
		service := server.ServiceMap[methodStr[0]][methodStr[1]]

		// 构造argv
		argv, err := MakeArgs(&req, edcode, *service)

		// 构造reply
		reply := reflect.New(service.ReplyType.Elem())

		// 调用对应的函数
		function := service.Method.Func
		out := function.Call([]reflect.Value{reflect.New(server.ServerType.Elem()), argv, reply})
		if out[0].Interface() != nil {
			fmt.Println(out[0].Interface())
			return
		}

		// encode
		replyData, err := edcode.Encode(reply.Elem().Interface())
		if err != nil {
			return
		}

		// 向conn写数据
		_, err = trans.WriteData(replyData)
		if err != nil {
			return
		}
	}
}

// Server 开启服务者
func (server *Server) Server(network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("net.Listen tcp :0: %v", err)
		return err
	}

	for {
		// 阻塞直到收到一个网络连接
		conn, e := l.Accept()
		if e != nil {
			log.Fatalf("l.Accept: %v", e)
		}

		//开始工作
		//go server.ServeConn(conn)
		go server.NewServerConn(conn) // 新处理函数
	}
}

// NewServer 实例化一个服务者
func NewServer() *Server {
	return &Server{
		ServiceMap:  make(map[string]map[string]*Service),
		serviceLock: sync.Mutex{}}
}

func (server *Server) NewServerConn(conn net.Conn)  {
	trans := NewTransfer(conn)
	for {
		// 长链接轮训处理
		requestID, data, err := trans.ServerReadDataByProtocol()
		if err != nil {
			log.Println("read each request err")
			continue
		}
		server.handlerEachDataBytes(requestID, data, trans)
	}
}
func (server *Server) handlerEachDataBytes(requestID uint64, data []byte, trans *Transfer)  {
	var req common.Request
	edcode, err := common.GetEdcode()
	if err != nil {
		return
	}
	err = edcode.Decode(data, &req)
	if err != nil {
		return
	}

	// 根据MethodName拿到service
	methodStr := strings.Split(req.MethodName, ".")
	if len(methodStr) != 2 {
		return
	}
	service := server.ServiceMap[methodStr[0]][methodStr[1]]

	// 构造argv
	argv, err := MakeArgs(&req, edcode, *service)

	// 构造reply
	reply := reflect.New(service.ReplyType.Elem())

	// 调用对应的函数
	function := service.Method.Func
	out := function.Call([]reflect.Value{reflect.New(server.ServerType.Elem()), argv, reply})
	if out[0].Interface() != nil {
		fmt.Println(out[0].Interface())
		return
	}

	// encode
	replyData, err := edcode.Encode(reply.Elem().Interface())
	if err != nil {
		return
	}

	// 向conn写数据
	trans.ServerWriteDataByProtocol(requestID, replyData)
}