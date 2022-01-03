package main

import (
	"awesomeProject/src/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
 模拟客户端
*/
func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	// 直接连接远程服务器  得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}
	for {
		////连接调用write 写data
		//_, err := conn.Write([]byte("Hello Zinx V0.1 ... "))
		//if err != nil {
		//	fmt.Println("write conn err")
		//	return
		//}
		//buf := make([]byte, 512)
		//cnt, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf error")
		//	return
		//}
		//fmt.Printf("server call back:%s,cnt=%d\n\n", buf, cnt)

		// 发送封包消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinxv0.5 client test message")))
		if err != nil {
			fmt.Println("pack error", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//服务器应该给客户端一个message数据  msgid： pingpingping
		//先读取流中的head部分
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}
		//将二进制的head拆包到msg结构中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msghead error", err)
			break
		}
		//再次读取数据
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("-----> recv server msg : id=", msg.Id, "len = ", msg.DataLen, "data = ", string(msg.Data))
		}
		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}

}
