package main

/*
使用zinx矿建实现一个应用
*/
import (
	"fmt"
	"zinxProject/ziface"
	"zinxProject/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	//先读取客户端的数据，再回写ping。。。
	fmt.Println("recv from client :msgID:", request.GetMsgId(), "msgData:", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}

}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	//先读取客户端的数据，再回写ping。。。
	fmt.Println("recv from client :msgID:", request.GetMsgId(), "msgData:", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx V0.6"))
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	//创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.2]") //新建了一个服务器，但是路由为空
	//当用户发送id为0的消息时，添加Pingrouter
	s.AddRouter(0, &PingRouter{})
	//当用户发送id为1的消息，添加hellorouter
	s.AddRouter(1, &HelloZinxRouter{})
	//启动server
	s.Server()
}
