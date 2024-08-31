package database

import (
	"errors"
	"sync"

	"library-manager/shared/api/bib"
)

var ConcreteUserRepo UserRepo = NewInMemoryUserRepo()

var (
	ErrUserAlreadyExists = errors.New("usuário já existe")
	ErrUserNotFound      = errors.New("usuário não encontrado")
)

type User struct {
	CPF     string
	Name    string
	Blocked bool
}

type UserRepo interface {
	Create(User) error
	Update(User) error
	Remove(string) error
	GetById(string) (User, error)
	GetAll() []User
}

func UserToProto(user User) api_bib.Usuario {
	return api_bib.Usuario{
		Cpf:  user.CPF,
		Nome: user.Name,
	}
}

func ProtoToUser(protoUser *api_bib.Usuario) User {
	return User{
		CPF:     (protoUser.Cpf),
		Name:    protoUser.Nome,
		Blocked: protoUser.Bloqueado,
	}
}

type InMemoryUserRepo struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewInMemoryUserRepo() UserRepo {
	return &InMemoryUserRepo{users: make(map[string]User)}
}

// Função para criar um usuário
func (repo *InMemoryUserRepo) Create(user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.users[user.CPF]; ok {
		return ErrUserAlreadyExists
	}

	repo.users[user.CPF] = user
	return nil
}

// Funcão para atualizar um usuário
func (repo *InMemoryUserRepo) Update(user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.users[user.CPF]; !ok {
		return ErrUserNotFound
	}

	repo.users[user.CPF] = user
	return nil
}

// Função para remover um usuário
func (repo *InMemoryUserRepo) Remove(userCPF string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.users[userCPF]; !ok {
		return ErrUserNotFound
	}

	delete(repo.users, userCPF)
	return nil
}

// Função para obter um usuário pelo CPF
func (repo *InMemoryUserRepo) GetById(userCPF string) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if user, ok := repo.users[userCPF]; ok {
		return user, nil
	}

	return User{}, ErrUserNotFound
}

// Função para obter todos os usuários
func (repo *InMemoryUserRepo) GetAll() []User {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var users []User
	for _, user := range repo.users {
		users = append(users, user)
	}

	return users
}
