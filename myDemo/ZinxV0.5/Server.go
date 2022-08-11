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
	fmt.Println("Call Router Handle")
	//先读取客户端的数据
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	//再回写ping..ping..ping..
	err := request.GetConn().SendMsg(1, []byte("ping..ping..ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()

	//自定义router方法
	router := PingRouter{}
	s.AddRouter(&router)

	s.Serve()
}
