package handlers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
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

	res, err := client.RealizaEmprestimo(ctx, &api.UsuarioLivro{user: user, livro: isbn})
	if err != nil {
		return fmt.Errorf("erro ao emprestar livro: %v", err)
	}

	fmt.Printf("Livro %v emprestado com sucesso: %v\n", user, isbn, res)
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

	ctx, cancel := context.RealizaDevolucao(context.Background(), time.Second)
	defer cancel()

	res, err := client.RealizaEmprestimo(ctx, &api.UsuarioLivro{user: user, livro: isbn})
	if err != nil {
		return fmt.Errorf("erro ao devolver livro: %v", err)
	}

	fmt.Printf("Livro %v devolvido com sucesso: %v\n", user, isbn, res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func listBorrowedBooks(client api.PortalBibliotecaClient) errr {
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
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", book.Isbn, book.Titulo, book.Autor)
	}

	w.Flush()
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func listMissingBooks(client api.PortalBibliotecaClient) err {
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
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", book.Isbn, book.Titulo, book.Autor)
	}

	w.Flush()
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func searchBook(client api.PortalBibliotecaClient) err {
	var crit string
	fmt.Println("Digite o critério de busca do livro:")
	fmt.Scan(&crit)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.PesquisaLivro(ctx, &api.Criterio{criterio: crit})

	if err != nil {
		return fmt.Errorf("erro ao buscar livro: %v", err)
	}
	fmt.Printf("Livro encontrado: %v\n", res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}