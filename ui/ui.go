package ui

// ui.go 包含了控制台界面相关的函数，用于清屏、打印日志、等待用户输入等功能

import (
	"bufio"
	"fmt"
	"os"
)

var serverPromptEnabled = true

// 设置服务端日志后是否重新显示命令提示符
func SetServerPromptEnabled(enabled bool) {
	serverPromptEnabled = enabled
}

// 控制台清屏
func ClearScreen() {
	fmt.Print("\033[H\033[2J") // ANSI 清屏
}

// 清除当前行
func ClearCurrentLine() {
	fmt.Print("\033[1A")
	fmt.Print("\r\033[2K")
}

// 打印服务器日志
func PrintServerLog(message string) { // 主动出发指令输出
	fmt.Printf("\r\033[2K")
	fmt.Println(message)
}
func PrintServerAsyncLog(message string) { // 异步函数(系统事件触发)
	if serverPromptEnabled {
		fmt.Printf("\r\033[2K")
	}
	fmt.Println(message)
	if serverPromptEnabled {
		fmt.Print("server > ")
	}
}

// 等待用户按回车继续
func WaitEnter() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n按回车继续...")
	reader.ReadString('\n')
}

// 打印标题、错误、系统消息和聊天消息
func PrintTitle(title string) {
	fmt.Println("--------------------------------")
	fmt.Println(" " + title)
	fmt.Println("--------------------------------")
}
func PrintSystem(message string) {
	fmt.Println("[系统] " + message)
}
