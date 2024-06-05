package screens

import tea "github.com/charmbracelet/bubbletea"

type UserCommands struct {
	nextScreen    map[string]tea.Model
	screenOptions []string
}

func NewUserCommands() UserCommands {
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

	return UserCommands{
		nextScreen:    commands,
		screenOptions: options,
	}
}

func (m UserCommands) Init() tea.Cmd {
	return nil
}

func (m UserCommands) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m UserCommands) View() string {
	return "\n\n Choose an operation to perform on User: \n\n" +
		"  (c)reate\n" +
		"  (r)ead\n" +
		"  (u)pdate\n" +
		"  (d)elete\n"
}
