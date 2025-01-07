package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	loginState = iota
	passwordListState
	passwordEditState
)

var (
	titleStyle          = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedLoginButton = focusedStyle.Render("  [Enter]")
	blurredLoginButton = fmt.Sprintf("  [%s]", blurredStyle.Render("Enter"))
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
)

type Password struct {
	ID       int
	Site     string
	UserName string
	Password string
}

type model struct {
	state              int
	loginFocusIndex    int
	passwordsDB        []Password
	passwordsList      list.Model
	passwordFocusIndex int
	loginInputs        []textinput.Model
	passwordInputs     []textinput.Model
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.passwordsList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	}

	switch m.state {
	case loginState:
		return LoginUpdate(msg, m)
	case passwordListState:
		return PasswordsListUpdate(msg, m)
	case passwordEditState:
		return PasswordsEditUpdate(msg, m)
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case loginState:
		return LoginView(m)
	case passwordListState:
		return PasswordsListView(m)
	case passwordEditState:
		return PasswordsEditView(m)
	}
	return fmt.Sprintln("Not implemented")
}

type item struct {
	site, userName, password string
}

func (i item) Title() string       { return i.site }
func (i item) Description() string { return i.userName }
func (i item) FilterValue() string { return i.site }

func initialModel() model {
	m := model{
		state: loginState,
		//loginFocusIndex: 0,
		loginInputs:    make([]textinput.Model, 2),
		passwordInputs: make([]textinput.Model, 3),
	}

	passwords := []Password{
		{ID: 1, Site: "yandex.ru", UserName: "admin", Password: "admin"},
		{ID: 2, Site: "google.com", UserName: "user", Password: "user"},
	}
	m.passwordsDB = passwords
	items := make([]list.Item, len(passwords))
	for i := range passwords {
		items[i] = item{site: passwords[i].Site, userName: passwords[i].UserName, password: passwords[i].Password}
	}

	m.passwordsList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.passwordsList.Title = "My Passwords"

	var l textinput.Model
	for i := range m.loginInputs {
		l = textinput.New()
		switch i {
		case 0:
			l.Placeholder = "Login"
			l.Focus()
			l.PromptStyle = focusedStyle
			l.TextStyle = focusedStyle
		case 1:
			l.Placeholder = "Password"
			l.EchoMode = textinput.EchoPassword
			l.EchoCharacter = '•'
		}

		m.loginInputs[i] = l
	}

	var t textinput.Model
	for i := range m.passwordInputs {
		t = textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Site"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "UserName"
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.passwordInputs[i] = t
	}

	return m
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
