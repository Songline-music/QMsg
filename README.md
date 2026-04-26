# QMsg

QMsg 是一个正在开发中的 Go 语言简易聊天项目，目标是从一个命令行聊天程序逐步扩展成轻量级 QQ 风格应用。

当前项目还不是正式成品，主要用于学习 TCP 通信、JSON 消息协议、客户端/服务端分离、云服务器部署和基础项目工程化。

## 当前功能

- TCP 客户端和服务端通信
- JSON 格式消息协议
- 房间密码校验
- 群聊消息
- 私聊消息
- 在线用户列表
- 多行消息输入
- 本地配置文件保存
- 服务端基础日志

## 项目结构

```text
QMsg/
  client/      客户端连接、输入、消息显示
  config/      配置文件读取和保存
  protocol/    消息结构体和 JSON 读写
  server/      服务端监听、登录校验、会话管理、命令处理
  text/        菜单和帮助文本
  ui/          控制台显示相关函数
  main.go      程序入口
```

## 本地运行

先编译：

```bash
go build -o qmsg
```

运行：

```bash
./qmsg
```

Windows 下可以是：

```powershell
go build -o qmsg.exe
.\qmsg.exe
```

启动后可以选择：

```text
1. 服务端
2. 客户端
```

服务端监听地址通常使用：

```text
:9000
```

客户端连接地址通常使用：

```text
服务器公网IP:9000
```

## 配置文件

项目会读取本地的 `config.json`，用于保存客户端连接信息。

请不要把真实的 `config.json` 提交到公开仓库。建议只保留示例文件：

```json
{
  "server_address": "127.0.0.1:9000",
  "username": "your-name",
  "token": "your-room-token"
}
```

推荐 `.gitignore` 中包含：

```gitignore
config.json
qmsg
qmsg.exe
*.log
```

## 云服务器部署

服务器系统推荐使用 Ubuntu。

第一次部署可以在服务器上执行：

```bash
cd /root
git clone https://github.com/Songline-music/QMsg.git
cd QMsg
go build -o qmsg
./qmsg
```

启动服务端时：

```text
监听地址输入 :9000
房间密码自行设置
```

同时需要在云服务器安全组中放行：

```text
TCP 9000
```

## 更新服务器版本

本地修改代码后，先提交并推送到 GitHub：

```bash
git add .
git commit -m "update qmsg"
git push
```

然后登录服务器，更新代码并重新编译：

```bash
cd /root/QMsg
git pull
go build -o qmsg
```

如果旧服务端正在运行，需要先停止旧进程，再启动新版本。

当前阶段可以手动运行：

```bash
./qmsg
```

后续计划改成 `systemd` 后台运行，实现开机自启和崩溃自动重启。

## 未完成计划

- 命令行参数启动，例如 `./qmsg -mode server -addr :9000 -token xxx`
- systemd 后台运行
- 更安全的密码输入方式
- 更完善的断线重连
- 服务端踢人和管理命令
- 聊天记录保存
- 账号系统
- 更完整的图形界面或 Web 界面

## 注意

这是一个学习中的项目，目前不建议直接用于正式生产环境。

如果项目公开，请不要提交真实服务器 IP、房间密码、SSH 密钥或其他敏感信息。
