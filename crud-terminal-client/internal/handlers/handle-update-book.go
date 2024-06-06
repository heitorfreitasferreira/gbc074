package handlers

import tea "github.com/charmbracelet/bubbletea"

type UpdateBook struct{}

func (c UpdateBook) Init() tea.Cmd {
	return nil
}

func (c UpdateBook) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c UpdateBook) View() string {
	return "Create Book"
}
