package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

//定义value结构体
type ClientInfo struct {
	C    chan string //给客户端传送消息
	Name string
	Addr string
}

//声明存放客户的map
var ClientOnMap map[string]ClientInfo

//建立一个message管道专门存放谁上线了
var message = make(chan string)

//新开一个协程，每来一个客户，遍历map，给map每个成员发送消息
func Manger() {
	//给map分配空间
	ClientOnMap = make(map[string]ClientInfo)
	for {
		//message为空时会阻塞
		msg := <-message

		//将谁上线的消息发送给每个客户
		for _, cli := range ClientOnMap {
			cli.C <- msg
		}
	}

}

//返回客户端的消息
func WriteToMessage(cli ClientInfo, msg string) (result string) {
	result = "[" + cli.Addr + "]" + cli.Name + " " + msg
	return
}

//广播客户端消息
func WriteToClient(cli ClientInfo, conn net.Conn) {
	for msg := range cli.C {
		conn.Write([]byte(msg + "\n"))
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	//获取客户端的地址
	Addr := conn.RemoteAddr().String()

	//将客户端用户加入到map中，这时候我们要建立一个map,因为value是一个结构体，我们要创建一个结构体
	cli := ClientInfo{make(chan string), Addr, Addr}
	ClientOnMap[Addr] = cli

	//将谁上线的消息传递给管道，以便将消息传递给每个客户端的管道
	msg := "login"
	message <- WriteToMessage(cli, msg)

	//给客户端广播哪个客户上线
	go WriteToClient(cli, conn)

	isQuit := make(chan bool)  //用户主动退出用
	hasData := make(chan bool) //超时处理
	//新建一个协程，接收客户发过来的消息
	go func() {
		m := make([]byte, 2048)
		for {
			n, err4 := conn.Read(m)
			if err4 != nil {
				fmt.Println("conn.Read err4 = ", err4)
				isQuit <- true
				return
			}
			msg := string(m[:n-1])

			if len(msg) == 3 && msg == "who" { //查询在线用户
				conn.Write([]byte("User list:\n"))
				for _, temp := range ClientOnMap {
					msg = ("[" + temp.Addr + "]" + temp.Name)
					conn.Write([]byte(msg + "\n"))
				}
			} else if len(msg) >= 8 && msg[:6] == "rename" { //修改用户名
				//rename|mike
				name := strings.Split(msg, "|")[1]
				cli.Name = name
				ClientOnMap[Addr] = cli
				conn.Write([]byte("successful\n"))
			} else {
				message <- WriteToMessage(cli, msg)
			}

			hasData <- true

		}
	}()

	//查询在线用户

	for {
		select {
		case <-isQuit:
			delete(ClientOnMap, Addr)
			message <- WriteToMessage(cli, "logout")
			return
		case <-hasData:

		case <-time.After(60 * time.Second):
			delete(ClientOnMap, Addr)
			message <- WriteToMessage(cli, "timeout")
			return

		}
	}

}

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.Listener err = ", err)
		return
	}

	defer listener.Close()

	//新开一个协程，每来一个客户，遍历map，给map每个成员发送消息
	go Manger()

	for {

		//并发聊天，肯定是可以不断连接客户端的
		conn, err1 := listener.Accept()
		if err1 != nil {
			fmt.Println("listener.Accept err1 = ", err1)
			continue
		}

		//处理连接
		go HandleConn(conn)
	}
}
