package utils

import (
	"awesomeProject/src/zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

// 存储一切有关zinx的框架的全局参数的对象  一些参数是可以通过zinx.json由用户进行配置

type GlobalObj struct {
	// server
	TcpServer ziface.IServer // 全局zinx的server对象
	Host      string         // 当前服务器主机监听的ip
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	// zinx
	Version          string //当前 zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前zinx框架数据包的最大值
	WorkerPoolSize   uint32 //工作池的大小
	MaxWorkerTaskLen uint32 //允许用户最多开辟的worker

}

// GlobalObject 定义一个全局的对外对象 globalogj
var GlobalObject *GlobalObj

// Reload 从zinx.json 中加载用户自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法用于初始化 当前的参数
func init() {
	//如果配置文件没有加载  默认的值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}
