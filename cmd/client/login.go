package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func LoginUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
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
				m.loginFocusIndex--
			} else {
				m.loginFocusIndex++
			}

			if m.loginFocusIndex > len(m.loginInputs) {
				m.loginFocusIndex = 0
			} else if m.loginFocusIndex < 0 {
				m.loginFocusIndex = len(m.loginInputs)
			}

			cmds := make([]tea.Cmd, len(m.loginInputs))
			for i := 0; i <= len(m.loginInputs)-1; i++ {
				if i == m.loginFocusIndex {
					// Set focused state
					cmds[i] = m.loginInputs[i].Focus()
					m.loginInputs[i].PromptStyle = focusedStyle
					m.loginInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.loginInputs[i].Blur()
				m.loginInputs[i].PromptStyle = noStyle
				m.loginInputs[i].TextStyle = noStyle
			}

			if s == "enter" && m.loginFocusIndex == len(m.loginInputs) {
				m.state = passwordListState
			}
			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking

	cmds := make([]tea.Cmd, len(m.loginInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
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

	button := &blurredLoginButton
	if m.loginFocusIndex == len(m.loginInputs) {
		button = &focusedLoginButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
