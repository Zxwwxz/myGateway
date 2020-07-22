package main

import (
	"fmt"
	"net"
)

func main()  {
	//连接服务器
	conn,err := net.DialTCP("tcp",nil,&net.TCPAddr{
		IP:net.IPv4(127,0,0,1),
		Port: 8001,
	})
	if err != nil {
		fmt.Println("Dial err:",err)
		return
	}
	//进行发送
	n,err := conn.Write([]byte("tcp req"))
	if err != nil {
		fmt.Println("Client Write err:",err)
		return
	}
	fmt.Println("Client Write num:",n)
	temp := [1024]byte{}
	//进行读取
	n,err = conn.Read(temp[:])
	if err != nil {
		fmt.Println("Client Read err:",err)
	}
	fmt.Println("Client Read str:",string(temp[:n]),n)
	_ = conn.Close()
}
