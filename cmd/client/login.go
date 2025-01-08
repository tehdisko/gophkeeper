package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func LoginUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down", "j", "k":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.loginInputsIndex--
			} else {
				m.loginInputsIndex++
			}

			if m.loginInputsIndex > len(m.loginInputs) {
				m.loginInputsIndex = 0
			} else if m.loginInputsIndex < 0 {
				m.loginInputsIndex = len(m.loginInputs)
			}

			cmds := make([]tea.Cmd, len(m.loginInputs))
			for i := 0; i <= len(m.loginInputs)-1; i++ {
				if i == m.loginInputsIndex {
					cmds[i] = m.loginInputs[i].Focus()
					m.loginInputs[i].PromptStyle = focusedStyle
					m.loginInputs[i].TextStyle = focusedStyle
					continue
				}
				m.loginInputs[i].Blur()
				m.loginInputs[i].PromptStyle = noStyle
				m.loginInputs[i].TextStyle = noStyle
			}

			if s == "enter" && m.loginInputsIndex == len(m.loginInputs) {
				m.state = passwordListState
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(m.loginInputs))

	for i := range m.loginInputs {
		m.loginInputs[i], cmds[i] = m.loginInputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func LoginView(m model) string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("Login View"))
	b.WriteString("\n\n")
	for i := range m.loginInputs {
		b.WriteString(m.loginInputs[i].View())
		if i < len(m.loginInputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
