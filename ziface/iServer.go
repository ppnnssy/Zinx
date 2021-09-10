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
}
