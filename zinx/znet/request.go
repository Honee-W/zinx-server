package znet

import "ZinxDemo/zinx/ziface"

type Request struct {
	//已经和客户端建立好的连接
	conn ziface.IConnection
	//二进制数据封装为消息
	msg ziface.IMessage
	//客户端请求的数据
	//data []byte
}

func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
