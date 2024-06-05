package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/screens"
)

func main() {
	p := tea.NewProgram(screens.EntityPicker{
		Choices: []string{
			"User",
			"Book",
		},
		SubModel: []tea.Model{
			screens.UserCommands{},
		},
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running client: %v", err)
		os.Exit(1)
	}
}
