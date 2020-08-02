package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type ForwardProxy struct {

}

func (f *ForwardProxy)ServeHTTP(rw http.ResponseWriter,req *http.Request)  {
	fmt.Println("Requsest :",req.RemoteAddr,req.Host,req.Method)
	//默认转发器
	transport := http.DefaultTransport
	//复制请求内容，防止干扰之前部分
	nextReq := new(http.Request)
	*nextReq = *req
	//头部添加经过ip
	if clientIp,_,err := net.SplitHostPort(nextReq.RemoteAddr);err == nil {
		if oldForward,ok := nextReq.Header["X-Forwarded-For"];ok == true {
			clientIp = strings.Join(oldForward,",") + ", " + clientIp
		}
		nextReq.Header.Set("X-Forwarded-For",clientIp)
	}
	//调用下游得到返回值
	rsp,err := transport.RoundTrip(nextReq)
	if err != nil {
		fmt.Println("RoundTrip err:",err)
	}
	//将返回的头部和体部再次返回上游
	for key,value := range rsp.Header{
		for _,v := range value{
			rw.Header().Add(key,v)
		}
	}
	fmt.Println("resp :",rsp.StatusCode)
	rw.WriteHeader(rsp.StatusCode)
	_,_ = io.Copy(rw,rsp.Body)
	_ = rsp.Body.Close()
}

func main()  {
	//只能转发http，需要浏览器设置web代理
	http.Handle("/",&ForwardProxy{})
	err := http.ListenAndServe("0.0.0.0:8010",nil)
	if err != nil {
		fmt.Println("ListenAndServe err:",err)
	}
}
