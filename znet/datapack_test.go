package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//测试的时候先注释掉	GlobalObject.Reload()

func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	listenner, err := net.Listen("tcp", "0.0.0.0:8989")
	if err != nil {
		fmt.Println("server listenner err", err)
		return
	}

	//创建go承载从客户端处理业务
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//拆包过程
				//定义一个拆包的对象
				dp := NewDataPack()
				for {
					//1.第一次从conn中读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err", err)
						return
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("unPack head err", err)
						return
					}

					if msgHead.GetMsgLen() > 0 { //判断数据包的长度是否大于0，也就是是否有数据
						//如果有数据，需要二次读取
						//2.第二次从conn中读，根据head中的datalen，在读取data内容
						msg := msgHead.(*Message) //通过断言，接口类型转换成结构体
						msg.Data = make([]byte, msg.DataLen)

						//根据datalen的长度再次从io流中读取
						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("read data err", err)
							return
						}

						//完整的消息已经读取完毕
						fmt.Println("--->Recv MsgID:", msg.Id, "datalen:", msg.DataLen, "data:", string(msg.Data))

					}

				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "0.0.0.0:8989")
	if err != nil {
		fmt.Println("链接出错", err)
		return
	}

	dp := NewDataPack()

	//模拟粘包过程，封装两个mes一期发送
	//封装第一个包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack(msg1)出错", err)
		return
	}

	//封装第二个包
	msg2 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack(msg2)出错", err)
		return
	}

	//两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性发送给服务端
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}
