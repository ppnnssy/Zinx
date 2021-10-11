package znet

import (
	"fmt"
	"strconv"
	"zinxProject/utils"
	"zinxProject/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	//存放每个msgId所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的Worker数量
	WorkerPoolSize uint32
}

//初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:      make(map[uint32]ziface.IRouter),
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		//业务工作Worker池的Worker数量
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

//调度、执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//检查是否存在对应的调度方法
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId", request.GetMsgId(), "is not found.need register!")
		return
	}

	//存在的话就调用
	handler.PreHandle(request)
	handler.Handle(request) //目前只重写并调用这个
	handler.PostHandle(request)
}

//为消息添加对应的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api,msgId=" + strconv.Itoa(int(msgId)))
	}

	mh.Apis[msgId] = router

}

//启动一个Worker工作池
//开启工作池的动作只能发生一次
func (mh *MsgHandle) StarWorkerPool() {
	//根据WorkerPoolSize，分别开启Worker，每个Worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1.给当前的Worker对应的channel消息队列开辟空间，第0个worker就用第0个channel。。。
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2.启动当前的worker，阻塞等待消息从chan中传递过来
		go mh.StarOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个Worker工作流程
//对外不暴露
func (mh *MsgHandle) StarOneWorker(workerId int, taskQueue chan ziface.IRequest) { //workerId记录消息队列中第几个消息。
	fmt.Println("WorkerID:", workerId, "is started")
	//阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}

}

// SendMsgToTaskQueue 将消息传递给TaskQueue，由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 将消息平均分配给不同的Worker
	//根据客户端建立的ConnID分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize //取余
	fmt.Println("Add ConnID:", request.GetConnection().GetConnID(),
		" request MsgID:", request.GetMsgId(),
		"to WorkID:", workerID)

	//2 将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
