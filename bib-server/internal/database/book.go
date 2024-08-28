package database

import (
	"errors"
	"library-manager/bib-server/api"
	"sync"
)

var ConcreteBookRepo BookRepo = NewInMemoryBookRepo()

type ISBN string

func (isbn ISBN) Validate() bool {
	return true
}

type Book struct {
	Isbn   ISBN
	Titulo string
	Autor  string
	Total  int32
}

type BookRepo interface {
	CreateBook(Book) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error)
	EditaLivro(Book) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error)
	RemoveLivro(ISBN) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error)
	ObtemLivro(ISBN) (Book, error)
	ObtemTodosLivros() ([]Book, error)
}

func BookToProto(book Book) br_ufu_facom_gbc074_projeto_biblioteca.Livro {
	return br_ufu_facom_gbc074_projeto_biblioteca.Livro{
		Isbn:   string(book.Isbn),
		Titulo: book.Titulo,
		Autor:  book.Autor,
		Total:  book.Total,
	}
}

func ProtoToBook(protoBook *br_ufu_facom_gbc074_projeto_biblioteca.Livro) Book {
	return Book{
		Isbn:   ISBN(protoBook.Isbn),
		Titulo: protoBook.Titulo,
		Autor:  protoBook.Autor,
		Total:  protoBook.Total,
	}
}

type InMemoryBookRepo struct {
	books map[ISBN]Book
	mu    sync.RWMutex
}

func NewInMemoryBookRepo() BookRepo {
	return &InMemoryBookRepo{books: make(map[ISBN]Book)}
}

func (repo *InMemoryBookRepo) CreateBook(book Book) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.books[book.Isbn]; ok {
		return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 1, Msg: "livro já existe"}, nil
	}
	repo.books[book.Isbn] = book
	return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 0, Msg: "livro criado com sucesso"}, nil
}

func (repo *InMemoryBookRepo) EditaLivro(book Book) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.books[book.Isbn]; !ok {
		return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 1, Msg: "livro não existe"}, nil
	}
	repo.books[book.Isbn] = book
	return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 0, Msg: "livro atualizado com sucesso"}, nil
}

func (repo *InMemoryBookRepo) RemoveLivro(id ISBN) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.books[id]; !ok {
		return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 1, Msg: "livro não existe"}, nil
	}
	delete(repo.books, id)
	return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 0, Msg: "livro removido com sucesso"}, nil
}

func (repo *InMemoryBookRepo) ObtemLivro(id ISBN) (Book, error) {
	if !id.Validate() {
		return Book{}, errors.New("isbn inválido")
	}
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	book, ok := repo.books[id]
	if !ok {
		return Book{}, errors.New("livro não encontrado")
	}
	return book, nil
}

func (repo *InMemoryBookRepo) ObtemTodosLivros() ([]Book, error) {
	var books []Book
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	for _, book := range repo.books {
		books = append(books, book)
	}
	return books, nil
}
