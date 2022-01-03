package ziface

// irequest 接口 实际上将客户端请求的连接信息和数据包装到一个request中

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到请求的数据
	GetData() []byte

	GetMsgId() uint32
}
