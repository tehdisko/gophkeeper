package tui2

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	LoginInputs      []textinput.Model
	LoginInputsIndex int
	Parent           *Model
}

func NewLoginModel() *LoginModel {
	return &LoginModel{}
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down", "j", "k":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.LoginInputsIndex--
			} else {
				m.LoginInputsIndex++
			}

			if m.LoginInputsIndex > len(m.LoginInputs) {
				m.LoginInputsIndex = 0
			} else if m.LoginInputsIndex < 0 {
				m.LoginInputsIndex = len(m.LoginInputs)
			}

			cmds := make([]tea.Cmd, len(m.LoginInputs))
			for i := 0; i <= len(m.LoginInputs)-1; i++ {
				if i == m.LoginInputsIndex {
					cmds[i] = m.LoginInputs[i].Focus()
					m.LoginInputs[i].PromptStyle = focusedStyle
					m.LoginInputs[i].TextStyle = focusedStyle
					continue
				}
				m.LoginInputs[i].Blur()
				m.LoginInputs[i].PromptStyle = noStyle
				m.LoginInputs[i].TextStyle = noStyle
			}

			if s == "enter" && m.LoginInputsIndex == len(m.LoginInputs) {
				m.Parent.state = passwordListState
				return m.Parent.Update(msg)
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(m.LoginInputs))

	for i := range m.LoginInputs {
		m.LoginInputs[i], cmds[i] = m.LoginInputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *LoginModel) View() string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("Login View"))
	b.WriteString("\n\n")
	for i := range m.LoginInputs {
		b.WriteString(m.LoginInputs[i].View())
		if i < len(m.LoginInputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
