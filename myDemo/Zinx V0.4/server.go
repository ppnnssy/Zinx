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

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping"))
	if err != nil {
		fmt.Println("Call Router PreHandle err", err)
	}
}
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("pinging..."))
	if err != nil {
		fmt.Println("Call Router Handle err", err)
	}

}
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping"))
	if err != nil {
		fmt.Println("Call Router PostHandle err", err)
	}
}

func main() {
	//创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.2]") //新建了一个服务器，但是路由为空
	//给服务器添加一个路由方法添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//启动server
	s.Server()
}
