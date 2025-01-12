package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func PasswordsListUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.state = loginState
			return m, nil
		}
		if msg.String() == "q" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			selectedItem, ok := m.passwordsList.SelectedItem().(item)
			if ok {
				m.currentPasswordID = selectedItem.id
				password, err := m.passwordRepo.GetPassword(m.currentPasswordID)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				decryptedPassword, err := m.crypto.Decrypt(password.Password)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				m.passwordUpdateInputs[0].SetValue(password.Site)
				m.passwordUpdateInputs[1].SetValue(password.UserName)
				m.passwordUpdateInputs[2].SetValue(decryptedPassword)
				m.state = passwordEditState
			}
			return m, nil
		}
		if msg.String() == "c" {
			m.state = passwordCreateState
			return m, nil
		}
		if msg.String() == "d" {
			selectedItem, ok := m.passwordsList.SelectedItem().(item)
			if ok {
				m.currentPasswordID = selectedItem.id
				err := m.passwordRepo.DeletePassword(m.currentPasswordID)
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				m.passwordsList.RemoveItem(m.passwordsList.Index())
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.passwordsList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	}

	var cmd tea.Cmd
	m.passwordsList, cmd = m.passwordsList.Update(msg)

	return m, cmd
}

func PasswordsListView(m Model) string {
	return m.passwordsList.View()
}
