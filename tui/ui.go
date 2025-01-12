package tui

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/tehdisko/gophkeeper/domain"
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

type CryptoManager interface {
	Encrypt(msg string) (string, error)
	Decrypt(msg string) (string, error)
}

type Model struct {
	dump              io.Writer
	state             int
	masterPassword    string
	crypto            CryptoManager
	passwordRepo      PasswordRepository
	currentPasswordID int

	passwordsList list.Model

	loginInputs      []textinput.Model
	loginInputsIndex int

	passwordUpdateInputs           []textinput.Model
	passwordUpdateInputsFocusIndex int

	passwordCreateInputs           []textinput.Model
	passwordCreateInputsFocusIndex int
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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
	case passwordCreateState:
		return PasswordsCreateUpdate(msg, m)
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case loginState:
		return LoginView(m)
	case passwordListState:
		return PasswordsListView(m)
	case passwordEditState:
		return PasswordsEditView(m)
	case passwordCreateState:
		return PasswordsCreateView(m)
	}
	return fmt.Sprintln("Not implemented")
}

type item struct {
	id                       int
	site, userName, password string
}

func (i item) Title() string       { return i.site }
func (i item) Description() string { return i.userName }
func (i item) FilterValue() string { return i.site }

func InitialModel() Model {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	//passwordRepo := repository.New("passwords")

	m := Model{
		dump:  dump,
		state: loginState,
		//passwordRepo: passwordRepo,
	}

	initLoginInputs(&m)
	initPasswordList(&m)
	initPasswordUpdateInputs(&m)
	initPasswordCreateInputs(&m)

	return m
}

func initLoginInputs(m *Model) {
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
	m.loginInputs = loginInputs
}

func initPasswordList(m *Model) {
	allPasswords := []domain.Password{{}}

	items := make([]list.Item, len(allPasswords))
	for i := range allPasswords {
		items[i] = item{id: allPasswords[i].ID, site: allPasswords[i].Site, userName: allPasswords[i].UserName, password: allPasswords[i].Password}
	}

	m.passwordsList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.passwordsList.Title = "My Passwords"
	m.passwordsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeymap.Create,
			DefaultKeymap.Enter,
			DefaultKeymap.Delete,
			DefaultKeymap.Back,
		}
	}
}

func initPasswordListBackup(m *Model) {
	allPasswords, err := m.passwordRepo.GetAllPasswords()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	items := make([]list.Item, len(allPasswords))
	for i := range allPasswords {
		items[i] = item{id: allPasswords[i].ID, site: allPasswords[i].Site, userName: allPasswords[i].UserName, password: allPasswords[i].Password}
	}

	m.passwordsList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.passwordsList.Title = "My Passwords"
	m.passwordsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeymap.Create,
			DefaultKeymap.Enter,
			DefaultKeymap.Delete,
			DefaultKeymap.Back,
		}
	}
}

func initPasswordUpdateInputs(m *Model) {
	passwordUpdateInputs := make([]textinput.Model, 3)
	var u textinput.Model
	for i := range passwordUpdateInputs {
		u = textinput.New()
		switch i {
		case 0:
			u.Placeholder = "Site"
			u.Focus()
			u.PromptStyle = focusedStyle
			u.TextStyle = focusedStyle
		case 1:
			u.Placeholder = "UserName"
		case 2:
			u.Placeholder = "Password"
		}
		passwordUpdateInputs[i] = u
	}
	m.passwordUpdateInputs = passwordUpdateInputs
}

func initPasswordCreateInputs(m *Model) {
	passwordCreateInputs := make([]textinput.Model, 3)
	var c textinput.Model
	for i := range passwordCreateInputs {
		c = textinput.New()
		switch i {
		case 0:
			c.Placeholder = "Site"
			c.Focus()
			c.PromptStyle = focusedStyle
			c.TextStyle = focusedStyle
		case 1:
			c.Placeholder = "UserName"
		case 2:
			c.Placeholder = "Password"
			c.EchoMode = textinput.EchoPassword
			c.EchoCharacter = '•'
		}
		passwordCreateInputs[i] = c
	}
	m.passwordCreateInputs = passwordCreateInputs
}
