package znet

import (
	"ZinxDemo/zinx/utils"
	"ZinxDemo/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

/*
	连接模块
*/
type Connection struct {
	//当前连接隶属于哪一个server
	TcpServer ziface.IServer
	//当前连接的TCP套接字
	Conn *net.TCPConn

	//当前连接的ID
	ConnID uint32

	//当前连接是否关闭
	isClosed bool

	//当前连接绑定的方法
	//handleAPI ziface.HandleFunc
	//该链接的处理方法Router
	//Router ziface.IRouter
	//单个Router优化为MsgHandler
	MsgHandler ziface.IMsgHandler

	//告知当前连接已经退出/停止的 channel (由Reader告知)
	ExitChan chan bool

	//无缓冲的管道，用于读，写Goroutine之间的消息通信
	MsgChan chan []byte

	//连接属性集合
	properties map[string]interface{} //interface{}万能类型
	//保护连接属性的锁
	propertiesLock sync.RWMutex
}

//初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		properties: make(map[string]interface{}),
	}

	//将conn加入ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

//读Goroutine 读取客户端发送的数据，将数据放入管道中供写Goroutine消费
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, " [Reader is exit], remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		//改为读取并封装消息
		//读取客户端数据到缓冲中 最大为512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取客户端消息的Msg Head  二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read Msg Head Err", err)
			break
		}
		//拆包，得到 msgID 和 msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack Err", err)
			break
		}
		///根据 msgDataLen 再次读取 Data 放在msg消息中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read Msg Data Err", err)
				break
			}
		}
		msg.SetMsgData(data)

		//得到当前连接的Request
		req := Request{
			conn: c,
			msg:  msg,
		}

		//再次优化 使用工作池和消息队列
		if utils.GlobalObject.WorkPoolSize > 0 {
			//已经开启了工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//未开启则一个连接对应一个handler协程
			go c.MsgHandler.DoMsgHandler(&req)
		}
		//开启新go程，执行注册的路由方法，取代原有的HandleAPI
		//go func(request ziface.IRequest) {
		//找到注册绑定的Conn对应的Router，调用方法
		//c.Router.PreHandle(request)
		//c.Router.Handle(request)
		//c.Router.PostHandle(request)

		//取代单一Router，根据request中的MsgID获取到Router，并执行业务方法
		//	c.MsgHandler.DoMsgHandler(request)
		//}(&req)
	}
}

//写Goroutine 读取管道中的内容，向客户端写数据
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.GetRemoteAddr().String(), "[conn Writer exit!]")

	//不断阻塞等待channel中的消息，读取后写回客户端
	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data Err, ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也退出
			return
		}
	}
}

//发送消息
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(&Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	})
	if err != nil {
		fmt.Println("Pack Data Err, ", err)
		return errors.New("Pack Msg Error")
	}

	//将数据写回客户端 ----> 改为写入管道，由写Goroutine写回客户端
	c.MsgChan <- binaryMsg
	//if _, err := c.GetTCPConnection().Write(binaryMsg); err != nil {
	//	fmt.Println("Write Msg id ", msgId, "error: ", err)
	//	return errors.New("Conn Write Error")
	//}
	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	//启动当前连接的读业务
	go c.StartReader()
	//启动当前连接的写业务
	go c.StartWriter()

	//按照开发者传递进来的方法，执行连接创建后对应的钩子函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	//按照开发者传递进来的方法，执行连接销毁前对应的钩子函数
	c.TcpServer.CallOnConnStop(c)

	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}

	//关闭连接
	c.isClosed = true
	for {
		err := c.Conn.Close()
		if err == nil {
			break
		}
	}

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从ConnMgr中删除
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭资源
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertiesLock.Lock() //Lock()写锁 RLock()读锁
	defer c.propertiesLock.Unlock()
	//添加属性
	c.properties[key] = value
}

//获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertiesLock.RLock() //加读锁
	defer c.propertiesLock.RUnlock()
	if v, ok := c.properties[key]; !ok {
		return nil, errors.New("no property found")
	} else {
		return v, nil
	}
}

//移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()
	delete(c.properties, key)
}
