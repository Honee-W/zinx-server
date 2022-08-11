package znet

import (
	"ZinxDemo/zinx/utils"
	"ZinxDemo/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

/*
	直接面向TCP连接中的数据流
	封包 拆包的具体实现
*/

type DataPack struct{} //无需属性 实现方法即可

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节) + ID uint32(4字节)
	return 8
}

func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	//创建一个存放byte的缓存
	buffer := bytes.NewBuffer([]byte{})

	//将Datalen写入buffer中 (二进制写入，小端(低地址存低字节)编码)
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MessageID写入buffer中
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgId()); err != nil {
		return nil, err
	}
	//将Data数据写入buffer中
	if err := binary.Write(buffer, binary.LittleEndian, message.GetMsgData()); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

//拆包方法 将包的Head信息读出即可，再根据Head信息中data的长度，再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个读取二进制数据的reader
	reader := bytes.NewReader(binaryData)

	//只读取head信息
	msg := &Message{}

	//读dataLen (对应小端读)
	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读MsgID
	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否超出允许的最大包大小，过长报错
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
