package main

// main.go 负责处理用户输入以及控制台显示

import (
	"QMsg/client"
	"QMsg/config"
	"QMsg/server"
	"QMsg/text"
	"QMsg/ui"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// 程序入口
func main() {
	modeFlag := flag.String("mode", "", "启动模式: server 或 client")
	addrFlag := flag.String("addr", "", "服务端监听地址或客户端连接地址")
	nameFlag := flag.String("name", "", "客户端用户名")
	tokenFlag := flag.String("token", "", "房间密码")
	flag.Parse()

	if *modeFlag != "" {
		runFlagMode(*modeFlag, *addrFlag, *nameFlag, *tokenFlag)
		return
	}

	// reader 对象读取用户输入
	reader := bufio.NewReader(os.Stdin)

	for {
		// 打印菜单
		ui.ClearScreen()
		fmt.Println(text.BuildMenuText())

		mode := readLine(reader, "请选择操作(1/2/3/0): ")

		switch strings.ToLower(mode) {
		case "1", "client":
			runClientMode(reader)
		case "2", "server":
			runServerMode(reader)
		case "3", "help", "/help":
			ui.ClearScreen()
			fmt.Println(text.BuildMenuHelpText())
			ui.WaitEnter()
		case "0", "exit", "quit":
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

// 通过命令行参数启动
func runFlagMode(mode string, address string, name string, token string) {
	switch strings.ToLower(mode) {
	case "server", "2":
		if address == "" {
			address = ":9000"
		}

		if strings.TrimSpace(token) == "" {
			fmt.Println("服务端启动失败: 房间密码不能为空")
			return
		}

		server.RunServer(address, token, false)

	case "client", "1":
		cfg, hasConfig := config.LoadConfig()
		if hasConfig {
			cfg = config.NormalizeConfig(cfg)
		} else {
			cfg = config.Config{
				ServerAddress: config.DefaultServerAddress,
				Username:      "匿名用户",
			}
		}

		if address != "" {
			cfg.ServerAddress = address
		}
		if name != "" {
			cfg.Username = name
		}
		if token != "" {
			cfg.Token = token
		}

		cfg = config.NormalizeConfig(cfg)

		if cfg.Username == "" {
			cfg.Username = "匿名用户"
		}

		if !config.IsClientConfigReady(cfg) {
			fmt.Println("客户端启动失败: 缺少服务端地址或房间密码")
			return
		}

		connectClient(cfg)

	default:
		fmt.Println("未知启动模式: " + mode)
	}
}

// 运行服务端模式
func runServerMode(reader *bufio.Reader) {
	ui.ClearScreen()
	address := readLineWithDefault(reader, "请输入监听地址: ", ":9000")

	token := readRequiredLine(reader, "请设置房间密码: ")

	server.RunServer(address, token, true)
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
		return
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
