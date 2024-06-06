package handlers

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
message Livro {
  // ISBN do livro (chave)
  string isbn    = 1;
  string titulo  = 2;
  string autor   = 3;
  int32 total    = 4;
  // campo presente apenas no portal biblioteca
  int32 restante = 5;
}*/

type ReadAllBooks struct {
	table table.Model
}

type myMsg []table.Row

func (r ReadAllBooks) Init() tea.Cmd {
	return func() tea.Msg {
		books := []book{
			{
				Isbn:     "978-0-306-40615-7",
				Titulo:   "The Art of Computer Programming",
				Autor:    "Donald Knuth",
				Total:    7,
				Restante: 3,
			},
			{
				Isbn:     "978-0-13-110362-7",
				Titulo:   "Clean Code",
				Autor:    "Robert C. Martin",
				Total:    5,
				Restante: 2,
			},
		}

		rows := make([]table.Row, len(books))
		for i, b := range books {
			rows[i] = b.toBubbleRow()
		}

		return rows

		// return ReadAllBooks{table: table, loaded: true}
	}
}

func (r ReadAllBooks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case myMsg:
		t := table.New(
			table.WithColumns([]table.Column{
				{Title: "ISBN", Width: 6},
				{Title: "Título", Width: 30},
				{Title: "Autor", Width: 20},
				{Title: "Total", Width: 4},
				{Title: "Restante", Width: 4},
			}),
			table.WithRows(msg),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)
		return ReadAllBooks{table: t}, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if r.table.Focused() {
				r.table.Blur()
			} else {
				r.table.Focus()
			}
		case "q", "ctrl+c":
			return r, tea.Quit
		case "enter":
			return r, tea.Batch(
				tea.Printf("Let's go to %s!", r.table.SelectedRow()[1]),
			)
		}
	}

	r.table, cmd = r.table.Update(msg)
	return r, cmd
}

func (r ReadAllBooks) View() string {
	quitMessage := "\n\nPress 'q' to quit or 'esc'."
	return r.table.View() + quitMessage
}
