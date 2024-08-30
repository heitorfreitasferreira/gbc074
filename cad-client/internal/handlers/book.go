package handlers

import (
	"bufio"
	"context"
	"fmt"
	api_cad "library-manager/shared/api/cad"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func removeBook(client api_cad.PortalCadastroClient) error {
	var isbn string
	fmt.Println("Digite o ISBN do livro:")
	fmt.Scan(&isbn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RemoveLivro(ctx, &api_cad.Identificador{Id: isbn})
	if err != nil {
		return fmt.Errorf("erro ao remover livro: %v", err)
	}
	fmt.Printf("Livro %v removido com sucesso: %v\n", isbn, res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func addBook(client api_cad.PortalCadastroClient) error {
	book, err := readBookData()
	if err != nil {
		return fmt.Errorf("erro ao ler dados do livro: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.NovoLivro(ctx, book)
	if err != nil {
		return fmt.Errorf("erro ao adicionar livro: %v", err)
	}
	fmt.Printf("Livro %s adicionado com sucesso: %v\n", book.Titulo, res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func editBook(client api_cad.PortalCadastroClient) error {
	book, err := readBookData()
	if err != nil {
		return fmt.Errorf("erro ao ler dados do livro: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.EditaLivro(ctx, book)
	if err != nil {
		return fmt.Errorf("erro ao editar livro: %v", err)
	}
	fmt.Printf("Livro %v editado com sucesso %v\n", book, res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func getBook(client api_cad.PortalCadastroClient) error {
	var isbn string
	fmt.Println("Digite o ISBN do livro:")
	fmt.Scan(&isbn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.ObtemLivro(ctx, &api_cad.Identificador{Id: isbn})
	if err != nil {
		return fmt.Errorf("erro ao obter livro: %v", err)
	}
	fmt.Printf("Livro obtido: %v\n", res)
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func getAllBooks(client api_cad.PortalCadastroClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := client.ObtemTodosLivros(ctx, &api_cad.Vazia{})
	if err != nil {
		return fmt.Errorf("erro ao obter todos os livros: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "ISBN\tTítulo\tAutor\tTotal\n")
	for {
		book, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", book.Isbn, book.Titulo, book.Autor, book.Total)
	}
	w.Flush()
	fmt.Println("Pressione ENTER...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func readBookData() (*api_cad.Livro, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Digite o ISBN do livro:")
	isbn, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("erro ao ler ISBN: %v", err)
	}
	isbn = strings.TrimSpace(isbn)

	fmt.Println("Digite o título do livro:")
	titulo, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("erro ao ler título: %v", err)
	}
	titulo = strings.TrimSpace(titulo)

	fmt.Println("Digite o autor do livro:")
	autor, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("erro ao ler autor: %v", err)
	}
	autor = strings.TrimSpace(autor)

	fmt.Println("Digite o total de exemplares do livro:")
	var total int32
	_, err = fmt.Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler total de exemplares: %v", err)
	}

	return &api_cad.Livro{
		Isbn:   isbn,
		Titulo: titulo,
		Autor:  autor,
		Total:  total,
	}, nil
}
