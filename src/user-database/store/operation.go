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
	TableUser     Table = "user"
	TableUserBook Table = "user-book"
)

// Define an interface that both database.User and database.UserBook implement
type OperationValue interface{}

// Ensure that database.User and database.UserBook implement the OperationValue interface
var _ OperationValue = (*database.User)(nil)
var _ OperationValue = (*database.UserBook)(nil)

type Operation struct {
	Table     Table
	OpType    OperationType
	Key       utils.CPF
	Value     OperationValue
	TimeStamp time.Time
}
