package handlers

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
)

type book struct {
	Isbn     string
	Titulo   string
	Autor    string
	Total    int
	Restante int
}

func (b book) toBubbleRow() table.Row {
	return table.Row{
		b.Isbn,
		b.Titulo,
		b.Autor,
		strconv.Itoa(b.Total),
		strconv.Itoa(b.Restante),
	}
}
