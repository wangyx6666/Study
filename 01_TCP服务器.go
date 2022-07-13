package main

import (
	"fmt"
	"net"
)

func main() {
	//监听   net.listen()
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	//阻塞等待接收   net.Accept()
	conn, err1 := listener.Accept()
	if err1 != nil {
		fmt.Println("err1 = ", err1)
		return
	}
	fmt.Println("小王上线了，我门尊贵的VIP客户，好好伺候。")

	for {

		//读数据      net.read()
		m := make([]byte, 1024)
		n, err2 := conn.Read(m)
		if err2 != nil {
			fmt.Println("err2 = ", err2)
			return
		}
		fmt.Println("服务器端接收数据是： = ", string(m[:n]))

		//应答        net.write()
		var m1 string
		fmt.Printf("到你发言了：")
		fmt.Scan(&m1)
		conn.Write([]byte(m1))

	}

	//关闭        close()
	defer listener.Close()
	defer conn.Close()

}
