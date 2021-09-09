package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinxProject/ziface"
)

type Connection struct {
	//绑定的链接
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前链接的状态
	IsClosed bool

	//告知当前链接已经停止/退出
	ExitChan chan bool

	//该链接处理的方法Router 也就是绑定了连接和路由方法
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

/*
实现接口中的各种方法

*/

// StartReader 读数据
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID", c.ConnID, "Reader is exit,remote addr is\n", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//buf:=make([]byte,utils.GlobalObject.MaxPackageSize)
		//_,err:=c.Conn.Read(buf)
		//if err!=nil{
		//	fmt.Println("recv buf err:",err)
		//	break
		//}

		//创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的MsgHead 二进制流，8个字节
		headData := make([]byte, dp.GetHeadLen()) //就是定义headData是8字节
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		//拆包得到MsgID和msgDatalen放到msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}

		//根据datalen再次读取data，放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err = io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetMsgData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			Conn: c,
			msg:  msg,
		}

		//调用路由，从路由中找到绑定的conn对应的Router
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}

}

//启动连接，让当前的连接准备工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID:", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()

	//启动从当前链接的写数据业务

}

//停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop()...ConnID:\n", c.ConnID)
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	c.Conn.Close()
	close(c.ExitChan)
}

//获取当前链接绑定的socket conn（套接字）
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr() //RemoteAddr returns the remote network address.
}

//发送数据，将数据发送给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("conn closed when send msg")
	}
	//将data进行封包
	dp := NewDataPack()

	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack msg error,msgId:", msgId)
		return errors.New("Pack error msg")
	}
	_, err = c.Conn.Write(msg)
	if err != nil {
		fmt.Println("Write msg id:", msgId, "error", err)
		return errors.New("conn write error")
	}

	return nil

}
