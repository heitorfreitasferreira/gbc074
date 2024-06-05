package screens

import tea "github.com/charmbracelet/bubbletea"

type BookCommands struct {
	nextScreen    map[string]tea.Model
	screenOptions []string
}

func NewBookCommands() BookCommands {
	commands := map[string]tea.Model{
		// "c": HandleCreateUser{},
		// "r": HandleReadUser{},
		// "u": HandleUpdateUser{},
		// "d": HandleDeleteUser{},
	}
	options := make([]string, 0, len(commands))
	for key := range commands {
		options = append(options, key)
	}

	return BookCommands{
		nextScreen:    commands,
		screenOptions: options,
	}
}

func (m BookCommands) Init() tea.Cmd {
	return nil
}

func (m BookCommands) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		for _, key := range m.screenOptions {
			if key == msg.String() {
				return m.nextScreen[key], nil
			}
		}
	}

	return m, nil
}

func (m BookCommands) View() string {
	return "\n\n Choose an operation to perform on Book: \n\n" +
		"  (c)reate\n" +
		"  (r)ead\n" +
		"  (u)pdate\n" +
		"  (d)elete\n"
}
