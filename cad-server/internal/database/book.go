package database

import (
	"errors"
	"sync"

	"library-manager/cad-server/api"
	"library-manager/cad-server/internal/utils"
)

var ConcreteBookRepo BookRepo = NewInMemoryBookRepo()

type Book struct {
	Isbn   utils.ISBN
	Titulo string
	Autor  string
	Total  int32
}

type BookRepo interface {
	CreateBook(Book) (api.Status, error)
	EditaLivro(Book) (api.Status, error)
	RemoveLivro(utils.ISBN) (api.Status, error)
	ObtemLivro(utils.ISBN) (Book, error)
	ObtemTodosLivros() ([]Book, error)
}

func BookToProto(book Book) api.Livro {
	return api.Livro{
		Isbn:   string(book.Isbn),
		Titulo: book.Titulo,
		Autor:  book.Autor,
		Total:  book.Total,
	}
}

func ProtoToBook(protoBook *api.Livro) Book {
	return Book{
		Isbn:   utils.ISBN(protoBook.Isbn),
		Titulo: protoBook.Titulo,
		Autor:  protoBook.Autor,
		Total:  protoBook.Total,
	}
}

type InMemoryBookRepo struct {
	books map[utils.ISBN]Book
	mu    sync.RWMutex
}

func NewInMemoryBookRepo() BookRepo {
	return &InMemoryBookRepo{books: make(map[utils.ISBN]Book)}
}

func (repo *InMemoryBookRepo) CreateBook(book Book) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.books[book.Isbn]; ok {
		return api.Status{Status: 1, Msg: "livro já existe"}, nil
	}
	repo.books[book.Isbn] = book
	return api.Status{Status: 0, Msg: "livro criado com sucesso"}, nil
}

func (repo *InMemoryBookRepo) EditaLivro(book Book) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.books[book.Isbn]; !ok {
		return api.Status{Status: 1, Msg: "livro não existe"}, nil
	}
	repo.books[book.Isbn] = book
	return api.Status{Status: 0, Msg: "livro atualizado com sucesso"}, nil
}

func (repo *InMemoryBookRepo) RemoveLivro(id utils.ISBN) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.books[id]; !ok {
		return api.Status{Status: 1, Msg: "livro não existe"}, nil
	}
	delete(repo.books, id)
	return api.Status{Status: 0, Msg: "livro removido com sucesso"}, nil
}

func (repo *InMemoryBookRepo) ObtemLivro(id utils.ISBN) (Book, error) {
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
