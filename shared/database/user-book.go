package database

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"library-manager/shared/api/bib"
)

var (
	ErrIdNotNumeric      = errors.New("identificador não numérico")
	ErrBookAlreadyLoaned = errors.New("livro já emprestado")
	ErrValueAlreadyInSet = errors.New("valor já presente no conjunto")
	ErrValueIsNotInSet   = errors.New("valor não presente no conjunto")
)

type UserBook struct {
	UserId   string
	BookISNB string
}

type UserBookKey string

type LoanDateTime string

type UserBookRepo interface {
	LoanBook(UserBook) (api_bib.Status, error)
	DeleteLoan(UserBook) (api_bib.Status, error)
}

type UserBookSet [][3]int64

func (s UserBookSet) Add(userId, bookId int64, time time.Time) error {
	for _, pair := range s {
		if pair[0] == userId && pair[1] == bookId {
			return ErrValueAlreadyInSet
		}
	}

	s = append(s, [3]int64{userId, bookId, time.UnixMilli()})
	return nil
}

func (s UserBookSet) Contains(userId, bookId int64) bool {
	for _, pair := range s {
		if pair[0] == userId && pair[1] == bookId {
			return true
		}
	}
	return false
}

func (s UserBookSet) Remove(userId, bookId int64) error {
	for i, pair := range s {
		if pair[0] == userId && pair[1] == bookId {
			s = append(s[:i], s[i+1:]...)
			return nil
		}
	}
	return ErrValueIsNotInSet
}

type InMemoryUserBookRepo struct {
	data UserBookSet
	mu   sync.RWMutex
}

func NewInMemoryUserBookRepo() UserBookRepo {
	return &InMemoryUserBookRepo{data: make(UserBookSet, 0)}
}

var ConcreteUserBookRepo UserBookRepo = NewInMemoryUserBookRepo()

func convertUserBookIds(userIdStr, bookIdStr string) (int64, int64, error) {
	userId, err := strconv.Atoi(strings.TrimSpace(userIdStr))
	if err != nil {
		return 0, 0, ErrIdNotNumeric
	}
	bookId, err := strconv.Atoi(strings.TrimSpace(bookIdStr))
	if err != nil {
		return 0, 0, ErrIdNotNumeric
	}

	return int64(userId), int64(bookId), nil
}

func (repo *InMemoryUserBookRepo) LoanBook(userBook UserBook) (api_bib.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userId, bookId, err := convertUserBookIds(userBook.UserId, userBook.BookISNB)
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

	userId, bookId, err := convertUserBookIds(userBook.UserId, userBook.BookISNB)
	if err != nil {
		return api_bib.Status{Status: 1, Msg: err.Error()}, err
	}

	if err := repo.data.Remove(userId, bookId); err != nil {
		return api_bib.Status{Status: 1, Msg: "livro não emprestado"}, nil
	}
	return api_bib.Status{Status: 0, Msg: "livro devolvido com sucesso"}, nil
}
