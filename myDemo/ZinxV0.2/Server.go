package main

import "ZinxDemo/zinx/znet"

/*
	基于Zinx框架来开发的 服务器端应用程序
*/

func main() {
	s := znet.NewServer("[zinx V0.2]")

	s.Serve()
}
