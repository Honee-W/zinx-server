package core

import (
	"ZinxDemo/mmo_game/pb"
	"ZinxDemo/zinx/ziface"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

/*
	服务器端的玩家对象
*/
type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家与客户端的连接
	X    float32            //平面X坐标
	Y    float32            //高度
	Z    float32            //平面Y坐标
	V    float32            //旋转角度  0-360°
}

//玩家ID生成器  -- 真实开发应为数据库提供自增主键
var PidGen int32
var IdLock sync.Mutex //保护生成器的锁

//玩家初始化方法
func NewPlayer(connection ziface.IConnection) *Player {
	//生成玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建一个玩家对象
	return &Player{
		Pid:  id,
		Conn: connection,
		X:    float32(160 + rand.Intn(10)),
		Y:    0, //高度为0
		Z:    float32(140 + rand.Intn(20)),
		V:    0, //角度为0
	}
}

/*
	给客户端发送消息
	主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将proto Message结构体序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}
	//将二进制文件 通过zinx的SendMsg方法将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection to player is nil")
		return
	}
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("player send msg err: ", err)
		return
	}
}

//将已经生成的玩家ID发送给客户端
func (p *Player) SyncPid() {
	//创建MsgID=0的proto数据
	data := &pb.SyncPid{Pid: p.Pid}
	//将消息发送给客户端
	p.SendMsg(1, data)
}

//广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//创建MsgID=200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //tp=2 代表广播位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}
	//将数据发送给客户端
	p.SendMsg(200, proto_msg)
}

//玩家广播消息方法
func (p *Player) Talk(content string) {
	//创建MsgID=200的proto数据
	proto_msg := &pb.BroadCast{
		Pid:  p.Pid,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}
	//得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()
	//向所有的玩家(包括自己)发送MsgID=200的消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

}

//向周边玩家广播自己的位置信息
func (p *Player) SyncSurrounding() {
	///获取当前玩家周围玩家(九宫格)
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	//将当前玩家的位置信息通过MsgID=200 发送给周围玩家（让其他玩家看到自己）
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	//将周围全部玩家的位置信息发送给当前的玩家客户端 MsgID=202（让自己看到其他玩家）
	//创建MsgID=202的proto数据
	//制作pb.player切片
	players_proto := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//创建一个message player
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto = append(players_proto, p)
	}
	//创建SyncPlayer protobuf数据
	syncPlayers := &pb.SyncPlayers{Ps: players_proto[:]}
	//将创建好的数据发送给当前玩家的客户端
	p.SendMsg(202, syncPlayers)
}

//广播并更新当前玩家坐标 MsgID=200, TP=4
func (p *Player) UpdatePos(x, y, z, v float32) {
	//更新当前玩家的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	//组建广播协议 MsgID=200 TP=4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}
	//获取当前玩家的周边玩家
	players := p.GetSurroundPlayers()
	//给每个玩家对应的客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

//获取当前玩家周边九宫格的玩家
func (p *Player) GetSurroundPlayers() []*Player {
	//获取pid
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	//将所有玩家放入切片中
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

func (p *Player) Offline() {
	//封装下线消息
	proto_msg := &pb.SyncPid{Pid: p.Pid}
	//得到周围九宫格的玩家
	players := p.GetSurroundPlayers()
	//向周围玩家广播下线
	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	//删除玩家
	WorldMgrObj.RomovePlayerByPid(p.Pid)
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
}
