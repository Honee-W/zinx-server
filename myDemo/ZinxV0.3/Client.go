package main

import (
	"fmt"
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
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("Con Error, ", err)
		return
	}

	//2.写数据
	for {
		_, err := conn.Write([]byte("Hello Zinx V0.3"))
		if err != nil {
			fmt.Println("Write Error, ", err)
			return
		}

		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read Error, ", err)
			return
		}
		fmt.Printf("server call back: %s, length=%d\n", buf[:n], n)

		//cpu阻塞，防止进入死循环，浪费资源
		time.Sleep(1 * time.Second)
	}
}
