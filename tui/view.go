package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

func (m Model) View() string {
	if !m.ready {
		return "正在初始化 QMsg..."
	}

	header := titleStyle.Render("QMsg") + " " +
		statusStyle.Render(fmt.Sprintf("用户: %s  服务器: %s", m.username, m.address))

	footer := m.input.View() + "\n" + helpStyle.Render("Enter 发送 | Alt+Enter/Ctrl+J 换行 | /help /users /exit | 鼠标滚轮/PgUp/PgDn 滚动")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}
