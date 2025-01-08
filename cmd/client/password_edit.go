package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func PasswordsEditUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = passwordListState
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down", "j", "k":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.passwordUpdateInputsFocusIndex--
			} else {
				m.passwordUpdateInputsFocusIndex++
			}

			if m.passwordUpdateInputsFocusIndex > len(m.passwordUpdateInputs) {
				m.passwordUpdateInputsFocusIndex = 0
			} else if m.passwordUpdateInputsFocusIndex < 0 {
				m.passwordUpdateInputsFocusIndex = len(m.passwordUpdateInputs)
			}

			cmds := make([]tea.Cmd, len(m.passwordUpdateInputs))
			for i := 0; i <= len(m.passwordUpdateInputs)-1; i++ {
				if i == m.passwordUpdateInputsFocusIndex {
					cmds[i] = m.passwordUpdateInputs[i].Focus()
					m.passwordUpdateInputs[i].PromptStyle = focusedStyle
					m.passwordUpdateInputs[i].TextStyle = focusedStyle
					continue
				}
				m.passwordUpdateInputs[i].Blur()
				m.passwordUpdateInputs[i].PromptStyle = noStyle
				m.passwordUpdateInputs[i].TextStyle = noStyle
			}

			if s == "enter" && m.passwordUpdateInputsFocusIndex == len(m.passwordUpdateInputs) {
				err := updatePassword(
					m.db, m.currentPasswordID,
					m.passwordUpdateInputs[0].Value(),
					m.passwordUpdateInputs[1].Value(),
					m.passwordUpdateInputs[2].Value(),
				)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}

				updatedItem := item{
					id:       m.currentPasswordID,
					site:     m.passwordUpdateInputs[0].Value(),
					userName: m.passwordUpdateInputs[1].Value(),
					password: m.passwordUpdateInputs[2].Value(),
				}
				m.passwordsList.SetItem(m.passwordsList.Index(), updatedItem)
				m.state = passwordListState
			}
			return m, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(m.passwordUpdateInputs))

	for i := range m.passwordUpdateInputs {
		m.passwordUpdateInputs[i], cmds[i] = m.passwordUpdateInputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}
func PasswordsEditView(m model) string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("Password Edit View"))
	b.WriteString("\n\n")
	for i := range m.passwordUpdateInputs {
		b.WriteString(m.passwordUpdateInputs[i].View())
		if i < len(m.passwordUpdateInputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
