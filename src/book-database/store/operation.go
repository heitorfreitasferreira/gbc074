package store

import (
	"time"

	"library-manager/book-database/repository"
)

type OperationType int

const (
	OperationCreate OperationType = iota
	OperationEdit
	OperationDelete
)

type Table string

const (
	TableUser Table = "user"
	TableBook Table = "book"
)

type Operation struct {
	OpType    OperationType
	Key       repository.ISBN
	Value     repository.Book
	TimeStamp time.Time
}
