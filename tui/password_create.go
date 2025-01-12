package tui

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func PasswordsCreateUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = passwordListState
			return m, nil
		case "tab", "shift+tab", "enter", "up", "down", "j", "k":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.passwordCreateInputsFocusIndex--
			} else {
				m.passwordCreateInputsFocusIndex++
			}

			if m.passwordCreateInputsFocusIndex > len(m.passwordCreateInputs) {
				m.passwordCreateInputsFocusIndex = 0
			} else if m.passwordCreateInputsFocusIndex < 0 {
				m.passwordCreateInputsFocusIndex = len(m.passwordCreateInputs)
			}

			if s == "enter" && m.passwordCreateInputsFocusIndex == len(m.passwordCreateInputs) {
				m.state = passwordListState
				return m, handlePasswordCreate(&m)
			}

			cmds := make([]tea.Cmd, len(m.passwordCreateInputs))
			for i := 0; i <= len(m.passwordCreateInputs)-1; i++ {
				if i == m.passwordCreateInputsFocusIndex {
					// Set focused state
					cmds[i] = m.passwordCreateInputs[i].Focus()
					m.passwordCreateInputs[i].PromptStyle = focusedStyle
					m.passwordCreateInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.passwordCreateInputs[i].Blur()
				m.passwordCreateInputs[i].PromptStyle = noStyle
				m.passwordCreateInputs[i].TextStyle = noStyle
			}
			cmds = append(cmds)
			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := updateCreateInputs(m, msg)

	return m, cmd
}

func updateCreateInputs(m Model, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.passwordCreateInputs))

	for i := range m.passwordCreateInputs {
		m.passwordCreateInputs[i], cmds[i] = m.passwordCreateInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func handlePasswordCreate(m *Model) tea.Cmd {
	encryptedPassword, err := m.crypto.Encrypt(m.passwordCreateInputs[2].Value())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	id, err := m.passwordRepo.AddPassword(
		m.passwordCreateInputs[0].Value(),
		m.passwordCreateInputs[1].Value(),
		encryptedPassword)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	newItem := item{
		id:       int(id),
		site:     m.passwordCreateInputs[0].Value(),
		userName: m.passwordCreateInputs[1].Value(),
		password: encryptedPassword,
	}
	cmd := m.passwordsList.InsertItem(len(m.passwordsList.Items()), newItem)
	return cmd
}

func PasswordsCreateView(m Model) string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("Password Create View"))
	b.WriteString("\n\n")
	for i := range m.passwordCreateInputs {
		b.WriteString(m.passwordCreateInputs[i].View())
		if i < len(m.passwordCreateInputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
