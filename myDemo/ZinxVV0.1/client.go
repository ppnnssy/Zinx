package main

import (
	"fmt"
	"net"
	"time"
)

/*
模拟客户端
 */
func main(){
	fmt.Println("client start...")
	time.Sleep(1*time.Second)
	//1.直接链接远程服务器，得到一个conn链接
	conn,err:=net.Dial("tcp","0.0.0.0:8899")
	if err!=nil{
		fmt.Println("client start err",err)
		return
	}


	//2.链接直接调用write，写数据
	for  {
		//向服务端发送数据
		_,err:=conn.Write([]byte("Hello Zinx V0.2.."))
		if err!=nil{
			fmt.Println("Write conn err",err)
			continue
		}

		//从服务端接收数据
		buf:=make([]byte,512)
		cnt,err:=conn.Read(buf)
		fmt.Printf("server call back:%s,cnt:%d\n",buf[:cnt],cnt)

		//cpu阻塞
		time.Sleep(1*time.Second)
	}





}
