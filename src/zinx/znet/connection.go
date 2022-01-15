package znet

import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

/*
 连接模块
*/
type Connection struct {
	// 当前connection 隶属于哪个server
	TcpServer ziface.IServer
	//当前 连接的socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//告知当前连接已经推出的 停止的channel
	ExitChan chan bool

	// 无缓冲通道 用于 读 写goroutine 之间的消息通信
	msgChan chan []byte

	//消息的管理msgid 和对应的处理业务的api关系
	MsgHandler ziface.IMsgHandle

	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化 连接模块
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}
	// 将conn加入到connmanager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//连接的读业务的方法
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running...")
	defer fmt.Println("connid =", c.ConnID, "reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buff中  最大的512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		//创建一个拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的Msg head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包 得到msgid和datalen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据datalen 再次读取data 放在msg.data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
			}
		}
		msg.SetData(data)
		//得到当前conn数据的request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制 将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//根据绑定好的id找到对应处理的业务
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 写 消息goroutine 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer gortine is running....]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")
	//不断的阻塞的等待channel的消息  进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error,", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader已经退出  此时writer也要退出
			return
		}
	}
}

// Start 启动连接 让当前的连接 开始工作
func (c *Connection) Start() {
	fmt.Println("conn start()...connid = ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()
	go c.StartWriter()
	// 按照开发者传递进来的 创建连接之后 需要调用的业务  执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接 结束当前连接的工作

func (c *Connection) Stop() {
	fmt.Println("conn stop().. connid = ", c.ConnID)
	//如果当前已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 调用开发者注册的 销毁连接之前 需要执行的业务hook函数
	c.TcpServer.CallOnConnStop(c)
	//关闭 socket连接
	c.Conn.Close()

	//告知writer关闭
	c.ExitChan <- true
	// 将当前连接从connmgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的id

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的tcp状态 ip port

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 提供一个sendmsg方法 将我们要发送给客户端的数据  先进行封包 再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data进行封包 msgdatalen|msgid|data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
