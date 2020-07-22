package main

import (
	"fmt"
	"net"
)

func main() {
	//直接得到连接，无需接受连接请求
	conn,err := net.ListenUDP("udp",&net.UDPAddr{
		IP:net.IPv4(127,0,0,1),
		Port: 8002,
	})
	if err != nil {
		fmt.Println("Listen err:",err)
		return
	}
	//执行读取
	process(conn)
}

func process(conn *net.UDPConn) {
	//循环读取
	for {
		defer func() {
			//可以不关闭
			//_ = conn.Close()
		}()
		//每次读取1024字节
		temp := [1024]byte{}
		//进行读取,要通过UDP方式读取，才能回复
		n,clientAddr,err := conn.ReadFromUDP(temp[:])
		if err != nil {
			//前端连接关闭时，这里会报错
			//前端进程关闭：错误提示强制关闭
			//前端close关闭：提示EOF
			fmt.Println("Server Read err:",err)
			break
		}
		//打印读取内容，读多少打印多少
		fmt.Println("Server Read Str:",string(temp[:n]),clientAddr,n)
		go do(conn,temp[:],clientAddr)
	}
}

func do(conn *net.UDPConn,buf []byte,clientAddr *net.UDPAddr){
	//进行回复，必须要携带客户端地址，不然可能连接关闭了
	n,err := conn.WriteToUDP([]byte("udp resp"),clientAddr)
	if err != nil {
		fmt.Println("Server Write err:",err)
		return
	}
	fmt.Println("Server Write Num:",n)
}