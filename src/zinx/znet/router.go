package znet

import "awesomeProject/src/zinx/ziface"

// BaseRouter 实现router时， 先嵌入这个baserouter基类，然后根据需要对这个基类的方法 进行重写
type BaseRouter struct {
}

// PreHandle 在处理conn业务之前的钩子方法hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 在处理conn业务的主方法 hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理conn业务之后的钩子方法hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}