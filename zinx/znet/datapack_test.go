package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

/*
	test结尾，负责测试拆包，封包的单元测试
*/

//固定形参
func TestDataPack(t *testing.T) {
	//模拟服务器
	//1.创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Server Listen Err:", err)
		return
	}
	//创建一个协程，负责从客户端处理业务
	go func() {
		//阻塞等待客户端连接
		for {
			//2.从客户端读取数据，拆包处理
			conn, e := listener.Accept()
			if e != nil {
				fmt.Println("Server Accept Err: ", err)
				return
			}
			//开启协程，处理数据
			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					//1.第一次从conn读，读出包的头head
					headData := make([]byte, dp.GetHeadLen())
					_, e := io.ReadFull(conn, headData) //根据切片大小一次读满
					if e != nil {
						fmt.Println("Read Head Err: ", e)
						return
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("Server Unpack Err: ", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						//message中有数据，进行第二次读取
						//2.第二次从conn读，根据datalen，读取data内容
						data := make([]byte, msgHead.GetMsgLen())
						_, e := io.ReadFull(conn, data)
						if e != nil {
							fmt.Println("Server Unpack Data Err: ", e)
							return
						}
						msgHead.SetMsgData(data)

						/* 效果相同
						msg := msgHead.(*Message) //类型断言，向下转型
						msg.Data = make([]byte, msg.DataLen)
						_, e = io.ReadFull(conn, msg.Data)
						if e != nil {
							fmt.Println("Server Unpack Data Err: ", e)
							return
						}
						*/

						//一个完整的消息读取完毕
						fmt.Println("----> Recv MsgID: ", msgHead.GetMsgId(), "datalen: ", msgHead.GetMsgLen(),
							"data: ", string(msgHead.GetMsgData()))
					}

				}
			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client Dial Err: ", err)
		return
	}
	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	pack1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Data Pack Err: ", err)
		return
	}

	//封装第二个msg包
	msg2 := &Message{
		Id:      2,
		DataLen: 6,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', '!'},
	}
	pack2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Data Pack Err: ", err)
		return
	}

	//将两个包连在一起
	pack1 = append(pack1, pack2...) //... - 切片打散传递

	//一次性发送给服务端
	_, err = conn.Write(pack1)
	if err != nil {
		fmt.Println("Data Send Err: ", err)
	}
	//客户端阻塞
	select {}
}
