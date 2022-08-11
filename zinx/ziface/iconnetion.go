package ziface

import (
	"net"
)

/*
	定义连接模块的抽象层
*/

type IConnection interface {
	Start() //启动连接
	Stop()  //停止连接

	//获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前连接模块的连接ID
	GetConnID() uint32

	//获取远程客户端的TCP状态 IP Port
	GetRemoteAddr() net.Addr

	//发送数据
	SendMsg(msgId uint32, data []byte) error

	//设置连接属性
	SetProperty(key string, value interface{})
	//获取连接属性
	GetProperty(key string) (interface{}, error)
	//移除连接属性
	RemoveProperty(key string)
}

/*
	定义处理连接业务的方法
*/
type HandleFunc func(*net.TCPConn, []byte, int) error
