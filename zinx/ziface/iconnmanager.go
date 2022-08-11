package ziface

/*
	连接管理模块抽象层
*/

type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(Conn IConnection)
	//根据ID获取连接
	Get(ConnID uint32) (IConnection, error)
	//得到当前连接总数
	TotalConn() int
	//清除并终止所有连接
	Clear()
}
