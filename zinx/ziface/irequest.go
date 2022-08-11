package ziface

/*
	IRequest接口
	把客户端请求的连接和数据包装在了一个Request中
*/

type IRequest interface {
	//得到当前连接
	GetConn() IConnection

	//得到请求数据
	GetData() []byte

	//得到消息的ID
	GetMsgID() uint32
}
