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

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("===> doconnectionbegin is called...")
	if err := conn.SendMsg(202, []byte("doconnection begin!")); err != nil {
		fmt.Println(err)
	}
	// 给当前的连接 设置一些属性
	fmt.Println("set conn name,hoe....")
	conn.SetProperty("Name", "L.G.X")
	conn.SetProperty("Age", "23")
}

// DoConnectionLost 连接断开之前需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("===> doconnectionLOST is called...")
	fmt.Println("===> connID = ", conn.GetConnID(), "is lost...")

	// 获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("name : ", name)
	}
	if age, err := conn.GetProperty("Age"); err == nil {
		fmt.Println("name : ", age)
	}
}

func main() {
	//1 创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx v0.1]")

	// 注册连接 hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//2 给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//3 启动server
	s.Serve()
}
