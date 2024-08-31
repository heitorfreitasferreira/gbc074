package database

import (
	"errors"
	"sync"
	"time"

	"library-manager/bib-server/internal/database/structure"
	"library-manager/bib-server/internal/database/utils"
	"library-manager/shared/api/bib"
)

var (
	ErrBookAlreadyLoaned = errors.New("livro já emprestado")
)

type UserBook struct {
	UserId   string
	BookISNB string
}

type UserBookRepo interface {
	LoanBook(UserBook) (api_bib.Status, error)
	DeleteLoan(UserBook) (api_bib.Status, error)
}

type InMemoryUserBookRepo struct {
	data structure.Set
	mu   sync.RWMutex
}

func NewInMemoryUserBookRepo() UserBookRepo {
	return &InMemoryUserBookRepo{data: make(structure.Set, 0)}
}

var ConcreteUserBookRepo UserBookRepo = NewInMemoryUserBookRepo()

func (repo *InMemoryUserBookRepo) LoanBook(userBook UserBook) (api_bib.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userId, bookId, err := utils.ConvertUserBookIds(userBook.UserId, userBook.BookISNB)
	if err != nil {
		return api_bib.Status{Status: 1, Msg: err.Error()}, err
	}

	if err := repo.data.Add(int64(userId), int64(bookId), time.Now()); err != nil {
		return api_bib.Status{Status: 1, Msg: "livro já emprestado"}, nil
	}

	return api_bib.Status{Status: 0, Msg: "livro emprestado com sucesso"}, nil
}

func (repo *InMemoryUserBookRepo) DeleteLoan(userBook UserBook) (api_bib.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userId, bookId, err := utils.ConvertUserBookIds(userBook.UserId, userBook.BookISNB)
	if err != nil {
		return api_bib.Status{Status: 1, Msg: err.Error()}, err
	}

	if err := repo.data.Remove(userId, bookId); err != nil {
		return api_bib.Status{Status: 1, Msg: "livro não emprestado"}, nil
	}
	return api_bib.Status{Status: 0, Msg: "livro devolvido com sucesso"}, nil
}
