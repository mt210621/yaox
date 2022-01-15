package znet

import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	"fmt"
	"net"
)

//IServer 的接口实现  定义一个server的服务模块

type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//当前server 的消息管理模块，用来绑定msgid和对应的业务api业务
	MsgHandler ziface.IMsgHandle
	// 该server的连接管理器
	ConnMgr ziface.IConnManager
	// server创建连接之后自动调用hook函数 -- onconnstarrt
	OnConnStart func(conn ziface.IConnection)
	// 该server销毁连接之前自动调用的hook函数--onconnstop
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[zinx]server name : %s,listener ai ip: %s,port:%d is starting\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n\n", s.IP, s.Port)
	go func() {
		// 0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()
		// 获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error:", err)
			return
		}
		//监听这个服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start zinx server succ,", s.Name, "suc,Listenning...")
		var cid uint32
		cid = 0
		//阻塞等待客户端连接  处理客户端连接业务
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接个数的判断  如果超过 则关闭此链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给用户客户端 响应一个最大连接的错误包
				fmt.Println("to many connections , maxconn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法  和conn进行绑定 得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 将服务器资源 状态或者一些已经开辟的资源停止
	fmt.Println("[STOP] zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	//启动server的服务的功能
	s.Start()

	//阻塞状态
	select {}
}

// AddRouter 路由功能 给当前服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ !!")
}

//
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// NewServer 初始化server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// SetOnConnStart 注册onconnstart 钩子函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册onconnstop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用onconnstart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用onconnstop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----> call onconnstop...")
		s.OnConnStop(conn)
	}
}
