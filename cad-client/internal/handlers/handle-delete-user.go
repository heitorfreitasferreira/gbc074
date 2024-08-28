package handlers

import tea "github.com/charmbracelet/bubbletea"

type DeleteUser struct{}

func (c DeleteUser) Init() tea.Cmd {
	return nil
}

func (c DeleteUser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c DeleteUser) View() string {
	return "Create Book"
}
