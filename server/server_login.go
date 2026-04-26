package server

// server_login.go 包含了处理客户端登录的函数,用于验证登录信息和管理用户名列表

import (
	"QMsg/protocol"
	"QMsg/ui"
	"bufio"
	"fmt"
	"net"
	"strings"
)

// 读取用户登入信息
func readLogin(reader *bufio.Reader, token string, conn net.Conn) (string, bool) {
	remoteAddr := conn.RemoteAddr().String()

	loginMessage, err := protocol.ReadMessage(reader)
	if err != nil {
		ui.PrintServerAsyncLog("登录失败 " + remoteAddr + " 读取登录消息失败 " + err.Error())
		protocol.SendMessage(conn, protocol.Message{
			Type:    protocol.MessageError,
			Content: "读取登录消息失败",
		})
		return "", false
	}

	if loginMessage.Type != protocol.MessageLogin {
		ui.PrintServerAsyncLog("登录失败 " + remoteAddr + " 登录消息格式错误")
		protocol.SendMessage(conn, protocol.Message{
			Type:    protocol.MessageError,
			Content: "登录消息格式错误",
		})
		return "", false
	}

	if loginMessage.Token != token {
		ui.PrintServerAsyncLog("登录失败 " + remoteAddr + " 房间密码错误")
		protocol.SendMessage(conn, protocol.Message{
			Type:    protocol.MessageError,
			Content: "房间密码错误",
		})
		return "", false
	}

	name := strings.TrimSpace(loginMessage.Content)
	if isUsernameTooLong(name) {
		ui.PrintServerAsyncLog("登陆失败 " + remoteAddr + " 用户名过长")
		protocol.SendMessage(conn, protocol.Message{
			Type:    protocol.MessageError,
			Content: "用户名不能超过 " + fmt.Sprint(maxUsernameLength) + " 个字符",
		})
		return "", false
	}

	name, ok := prepareClientName(name)
	if !ok {
		ui.PrintServerAsyncLog("登陆失败 " + remoteAddr + " 用户名重复 " + name)
		protocol.SendMessage(conn, protocol.Message{
			Type:    protocol.MessageError,
			Content: "用户名已存在,请更换用户名",
		})
		return "", false
	}

	return name, true
}
