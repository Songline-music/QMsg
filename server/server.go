package server

// server.go 包含了服务端相关的函数,用于监听客户端连接.处理消息和管理用户列表

import (
	"QMsg/protocol"
	"QMsg/ui"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// 运行服务器
func RunServer(address string, token string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("监听失败 ", err)
		return
	}
	defer listener.Close()

	ui.ClearScreen()
	ui.PrintTitle("QMsg 服务端")
	ui.PrintSystem("服务端已启动,监听地址 " + address)
	ui.PrintSystem("等待客户端连接...")
	fmt.Println()

	shutdown := make(chan struct{})
	go listenServerCommands(listener, shutdown)

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-shutdown:
				closeAllClients()
				ui.PrintSystem("服务端已关闭")
				return
			default:
				ui.PrintServerAsyncLog("接收客户端连接失败 " + err.Error())
				continue
			}
		}

		go handleClient(conn, token)
	}
}

// 监听服务端命令
func listenServerCommands(listener net.Listener, shutdown chan struct{}) {
	input := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("server > ")
		command, err := input.ReadString('\n')
		if err != nil {
			return
		}

		ui.ClearCurrentLine()

		command = strings.TrimSpace(command)

		switch strings.ToLower(command) {
		case "exit", "/shutdown":
			close(shutdown)
			listener.Close()
			return
		case "help", "/help":
			ui.PrintServerLog("服务端命令: exit 或 /shutdown 关闭服务端,help 或 /help 查看帮助")
		case "":
			continue
		default:
			fmt.Println("未知服务端指令,输入 help 查看帮助")
		}
	}
}

// 处理客户端连接
func handleClient(conn net.Conn, token string) {
	reader := bufio.NewReader(conn)
	remoteAddr := conn.RemoteAddr().String()

	name, ok := readLogin(reader, token, conn)
	if !ok {
		conn.Close()
		return
	}

	protocol.SendMessage(conn, protocol.Message{
		Type:    protocol.MessageOK,
		Content: name,
	})

	client := &Client{
		name: name,
		conn: conn,
		send: make(chan protocol.Message, 16),
	}

	addClient(client)
	defer removeClient(client)

	go writeToClient(client)
	broadcast(protocol.Message{
		Type:    protocol.MessageSystem,
		Content: client.name + " 加入了群聊",
	})
	ui.PrintServerAsyncLog("用户上线 " + client.name + " " + remoteAddr)

	readClientMessage(reader, client, remoteAddr)
}

// 消息循环时刻监听客户端消息
func readClientMessage(reader *bufio.Reader, client *Client, remoteAddr string) {
	for {
		message, err := protocol.ReadMessage(reader)
		if err != nil {
			ui.PrintServerAsyncLog("连接断开 " + client.name + " " + remoteAddr + " " + err.Error())
			return
		}

		message.Content = strings.TrimSpace(message.Content)
		if message.Content == "" {
			continue
		}
		if isMessageTooLong(message.Content) {
			ui.PrintServerAsyncLog("消息过长 " + client.name + " " + remoteAddr)
			sendToClient(client, protocol.Message{
				Type:    protocol.MessageError,
				Content: "消息不能超过 " + fmt.Sprint(maxMessageLength) + " 个字符",
			})
			continue
		}

		handleMessage(client, message)
	}
}
