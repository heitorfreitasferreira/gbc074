package database

import (
	"errors"
	"sync"
	"time"

	"library-manager/shared/api/bib"
)

var ConcreteUserBookRepo UserBookRepo = NewInMemoryUserBookRepo()

var (
	ErrUserDoesNotHaveLoans = errors.New("usuário não tem empréstimos")
	ErrUserDidNotLoanBook   = errors.New("usuário não emprestou este livro")
	ErrBookAlreadyLoaned    = errors.New("livro já emprestado")
)

type UserBook struct {
	UserCPF  string
	BookISBN string
}

func NewUserBook(userCPF, bookISBN string) UserBook {
	return UserBook{UserCPF: userCPF, BookISBN: bookISBN}
}

func UserBookToProto(userBook UserBook) api_bib.UsuarioLivro {
	return api_bib.UsuarioLivro{
		Usuario: &api_bib.Identificador{
			Id: userBook.UserCPF,
		},
		Livro: &api_bib.Identificador{
			Id: userBook.BookISBN,
		},
	}
}

func ProtoToUserBook(protoUserBook *api_bib.UsuarioLivro) UserBook {
	return UserBook{
		UserCPF:  protoUserBook.Usuario.Id,
		BookISBN: protoUserBook.Livro.Id,
	}
}

type UserBookRepo interface {
	Create(UserBook)
	Delete(UserBook) error
	GetAll() []UserBook
	GetUserLoans(string) []LoanBookAndTime
	RemoveUserLoan(UserBook) error
}

type LoanBookAndTime struct {
	BookISBN  string
	Timestamp int64
}

type InMemoryUserBookRepo struct {
	data map[string][]LoanBookAndTime
	mu   sync.RWMutex
}

func NewInMemoryUserBookRepo() UserBookRepo {
	return &InMemoryUserBookRepo{data: make(map[string][]LoanBookAndTime)}
}

func (repo *InMemoryUserBookRepo) Create(userBook UserBook) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userLoans := repo.data[userBook.UserCPF]
	if userLoans == nil || len(userLoans) == 0 {
		newLoans := make([]LoanBookAndTime, 0)
		newLoans = append(newLoans, LoanBookAndTime{BookISBN: userBook.BookISBN, Timestamp: time.Now().UnixMilli()})
		repo.data[userBook.UserCPF] = newLoans

		return
	}

	repo.data[userBook.UserCPF] = append(userLoans, LoanBookAndTime{BookISBN: userBook.BookISBN, Timestamp: time.Now().UnixMilli()})
}

func (repo *InMemoryUserBookRepo) Delete(userBook UserBook) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userLoans := repo.data[userBook.UserCPF]
	if userLoans == nil || len(userLoans) == 0 {
		return ErrUserDoesNotHaveLoans
	}

	// Remove livro do empréstimo
	for i, loan := range userLoans {
		if loan.BookISBN == userBook.BookISBN {
			repo.data[userBook.UserCPF] = append(userLoans[:i], userLoans[i+1:]...)
			return nil
		}
	}

	return ErrUserDidNotLoanBook
}

func (repo *InMemoryUserBookRepo) GetAll() []UserBook {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	userBookList := make([]UserBook, 0)
	for userId, loans := range repo.data {
		for _, loan := range loans {
			userkBook := NewUserBook(userId, loan.BookISBN)
			userBookList = append(userBookList, userkBook)
		}
	}

	return userBookList
}

func (repo *InMemoryUserBookRepo) GetUserLoans(userId string) []LoanBookAndTime {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	return repo.data[userId]
}

func (repo *InMemoryUserBookRepo) RemoveUserLoan(userBook UserBook) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	userLoans, ok := repo.data[userBook.UserCPF]
	if !ok {
		return ErrUserNotFound
	}

	for i, loan := range userLoans {
		if loan.BookISBN == userBook.BookISBN {
			repo.data[userBook.UserCPF] = append(userLoans[:i], userLoans[i+1:]...)
			return nil
		}
	}

	return ErrUserDidNotLoanBook

}
