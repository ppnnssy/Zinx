package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinxProject/utils"
	"zinxProject/ziface"
)

type Connection struct {
	//以后框架变大了之后会有多个服务器，所以需要一个属性判断当前链接属于哪个server
	//目前可以通过本字段访问connManager 本条注释会写多次因为我记不住
	TcpServer ziface.IServer
	//绑定的链接
	Conn *net.TCPConn
	//链接ID
	ConnID uint32
	//当前链接的状态
	IsClosed bool
	//告知当前链接已经停止/退出
	ExitChan chan bool
	//消息管理MsgId和对应的业务处理API
	MsgHandle ziface.IMsgHandle
	//添加一个无缓冲的管道，用于读写之间的通信
	MsgChan chan []byte
	//链接属性集合string可以是链接名，空接口可以是任何属性
	property map[string]interface{}

	//保护链接属性的锁
	protectLock sync.RWMutex
}

//初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		//以后框架变大了之后会有多个服务器，所以需要一个属性判断当前链接属于哪个server
		//目前可以通过本字段访问connManager 本条注释会写多次因为我记不住
		TcpServer: server, //传个形参进来
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandle,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
		MsgChan:   make(chan []byte),
		property:  make(map[string]interface{}),
	}

	c.TcpServer.GetConnMgr().Add(c) //有种左脚踩右脚的感觉

	return c
}

/*
实现接口中的各种方法

*/

// StartReader 读数据
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "connID", c.ConnID, "Reader is exit,remote addr is")

	//只要Reader停止，就会调用stop函数，所以可以在stop函数中写吧关闭的消息传递给Exitchan
	defer c.Stop()

	for {
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

		//发送request到TaskQueue
		if utils.GlobalObject.WorkerPoolSize > 0 { //说明已经开启了工作池
			c.MsgHandle.SendMsgToTaskQueue(&req) //这里不需要用go，因为go已经在方法中了
		} else { //没有开启连接池，直接发送消息找到绑定的路由
			go c.MsgHandle.DoMsgHandler(&req)

		}

	}

}

//写消息的goroutine，专门发送给客户消息的方法
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]")

	//不断的阻塞等待chan的消息
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan: //可以在Exitchan中读到数据说明read那边说可以退出了
			return

		}
	}
}

//启动连接，让当前的连接准备工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID:", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()

	//启动从当前链接的写数据业务
	go c.StartWriter()

	//按照开发者传进来的，创建链接后需要调用的处理业务，执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop()...ConnID:\n", c.ConnID)
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true

	//在链接关闭之前处理一些业务，调用hook函数
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前链接从ConnManager中删除
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.MsgChan)

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

	//把消息发给管道
	c.MsgChan <- msg
	return nil
}

/*
一些管理属性的方法
*/
//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.protectLock.Lock()
	defer c.protectLock.Unlock()

	//添加一个属性
	c.property[key] = value

}

//获取链接属性
func (c *Connection) GetProperty(key string) (value interface{}, err error) {
	c.protectLock.RLock()
	defer c.protectLock.RUnlock()

	value, ok := c.property[key]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("no property FOUND!")
	}

}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.protectLock.Lock()
	defer c.protectLock.Unlock()

	delete(c.property, key)
}
