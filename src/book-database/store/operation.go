package store

import (
	"time"

	"library-manager/book-database/database"
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
	Key       database.ISBN
	Value     database.Book
	TimeStamp time.Time
}
