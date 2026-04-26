package protocol

// message.go 包含了消息相关的函数，用于定义消息结构体和发送/接收消息的函数

import (
	"bufio"
	"encoding/json"
	"net"
)

// 消息类型常量
const (
	MessageLogin   = "login"
	MessageChat    = "chat"
	MessagePrivate = "private"
	MessageCommand = "command"
	MessageSystem  = "system"
	MessageOK      = "ok"
	MessageError   = "error"
	MessageText    = "text"
)

// 消息结构体
type Message struct {
	Type    string `json:"type"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Content string `json:"content,omitempty"`
	Token   string `json:"token,omitempty"`
}

// 发送信息到服务器
func SendMessage(conn net.Conn, message Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = conn.Write(append(data, '\n'))
	return err
}

// 从服务器接收信息
func ReadMessage(reader *bufio.Reader) (Message, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return Message{}, err
	}

	var message Message
	err = json.Unmarshal([]byte(line), &message)
	if err != nil {
		return Message{}, err
	}

	return message, nil
}
