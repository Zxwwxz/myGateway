package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP (rsp http.ResponseWriter,req *http.Request) {
	f(rsp,req)
}

func main()  {
	hf := HandlerFunc(hello)
	rsp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte("test")))
	hf.ServeHTTP(rsp,req)
	bodyBuf, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(bodyBuf))
}

func hello(rsp http.ResponseWriter, req *http.Request) {
	rsp.Write([]byte("Hello world"))
}
