package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main()  {
	//连接池
	transport := &http.Transport{
		DialContext:(&net.Dialer{
			Timeout:   30 * time.Second, //连接超时
			KeepAlive: 30 * time.Second, //探活时间
		}).DialContext,
		MaxIdleConns:          100,              //最大空闲连接
		IdleConnTimeout:       90 * time.Second, //空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
	}
	client := http.Client{
		Timeout  :time.Second * 30, //请求超时时间
		Transport:transport,
	}
	//发送请求
	rsp,err := client.Get("http://127.0.0.1:8004/hello")
	if err != nil {
		fmt.Println("Get err:",err)
	}
	fmt.Println("Status:",rsp.Status)
	defer rsp.Body.Close()
	//读取内容
	buf,err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("ReadAll err:",err)
	}
	fmt.Println("body:",string(buf))
}
