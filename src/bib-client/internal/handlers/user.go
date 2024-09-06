package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	api "library-manager/shared/api/bib"
)

func blockUsers(client api.PortalBibliotecaClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.BloqueiaUsuarios(ctx, &api.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao bloquear usuários com devoluções pendentes: %v", err)
	}

	fmt.Printf("%v\n", res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func releaseUsers(client api.PortalBibliotecaClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.LiberaUsuarios(ctx, &api.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao liberar usuários com devoluções pendentes: %v", err)
	}

	fmt.Printf("%v\n", res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func listBlockedUsers(client api.PortalBibliotecaClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.ListaUsuariosBloqueados(ctx, &api.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao obter todos os usuários bloqueados: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Nome\tCPF")
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erro ao obter usuários: %v", err)
		}
		fmt.Fprintf(w, "Usuario %v com CPF: %v bloqueado pelo livro %v\n", res.Usuario.Nome, res.Usuario.Cpf, res.Livros[0].Titulo)
	}

	w.Flush()
	fmt.Printf("\n")
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}
