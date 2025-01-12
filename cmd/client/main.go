package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tehdisko/gophkeeper/tui"
)

func main() {
	if _, err := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
