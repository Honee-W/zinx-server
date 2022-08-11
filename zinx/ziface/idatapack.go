package ziface

/*
	解决TCP粘包问题的封包和拆包模块
	针对Massage进行封装和拆包
	格式为 Len(4字节)+ID(4字节)+Data(Len字节)
			Header				Body
*/

type IDataPack interface {
	//获取包的头的长度
	GetHeadLen() uint32
	//封包方法
	Pack(message IMessage) ([]byte, error)
	//拆包方法
	Unpack([]byte) (IMessage, error)
}
