package znet

/*
	Message 实现类
*/
type Message struct {
	Id      uint32 //消息ID
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) SetMsgId(u uint32) {
	m.Id = u
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) SetMsgLen(u uint32) {
	m.DataLen = u
}

func (m *Message) GetMsgData() []byte {
	return m.Data
}

func (m *Message) SetMsgData(bytes []byte) {
	m.Data = bytes
}

//创建Msg方法
func (m *Message) NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}
