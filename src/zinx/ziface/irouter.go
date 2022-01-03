package ziface

//定义一个 路由抽象接口
//路由里的数据  都是irequest

type IRouter interface {
	// PreHandle 在处理conn业务之前的钩子方法hook
	PreHandle(request IRequest)
	// Handle 在处理conn业务的主方法 hook
	Handle(request IRequest)
	// PostHandle 在处理conn业务之后的钩子方法hook
	PostHandle(request IRequest)
}
