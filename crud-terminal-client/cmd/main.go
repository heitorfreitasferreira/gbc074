package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/handlers"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/screens"
)

func main() {
	p := tea.NewProgram(screens.OptionPicker{
		Choices: []string{
			"User",
			"Book",
		},
		SubModel: []tea.Model{
			screens.OptionPicker{
				Choices: []string{
					"Create",
					"Read",
					"Update",
					"Delete",
				},
				SubModel: []tea.Model{
					handlers.NewCreateUser(),
					screens.OptionPicker{
						Choices: []string{
							"By term",
							"List all",
						},
						SubModel: []tea.Model{
							handlers.ReadUserByTerm{},
							handlers.ReadAllUsers{},
						},
					},
					handlers.UpdateUser{},
					handlers.DeleteUser{},
				},
			},
			// Handlers dos livros
			screens.OptionPicker{
				Choices: []string{
					"Create",
					"Read",
					"Update",
					"Delete",
				},
				SubModel: []tea.Model{
					handlers.CreateBook{},
					screens.OptionPicker{
						Choices: []string{
							"By term",
							"List all",
						},
						SubModel: []tea.Model{
							handlers.ReadBookByTerm{},
							handlers.ReadAllBooks{},
						},
					},
					handlers.UpdateBook{},
					handlers.DeleteBook{},
				},
			},
		},
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running client: %v", err)
		os.Exit(1)
	}
}
