package tui

import (
	"QMsg/protocol"
	"bufio"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.input.SetWidth(msg.Width)

		contentHeight := msg.Height - m.input.Height() - 3
		if contentHeight < 1 {
			contentHeight = 1
		}

		m.viewport = viewport.New(msg.Width, contentHeight)
		m.ready = true
		m.syncViewport()

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.stickBottom = false
			m.viewport.LineUp(3)
			return m, nil
		case tea.MouseWheelDown:
			m.viewport.LineDown(3)
			m.stickBottom = m.viewport.AtBottom()
			return m, nil
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "alt+up", "k":
			m.stickBottom = false
			m.viewport.LineUp(1)
			return m, nil
		case "alt+down", "j":
			m.viewport.LineDown(1)
			m.stickBottom = m.viewport.AtBottom()
			return m, nil
		case "pgup":
			m.stickBottom = false
			m.viewport.HalfViewUp()
			return m, nil
		case "pgdown":
			m.viewport.HalfViewDown()
			m.stickBottom = m.viewport.AtBottom()
			return m, nil
		case "home":
			m.stickBottom = false
			m.viewport.GotoTop()
			return m, nil
		case "end":
			m.stickBottom = true
			m.viewport.GotoBottom()
			return m, nil
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				return m, nil
			}

			lowerText := strings.ToLower(text)
			if lowerText == "/exit" {
				return m, tea.Quit
			}

			m.input.Reset()
			cmds = append(cmds, sendOutgoingMessage(m.conn, text))
			return m, tea.Batch(cmds...)
		}

	case serverMessageMsg:
		message := protocol.Message(msg)
		m.messages = append(m.messages, FormatIncomingMessage(message))
		m.syncViewport()
		cmds = append(cmds, waitForServerMessage(m.reader))

	case serverErrorMsg:
		m.messages = append(m.messages, "[系统] 连接已断开: "+error(msg).Error())
		m.syncViewport()
		return m, tea.Quit

	case sendErrorMsg:
		m.messages = append(m.messages, "[错误] 发送失败: "+error(msg).Error())
		m.syncViewport()

	case sentMessageMsg:
		// 发送成功不需要额外显示，服务端广播回来后会显示。
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func waitForServerMessage(reader *bufio.Reader) tea.Cmd {
	return func() tea.Msg {
		message, err := protocol.ReadMessage(reader)
		if err != nil {
			return serverErrorMsg(err)
		}

		return serverMessageMsg(message)
	}
}

func sendOutgoingMessage(conn net.Conn, text string) tea.Cmd {
	return func() tea.Msg {
		message := BuildOutgoingMessage(text)
		err := protocol.SendMessage(conn, message)
		if err != nil {
			return sendErrorMsg(err)
		}

		return sentMessageMsg{}
	}
}

func (m *Model) syncViewport() {
	if !m.ready {
		return
	}

	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	if m.stickBottom {
		m.viewport.GotoBottom()
	}
}
