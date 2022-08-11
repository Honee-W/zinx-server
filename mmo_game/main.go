package main

import (
	"ZinxDemo/mmo_game/apis"
	"ZinxDemo/mmo_game/core"
	"ZinxDemo/zinx/ziface"
	"ZinxDemo/zinx/znet"
	"fmt"
)

//当前客户端建立连接之后的hook函数
func OnConnAdd(connection ziface.IConnection) {
	//创建一个player对象
	player := core.NewPlayer(connection)
	//给客户端发送MsgID=1的消息 --- 同步玩家ID给客户端
	player.SyncPid()
	//给客户端发送MsgID=200的消息 --- 广播上线位置
	player.BroadCastStartPosition()

	//将当前新上线的玩家添加到worldmanager中
	core.WorldMgrObj.AddPlayer(player)

	//将该连接绑定一个Pid(玩家ID)的属性
	connection.SetProperty("pid", player.Pid)

	//同步周边玩家，告知他们当前玩家已经上线，广播当前玩家的位置信息
	player.SyncSurrounding()

	fmt.Println("=====> Player pid = ", player.Pid, " is online <=====")
}

//当前客户端断开连接之前的hook函数
func OnConnStop(connection ziface.IConnection) {
	//根据连接属性获取pid
	pid, err := connection.GetProperty("pid")
	if err != nil {
		fmt.Println("get property pid err: ", err)
		return
	}
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//得到玩家句柄，执行下线业务
	player.Offline()
}

func main() {
	//创建zinx server句柄
	s := znet.NewServer()

	//连接创建和销毁的hook函数
	s.SetOnConnStart(OnConnAdd)
	s.SetOnConnStop(OnConnStop)

	//注册路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.Move{})

	//启动服务
	s.Serve()
}
