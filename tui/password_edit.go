package tui

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func PasswordsEditUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
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

			if s == "enter" && m.passwordUpdateInputsFocusIndex == len(m.passwordUpdateInputs) {
				m.state = passwordListState
				return m, handlePasswordUpdate(&m)
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

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := updateUpdateInputs(m, msg)

	return m, cmd
}

func updateUpdateInputs(m Model, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.passwordUpdateInputs))

	for i := range m.passwordUpdateInputs {
		m.passwordUpdateInputs[i], cmds[i] = m.passwordUpdateInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func handlePasswordUpdate(m *Model) tea.Cmd {
	encryptedPassword, err := m.crypto.Encrypt(m.passwordUpdateInputs[2].Value())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = m.passwordRepo.UpdatePassword(
		m.currentPasswordID,
		m.passwordUpdateInputs[0].Value(),
		m.passwordUpdateInputs[1].Value(),
		encryptedPassword,
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	updatedItem := item{
		id:       m.currentPasswordID,
		site:     m.passwordUpdateInputs[0].Value(),
		userName: m.passwordUpdateInputs[1].Value(),
		password: encryptedPassword,
	}
	m.passwordsList.SetItem(m.passwordsList.Index(), updatedItem)
	return nil
}

func PasswordsEditView(m Model) string {
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
