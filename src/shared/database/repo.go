package database

type Repository interface {
	Close() error
}
