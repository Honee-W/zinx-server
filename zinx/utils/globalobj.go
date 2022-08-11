package utils

import (
	"ZinxDemo/zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
	读取用户编写的json格式配置文件，加载到全局对象(globalobj)中
	将之前的硬编码采用全局对象(globalobj)中的参数替换
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前全局Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/
	Version          string //当前版本号
	MaxConn          int    //服务器主机允许的最大链接数
	MaxPackageSize   uint32 //当前数据包的最大值
	WorkPoolSize     uint32 //当前工作池中goroutine的数量 -- 与消息队列个数一一对应
	MaxWorkerTaskLen uint32 //每个worker对应的消息队列的任务数量的最大值
}

/*
	定义全局的对外GlobalObj对象
*/
var GlobalObject *GlobalObj

/*
	从conf/zinx.json加载配置
*/
func (g *GlobalObj) Reload() {
	//fmt.Println(os.Getwd())
	data, err := ioutil.ReadFile("mmo_game/conf/zinx.json") //myDemo/ZinxV0.4/conf/zinx.json
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	e := json.Unmarshal(data, &GlobalObject)
	if e != nil {
		panic(e)
	}
}

/*
	init()方法 只要被导入，就会执行
	提供init()，初始化当前的全局配置对象，提供默认配置
*/
func init() {
	GlobalObject = &GlobalObj{
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "ZinxServerApp",
		Version:          "v0.8",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkPoolSize:     10,
		MaxWorkerTaskLen: 1024,
	}

	//尝试从conf/zinx.json加载配置
	GlobalObject.Reload()
}
