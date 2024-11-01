package repository

import (
	"encoding/json"
	"sync"
	"time"

	api_cad "library-manager/shared/api/cad"
	"library-manager/shared/database"

	"github.com/syndtr/goleveldb/leveldb"
)

type ISBN string

func (isbn ISBN) Validate() bool {
	return true
}

type Book struct {
	Isbn   ISBN
	Titulo string
	Autor  string
	Total  int32
}

type BookRepo interface {
	CreateBook(Book) (api_cad.Status, error)
	EditBook(Book) (api_cad.Status, error)
	RemoveBook(ISBN) (api_cad.Status, error)
	GetBook(ISBN) (Book, error)
	GetAllBooks() ([]Book, error)
	database.Repository
}

func BookToProto(book Book) api_cad.Livro {
	return api_cad.Livro{
		Isbn:   string(book.Isbn),
		Titulo: book.Titulo,
		Autor:  book.Autor,
		Total:  book.Total,
	}
}

func ProtoToBook(protoBook *api_cad.Livro) Book {
	return Book{
		Isbn:   ISBN(protoBook.Isbn),
		Titulo: protoBook.Titulo,
		Autor:  protoBook.Autor,
		Total:  protoBook.Total,
	}
}

type LevelDBBookRepo struct {
	cache map[ISBN]database.Timed[Book]
	mu    sync.RWMutex
	db    *leveldb.DB
}

func New(dbPath string) (BookRepo, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	repo := &LevelDBBookRepo{cache: make(map[ISBN]database.Timed[Book]), db: db}
	// TUTORIAL: Inicia a goroutine que limpa o cache de tempo em tempo
	go repo.cleanupCache()
	return repo, nil
}
func (repo *LevelDBBookRepo) Close() error {
	return repo.db.Close()
}

// CRUD =======================================================================
func (repo *LevelDBBookRepo) CreateBook(book Book) (api_cad.Status, error) {
	exists, err := repo.db.Has([]byte(book.Isbn), nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao verificar se o livro já existe"}, err
	}
	if exists {
		return api_cad.Status{Status: 1, Msg: "livro já existe"}, nil
	}

	data, err := json.Marshal(book)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao serializar livro"}, err
	}
	err = repo.db.Put([]byte(book.Isbn), data, nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao salvar livro"}, err
	}
	repo.putInCache(book)
	return api_cad.Status{Status: 0, Msg: "livro criado com sucesso"}, nil

}

func (repo *LevelDBBookRepo) EditBook(book Book) (api_cad.Status, error) {
	exists, err := repo.db.Has([]byte(book.Isbn), nil)
	if err != nil {
		return api_cad.Status{}, err
	}
	if !exists {
		return api_cad.Status{Status: 1, Msg: "livro não existe"}, nil
	}

	bookData, err := json.Marshal(book)
	if err != nil {
		return api_cad.Status{}, err
	}

	err = repo.db.Put([]byte(book.Isbn), bookData, nil)
	if err != nil {
		return api_cad.Status{}, err
	}

	repo.putInCache(book)

	return api_cad.Status{Status: 0, Msg: "livro atualizado com sucesso"}, nil
}

func (repo *LevelDBBookRepo) RemoveBook(id ISBN) (api_cad.Status, error) {
	exists, err := repo.db.Has([]byte(id), nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao verificar livro"}, err
	}
	if !exists {
		return api_cad.Status{Status: 1, Msg: "livro não existe"}, nil
	}
	err = repo.db.Delete([]byte(id), nil)
	if err != nil {
		return api_cad.Status{Status: 1, Msg: "erro ao remover livro"}, err
	}
	// TODO: Conferir se tem que remover no cache ou é para usar consistência eventual, código abaixo remove do cache
	// repo.mu.Lock()
	// delete(repo.cache, id)
	// repo.mu.Unlock()
	return api_cad.Status{Status: 0, Msg: "livro removido com sucesso"}, nil
}

func (repo *LevelDBBookRepo) GetBook(id ISBN) (Book, error) {
	// TUTORIAL: Leitura tem que ver no cache primeiro
	if book, found := repo.getFromCache(id); found {
		return book, nil
	}
	data, err := repo.db.Get([]byte(id), nil)
	if err != nil {
		return Book{}, err
	}
	var book Book
	err = json.Unmarshal(data, &book)
	if err != nil {
		return Book{}, err
	}
	// TUTORIAL: Coloca no cache
	repo.putInCache(book)
	return book, nil
}

func (repo *LevelDBBookRepo) GetAllBooks() ([]Book, error) {
	var books []Book
	iter := repo.db.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		var book Book
		err := json.Unmarshal(iter.Value(), &book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}
	return books, nil
}

// FIM CRUD ====================================================================

// CACHE HELPERS ============================================================

// De 0.5 em 0.5 segundos, remove os itens do cache que já expiraram
func (repo *LevelDBBookRepo) cleanupCache() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for range ticker.C {
		repo.mu.Lock()
		for isbn, timedBook := range repo.cache {
			// Percorre o cache e remove os itens que já expiraram
			if !timedBook.IsValid() {
				delete(repo.cache, isbn)
			}
		}
		repo.mu.Unlock()
	}
}

func (repo *LevelDBBookRepo) getFromCache(isbn ISBN) (Book, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	book, ok := repo.cache[isbn]
	if !ok {
		return Book{}, false
	}
	if !book.IsValid() {
		delete(repo.cache, isbn)
		return Book{}, false
	}
	return book.Item, true
}

func (repo *LevelDBBookRepo) putInCache(book Book) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.cache[book.Isbn] = database.NewTimed(book)
}
