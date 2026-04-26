package client

// client.go 包含了客户端相关的函数，用于连接服务器、发送消息和接收消息

import (
	"QMsg/protocol"
	"QMsg/ui"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	currentPrompt = "you > "
	promptLock    sync.Mutex
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
		ui.WaitEnter()
		return false, ""
	}

	if reply.Type == protocol.MessageOK {
		name = reply.Content
	} else {
		fmt.Println("服务端响应格式错误")
		ui.WaitEnter()
		return false, ""
	}

	ui.ClearScreen()
	ui.PrintTitle("QMsg 客户端")
	ui.PrintSystem("已连接到服务器 " + address)
	ui.PrintSystem("当前用户 " + name)
	ui.PrintSystem("输入 help 查看可用功能,输入 exit 退出聊天")
	fmt.Println()

	done := make(chan struct{})
	go readFromServer(serverReader, done) // 多线程读取服务器消息
	writeFromKeyboard(conn, done)
	return true, name
}

// 从服务器读取消息并且打印
func readFromServer(reader *bufio.Reader, done chan<- struct{}) {
	for {
		message, err := protocol.ReadMessage(reader)
		if err != nil {
			printIncoming("连接已断开 " + err.Error())
			close(done)
			return
		}

		content := strings.TrimSpace(message.Content)
		if content == "" {
			continue
		}

		printIncoming(formatIncomingMessage(message))
	}
}

// 格式化服务器发送消息
func formatIncomingMessage(message protocol.Message) string {
	switch message.Type {
	case protocol.MessageChat:
		return formatDisplayMessage(message.From, message.Content)
	case protocol.MessagePrivate:
		return formatDisplayMessage("[私聊] "+message.From, message.Content)
	case protocol.MessageSystem:
		return "[系统] " + message.Content
	case protocol.MessageError:
		return "[错误] " + message.Content
	case protocol.MessageText:
		return message.Content
	default:
		return message.Content
	}
}

// 从键盘读取输入并发送至服务器
func writeFromKeyboard(conn net.Conn, done <-chan struct{}) {
	input := bufio.NewReader(os.Stdin)
	for {
		printPrompt()
		text, err := input.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入失败 ", err)
			return
		}

		ui.ClearCurrentLine()

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		if text == "exit" {
			fmt.Println("聊天结束")
			return
		}

		if text == "/multi" {
			writeMultiLineMessage(conn, input)
			continue
		}

		message := buildOutgoingMessage(text)

		err = protocol.SendMessage(conn, message)
		if err != nil {
			fmt.Println("发送消息失败 ", err)
			return
		}

		select {
		case <-done:
			return
		default:
		}
	}
}

// 构建要发送的消息
func buildOutgoingMessage(text string) protocol.Message {
	messageType := protocol.MessageChat

	if strings.HasPrefix(text, "/") || strings.ToLower(text) == "help" {
		messageType = protocol.MessageCommand
	}
	if strings.HasPrefix(text, "@") {
		messageType = protocol.MessagePrivate
	}

	return protocol.Message{
		Type:    messageType,
		Content: text,
	}
}

// 处理客户端发送的多行消息
func writeMultiLineMessage(conn net.Conn, input *bufio.Reader) {
	setPrompt("| ")
	defer setPrompt("you > ")

	fmt.Println("---------进入多行输入模式---------")
	fmt.Println("输入 /send 发送, 输入 /cancel取消")

	lines := []string{}

	for {
		fmt.Print("| ")
		line, err := input.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入失败 ", err)
			return
		}

		line = strings.TrimRight(line, "\r\n")

		switch line {
		case "/send":
			if len(lines) == 0 {
				fmt.Println("没有输入任何内容")
				return
			}

			message := strings.Join(lines, "\\n")
			err = protocol.SendMessage(conn, protocol.Message{
				Type:    protocol.MessageChat,
				Content: message,
			})
			if err != nil {
				fmt.Println("发送消息失败 ", err)
			}
			return
		case "/cancel":
			fmt.Println("已取消多行输入")
			return
		default:
			lines = append(lines, line)
		}
	}
}

// 处理展示文本的格式
func formatDisplayMessage(from string, content string) string {
	prefix := from + ": "
	lines := strings.Split(content, "\n")

	if len(lines) == 0 {
		return prefix
	}

	for i := 1; i < len(lines); i++ {
		lines[i] = strings.Repeat(" ", len(prefix)) + lines[i]
	}

	return prefix + strings.Join(lines, "\n")
}

// 打印服务器发送的消息
func printIncoming(message string) {
	fmt.Printf("\r\033[2K")
	ui.PrintChat(message)
	printPrompt()
}

// 设置用户输入提示词
func setPrompt(prompt string) {
	promptLock.Lock()
	currentPrompt = prompt
	promptLock.Unlock()
}

// 显示用户输入提示词
func printPrompt() {
	promptLock.Lock()
	fmt.Print(currentPrompt)
	promptLock.Unlock()
}
