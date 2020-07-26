package main

import (
	"fmt"
	"net/http"
)

func main()  {
	//新建路由器
	mux := http.NewServeMux()
	//设置函数映射
	mux.HandleFunc("/hello",hello)
	//新建服务
	server := http.Server{
		Addr:":8004",
		Handler:mux,
	}
	//开启监听
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe err:",err)
	}
}

func hello(rsp http.ResponseWriter,req *http.Request)  {
	fmt.Println("RemoteAddr:",req.RemoteAddr)
	fmt.Println("RequestURI:",req.RequestURI)
	fmt.Println("Method:",req.Method)
	fmt.Println("Host:",req.Host)
	fmt.Println("Body:",req.Body)
	_,err := rsp.Write([]byte("hello rsp"))
	if err != nil {
		fmt.Println("Write err:",err)
	}
}
