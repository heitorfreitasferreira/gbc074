package handlers

import tea "github.com/charmbracelet/bubbletea"

type ReadAllUsers struct{}

func (c ReadAllUsers) Init() tea.Cmd {
	return nil
}

func (c ReadAllUsers) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c ReadAllUsers) View() string {
	return "Create Book"
}
