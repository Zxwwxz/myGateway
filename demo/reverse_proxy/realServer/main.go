package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//开启服务器1
	rs1 := &RealServer{Addr: "127.0.0.1:8012"}
	rs1.Run()
	//开启服务器2
	rs2 := &RealServer{Addr: "127.0.0.1:8013"}
	rs2.Run()
	//监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<- quit
}

type RealServer struct {
	Addr string
}

func (r *RealServer) Run(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		//协程开启服务器
		log.Fatal(server.ListenAndServe())
	}()
}

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	//http://127.0.0.1:8008/abc?def=11
	fmt.Println("Host:",req.Host)
	fmt.Println("RemoteAddr:",req.RemoteAddr)
	fmt.Println("RequestURI:",req.RequestURI)
	fmt.Println("Method:",req.Method)
	fmt.Println("Header:",req.Header)
	fmt.Println("URL:",*req.URL)
	fmt.Println("Body:",req.Body)
	upath := fmt.Sprintf("http://%s%s\n", r.Addr,req.RequestURI)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
}

func (r *RealServer) ErrorHandler(rw http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	rw.WriteHeader(500)
	io.WriteString(rw, upath)
}
