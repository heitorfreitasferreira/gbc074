package structure

import (
	"errors"
	"time"
)

var (
	ErrIdNotNumeric      = errors.New("identificador não numérico")
	ErrBookAlreadyLoaned = errors.New("livro já emprestado")
	ErrValueAlreadyInSet = errors.New("valor já presente no conjunto")
	ErrValueIsNotInSet   = errors.New("valor não presente no conjunto")
)

type Set [][3]int64

func (s Set) Add(userId, bookId int64, time time.Time) error {
	for _, pair := range s {
		if pair[0] == userId && pair[1] == bookId {
			return ErrValueAlreadyInSet
		}
	}

	s = append(s, [3]int64{userId, bookId, time.UnixMilli()})
	return nil
}

func (s Set) Contains(userId, bookId int64) bool {
	for _, pair := range s {
		if pair[0] == userId && pair[1] == bookId {
			return true
		}
	}
	return false
}

func (s Set) Remove(userId, bookId int64) error {
	for i, pair := range s {
		if pair[0] == userId && pair[1] == bookId {
			s = append(s[:i], s[i+1:]...)
			return nil
		}
	}
	return ErrValueIsNotInSet
}
