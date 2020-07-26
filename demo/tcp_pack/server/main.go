package main

import (
	"fmt"
	"myGateway/demo/tcp_pack/pack"
	"net"
)

func main()  {
	//监听
	listen, err := net.Listen("tcp",":8002")
	if err != nil {
		fmt.Println("Listen err:",err)
		return
	}
	//每个客户端一个连接
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept err:",err)
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn){
	defer conn.Close()
	//一个连接的每条消息
	for {
		//反序列化
		msgBuf,err := pack.Decode(conn)
		if err != nil {
			fmt.Println("Decode err:",err)
			break
		}
		do(msgBuf)
	}
}

func do(msgBuf []byte) {
	fmt.Println("server:",string(msgBuf))
}