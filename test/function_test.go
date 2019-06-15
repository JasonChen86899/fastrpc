package test

import (
	"fastrpc/client"
	"fastrpc/server"
	"fmt"
	"testing"
	"time"
)

type ServerTest struct {
}

type Args struct {

}

func (s *ServerTest) Print(args Args, reValue *string ) interface{} {

	*reValue = "success"
	return nil
}

func Test_Fastrpc(t *testing.T) {
	serverTest := server.NewServer()
	go serverTest.Server("tcp", "127.0.0.1:9999")
	time.Sleep(time.Second)
	serverTest.Register(new(ServerTest))
	time.Sleep(time.Second)
	clientTest, err := client.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	var s string
	clientTest.SendReqByProtocol("ServerTest.Print", Args{}, &s)
	fmt.Println(s)
}
