package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func PasswordsEditUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down", "j", "k":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.passwordFocusIndex--
			} else {
				m.passwordFocusIndex++
			}

			if m.passwordFocusIndex > len(m.passwordInputs) {
				m.passwordFocusIndex = 0
			} else if m.passwordFocusIndex < 0 {
				m.passwordFocusIndex = len(m.passwordInputs)
			}

			cmds := make([]tea.Cmd, len(m.passwordInputs))
			for i := 0; i <= len(m.passwordInputs)-1; i++ {
				if i == m.passwordFocusIndex {
					// Set focused state
					cmds[i] = m.passwordInputs[i].Focus()
					m.passwordInputs[i].PromptStyle = focusedStyle
					m.passwordInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.passwordInputs[i].Blur()
				m.passwordInputs[i].PromptStyle = noStyle
				m.passwordInputs[i].TextStyle = noStyle
			}

			if s == "enter" && m.passwordFocusIndex == len(m.passwordInputs) {
				m.state = passwordListState
			}
			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking

	cmds := make([]tea.Cmd, len(m.passwordInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.passwordInputs {
		m.passwordInputs[i], cmds[i] = m.passwordInputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}
func PasswordsEditView(m model) string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render("Password Edit View"))
	b.WriteString("\n\n")
	for i := range m.passwordInputs {
		b.WriteString(m.passwordInputs[i].View())
		if i < len(m.passwordInputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredLoginButton
	if m.passwordFocusIndex == len(m.passwordInputs) {
		button = &focusedLoginButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
