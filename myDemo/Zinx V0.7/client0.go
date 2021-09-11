package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinxProject/znet"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("client0 start...")
	time.Sleep(1 * time.Second)
	//1.直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}

	//2.链接直接调用write，写数据
	for {
		//发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx v0.7 client test message")))
		if err != nil {
			fmt.Println("封包失败", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("发送消息错误", err)
			break
		}

		//从服务端接收消息
		//先读取流中的head部分 得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err", err)
			break
		}

		//再根据id和datalen第二次读取
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("解包失败", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			_, err = io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data error:", err)
				break
			}

			fmt.Println("-->Recv Server Msg:ID:", msg.GetMsgId(), "Len:", msg.GetMsgLen(), "data:", string(msg.GetMsgData()))

		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
