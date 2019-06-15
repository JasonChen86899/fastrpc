package test

import (
	"fastrpc/client"
	"fastrpc/server"
	"fmt"
	"testing"
	"time"
)

func Test_Fastrpc(t *testing.T) {
	serverTest := server.NewServer()
	go serverTest.Server("tcp", "127.0.0.1:9999")
	time.Sleep(time.Second)
	clientTest, err := client.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	clientTest.SendReqByProtocol("aaaa", "", "")
}
