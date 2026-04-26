package text

import "strings"

// text.go 包含了文本相关的函数，用于构建菜单文本、帮助文本和读取用户输入

// 菜单页面打印
func BuildMenuText() string {
	return strings.TrimSpace(`
------------- QMsg -------------
提示:
  输入 help 查看帮助
--------------------------------
`)
}

// 帮助页面打印
func BuildHelpText() string {
	return strings.TrimSpace(`
------------- help -------------
启动阶段:
  1 或 server       启动服务端
  2 或 client       启动客户端
  help 或 /help     查看帮助

聊天阶段:
  普通文字          发送群聊消息
  /users            查看当前在线用户
  @用户名 消息      给指定用户发送私聊
  help 或 /help     查看帮助
  exit              退出聊天
  /multi            进入多行输入模式
  /send             多行输入模式中发送
  /cancel           多行输入模式中取消
--------------------------------
`)
}
