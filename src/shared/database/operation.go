package database

import (
	"encoding/json"
	"time"
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
	Table     Table
	Key       string
	Value     repository.Book
	TimeStamp time.Time
}
