package main

import (
	"fmt"
	"net"
)

func main (){
	//监听端口
	listen,err := net.ListenTCP("tcp",&net.TCPAddr{
		IP:net.IPv4(127,0,0,1),
		Port: 8001,
	})
	if err != nil {
		fmt.Println(" Server Listen err:",err)
		return
	}
	for {
		//建立连接,每个连接是一个客户端
		conn,err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("Server Accept err:",err)
			continue
		}
		//开启协程
		go process(conn)
	}
}

func process(conn net.Conn)  {
	//循环读取，每次收到是当前客户端的一次发送
	for {
		defer func() {
			//不关闭会出现close_wait
			_ = conn.Close()
		}()
		//每次读取1024字节
		temp := [1024]byte{}
		//进行读取
		n,err := conn.Read(temp[:])
		if err != nil {
			//前端连接关闭时，这里会报错
			//前端进程关闭：错误提示强制关闭
			//前端close关闭：提示EOF
			fmt.Println("Server Read err:",err)
			break
		}
		//打印读取内容，读多少打印多少
		fmt.Println("Server Read Str:",string(temp[:n]),n)
		//执行具体逻辑
		go do(conn,temp[:])
	}
}

func do(conn net.Conn,buf []byte){
	//进行回复
	n,err := conn.Write([]byte("tcp resp"))
	if err != nil {
		fmt.Println("Server Write err:",err)
		return
	}
	fmt.Println("Server Write Num:",n)
}
