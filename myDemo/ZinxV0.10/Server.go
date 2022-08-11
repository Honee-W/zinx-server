package main

import (
	"ZinxDemo/zinx/ziface"
	"ZinxDemo/zinx/znet"
	"fmt"
)

/*
	基于Zinx框架来开发的 服务器端应用程序
*/

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	//再回写ping..ping..ping..
	err := request.GetConn().SendMsg(200, []byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (h *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle")
	//先读取客户端的数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	//再回写hello
	err := request.GetConn().SendMsg(201, []byte("Hello Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

//连接创建之后的hook函数
func OnConnStartHook(connection ziface.IConnection) {
	fmt.Println("======> DoConnectionStart is Called ... ")
	err := connection.SendMsg(202, []byte("DoConnectionStart"))
	if err != nil {
		fmt.Println("DoConnectionStart Err", err)
	}

	//给当前的连接设置一些属性
	fmt.Println("Set conn Properties ...")
	connection.SetProperty("Name", "小王")
	connection.SetProperty("Home", "https://www.baidu.com")
}

//连接销毁前的hook函数 -- 回收资源等
func OnConnStopHook(connection ziface.IConnection) {
	fmt.Println("======> DoConnectionStop is Called ... ")
	fmt.Println("connID = ", connection.GetConnID(), " is Offline ...")

	if v, err := connection.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", v)
	}

	if v, err := connection.GetProperty("Home"); err == nil {
		fmt.Println("Home = ", v)
	}
}

func main() {
	//1.船舰一个server句柄
	s := znet.NewServer()

	//2.注册连接的Hook函数
	s.SetOnConnStart(OnConnStartHook)
	s.SetOnConnStop(OnConnStopHook)

	//3，自定义router方法 --- 0号消息，返回ping
	router := PingRouter{}
	s.AddRouter(0, &router)
	//1号消息，返回hello
	r := HelloZinxRouter{}
	s.AddRouter(1, &r)

	//4.启动server
	s.Serve()
}
