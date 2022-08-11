package znet

import "ZinxDemo/zinx/ziface"

/*
	路由接口实现基类，方法均为空，可选择性继承重写 (适配器设计模式)
*/

type BaseRouter struct{}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
