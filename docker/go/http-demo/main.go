package main

import (
	// 导入net/http包
	"context"
	"http-demo/demov1"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	// ------------------ 使用http包启动一个http服务 方式一 ------------------
	// *http.Request http请求内容实例的指针
	// http.ResponseWriter 写http响应内容的实例
	http.HandleFunc("/v1/demo", func(w http.ResponseWriter, r *http.Request) {
		name, err := demoGrpcReq()
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		// 写入响应内容
		w.Write([]byte(name))
	})

	go (func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":6061", nil)
	})()

	// 启动一个http服务并监听8888端口 这里第二个参数可以指定handler
	http.ListenAndServe(":6060", nil)
}

func demoGrpcReq() (string, error) {
	conn, err := grpc.Dial("grpc-demo:1010", grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	// 泄露
	defer conn.Close()

	client := demov1.NewGreeterClient(conn)
	resp, err := client.SayHello(context.TODO(), &demov1.HelloRequest{
		Name: "http request",
	})
	if err != nil {
		return "", err
	}
	return resp.GetMessage(), nil
}
