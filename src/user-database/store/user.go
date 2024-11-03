package store

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	api_cad "library-manager/shared/api/cad"
	"library-manager/shared/utils"
	"library-manager/user-database/database"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

// UserStore is a simple key-value store, where all changes are made via Raft consensus.
type UserStore struct {
	RaftDir  string
	RaftBind string
	inmem    bool

	mu           sync.Mutex
	userRepo     database.UserRepo
	userBookRepo database.UserBookRepo

	raft   *raft.Raft
	logger *log.Logger
}

// New returns a new Store (BookDatabase).
func New(inmem bool, raftPath string) *UserStore {
	db, err := database.New(raftPath)

	if err != nil {
		log.Fatalf("failed to create database: %s", err)
		return nil
	}

	return &UserStore{
		userRepo: db,
		inmem:    inmem,
		logger:   log.New(os.Stderr, "[store] ", log.LstdFlags),
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (s *UserStore) Open(enableSingle bool, localID string) error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	if s.inmem {
		logStore = raft.NewInmemStore()
		stableStore = raft.NewInmemStore()
	} else {
		boltDB, err := raftboltdb.New(raftboltdb.Options{
			Path: filepath.Join(s.RaftDir, "raft.db"),
		})
		if err != nil {
			return fmt.Errorf("new bbolt store: %s", err)
		}
		logStore = boltDB
		stableStore = boltDB
	}

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	s.raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}

	return nil
}

// GetUser returns the value for the given key.
func (s *UserStore) GetUser(cpf utils.CPF) (database.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.userRepo.ObtemUsuario(cpf)
}

// GetAllUsers returns all values.
func (s *UserStore) GetAllUsers() ([]database.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.userRepo.ObtemTodosUsuarios()
}

func (s *UserStore) CreateUser(value database.User) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	operation := &Operation{
		OpType:    OperationCreate,
		Value:     value,
		TimeStamp: time.Now(),
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return err
	}
	f := s.raft.Apply(operationJson, raftTimeout)
	return f.Error()
}

// Set sets the value for the given key.
func (s *UserStore) EditUser(key utils.CPF, value database.User) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	operation := &Operation{
		OpType:    OperationEdit,
		Key:       key,
		Value:     value,
		TimeStamp: time.Now(),
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return err
	}

	f := s.raft.Apply(operationJson, raftTimeout)
	return f.Error()
}

// DeleteUser deletes the given key.
func (s *UserStore) DeleteUser(key utils.CPF) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	operation := &Operation{
		OpType:    OperationDelete,
		Key:       key,
		TimeStamp: time.Now(),
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return err
	}

	f := s.raft.Apply(operationJson, raftTimeout)
	return f.Error()
}

func (s *UserStore) CreateUserBook(userBook database.UserBook) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	operation := &Operation{
		Table:     TableUserBook,
		OpType:    OperationCreate,
		Value:     userBook,
		TimeStamp: time.Now(),
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return err
	}
	f := s.raft.Apply(operationJson, raftTimeout)
	return f.Error()
}

func (s *UserStore) DeleteUserBook(userbook database.UserBook) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	operation := &Operation{
		Table:     TableUserBook,
		OpType:    OperationDelete,
		Value:     userbook,
		TimeStamp: time.Now(),
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return err
	}
	f := s.raft.Apply(operationJson, raftTimeout)
	return f.Error()
}

func (s *UserStore) GetAllUserBooks() []database.UserBook {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.userBookRepo.GetAll()
}

func (s *UserStore) GetUserLoans(userId string) []database.LoanBookAndTime {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.userBookRepo.GetUserLoans(userId)
}

func (s *UserStore) RemoveUserLoan(userbook database.UserBook) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}
	operation := &Operation{
		Table:     TableUserBook,
		OpType:    OperationDelete,
		Value:     userbook,
		TimeStamp: time.Now(),
	}
	operationJson, err := json.Marshal(operation)
	if err != nil {
		return err
	}
	f := s.raft.Apply(operationJson, raftTimeout)
	return f.Error()
}

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
func (s *UserStore) Join(nodeID, addr string) error {
	s.logger.Printf("received join request for remote node %s at %s", nodeID, addr)

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				s.logger.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}

			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	s.logger.Printf("node %s at %s joined successfully", nodeID, addr)
	return nil
}

type fsm UserStore

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var operation Operation
	if err := json.Unmarshal(l.Data, &operation); err != nil {
		panic(fmt.Sprintf("failed to unmarshal Operation: %s", err.Error()))
	}

	var status api_cad.Status
	if operation.Table == TableUser {
		switch operation.OpType {
		case OperationCreate:
			user, ok := operation.Value.(database.User)
			if !ok {
				panic(fmt.Sprintf("unexpected value type: %T", operation.Value))
			}

			status, _ = f.applyCreateUser(user)
			break
		case OperationEdit:
			user, ok := operation.Value.(database.User)
			if !ok {
				panic(fmt.Sprintf("unexpected value type: %T", operation.Value))
			}

			status, _ = f.applyEditUser(user)
			break
		case OperationDelete:
			status, _ = f.applyDeleteUser(operation.Key)
			break
		default:
			panic(fmt.Sprintf("unrecognized Operation op: %v", operation.OpType))
		}

	} else if operation.Table == TableUserBook {
		switch operation.OpType {
		case OperationCreate:
			userBook, ok := operation.Value.(database.UserBook)
			if !ok {
				panic(fmt.Sprintf("unexpected value type: %T", operation.Value))
			}
			f.applyCreateUserBook(userBook)
			break
		case OperationDelete:
			userBook, ok := operation.Value.(database.UserBook)
			if !ok {
				panic(fmt.Sprintf("unexpected value type: %T", operation.Value))
			}
			f.applyDeleteUserBook(userBook)
			break
		}
	}

	return &status
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Clone the map.
	udb := make(map[string]database.User)
	allUsers, err := f.userRepo.ObtemTodosUsuarios()
	if err != nil {
		return nil, err
	}
	for _, user := range allUsers {
		udb[string(user.Cpf)] = user
	}

	ubdb := make(map[string]database.UserBook)
	allUserBooks := f.userBookRepo.GetAll()
	for _, userBook := range allUserBooks {
		ubdb[userBook.UserCPF] = userBook
	}

	return &fsmSnapshot{store: StoreDatabases{
		UserRepo:     udb,
		UserBookRepo: ubdb,
	}}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	o := make(map[string]database.User)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	for _, value := range o {
		f.userRepo.CreateUser(value)
	}

	return nil
}

func (f *fsm) applyCreateUser(value database.User) (api_cad.Status, error) {
	return f.userRepo.CreateUser(value)
}

func (f *fsm) applyEditUser(value database.User) (api_cad.Status, error) {
	return f.userRepo.EditaUsuario(value)
}

func (f *fsm) applyDeleteUser(isbn utils.CPF) (api_cad.Status, error) {
	return f.userRepo.RemoveUsuario(isbn)
}

func (f *fsm) applyCreateUserBook(value database.UserBook) {
	f.userBookRepo.Create(value)
}

func (f *fsm) applyDeleteUserBook(value database.UserBook) error {
	return f.userBookRepo.Delete(value)
}

func (f *fsm) applyRemoveUserLoan(value database.UserBook) error {
	return f.userBookRepo.RemoveUserLoan(value)
}

type StoreDatabases struct {
	UserRepo     map[string]database.User
	UserBookRepo map[string]database.UserBook
}

type fsmSnapshot struct {
	store StoreDatabases
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
