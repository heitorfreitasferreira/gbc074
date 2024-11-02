package database

import (
	"encoding/json"
	"sync"
	"time"

	api_cad "library-manager/shared/api/cad"
	"library-manager/shared/database"
	"library-manager/shared/utils"

	"github.com/syndtr/goleveldb/leveldb"
)

type User struct {
	Cpf  utils.CPF
	Nome string
}

type UserRepo interface {
	CreateUser(User) (api_cad.Status, error)
	EditaUsuario(User) (api_cad.Status, error)
	RemoveUsuario(utils.CPF) (api_cad.Status, error)
	ObtemUsuario(utils.CPF) (User, error)
	ObtemTodosUsuarios() ([]User, error)
	database.Repository
}

func UserToProto(user User) api_cad.Usuario {
	return api_cad.Usuario{
		Cpf:  string(user.Cpf),
		Nome: user.Nome,
	}
}

func ProtoToUser(protoUser *api_cad.Usuario) User {
	return User{
		Cpf:  utils.CPF(protoUser.Cpf),
		Nome: protoUser.Nome,
	}
}

type InMemoryUserRepo struct {
	db    *leveldb.DB
	cache map[utils.CPF]database.Timed[User]
	mu    sync.RWMutex
}

func New(dbPath string) (UserRepo, error) {
	db, err := leveldb.OpenFile(dbPath, nil)

	if err != nil {
		return nil, err
	}

	repo := &InMemoryUserRepo{cache: make(map[utils.CPF]database.Timed[User]), db: db}
	go repo.cleanupCache()
	return repo, nil
}
func (repo *InMemoryUserRepo) Close() error {
	return repo.db.Close()
}

func (repo *InMemoryUserRepo) CreateUser(user User) (api_cad.Status, error) {
	exists, err := repo.db.Has([]byte(user.Cpf), nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao verificar se o usuário já existe"}, err
	}
	if exists {
		return api_cad.Status{Status: 1, Msg: "usuário já existe"}, nil
	}

	data, err := json.Marshal(user)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao serializar usuário"}, err
	}
	err = repo.db.Put([]byte(user.Cpf), data, nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao salvar usuário"}, err
	}
	repo.putInCache(user)
	return api_cad.Status{Status: 0, Msg: "usuário criado com sucesso"}, nil
}

func (repo *InMemoryUserRepo) EditaUsuario(user User) (api_cad.Status, error) {
	exists, err := repo.db.Has([]byte(user.Cpf), nil)
	if err != nil {
		return api_cad.Status{}, err
	}
	if !exists {
		return api_cad.Status{Status: 1, Msg: "usuário não existe"}, nil
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return api_cad.Status{}, err
	}

	err = repo.db.Put([]byte(user.Cpf), userData, nil)
	if err != nil {
		return api_cad.Status{}, err
	}

	repo.putInCache(user)

	return api_cad.Status{Status: 0, Msg: "usuário atualizado com sucesso"}, nil
}

func (repo *InMemoryUserRepo) RemoveUsuario(id utils.CPF) (api_cad.Status, error) {
	exists, err := repo.db.Has([]byte(id), nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao verificar usuário"}, err
	}
	if !exists {
		return api_cad.Status{Status: 1, Msg: "usuário não existe"}, nil
	}
	err = repo.db.Delete([]byte(id), nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao remover usuário"}, err
	}
	// TODO: Conferir se tem que remover no cache ou é para usar consistência eventual, código abaixo remove do cache
	// repo.mu.Lock()
	// delete(repo.cache, id)
	// repo.mu.Unlock()
	return api_cad.Status{Status: 0, Msg: "usuário removido com sucesso"}, nil
}

func (repo *InMemoryUserRepo) ObtemUsuario(id utils.CPF) (User, error) {
	// TUTORIAL: Leitura tem que ver no cache primeiro
	if user, found := repo.getFromCache(id); found {
		return user, nil
	}
	data, err := repo.db.Get([]byte(id), nil)
	if err != nil {
		return User{}, err
	}
	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}
	// TUTORIAL: Coloca no cache
	repo.putInCache(user)
	return user, nil
}

func (repo *InMemoryUserRepo) ObtemTodosUsuarios() ([]User, error) {
	var users []User
	iter := repo.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		var user User
		err := json.Unmarshal(iter.Value(), &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}
	return users, nil
}

// CACHE

func (repo *InMemoryUserRepo) cleanupCache() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for range ticker.C {
		repo.mu.Lock()
		for isbn, timeduser := range repo.cache {
			// Percorre o cache e remove os itens que já expiraram
			if !timeduser.IsValid() {
				delete(repo.cache, isbn)
			}
		}
		repo.mu.Unlock()
	}
}

func (repo *InMemoryUserRepo) getFromCache(cpf utils.CPF) (User, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	user, ok := repo.cache[cpf]
	if !ok {
		return User{}, false
	}
	if !user.IsValid() {
		delete(repo.cache, cpf)
		return User{}, false
	}
	return user.Item, true
}

func (repo *InMemoryUserRepo) putInCache(user User) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.cache[user.Cpf] = database.NewTimed(user)
}
