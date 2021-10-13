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
	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx V0.8"))
	if err != nil {
		fmt.Println(err)
	}

}

//开发者定义两个hook函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("===>DoConnBegin is Called!...")
	err := conn.SendMsg(202, []byte("DoConnBegin"))
	if err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	fmt.Println("Set conn Name...")
	conn.SetProperty("Name", "冰冰的小圆脸")
	conn.SetProperty("git Addr", "https://github.com/ppnnssy/Zinx.git")

}

//结束后的钩子函数
func DoConnLost(conn ziface.IConnection) {
	fmt.Println("===>DoConnLost is Called!...")
	fmt.Println("connId:", conn.GetConnID(), "is lost!")

	//结束时获取一下链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name:", name)
	}
	if add, err := conn.GetProperty("git Addr"); err == nil {
		fmt.Println("add:", add)
	}

}



func main() {
	//创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.7]") //新建了一个服务器，但是路由为空

	//注册hook钩子函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	//当用户发送id为0的消息时，添加Pingrouter
	s.AddRouter(0, &PingRouter{})
	//当用户发送id为1的消息，添加hellorouter
	s.AddRouter(1, &HelloZinxRouter{})


	//启动server
	s.Server()
}
