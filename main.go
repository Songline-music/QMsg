package main

// main.go 负责处理用户输入以及控制台显示

import (
	"QMsg/client"
	"QMsg/config"
	"QMsg/server"
	"QMsg/text"
	"QMsg/ui"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 程序入口
func main() {
	// reader 对象读取用户输入
	reader := bufio.NewReader(os.Stdin)

	for {
		// 打印菜单
		ui.ClearScreen()
		fmt.Println(text.BuildMenuText())

		mode := readLine(reader, "请选择模式(1.服务端, 2.客户端, help.帮助): ")

		switch strings.ToLower(mode) {
		case "1", "server":
			runServerMode(reader)
		case "2", "client":
			runClientMode(reader)
		case "help", "/help":
			ui.ClearScreen()
			fmt.Println(text.BuildHelpText())
			ui.WaitEnter()
		case "exit":
			fmt.Println("程序已退出")
			ui.WaitEnter()
			return
		default:
			ui.ClearScreen()
			fmt.Println("模式输入错误,请重新输入")
			ui.WaitEnter()
		}
	}
}

// 运行服务端模式
func runServerMode(reader *bufio.Reader) {
	ui.ClearScreen()
	address := readLineWithDefault(reader, "请输入监听地址: ", ":9000")

	token := readRequiredLine(reader, "请设置房间密码: ")

	server.RunServer(address, token)
	ui.WaitEnter()
}

// 运行客户端模式
func runClientMode(reader *bufio.Reader) {
	ui.ClearScreen()

	cfg, hasConfig := config.LoadConfig()
	if hasConfig {
		cfg = config.NormalizeConfig(cfg)

		if config.IsClientConfigReady(cfg) {
			fmt.Println("[系统] 已读取配置文件")
			fmt.Println("[系统] 服务器地址: " + cfg.ServerAddress)
			fmt.Println("[系统] 用户名: " + cfg.Username)

			useConfig := readLineWithDefault(reader, "是否直接连接(Y/N): ", "Y")
			if strings.ToLower(useConfig) == "y" {
				connectClient(cfg)
				return
			}
		}
	}

	cfg = readClientConfig(reader)
	connectClient(cfg)
}

// 读取客户端配置
func readClientConfig(reader *bufio.Reader) config.Config {
	address := readLineWithDefault(reader, "请输入服务端地址: ", config.DefaultServerAddress)
	name := readLine(reader, "请输入你的用户名: ")
	token := readRequiredLine(reader, "请输入房间密码: ")

	if name == "" {
		name = "匿名用户"
	}

	return config.Config{
		ServerAddress: address,
		Username:      name,
		Token:         token,
	}
}

// 连接服务器并进入聊天
func connectClient(cfg config.Config) {
	if cfg.Username == "" {
		cfg.Username = "匿名用户"
	}

	ok, name := client.RunClient(cfg.ServerAddress, cfg.Username, cfg.Token)
	if ok {
		cfg.Username = name
		config.SaveConfig(cfg)
	}

	ui.WaitEnter()
}

// 读取用户输入
func readLine(reader *bufio.Reader, tip string) string {
	fmt.Print(tip) // 传给用户的输入提示文本
	text, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(text) // 去掉首尾空格
}

// 读取用户输入,如果输入为空则返回默认值
func readLineWithDefault(reader *bufio.Reader, tip string, defaultValue string) string {
	text := readLine(reader, tip)
	if text == "" {
		return defaultValue
	}
	return text
}

// 读取用户输入,直到输入不为空
func readRequiredLine(reader *bufio.Reader, tip string) string {
	for {
		text := readLine(reader, tip)
		if text != "" {
			return text
		}

		fmt.Println("输入不能为空")
	}
}
