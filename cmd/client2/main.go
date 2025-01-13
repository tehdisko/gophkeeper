package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tehdisko/gophkeeper/tui2"
)

type App struct {
	tui    *tea.Program
	state  *AppState
	sync   *SyncManager
	crypto *CryptoManager
	view   map[string]View
}
type AppState interface {
}

type SyncManager interface {
	Sync() error
}

type CryptoManager interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type View interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type PasswordListView struct {
	app *App
}

func main() {
	if _, err := tea.NewProgram(tui2.InitialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
