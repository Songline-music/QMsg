package server

import (
	"QMsg/protocol"
	"QMsg/ui"
	"fmt"
	"net"
	"sync"
)

// session.go 包含了客户端会话相关的函数，用于处理用户输入、发送消息和接收消息等功能

// 用户结构体
type Client struct {
	name string                // 用户名
	conn net.Conn              // 网络连接
	send chan protocol.Message // 消息发送通道
}

var (
	clients = make(map[*Client]bool) // 用户列表
	lock    sync.Mutex               // 用户列表锁
)

// 关闭所有客户端连接
func closeAllClients() {
	lock.Lock()
	targets := make([]*Client, 0, len(clients))
	for client := range clients {
		targets = append(targets, client)
	}
	lock.Unlock()

	for _, client := range targets {
		sendToClient(client, protocol.Message{
			Type:    protocol.MessageSystem,
			Content: "服务端已关闭",
		})
		client.conn.Close()
	}
}

// 为用户准备用户名(防重名),如重名则返回false
func prepareClientName(name string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()

	if name == "" || name == "匿名用户" {
		return nextAnonymousNameLocked(), true
	}

	for client := range clients {
		if client.name == name {
			return "", false
		}
	}

	return name, true
}

// 获取下一个匿名用户名,如匿名用户1,匿名用户2
func nextAnonymousNameLocked() string {
	for i := 1; ; i++ {
		name := fmt.Sprintf("匿名用户%d", i)

		exists := false
		for client := range clients {
			if client.name == name {
				exists = true
				break
			}
		}

		if !exists {
			return name
		}
	}
}

// 添加一个用户
func addClient(client *Client) {
	lock.Lock()
	clients[client] = true
	lock.Unlock()
}

// 移除一个用户
func removeClient(client *Client) {
	lock.Lock()
	exists := clients[client]
	if exists {
		delete(clients, client)
		close(client.send)
		client.conn.Close()
	}
	lock.Unlock()

	if exists {
		broadcast(protocol.Message{
			Type:    protocol.MessageSystem,
			Content: client.name + " 离开了群聊",
		})
		ui.PrintServerAsyncLog("用户离线 " + client.name)
	}
}

// 广播
func broadcast(message protocol.Message) {
	lock.Lock()
	targets := make([]*Client, 0, len(clients))
	for client := range clients {
		targets = append(targets, client)
	}
	lock.Unlock()

	for _, client := range targets {
		sendToClient(client, message)
	}
}

// 向指定用户发送消息
func sendToClient(client *Client, message protocol.Message) {
	select {
	case client.send <- message:
	default:
		ui.PrintServerAsyncLog("消息队列已满 " + client.name + " 跳过一次发送")
	}
}

// 向用户发送消息的携程
func writeToClient(client *Client) {
	for message := range client.send {
		err := protocol.SendMessage(client.conn, message)
		if err != nil {
			ui.PrintServerAsyncLog("发送消息失败 " + client.name + " " + err.Error())
			return
		}
	}
}

// 获取在线用户切片
func getClientNames() []string {
	lock.Lock()
	defer lock.Unlock()

	names := make([]string, 0, len(clients))
	for client := range clients {
		names = append(names, client.name)
	}

	return names
}

// 通过用户名查找用户
func findClientByName(name string) *Client {
	lock.Lock()
	defer lock.Unlock()

	for client := range clients {
		if client.name == name {
			return client
		}
	}

	return nil
}
