package znet

import (
	"zinxProject/ziface"
)

/*
实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
*/

type BaseRouter struct{}

//这里方法都先空着，需要处理具体业务的时候再重写
//在处理conn业务之前的钩子方法
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//处理conn业务的主方法
func (br *BaseRouter) Handle(request ziface.IRequest) {}

//在处理conn业务之后的钩子方法
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
