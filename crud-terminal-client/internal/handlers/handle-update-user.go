package handlers

import tea "github.com/charmbracelet/bubbletea"

type UpdateUser struct{}

func (c UpdateUser) Init() tea.Cmd {
	return nil
}

func (c UpdateUser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c UpdateUser) View() string {
	return "Create Book"
}
