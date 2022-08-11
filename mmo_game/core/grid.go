package core

import (
	"fmt"
	"sync"
)

/*
	一个AOI(兴趣点)地图中的格子类型
*/
type Grid struct {
	//格子ID
	GID int
	//格子左边界边界坐标
	MinX int
	//格子右边界边界坐标
	MaxX int
	//格子上边界边界坐标
	MinY int
	//格子下边界边界坐标
	MaxY int
	//当前格子内玩家或物体成员的ID集合
	playerIDs map[int]bool
	//保护当前集合的锁
	pIDLock sync.RWMutex
}

//初始化当前格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
		pIDLock:   sync.RWMutex{},
	}
}

//向格子中添加一个玩家
func (g *Grid) Add(playerID int) {
	//加写锁
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

//从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	//加写锁
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

//得到当前格子中的所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	//加读锁
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

//调试使用 -- 打印出格子的基本信息
func (g *Grid) String() string {
	//重写string()方法 -- fmt.println()
	return fmt.Sprintf("Grid id:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
