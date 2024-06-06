package handlers

import tea "github.com/charmbracelet/bubbletea"

type CreateBook struct{}

func (c CreateBook) Init() tea.Cmd {
	return nil
}

func (c CreateBook) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c CreateBook) View() string {
	return "Create Book"
}
