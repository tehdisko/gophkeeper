package tui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tehdisko/gophkeeper/crypto"
	"github.com/tehdisko/gophkeeper/repository/password"
)

func LoginUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
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
				login := m.loginInputs[0].Value()
				password := m.loginInputs[1].Value()
				initPasswordRepo(login, &m)
				updatePasswordList(&m)
				initCryptoManager(password, &m)
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

func initPasswordRepo(login string, m *Model) {
	passwordRepo := password.New(login)
	m.passwordRepo = passwordRepo
}

func updatePasswordList(m *Model) {
	allPasswords, err := m.passwordRepo.GetAllPasswords()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	items := make([]list.Item, len(allPasswords))
	for i := range allPasswords {
		items[i] = item{id: allPasswords[i].ID, site: allPasswords[i].Site, userName: allPasswords[i].UserName, password: allPasswords[i].Password}
	}
	m.passwordsList.SetItems(items)
}

func initCryptoManager(password string, m *Model) {
	cryptoManager := crypto.NewCryptoManager(password)
	m.crypto = cryptoManager
}

func LoginView(m Model) string {
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
