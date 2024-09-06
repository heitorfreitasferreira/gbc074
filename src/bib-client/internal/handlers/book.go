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

func borrowBook(client api.PortalBibliotecaClient) error {
	var user, isbn string

	fmt.Println("Digite a chave do usuário:")
	fmt.Scan(&user)
	fmt.Println("Digite o ISBN do livro:")
	fmt.Scan(&isbn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &api.UsuarioLivro{
		Usuario: &api.Identificador{
			Id: user,
		},
		Livro: &api.Identificador{
			Id: isbn,
		},
	}

	stream, err := client.RealizaEmprestimo(ctx)
	if err != nil {
		return fmt.Errorf("erro ao emprestar livro: %v", err)
	}

	if err := stream.Send(req); err != nil {
		return fmt.Errorf("erro ao enviar o livro para emprestimo: %v", err)
	}
	status, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("erro ao receber o status do emprestimo: %v", err)
	}
	if status.Status == 1 {
		return fmt.Errorf("erro ao emprestar o livro: %v", status)
	}

	fmt.Printf("Livro %s emprestado para usuário %s com sucesso\n", isbn, user)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func returnBook(client api.PortalBibliotecaClient) error {
	var user, isbn string

	fmt.Println("Digite a chave do usuário:")
	fmt.Scan(&user)
	fmt.Println("Digite o ISBN do livro:")
	fmt.Scan(&isbn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.RealizaDevolucao(ctx)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao servidor: %v", err)
	}
	req := &api.UsuarioLivro{
		Usuario: &api.Identificador{Id: user},
		Livro:   &api.Identificador{Id: isbn},
	}
	if err := stream.Send(req); err != nil {
		return fmt.Errorf("erro ao enviar o livro para devolução: %v", err)
	}
	status, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("erro ao receber o status da devolução: %v", err)
	}
	if status.Status == 1 {
		return fmt.Errorf("erro ao devolver o livro: %v", status)
	}

	fmt.Printf("Livro %s devolvido pelo usuário %s com sucesso\n", isbn, user)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func listBorrowedBooks(client api.PortalBibliotecaClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.ListaLivrosEmprestados(ctx, &api.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao obter todos os livros emprestados: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "ISBN\tTítulo\tAutor\n")
	for {
		book, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Fprintf(w, "%v\t%v\t%v\n", book.Isbn, book.Titulo, book.Autor)
	}

	w.Flush()
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func listMissingBooks(client api.PortalBibliotecaClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.ListaLivrosEmFalta(ctx, &api.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao obter todos os livros em falta: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "ISBN\tTítulo\tAutor\n")
	for {
		book, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Fprintf(w, "%v\t%v\t%v\n", book.Isbn, book.Titulo, book.Autor)
	}

	w.Flush()
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func searchBook(client api.PortalBibliotecaClient) error {
	var crit string
	fmt.Println("Digite o critério de busca do livro:")
	fmt.Scan(&crit)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.PesquisaLivro(ctx, &api.Criterio{Criterio: crit})
	if err != nil {
		return fmt.Errorf("erro ao pesquisar livro: %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {

			return fmt.Errorf("erro ao buscar livro: %v", err)
		}

		fmt.Printf("Livro encontrado: %v\n", res)
	}
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}
