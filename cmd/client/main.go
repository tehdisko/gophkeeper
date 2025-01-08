package main

import (
	"database/sql"
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

type Password struct {
	ID       int
	Site     string
	UserName string
	Password string
}

type model struct {
	dump              io.Writer
	state             int
	db                *sql.DB
	currentPasswordID int

	passwordsList list.Model

	loginInputs                    []textinput.Model
	loginInputsIndex               int
	passwordUpdateInputs           []textinput.Model
	passwordUpdateInputsFocusIndex int
	passwordCreateInputs           []textinput.Model
	passwordCreateInputsFocusIndex int
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
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

func initialModel() model {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	db, err := sql.Open("sqlite3", "./passwords.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	m := model{
		dump:                 dump,
		state:                passwordListState,
		db:                   db,
		loginInputs:          make([]textinput.Model, 2),
		passwordUpdateInputs: make([]textinput.Model, 3),
		passwordCreateInputs: make([]textinput.Model, 3),
	}

	allPasswords, err := getAllPasswords(m.db)

	items := make([]list.Item, len(allPasswords))
	for i := range allPasswords {
		items[i] = item{id: allPasswords[i].ID, site: allPasswords[i].Site, userName: allPasswords[i].UserName, password: allPasswords[i].Password}
	}

	m.passwordsList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.passwordsList.Title = "My Passwords"
	m.passwordsList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Create,
			Keymap.Enter,
			Keymap.Delete,
			Keymap.Back,
		}
	}
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

	var u textinput.Model
	for i := range m.passwordUpdateInputs {
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
			u.EchoMode = textinput.EchoPassword
			u.EchoCharacter = '•'
		}

		m.passwordUpdateInputs[i] = u
	}

	var c textinput.Model
	for i := range m.passwordCreateInputs {
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

		m.passwordCreateInputs[i] = c
	}

	return m
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
