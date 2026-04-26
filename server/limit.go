package server

// limit.go 包含了输入限制相关的函数，用于限制用户名长度和消息长度

const (
	maxUsernameLength = 20
	maxMessageLength  = 200
)

// 判断用户名是否过长
func isUsernameTooLong(name string) bool {
	return len([]rune(name)) > maxUsernameLength
}

// 判断消息是否过长
func isMessageTooLong(message string) bool {
	return len([]rune(message)) > maxMessageLength
}
