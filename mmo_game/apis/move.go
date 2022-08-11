package apis

import (
	"ZinxDemo/mmo_game/core"
	"ZinxDemo/mmo_game/pb"
	"ZinxDemo/zinx/ziface"
	"ZinxDemo/zinx/znet"
	"fmt"
	"google.golang.org/protobuf/proto"
)

//玩家移动 业务路由
type Move struct {
	znet.BaseRouter
}

func (m *Move) Handle(request ziface.IRequest) {
	//解析客户端传来的proto数据
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: position unmarshal err: ", err)
		return
	}

	//得到当前发送位置的玩家
	pid, err := request.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("get property pid err: ", err)
		return
	}

	fmt.Printf("Player pid = %d, move(%f,%f,%f,%f)\n", pid, proto_msg.X,
		proto_msg.Y, proto_msg.Z, proto_msg.V)

	//向其他玩家广播当前玩家的位置信息
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//广播并更新当前玩家坐标
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
