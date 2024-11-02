package store

import (
	"library-manager/shared/utils"
	"library-manager/user-database/database"
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
	TableBook Table = "user-book"
)

type Operation struct {
	Table     Table
	OpType    OperationType
	Key       utils.CPF
	Value     database.User
	TimeStamp time.Time
}
