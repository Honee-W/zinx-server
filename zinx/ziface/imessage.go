package ziface

/*
	将请求的消息封装到一个Message中 抽象接口
*/
type IMessage interface {
	//字段的Getter和Setter

	GetMsgId() uint32
	SetMsgId(uint32)

	GetMsgLen() uint32
	SetMsgLen(uint32)

	GetMsgData() []byte
	SetMsgData([]byte)
}
