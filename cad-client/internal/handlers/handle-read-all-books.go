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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func NewReadAllBooks() ReadAllBooks {
	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "ISBN", Width: 6},
			{Title: "Título", Width: 30},
			{Title: "Autor", Width: 20},
			{Title: "Total", Width: 4},
			{Title: "Restante", Width: 4},
		}),
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

	return ReadAllBooks{
		table: t,
	}

}

type listAllBooksMsg []table.Row

func (r ReadAllBooks) Init() tea.Cmd {
	return func() tea.Msg {
		books := []table.Row{
			{
				"978-0-306-40615-7",
				"The Art of Computer Programming",
				"Donald Knuth",
				"7",
				"3",
			},
			{
				"978-0-13-110362-7",
				"Clean Code",
				"Robert C. Martin",
				"5",
				"2",
			},
		}
		return listAllBooksMsg(books)
		// rows := make([]table.Row, len(books))
		// for i, b := range books {
		// 	rows[i] = b
		// }
		// return listAllBooksMsg(rows)
	}
}

func (r ReadAllBooks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case listAllBooksMsg:
		r.table.SetRows([]table.Row(msg))
		return r, tea.Quit
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
	// r.table, cmd = r.table.Update(msg)
	return r, cmd
}

func (r ReadAllBooks) View() string {
	return baseStyle.Render(r.table.View()) + "\n"
}
