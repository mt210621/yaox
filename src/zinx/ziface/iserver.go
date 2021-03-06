package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 路由功能 给当前的服务注册一个路由方法 供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	// GetConnMgr 获取当前server的连接管理器
	GetConnMgr() IConnManager
	// SetOnConnStart 注册onconnstart 钩子函数
	SetOnConnStart(func(connection IConnection))
	// SetOnConnStop 注册onconnstop钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	// CallOnConnStart 调用onconnstart钩子函数的方法
	CallOnConnStart(connection IConnection)
	// CallOnConnStop 调用onconnstop钩子函数的方法
	CallOnConnStop(connection IConnection)
}
