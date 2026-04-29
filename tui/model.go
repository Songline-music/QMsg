package tui

import (
	"QMsg/protocol"
	"bufio"
	"net"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type serverMessageMsg protocol.Message
type serverErrorMsg error
type sendErrorMsg error
type sentMessageMsg struct{}

type Model struct {
	conn     net.Conn
	reader   *bufio.Reader
	username string
	address  string

	input    textarea.Model
	viewport viewport.Model
	messages []string

	width       int
	height      int
	ready       bool
	stickBottom bool
}

func NewModel(conn net.Conn, reader *bufio.Reader, username string, address string) Model {
	input := textarea.New()
	input.Placeholder = "输入消息或 /help..."
	input.Prompt = "> "
	input.ShowLineNumbers = false
	input.SetHeight(3)
	input.KeyMap.InsertNewline = key.NewBinding(key.WithKeys("alt+enter", "ctrl+j"))
	input.Focus()

	return Model{
		conn:        conn,
		reader:      reader,
		username:    username,
		address:     address,
		input:       input,
		messages:    []string{"[系统] 已连接到服务器 " + address},
		stickBottom: true,
	}
}

func Run(conn net.Conn, reader *bufio.Reader, username string, address string) error {
	model := NewModel(conn, reader, username, address)
	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		waitForServerMessage(m.reader),
	)
}
