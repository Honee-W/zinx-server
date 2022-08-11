package znet

import (
	"ZinxDemo/zinx/utils"
	"ZinxDemo/zinx/ziface"
	"fmt"
	"strconv"
)

/*
 消息处理模块实现
*/
type MsgHandler struct {
	//存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	//Worker取任务的消息队列 ---- 消息队列与工作Goroutine一对一
	TaskQueue []chan ziface.IRequest
	//业务工作池的数量
	WorkerPoolSize uint32
}

//构造方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkPoolSize, //全局配置文件中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkPoolSize),
	}
}

//执行对应Router消息处理方法
func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//获取对应MsgID的路由
	router := m.Apis[request.GetMsgID()]
	if router == nil {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! need register!")
		return
	}
	//执行对应业务方法
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

//添加路由
func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前Msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		//已经存在
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//添加msg与API的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}

//启动一个Worker工作池 只能发生一次 一个zinx框架只有一个worker工作池
func (m *MsgHandler) StartWorkerPool() {
	//根据workerPoolSize分别开启Worker，每个Worker由一个协程承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个worker被启动
		//1.给当前的worker对应的channel开辟空间 第i个worker用第i个channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2.启动当前的worker，阻塞等待消息从channel中到来，处理业务
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

//启动一个Worker工作流程
func (m *MsgHandler) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkID = ", workID, "is started ...")
	//不断阻塞等待对应的消息队列的消息
	for {
		select {
		//有消息传来，出列的就是一个客户端的request，执行当前request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue，由worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	//1。将消息平均分配给不同的worker --- 轮询分配
	//根据客户端建立的ConnID来分配
	workerID := request.GetConn().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConn().GetConnID(),
		" request MsgID = ", request.GetMsgID(), " to WorkerID = ", workerID)

	//2.将消息发送给对应worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}
