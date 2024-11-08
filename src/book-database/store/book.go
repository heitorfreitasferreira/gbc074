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

	"library-manager/book-database/database"
	api_cad "library-manager/shared/api/cad"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

// BookStore is a simple key-value store, where all changes are made via Raft consensus.
type BookStore struct {
	RaftDir  string
	RaftBind string
	inmem    bool

	mu   sync.Mutex
	repo database.BookRepo

	raft   *raft.Raft
	logger *log.Logger
}

// New returns a new Store (BookDatabase).
func New(inmem bool, raftPath string) *BookStore {
	db, err := database.New(raftPath)

	if err != nil {
		log.Fatalf("failed to create database: %s", err)
		return nil
	}

	return &BookStore{
		repo:   db,
		inmem:  inmem,
		logger: log.New(os.Stderr, "[store] ", log.LstdFlags),
	}
}

// Open opens the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (s *BookStore) Open(enableSingle bool, localID string) error {
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

// Get returns the value for the given key.
func (s *BookStore) Get(isbn database.ISBN) (database.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.GetBook(isbn)
}

// GetAll returns all values.
func (s *BookStore) GetAll() ([]database.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.repo.GetAllBooks()
}

func (s *BookStore) Create(value database.Book) error {
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
func (s *BookStore) Edit(key database.ISBN, value database.Book) error {
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

// Delete deletes the given key.
func (s *BookStore) Delete(key database.ISBN) error {
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

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
func (s *BookStore) Join(nodeID, addr string) error {
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

type fsm BookStore

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var operation Operation
	if err := json.Unmarshal(l.Data, &operation); err != nil {
		panic(fmt.Sprintf("failed to unmarshal Operation: %s", err.Error()))
	}

	var status api_cad.Status
	switch operation.OpType {
	case OperationCreate:
		status, _ = f.applyCreate(operation.Value)
		break
	case OperationEdit:
		status, _ = f.applyEdit(operation.Value)
		break
	case OperationDelete:
		status, _ = f.applyDelete(operation.Key)
		break
	default:
		panic(fmt.Sprintf("unrecognized Operation op: %v", operation.OpType))
	}

	return &status
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Clone the map.
	o := make(map[string]database.Book)
	allBooks, err := f.repo.GetAllBooks()
	if err != nil {
		return nil, err
	}

	for _, book := range allBooks {
		o[string(book.Isbn)] = book
	}
	return &fsmSnapshot{store: o}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	o := make(map[string]database.Book)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	for _, value := range o {
		f.repo.CreateBook(value)
	}

	return nil
}

func (f *fsm) applyCreate(value database.Book) (api_cad.Status, error) {
	return f.repo.CreateBook(value)
}

func (f *fsm) applyEdit(value database.Book) (api_cad.Status, error) {
	return f.repo.EditBook(value)
}

func (f *fsm) applyDelete(isbn database.ISBN) (api_cad.Status, error) {
	return f.repo.RemoveBook(isbn)
}

type fsmSnapshot struct {
	store map[string]database.Book
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
