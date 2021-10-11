package ziface

//定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Server()
	//给当前的服务注册一个路由方法，供当前的客户端链接使用
	AddRouter(msgID uint32, router IRouter)
	//定义一个获得connManager的方法
	GetConnMgr() IConnManager

	//注册OnConnStart函数的方法
	SetOnConnStart(func(connection IConnection))
	//注册OnConnStop函数的方法
	SetOnConnStop(func(connection IConnection))

	//调用OnConnStop函数的方法
	CallOnConnStart(connection IConnection)
	//调用OnConnStart函数的方法
	CallOnConnStop(connection IConnection)
}
