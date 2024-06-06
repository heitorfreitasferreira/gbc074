package handlers

import tea "github.com/charmbracelet/bubbletea"

type DeleteBook struct{}

func (c DeleteBook) Init() tea.Cmd {
	return nil
}

func (c DeleteBook) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c DeleteBook) View() string {
	return "Create Book"
}
