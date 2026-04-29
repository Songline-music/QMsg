package text

import "strings"

// text.go 包含了文本相关的函数，用于构建菜单文本、帮助文本和读取用户输入

// 菜单页面打印
func BuildMenuText() string {
	return strings.TrimSpace(`
------------- QMsg -------------
1. 启动客户端
2. 启动服务端
3. 查看帮助
0. 退出程序
--------------------------------
`)
}

// 选择界面帮助
func BuildMenuHelpText() string {
	return strings.TrimSpace(`
------------- help -------------
启动阶段:
  1 或 client       启动客户端
  2 或 server       启动服务端
  3 或 /help        查看帮助
  0 或 exit         退出程序
--------------------------------
`)
}

// 客户端帮助
func BuildClientHelpText() string {
	return strings.TrimSpace(`
------------- client help -------------
聊天阶段:
  普通文字          发送群聊消息
  /users            查看当前在线用户
  @用户名 消息      给指定用户发送私聊
  /help             查看帮助
  Esc 或 Ctrl+C     退出聊天
  /exit             退出聊天
--------------------------------
`)
}

// 服务端帮助
func BuildServerHelpText() string {
	return strings.TrimSpace(`
------------- server help -------------
服务端命令:
  /help             查看帮助
  /shutdown         关闭服务端
--------------------------------
`)
}
