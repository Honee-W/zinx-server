package ziface

/*
	IRouter 抽象接口
	数据格式为IRequest
*/

type IRouter interface {
	//处理业务前的方法
	PreHandle(request IRequest)

	//处理业务方法
	Handle(request IRequest)

	//处理业务后的方法
	PostHandle(request IRequest)
}
