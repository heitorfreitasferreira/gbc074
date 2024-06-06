package handlers

import tea "github.com/charmbracelet/bubbletea"

type ReadUserByTerm struct{}

func (c ReadUserByTerm) Init() tea.Cmd {
	return nil
}

func (c ReadUserByTerm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c ReadUserByTerm) View() string {
	return "Create Book"
}
