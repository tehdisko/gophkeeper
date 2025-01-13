package tui2

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	id                       int
	site, userName, password string
}

func (i item) Title() string       { return i.site }
func (i item) Description() string { return i.userName }
func (i item) FilterValue() string { return i.site }

type PasswordListModel struct {
	PasswordsList list.Model
}

func (m *PasswordListModel) Init() tea.Cmd {
	return nil
}

func (m *PasswordListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.PasswordsList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	}

	m.PasswordsList, cmd = m.PasswordsList.Update(msg)
	return m, cmd
}

func (m *PasswordListModel) View() string {
	return m.PasswordsList.View()
}

func NewPasswordListModel() *PasswordListModel {
	return &PasswordListModel{}
}
