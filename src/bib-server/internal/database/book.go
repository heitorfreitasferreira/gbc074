package database

import (
	"errors"
	"sync"

	"library-manager/shared/api/bib"
)

var ConcreteBookRepo BookRepo = NewInMemoryBookRepo()

var (
	ErrBookAlreadyExists = errors.New("livro já existe")
	ErrBookNotFound      = errors.New("livro não encontrado")
)

type Book struct {
	ISBN      string
	Title     string
	Author    string
	Total     int32
	Remaining int32
}

func BookToProto(book Book) api_bib.Livro {
	return api_bib.Livro{
		Isbn:     string(book.ISBN),
		Titulo:   book.Title,
		Autor:    book.Author,
		Total:    book.Total,
		Restante: book.Remaining,
	}
}

func ProtoToBook(protoBook *api_bib.Livro) Book {
	return Book{
		ISBN:      protoBook.Isbn,
		Title:     protoBook.Titulo,
		Author:    protoBook.Autor,
		Total:     protoBook.Total,
		Remaining: protoBook.Restante,
	}
}

type BookRepo interface {
	Create(Book) error
	Update(Book) error
	Remove(string) error
	GetById(string) (Book, error)
	GetAll() []Book
}

type InMemoryBookRepo struct {
	books map[string]Book
	mu    sync.RWMutex
}

func NewInMemoryBookRepo() BookRepo {
	return &InMemoryBookRepo{books: make(map[string]Book)}
}

func (repo *InMemoryBookRepo) Create(book Book) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, err := repo.books[book.ISBN]; err {
		return ErrBookAlreadyExists
	}

	repo.books[book.ISBN] = book

	return nil
}

func (repo *InMemoryBookRepo) Update(book Book) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, err := repo.books[book.ISBN]; !err {
		return ErrBookNotFound
	}

	repo.books[book.ISBN] = book

	return nil
}

func (repo *InMemoryBookRepo) Remove(bookISBN string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, err := repo.books[bookISBN]; !err {
		return ErrBookNotFound
	}

	delete(repo.books, bookISBN)

	return nil
}

func (repo *InMemoryBookRepo) GetById(bookISBN string) (Book, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	book, err := repo.books[bookISBN]
	if err {
		return Book{}, ErrBookNotFound
	}
	return book, nil
}

func (repo *InMemoryBookRepo) GetAll() []Book {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var books []Book
	for _, book := range repo.books {
		books = append(books, book)
	}

	return books
}
