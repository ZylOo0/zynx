package znet

import "zyloo.com/zinx/ziface"

// 实现 Router 时，先嵌入这个 BaseRouter 基类，然后根据需要对这个基类的方法重写
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

func (br *BaseRouter) Handle(request ziface.IRequest) {}

func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
