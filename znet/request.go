package znet

import "zinxProject/ziface"

type Request struct {
	Conn ziface.IConnection //这个属性是一个接口
	msg  ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
