package znet

import (
	"ZinxDemo/zinx/utils"
	"ZinxDemo/zinx/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//Router ziface.IRouter
	//当前server的消息管理模块，绑定MsgID和对应的处理业务API
	MsgHandler ziface.IMsgHandler
	//当前server的连接管理模块
	ConnManager ziface.IConnManager
	//创建连接后自动调用的Hook函数
	OnConnStart func(conn ziface.IConnection)
	//销毁连接前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

//----由开发者定义
//定义客户端的回调方法，目前写死
//func CallBackToClient (conn *net.TCPConn, data []byte, cnt int) error {
//	//回写业务
//	fmt.Println("[Conn Handle] CallBackToClient")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

func (s *Server) Start() {
	fmt.Printf("[zinx] Server Name : %s, listenner at IP : %s, Port : %d is starting", utils.GlobalObject.Name,
		utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("Resolve Error, ", err)
		return
	}

	//开启消息队列
	s.MsgHandler.StartWorkerPool()

	//接受客户端连接
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("Listen Error, ", err)
		return
	}
	var cid uint32
	cid = 0

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept Error, ", err)
			continue
		}

		//处理连接前，先判断是否已经到达最大连接，若超过，则关闭
		if s.ConnManager.TotalConn() >= utils.GlobalObject.MaxConn {
			//TODO 给客户端响应一个超出最大连接的错误方法
			fmt.Println("Too Many Connection !")
			err := conn.Close()
			if err != nil {
				fmt.Println("conn close failed, ", err)
			}
			continue
		}

		////将处理连接的业务方法和连接绑定
		//dealConn := NewConnection(conn, cid, CallBackToClient)
		dealConn := NewConnection(s, conn, cid, s.MsgHandler)
		cid++

		//启动业务处理协程
		go dealConn.Start()
	}
}

func (s *Server) Stop() {
	//关闭连接，回收资源
	s.ConnManager.Clear()
	fmt.Println("[STOP] Zinx Serve name: ", s.Name)
}

func (s *Server) Serve() {
	s.Start()

	//TODO 启动服务器后的额外业务

	//阻塞等待
	select {}
}

//注册绑定路由方法
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}

	return s
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnManager
}

//注册Hook方法
func (s *Server) SetOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

//调用Hook方法
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---->Call OnConnStart() ... ")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---->Call OnConnStop() ... ")
		s.OnConnStop(connection)
	}
}
