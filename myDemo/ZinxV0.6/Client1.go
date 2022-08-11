package main

import (
	"ZinxDemo/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
	模拟客户端
*/

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	//1.连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Con Error, ", err)
		return
	}

	//2.写数据
	for {
		//发送封包消息
		dp := znet.NewDataPack()
		data := []byte("zinx  client1 test")
		binaryData, err := dp.Pack(&znet.Message{
			Id:      1,
			DataLen: uint32(len(data)),
			Data:    data,
		})
		if err != nil {
			fmt.Println("Pack Data Err, ", err)
			return
		}
		_, err = conn.Write(binaryData)
		if err != nil {
			fmt.Println("Send Data Err, ", err)
			return
		}

		//对服务器端返回的消息进行解包
		//1.先读取流中的Msg Head部分 得到ID和dataLen
		MsgHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, MsgHead)
		if err != nil {
			fmt.Println("Read Head Err", err)
			return
		}
		msg_head, err := dp.Unpack(MsgHead)
		if err != nil {
			fmt.Println("Unpack Head Err", err)
			return
		}
		//2.再读取Msg Body
		if msg_head.GetMsgLen() > 0 {
			msg := msg_head.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("Read Msg Err", err)
				return
			}
			fmt.Println("---->Recv Server Msg: ID = ", msg.Id,
				",len = ", msg.DataLen, ",data = ", string(msg.Data))
		}
		//cpu阻塞，防止进入死循环，浪费资源
		time.Sleep(1 * time.Second)
	}
}
