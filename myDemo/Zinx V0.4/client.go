package main

import (
	"fmt"
	"net"
	"time"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	//1.直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("client start err", err)
		return
	}

	//2.链接直接调用write，写数据
	for {
		//向服务端发送数据
		_, err := conn.Write([]byte("000000000000000")) //直接发送字节过去，没有封装成数据包。服务端读取数据时会把前8个字节读取成数据头，然后显示数据包过大
		if err != nil {
			fmt.Println("Write conn err", err)
			break
		}

		//从服务端接收数据
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		fmt.Printf("server call back:%s,cnt:%d\n", buf[:cnt], cnt)

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
