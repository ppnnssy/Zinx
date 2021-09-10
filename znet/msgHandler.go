package znet

import (
	"fmt"
	"strconv"
	"zinxProject/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	//存放每个msgIdsuo对应的处理方法
	Apis map[uint32]ziface.IRouter
}

//初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
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
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加对应的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api,msgId=" + strconv.Itoa(int(msgId)))
	}

	mh.Apis[msgId] = router

}
