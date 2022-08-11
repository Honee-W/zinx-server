package core

import "fmt"

//定义一些AOI的边界值 -- 宏
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*
	AOI区域管理模块
*/

type AOIManager struct {
	//区域左边界的坐标
	MinX int
	//区域右边界的坐标
	MaxX int
	//X轴方向上格子的数量
	CntsX int
	//区域上边界的坐标
	MinY int
	//区域下边界的左右
	MaxY int
	//Y轴方向上各自的数量
	CntsY int
	//当前区域的格子集合 key=格子id，value=格子对象
	grids map[int]*Grid
}

/*
初始化一个AOI区域管理模块
*/
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	//给初始化区域中的所有格子进行编号和初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//根据x，y编号，计算格子的坐标
			gid := y + cntsX*x

			//初始化gid对应的格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight())
		}
	}
	return aoiMgr
}

/*
得到每个格子在X轴方向的宽度
*/
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

/*
得到每个格子在Y轴方向的高度
*/
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

/*
打印格子信息---调试
*/
func (m *AOIManager) String() string {
	//打印AOIManager信息
	s := fmt.Sprintf("AOIManager:\nMinX:%d, MaxX:%d, cntsX:%d, "+
		"minY:%d, maxY:%d, cntsY:%d\n", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	//打印所有格子信息
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

/*
根据格子GID得到周边九宫格格子的ID集合
*/
func (m *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	//判断当前gid对应的格子是否在AOIManager中
	if _, ok := m.grids[gid]; !ok {
		return
	}
	//初始化grids返回值 将当前gid本身加入九宫格切片中
	grids = append(grids, m.grids[gid])

	//判断gid的左边是否有格子？右边是否有格子？
	idx := gid % m.CntsX //通过gid得到当前格子x轴的编号

	//判断idx编号左边是否还有格子  有则放入grids中
	//判断idx编号右边是否还有格子  有则放入grids中
	if idx > 0 {
		//左边有格子
		grids = append(grids, m.grids[gid-1])
	}
	if idx < m.CntsX-1 {
		//右边有格子
		grids = append(grids, m.grids[gid+1])
	}

	//遍历grids中的每个格子的gid
	//gid上边是否还有格子？ 有则放入grids中
	//gid下边是否还有格子？ 有则放入grids中
	for _, grid := range grids {
		idy := grid.GID / m.CntsY
		if idy > 0 {
			//上边有格子
			grids = append(grids, m.grids[grid.GID-m.CntsX])
		}
		if idy < m.CntsY-1 {
			//下边有格子
			grids = append(grids, m.grids[grid.GID+m.CntsX])
		}
	}
	return grids
}

/*
根据玩家坐标得到所在格子的gid
*/
func (m *AOIManager) GetGidByPos(x, y float32) int {

	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()

	return idy*m.CntsX + idx
}

/*
通过横纵坐标得到周边九宫格内全部的PlayerIDs
*/
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	//得到当前玩家的所在格子的gid
	gid := m.GetGidByPos(x, y)
	//通过gid得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gid)
	//将九宫格的中所有playerID的累加到返回值中
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...) //数组拼接 ...打散处理
		//fmt.Printf("====> grid ID : %d, pids : %v ====", grid.GID, grid.playerIDs)
	}
	return playerIDs
}

//添加一个playerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

//移除一个格子的中playerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

//通过GID获取全部的playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	return m.grids[gID].GetPlayerIDs()
}

//通过坐标将player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.grids[gID].Add(pID)
}

//通过坐标把一个player从格子中移除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.grids[gID].Remove(pID)
}
