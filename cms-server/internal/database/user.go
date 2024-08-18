package database

import (
	"errors"
	"github.com/rpc-mqtt-library-manager/cms-server/api"
	"sync"
)

var ConcreteUserRepo UserRepo = NewInMemoryUserRepo()

type CPF string

func (cpf CPF) Validate() bool {
	return true
}

type User struct {
	Cpf  CPF
	Nome string
}

type UserRepo interface {
	CreateUser(User) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error)
	EditaUsuario(User) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error)
	RemoveUsuario(CPF) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error)
	ObtemUsuario(CPF) (User, error)
	ObtemTodosUsuarios() ([]User, error)
}

func UserToProto(user User) br_ufu_facom_gbc074_projeto_biblioteca.Usuario {
	return br_ufu_facom_gbc074_projeto_biblioteca.Usuario{
		Cpf:  string(user.Cpf),
		Nome: user.Nome,
	}
}

func ProtoToUser(protoUser *br_ufu_facom_gbc074_projeto_biblioteca.Usuario) User {
	return User{
		Cpf:  CPF(protoUser.Cpf),
		Nome: protoUser.Nome,
	}
}

type InMemoryUserRepo struct {
	users map[CPF]User
	mu    sync.RWMutex
}

func NewInMemoryUserRepo() UserRepo {
	return &InMemoryUserRepo{users: make(map[CPF]User)}
}

func (repo *InMemoryUserRepo) CreateUser(user User) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.users[user.Cpf]; ok {
		return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 1, Msg: "usuário já existe"}, nil
	}
	repo.users[user.Cpf] = user
	return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 0, Msg: "usuário criado com sucesso"}, nil
}

func (repo *InMemoryUserRepo) EditaUsuario(user User) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.users[user.Cpf]; !ok {
		return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 1, Msg: "usuário não existe"}, nil
	}
	repo.users[user.Cpf] = user
	return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 0, Msg: "usuário atualizado com sucesso"}, nil
}

func (repo *InMemoryUserRepo) RemoveUsuario(id CPF) (br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, ok := repo.users[id]; !ok {
		return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 1, Msg: "usuário não existe"}, nil
	}
	delete(repo.users, id)
	return br_ufu_facom_gbc074_projeto_biblioteca.Status{Status: 0, Msg: "usuário removido com sucesso"}, nil
}

func (repo *InMemoryUserRepo) ObtemUsuario(id CPF) (User, error) {
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
