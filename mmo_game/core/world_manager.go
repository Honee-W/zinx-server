package core

import "sync"

/*
	世界管理模块
*/

type WorldManager struct {
	//当前AOI模块
	AoiMgr *AOIManager
	//玩家集合
	Players map[int32]*Player
	//保护集合的锁
	pLock sync.RWMutex
}

//初始化方法 -- 导包只调用一次
//提高一个对外的全局的世界管理模块句柄
var WorldMgrObj *WorldManager

func init() {
	WorldMgrObj = &WorldManager{
		//创建AOI地图
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
		pLock:   sync.RWMutex{},
	}
}

//添加一个玩家
func (m *WorldManager) AddPlayer(p *Player) {
	//将player添加到玩家集合中
	m.pLock.Lock()
	m.Players[p.Pid] = p
	m.pLock.Unlock()
	//将player添加到AOIManager中
	m.AoiMgr.AddToGridByPos(int(p.Pid), p.X, p.Z)
}

//删除玩家
func (m *WorldManager) RomovePlayerByPid(pid int32) {
	//通过id获得玩家，将玩家从AOIManager中删除
	player := m.Players[pid]
	m.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)
	//将玩家从玩家集合中删除
	m.pLock.Lock()
	delete(m.Players, pid)
	m.pLock.Unlock()
}

//通过玩家ID查询player对象
func (m *WorldManager) GetPlayerByPid(pid int32) *Player {
	//加读锁
	m.pLock.RLock()
	defer m.pLock.RUnlock()
	return m.Players[pid]
}

//获取全部的在线玩家
func (m *WorldManager) GetAllPlayers() []*Player {
	//加读锁
	m.pLock.RLock()
	defer m.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, player := range m.Players {
		players = append(players, player)
	}

	return players
}
