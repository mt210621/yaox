package znet

import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	"fmt"
	"strconv"
)

// MsgHandle 消息处理模块的实现
type MsgHandle struct {
	// 存放每一个msgid 所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
		//从全局配置中 获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度 /执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从request中找到msgid
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgid = ", request.GetMsgId(), "is not found ! need register!")
	}
	// 2 根据msgid 调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断 当前msg绑定的api处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册了
		panic("repeat api,msgId = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与api的关系
	mh.Apis[msgID] = router
	fmt.Println("add api msgid = ", msgID, "succ!")
}

// StartWorkerPool 启动一个worker 工作池 (开启工作池的动作只能发生一次，一个zinx框架只能有一个work工作池）
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerpoolsize 分别开启worker  每个worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 当前的worker对应的channel消息队列 开辟空间 第0个worker 就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 尝试启动当前的worker 阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启用一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker ID = ", workerID, "is started ...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息 过来，出列的就是一个客户端的request，执行当前的request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给taskqueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不通过的worker

	//根据客户端建立的connid来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connid = ", request.GetConnection().GetConnID(), "request msgid = ", request.GetMsgId(), "to workerid = ", workerID)

	// 2 将消息发送给对应的worker的taskqueue即可
	mh.TaskQueue[workerID] <- request
}
