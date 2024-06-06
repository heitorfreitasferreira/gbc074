package handlers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type CreateUser struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func NewCreateUser() CreateUser {
	m := CreateUser{
		inputs: make([]textinput.Model, 2),
	}

	m.inputs[0] = textinput.New()
	m.inputs[0].Cursor.Style = cursorStyle
	m.inputs[0].Placeholder = "CPF"
	m.inputs[0].CharLimit = 11
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.inputs[0].TextStyle = focusedStyle

	m.inputs[1] = textinput.New()
	m.inputs[1].Cursor.Style = cursorStyle
	m.inputs[1].Placeholder = "Nome"
	m.inputs[1].CharLimit = 64

	return m
}

func (c CreateUser) Init() tea.Cmd {
	return textinput.Blink
}

func (c CreateUser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return c, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			c.cursorMode++
			if c.cursorMode > cursor.CursorHide {
				c.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(c.inputs))
			for i := range c.inputs {
				cmds[i] = c.inputs[i].Cursor.SetMode(c.cursorMode)
			}
			return c, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && c.focusIndex == len(c.inputs) {
				newUserRequest := user{
					cpf:  c.inputs[0].Value(),
					nome: c.inputs[1].Value(),
				}
				_ = newUserRequest // TODO: send request to create user via gRPC
				return c, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				c.focusIndex--
			} else {
				c.focusIndex++
			}

			if c.focusIndex > len(c.inputs) {
				c.focusIndex = 0
			} else if c.focusIndex < 0 {
				c.focusIndex = len(c.inputs)
			}

			cmds := make([]tea.Cmd, len(c.inputs))
			for i := 0; i <= len(c.inputs)-1; i++ {
				if i == c.focusIndex {
					// Set focused state
					cmds[i] = c.inputs[i].Focus()
					c.inputs[i].PromptStyle = focusedStyle
					c.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				c.inputs[i].Blur()
				c.inputs[i].PromptStyle = noStyle
				c.inputs[i].TextStyle = noStyle
			}

			return c, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := c.updateInputs(msg)

	return c, cmd
}

func (m *CreateUser) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m CreateUser) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
