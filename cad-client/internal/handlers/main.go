package handlers

import (
	"fmt"
	"os"

	api "library-manager/shared/api/cad"
)

type MyHandler func(api.PortalCadastroClient) error

func Choose() MyHandler {
	options := []struct {
		id      int
		desc    string
		handler MyHandler
	}{
		{
			id:      1,
			desc:    "Novo Usuário",
			handler: createUser,
		},
		{
			id:      2,
			desc:    "Edita Usuário",
			handler: updateUser,
		},
		{
			id:      3,
			desc:    "Remove Usuário",
			handler: removeUser,
		},
		{
			id:      4,
			desc:    "Obtem Usuário",
			handler: readUser,
		},
		{
			id:      5,
			desc:    "Obtem Todos Usuários",
			handler: readAllUsers,
		},
		{
			id:      6,
			desc:    "Novo Livro",
			handler: nil,
		},
		{
			id:      7,
			desc:    "Edita Livro",
			handler: nil,
		},
		{
			id:      8,
			desc:    "Remove Livro",
			handler: nil,
		},
		{
			id:      9,
			desc:    "Obtem Livro",
			handler: nil,
		},
		{
			id:      10,
			desc:    "Obtem Todos Livros",
			handler: nil,
		},
	}
	fmt.Print("\033[H\033[2J")
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
