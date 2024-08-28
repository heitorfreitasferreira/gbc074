package database

import (
	"sync"
	"time"

	api "library-manager/bib-server/api"
)

type UserBook struct {
	UserId   string
	BookISNB string
}

type UserBookKey string

type LoanDateTime string

type UserBookRepo interface {
	LoanBook(UserBook) (api.Status, error)
	DeleteLoan(UserBook) (api.Status, error)
}

type InMemoryUserBookRepo struct {
	data map[UserBookKey]LoanDateTime
	mu   sync.RWMutex
}

func NewInMemoryUserBookRepo() UserBookRepo {
	return &InMemoryUserBookRepo{data: make(map[UserBookKey]LoanDateTime)}
}

var ConcreteUserBookRepo UserBookRepo = NewInMemoryUserBookRepo()

func (repo *InMemoryUserBookRepo) LoanBook(userBook UserBook) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	key := UserBookKey(userBook.UserId + userBook.BookISNB)
	if _, ok := repo.data[key]; ok {
		return api.Status{Status: 1, Msg: "livro já emprestado"}, nil
	}
	repo.data[key] = LoanDateTime(time.Now().UTC().Format(time.RFC3339))
	return api.Status{Status: 0, Msg: "livro emprestado com sucesso"}, nil
}

func (repo *InMemoryUserBookRepo) DeleteLoan(userBook UserBook) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	key := UserBookKey(userBook.UserId + userBook.BookISNB)
	if _, ok := repo.data[key]; !ok {
		return api.Status{Status: 1, Msg: "livro não emprestado"}, nil
	}
	delete(repo.data, key)
	return api.Status{Status: 0, Msg: "livro devolvido com sucesso"}, nil
}
