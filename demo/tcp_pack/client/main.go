package main

import (
	"fmt"
	"myGateway/demo/tcp_pack/pack"
	"net"
)

func main()  {
	//连接
	conn,err := net.Dial("tcp",":8002")
	if err != nil {
		fmt.Println("Dial err:",err)
		return
	}
	sendMsg := "hello"
	//序列化
	err = pack.Encode(conn,[]byte(sendMsg))
	if err != nil {
		fmt.Println("Encode err:",err)
		return
	}
	_ = conn.Close()
}
