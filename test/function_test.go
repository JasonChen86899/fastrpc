package test

import (
	"fastrpc/client"
	"fastrpc/server"
	"fmt"
	"sync"
	"testing"
	"time"
)

type ServerTest struct {
}

type Args struct {

}

func (s *ServerTest) Print1(args Args, reValue *string ) interface{} {

	*reValue = "success_1"
	return nil
}

func (s *ServerTest) Print2(args Args, reValue *string ) interface{} {

	*reValue = "success_2"
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

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(7)
	go showReturnValue1(clientTest, &waitGroup)
	go showReturnValue2(clientTest, &waitGroup)
	go showReturnValue1(clientTest, &waitGroup)
	go showReturnValue2(clientTest, &waitGroup)
	go showReturnValue1(clientTest, &waitGroup)
	go showReturnValue2(clientTest, &waitGroup)

	var s string
	clientTest.SendReqByProtocol("ServerTest.Print1", Args{}, &s)
	waitGroup.Done()
	fmt.Println(s)

	waitGroup.Wait()
}

func showReturnValue1(client *client.Client, waitGroup *sync.WaitGroup)  {
	var s string
	client.SendReqByProtocol("ServerTest.Print1", Args{}, &s)
	fmt.Println(s)
	waitGroup.Done()
}

func showReturnValue2(client *client.Client, waitGroup *sync.WaitGroup)  {
	var s string
	client.SendReqByProtocol("ServerTest.Print2", Args{}, &s)
	fmt.Println(s)
	waitGroup.Done()
}