package tui

import (
	"QMsg/protocol"
	"fmt"
	"strings"
	"time"
)

func FormatIncomingMessage(message protocol.Message) string {
	now := time.Now().Format("15:04")

	switch message.Type {
	case protocol.MessageChat:
		return formatDisplayMessage(fmt.Sprintf("[%s] %s", now, message.From), message.Content)
	case protocol.MessagePrivate:
		return formatDisplayMessage(fmt.Sprintf("[%s] [私聊] %s", now, message.From), message.Content)
	case protocol.MessageSystem:
		return fmt.Sprintf("[%s] [系统] %s", now, message.Content)
	case protocol.MessageError:
		return fmt.Sprintf("[%s] [错误] %s", now, message.Content)
	case protocol.MessageText:
		return message.Content
	default:
		return message.Content
	}
}

func BuildOutgoingMessage(text string) protocol.Message {
	messageType := protocol.MessageChat
	commandText := strings.TrimSpace(text)

	if strings.HasPrefix(commandText, "/") {
		messageType = protocol.MessageCommand
	}

	if strings.HasPrefix(commandText, "@") {
		messageType = protocol.MessagePrivate
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\n", "\\n")

	return protocol.Message{
		Type:    messageType,
		Content: text,
	}
}

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
