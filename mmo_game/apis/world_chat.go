package apis

import (
	"ZinxDemo/mmo_game/core"
	"ZinxDemo/mmo_game/pb"
	"ZinxDemo/zinx/ziface"
	"ZinxDemo/zinx/znet"
	"fmt"
	"google.golang.org/protobuf/proto"
)

/*
世界聊天 路由业务
*/
type WorldChatApi struct {
	znet.BaseRouter //继承
}

func (w *WorldChatApi) Handle(request ziface.IRequest) {
	//1.解析客户端传送的proto
	proto_msg := &pb.Talk{}
	e := proto.Unmarshal(request.GetData(), proto_msg)
	if e != nil {
		fmt.Println("Talk Unmarshal err: ", e)
		return
	}
	//获取发送当前聊天的玩家id
	pid, err := request.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("get pid err: ", err)
		return
	}
	//根据pid得到对应的player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//将这个消息广播给其他所有玩家
	player.Talk(proto_msg.Content)

}
