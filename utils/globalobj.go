package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinxProject/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer //全局的Server对象
	Host      string         //当前服务器的主机监听IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //服务器名称

	/*
		zinx
	*/
	Version        string //版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //数据包的最大值
}

//提供一个全局对外的GlobalObj对象
var GlobalObject *GlobalObj

//从Zinx。Json文件中读取配置参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("myDemo\\Zinx V0.6\\conf\\zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject) //将配置文件的数据加载到全局变量GlobalObject中
	if err != nil {
		panic(err)
	}

}

//提供一个init方法，初始化当前的GlobalObject
func init() {
	//若果配置文件没有加载，默认值如下
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "Zinx V0.4",
		TcpPort:        8989,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//尝试加载自定义的值

	GlobalObject.Reload()

}
