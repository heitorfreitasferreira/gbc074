package screens

import tea "github.com/charmbracelet/bubbletea"

type OptionPicker struct {
	Choices  []string
	Cursor   int
	SubModel []tea.Model
}

func (m OptionPicker) Init() tea.Cmd {
	return nil
}

func (m OptionPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter", " ":
			return m.SubModel[m.Cursor], nil
		}
	}
	return m, nil
}

func (m OptionPicker) View() string {
	baseText := "Choose an option: \n\n"
	for i, choice := range m.Choices {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}
		baseText += cursor + " " + choice + "\n"
	}
	return baseText
}
