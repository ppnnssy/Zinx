package znet

import (
	"fmt"
	"net"
	"zinxProject/utils"
	"zinxProject/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器绑定的监听ip
	IP string
	//服务器监听的端口
	Port int
	//添加一个Router对象
	Router ziface.IRouter
}

//给当前的服务注册一个路由方法，供当前的客户端链接使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Success")
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[zinx] Server Name:%s;listenner at IP:%s;Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[zinx]Version:%s;MaxConn:%d,MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	//把所有业务放到go中做，不会无限阻塞
	go func() {
		//1.获取一个tcp的addr
		//ResolveIPAddr将addr作为一个格式为"host"或"ipv6-host%zone"的IP地址来解析。
		//函数会在参数net指定的网络类型上解析，net必须是"ip"、"ip4"或"ip6"。
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error:", err)
			return
		}
		//2.尝试监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen:", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start zinx server success,", s.Name, "success,Listenning...")

		//定义一个变量，记录链接的编号
		var cid uint32
		cid = 0

		//3.阻塞式的等待客户端连接，处理客户端的连接业务
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//到此为止客户端已经建立连接，做一些业务。暂时做一个最基本的最大512字节的回显业务
			//初始化链接，绑定链接conn和业务CallBackToClient
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动当前的链接业务
			go dealConn.Start()

		}

	}()

}

//停止服务器
//本函数用于将一些服务器的资源，状态或者已经开辟的链接信息进行回收或停止
func (s *Server) Stop() {

}

//运行服务器
func (s *Server) Server() {
	//启动server的服务功能
	s.Start()

	//因为Start（）所有的服务都在go中执行，main进程结束后go也会提前结束，所以需要阻塞一下主进程
	select {}
}

//提供一个初始化Server模块的方法(工厂模式）
func NewServer(name string) *Server {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil, //先默认为空，实际项目中调用自己重写的路由方法
	}
	return s
}
