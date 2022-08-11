package znet

import (
	"ZinxDemo/zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

//构造方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map  加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//将conn加入map中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("conn added to ConnManager successfully : conn num = ", c.TotalConn())
}

func (c *ConnManager) Remove(Conn ziface.IConnection) {
	//保护共享资源map  加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除map中的连接
	delete(c.connections, Conn.GetConnID())
	fmt.Println("conn removed from ConnManager successfully : connID = ", Conn.GetConnID())
}

func (c *ConnManager) Get(ConnID uint32) (ziface.IConnection, error) {
	//保护资源 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[ConnID]; !ok {
		return nil, errors.New("connection not FOUND !")
	} else {
		return conn, nil
	}
}

func (c *ConnManager) TotalConn() int {
	return len(c.connections) //返回map中key的数量
}

func (c *ConnManager) Clear() {
	//保护共享资源map  加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	//删除conn并停止工作
	for connID, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
	}

	fmt.Println("clear all connections succ! conn num = ", c.TotalConn())
}
