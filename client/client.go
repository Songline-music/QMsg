package client

// client.go 包含了客户端相关的函数，用于连接服务器、发送消息和接收消息

import (
	"QMsg/protocol"
	"QMsg/tui"
	"bufio"
	"fmt"
	"net"
)

// 运行客户端
func RunClient(address string, name string, token string) (bool, string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("连接服务端失败 ", err)
		return false, ""
	}
	defer conn.Close()

	err = protocol.SendMessage(conn, protocol.Message{
		Type:    protocol.MessageLogin,
		Content: name,
		Token:   token,
	})
	if err != nil {
		fmt.Println("发送用户名失败 ", err)
		return false, ""
	}

	serverReader := bufio.NewReader(conn)
	reply, err := protocol.ReadMessage(serverReader)
	if err != nil {
		fmt.Println("读取服务端响应失败 ", err)
		return false, ""
	}

	if reply.Type == protocol.MessageError {
		fmt.Println("[系统] " + reply.Content)
		return false, ""
	}

	if reply.Type == protocol.MessageOK {
		name = reply.Content
	} else {
		fmt.Println("服务端响应格式错误")
		return false, ""
	}

	err = tui.Run(conn, serverReader, name, address)
	if err != nil {
		fmt.Println("启动 TUI 失败 ", err)
	}
	return true, name
}
