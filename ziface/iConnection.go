package ziface

import "net"

//
type IConnection interface {
	//启动连接，让当前的连接准备工作
	Start()
	//停止链接，结束当前链接的工作
	Stop()
	//获取当前链接绑定的socket conn（套接字）
	GetTCPConnection() *net.TCPConn
	//获取当前链接模块的链接ID
	GetConnID() uint32
	//获取远程客户端的TCP状态ip port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给客户端
	SendMsg(msgId uint32, data []byte) error
}

//定义一个处理链接业务的函数指针
type HandleFunc func(*net.TCPConn, []byte, int) error
