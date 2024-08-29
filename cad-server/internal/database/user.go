package database

import (
	"errors"
	"sync"

	"library-manager/cad-server/api"
	"library-manager/cad-server/internal/utils"
)

var ConcreteUserRepo UserRepo = NewInMemoryUserRepo()

type User struct {
	Cpf  utils.CPF
	Nome string
}

type UserRepo interface {
	CreateUser(User) (api.Status, error)
	EditaUsuario(User) (api.Status, error)
	RemoveUsuario(utils.CPF) (api.Status, error)
	ObtemUsuario(utils.CPF) (User, error)
	ObtemTodosUsuarios() ([]User, error)
}

func UserToProto(user User) api.Usuario {
	return api.Usuario{
		Cpf:  string(user.Cpf),
		Nome: user.Nome,
	}
}

func ProtoToUser(protoUser *api.Usuario) User {
	return User{
		Cpf:  utils.CPF(protoUser.Cpf),
		Nome: protoUser.Nome,
	}
}

type InMemoryUserRepo struct {
	users map[utils.CPF]User
	mu    sync.RWMutex
}

func NewInMemoryUserRepo() UserRepo {
	return &InMemoryUserRepo{users: make(map[utils.CPF]User)}
}

func (repo *InMemoryUserRepo) CreateUser(user User) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.users[user.Cpf]; ok {
		return api.Status{Status: 1, Msg: "usuário já existe"}, nil
	}
	repo.users[user.Cpf] = user
	return api.Status{Status: 0, Msg: "usuário criado com sucesso"}, nil
}

func (repo *InMemoryUserRepo) EditaUsuario(user User) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.users[user.Cpf]; !ok {
		return api.Status{Status: 1, Msg: "usuário não existe"}, nil
	}
	repo.users[user.Cpf] = user
	return api.Status{Status: 0, Msg: "usuário atualizado com sucesso"}, nil
}

func (repo *InMemoryUserRepo) RemoveUsuario(id utils.CPF) (api.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.users[id]; !ok {
		return api.Status{Status: 1, Msg: "usuário não existe"}, nil
	}
	delete(repo.users, id)
	return api.Status{Status: 0, Msg: "usuário removido com sucesso"}, nil
}

func (repo *InMemoryUserRepo) ObtemUsuario(id utils.CPF) (User, error) {
	if !id.Validate() {
		return User{}, errors.New("cpf inválido")
	}
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	user, ok := repo.users[id]
	if !ok {
		return User{}, errors.New("usuário não encontrado")
	}
	return user, nil
}

func (repo *InMemoryUserRepo) ObtemTodosUsuarios() ([]User, error) {
	var users []User
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	for _, user := range repo.users {
		users = append(users, user)
	}
	return users, nil
}
