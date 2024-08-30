package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	api "library-manager/shared/api/cad"
)

func createUser(client api.PortalCadastroClient) error {
	var name, cpf string
	fmt.Println("Digite o nome do usuário:")
	fmt.Scan(&name)
	fmt.Println("Digite o CPF do usuário:")
	fmt.Scan(&cpf)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user := &api.Usuario{
		Nome: name,
		Cpf:  cpf,
	}
	res, err := client.NovoUsuario(ctx, user)
	if err != nil || res.Status == 1 {
		return fmt.Errorf("erro ao criar usuario: %v", err)
	}
	fmt.Printf("Usuário %v criado com sucesso: %v\n", user, res)

	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func readUser(client api.PortalCadastroClient) error {
	var cpf string
	fmt.Println("Digite o CPF do usuário:")
	fmt.Scan(&cpf)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.ObtemUsuario(ctx, &api.Identificador{Id: cpf})
	if err != nil {
		return fmt.Errorf("erro ao obter usuario: %v", err)
	}
	fmt.Printf("Usuário encontrado: %v\n", res)
	fmt.Println("Pressione qualquer tecla para continuar...")
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func readAllUsers(client api.PortalCadastroClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.ObtemTodosUsuarios(ctx, &api.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao obter usuarios: %v", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Nome\tCPF")
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erro ao obter usuarios: %v", err)
		}
		fmt.Fprintf(w, "%v\t%v\n", user.Nome, user.Cpf)
	}
	w.Flush()
	fmt.Printf("\n")
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func updateUser(client api.PortalCadastroClient) error {
	var name, cpf string

	fmt.Println("Digite o CPF do usuário:")
	fmt.Scan(&cpf)
	fmt.Println("Digite o novo nome do usuário:")
	fmt.Scan(&name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user := &api.Usuario{
		Nome: name,
		Cpf:  cpf,
	}
	res, err := client.EditaUsuario(ctx, user)
	if err != nil {
		return fmt.Errorf("erro ao atualizar usuario: %v", err)
	}
	fmt.Printf("Usuário %v atualizado com sucesso: %v\n", user, res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func removeUser(client api.PortalCadastroClient) error {
	var cpf string
	fmt.Println("Digite o CPF do usuário:")
	fmt.Scan(&cpf)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RemoveUsuario(ctx, &api.Identificador{Id: cpf})
	if err != nil || res.Status == 1 {
		return fmt.Errorf("erro ao remover usuario: %v", err)
	}
	fmt.Printf("Usuário %v removido com sucesso: %v\n", cpf, res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}
