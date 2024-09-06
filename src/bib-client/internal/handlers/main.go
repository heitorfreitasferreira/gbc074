package handlers

import (
	"fmt"
	"os"

	api "library-manager/shared/api/bib"
)

type MyHandler func(api.PortalBibliotecaClient) error

func Choose(host, port string) MyHandler {
	options := []struct {
		id      int
		desc    string
		handler MyHandler
	}{
		{
			id:      1,
			desc:    "Realizar Empréstimo",
			handler: borrowBook,
		},
		{
			id:      2,
			desc:    "Realizar Devolução",
			handler: returnBook,
		},
		{
			id:      3,
			desc:    "Bloquear Usuários",
			handler: blockUsers,
		},
		{
			id:      4,
			desc:    "Liberar Usuários",
			handler: releaseUsers,
		},
		{
			id:      5,
			desc:    "Listar Usuários Bloqueados",
			handler: listBlockedUsers,
		},
		{
			id:      6,
			desc:    "Listar Livros Emprestados",
			handler: listBorrowedBooks,
		},
		{
			id:      7,
			desc:    "Listar Livros Em Falta",
			handler: listMissingBooks,
		},
		{
			id:      8,
			desc:    "Pesquisar Livros",
			handler: searchBook,
		},
	}

	fmt.Print("\033[H\033[2J")
	fmt.Printf("Cliente cadastro %d conectado ao servidor %s:%s\n\n", os.Getpid(), host, port)
	for _, opt := range options {
		fmt.Printf("%d - %s\n", opt.id, opt.desc)
	}
	var opt int
	fmt.Scan(&opt)
	if opt < 1 || opt > len(options) {
		fmt.Printf("Opção invalida, tchau")
		os.Exit(0)
	}
	fmt.Print("\033[H\033[2J")
	return options[opt-1].handler
}
