package ziface

type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Serve()
	//路由功能，注册路由方法，供客户端连接使用
	AddRouter(uint32, IRouter)
	//获取当前server的连接管理模块
	GetConnMgr() IConnManager
	//注册onConnStart 钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	//注册onConnStop 钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	//调用onConnStart 钩子函数的方法
	CallOnConnStart(connection IConnection)
	//调用onConnStop 钩子函数的方法
	CallOnConnStop(connection IConnection)
}
