package handlers

import tea "github.com/charmbracelet/bubbletea"

type ReadBookByTerm struct{}

func (c ReadBookByTerm) Init() tea.Cmd {
	return nil
}

func (c ReadBookByTerm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c ReadBookByTerm) View() string {
	return "Create Book"
}
