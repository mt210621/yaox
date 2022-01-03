package main

import (
	"awesomeProject/src/zinx/ziface"
	"awesomeProject/src/zinx/znet"
	"fmt"
)

/*
基于 zinx框架开发的服务端应用程序
*/

// ping test 自定义路由

type PingRouter struct {
	znet.BaseRouter
}

// PreHandle test prehandle
//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call router prehandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}

// Handle  test handle
//func (this *PingRouter) Handle(request ziface.IRequest) {
//	fmt.Println("Call router handle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping-ping...\n"))
//	if err != nil {
//		fmt.Println("call back ping error")
//	}
//}

// PostHandle test posthandle
//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call router posthandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping\n"))
//	if err != nil {
//		fmt.Println("call post ping error")
//	}
//}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call PingRouter-router handle...")
	//先读取客户端的数据 再回写ping...ping...ping
	fmt.Println("recv from client: msgId:", request.GetMsgId(), "data:", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("call HelloZinxRouter-router handle...")
	//先读取客户端的数据 再回写ping...ping...ping
	fmt.Println("recv from client: msgId:", request.GetMsgId(), "data:", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello zinx welcome to zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1 创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx v0.1]")

	//2 给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//3 启动server
	s.Serve()
}
