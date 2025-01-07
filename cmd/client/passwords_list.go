package main

import tea "github.com/charmbracelet/bubbletea"

func PasswordsListUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			_, ok := m.passwordsList.SelectedItem().(item)
			if ok {
				m.state = passwordEditState
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

func PasswordsListView(m model) string {
	//return docStyle.Render(m.passwordsList.View())
	return m.passwordsList.View()
}
