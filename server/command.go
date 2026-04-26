package server

// command.go 包含了处理客户端消息的函数，包括帮助信息、用户列表和私聊功能

import (
	"QMsg/protocol"
	"QMsg/text"
	"strings"
)

// 处理客户端发送信息
func handleMessage(client *Client, message protocol.Message) {
	content := strings.TrimSpace(message.Content)
	lowerMessage := strings.ToLower(content)

	switch message.Type {
	case protocol.MessageCommand:
		handleCommand(client, lowerMessage)
	case protocol.MessagePrivate:
		sendPrivateMessage(client, content)
	case protocol.MessageChat:
		broadcastChat(client, content)
	default:
		sendToClient(client, protocol.Message{
			Type:    protocol.MessageError,
			Content: "未知消息类型",
		})
	}
}

// 发送用户列表给客户端
func sendUserList(client *Client) {
	names := getClientNames()

	sendToClient(client, protocol.Message{
		Type:    protocol.MessageSystem,
		Content: "当前在线用户 " + strings.Join(names, ", "),
	})
}

// 发送私聊信息
func sendPrivateMessage(from *Client, message string) {
	parts := strings.SplitN(message, " ", 2) // 分割成两部分：@用户名 和 消息内容
	if len(parts) < 2 {
		sendToClient(from, protocol.Message{
			Type:    protocol.MessageError,
			Content: "私聊格式错误,请使用 @用户名 消息内容",
		})
		return
	}

	targetName := strings.TrimPrefix(parts[0], "@")
	content := strings.TrimSpace(parts[1])
	if targetName == "" || content == "" {
		sendToClient(from, protocol.Message{
			Type:    protocol.MessageError,
			Content: "私聊格式错误,请使用 @用户名 消息内容",
		})
		return
	}

	// 获取目标用户
	target := findClientByName(targetName)
	if target == nil {
		sendToClient(from, protocol.Message{
			Type:    protocol.MessageError,
			Content: "用户 " + targetName + " 不在线",
		})
		return
	}

	content = strings.ReplaceAll(content, "\\n", "\n")

	sendToClient(target, protocol.Message{
		Type:    protocol.MessagePrivate,
		From:    from.name,
		Content: content,
	})
	sendToClient(from, protocol.Message{
		Type:    protocol.MessagePrivate,
		From:    "你 -> " + target.name,
		Content: content,
	})
}

// 处理命令
func handleCommand(client *Client, command string) {
	switch command {
	case "help", "/help":
		sendToClient(client, protocol.Message{
			Type:    protocol.MessageText,
			Content: text.BuildHelpText(),
		})
	case "/users":
		sendUserList(client)
	default:
		sendToClient(client, protocol.Message{
			Type:    protocol.MessageError,
			Content: "未知命令",
		})
	}
}

// 广播聊天
func broadcastChat(client *Client, content string) {
	content = strings.ReplaceAll(content, "\\n", "\n")

	broadcast(protocol.Message{
		Type:    protocol.MessageChat,
		From:    client.name,
		Content: content,
	})
}
