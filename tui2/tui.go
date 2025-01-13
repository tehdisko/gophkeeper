package tui2

import (
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tehdisko/gophkeeper/domain"
	"github.com/tehdisko/gophkeeper/repository/password"
)

const (
	loginState = iota
	passwordListState
	passwordEditState
	passwordCreateState
)

var (
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
)

type PasswordRepository interface {
	GetPassword(id int) (domain.Password, error)
	AddPassword(site, username, password string) (int64, error)
	UpdatePassword(id int, site, username, password string) error
	DeletePassword(id int) error
	GetAllPasswords() ([]domain.Password, error)
}

type Model struct {
	dump              io.Writer
	state             int
	passwordRepo      PasswordRepository
	currentPasswordID int

	passwordListView *PasswordListModel
	loginView        *LoginModel
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.passwordListView.PasswordsList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	}

	switch m.state {
	case loginState:
		return m.loginView.Update(msg)
	case passwordListState:
		return m.passwordListView.Update(msg)
	}

	return m, nil
}

func (m *Model) View() string {
	switch m.state {
	case loginState:
		return m.loginView.View()
	case passwordListState:
		return m.passwordListView.View()
	}
	return "empty"
}
func InitialModel() *Model {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	passwordRepo := password.New("passwords")

	m := Model{
		dump:         dump,
		state:        loginState,
		passwordRepo: passwordRepo,
	}

	m.initLoginInputs()
	m.initPasswordList()

	return &m
}

func (m *Model) initLoginInputs() {
	loginInputs := make([]textinput.Model, 2)
	var l textinput.Model
	for i := range loginInputs {
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

		loginInputs[i] = l
	}
	loginView := NewLoginModel()
	loginView.LoginInputs = loginInputs
	m.loginView = loginView
	m.loginView.Parent = m
}

func (m *Model) initPasswordList() {
	allPasswords, err := m.passwordRepo.GetAllPasswords()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	items := make([]list.Item, len(allPasswords))
	for i := range allPasswords {
		items[i] = item{id: allPasswords[i].ID, site: allPasswords[i].Site, userName: allPasswords[i].UserName, password: allPasswords[i].Password}
	}

	passwordsList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	passwordsList.Title = "My Passwords"
	passwordsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeymap.Create,
			DefaultKeymap.Enter,
			DefaultKeymap.Delete,
			DefaultKeymap.Back,
		}
	}
	passwordListView := NewPasswordListModel()
	passwordListView.PasswordsList = passwordsList
	m.passwordListView = passwordListView
}
