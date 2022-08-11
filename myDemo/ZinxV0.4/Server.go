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

//PreHandle
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	conn := request.GetConn().GetTCPConnection()
	_, err := conn.Write([]byte("before ping..\n"))
	if err != nil {
		fmt.Println("Call Back Error", err)
	}
}

//Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	conn := request.GetConn().GetTCPConnection()
	_, err := conn.Write([]byte("ping..ping..ping..\n"))
	if err != nil {
		fmt.Println("Call Back Error", err)
	}
}

//PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	conn := request.GetConn().GetTCPConnection()
	_, err := conn.Write([]byte("after ping..\n"))
	if err != nil {
		fmt.Println("Call Back Error", err)
	}
}

func main() {
	s := znet.NewServer()

	//自定义router方法
	router := PingRouter{}
	s.AddRouter(&router)

	s.Serve()
}
