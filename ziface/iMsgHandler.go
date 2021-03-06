package ziface

type IMsgHandle interface {
	//调度、执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加对应的处理逻辑
	AddRouter(msgId uint32, router IRouter)

	//开启一个工作池
	StarWorkerPool()
	//发送request到消息队列
	SendMsgToTaskQueue(request IRequest)
}
