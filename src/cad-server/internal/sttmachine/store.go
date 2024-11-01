package sttmachine

// Inspirado na documentação https://github.com/otoolep/hraftd/blob/master/store/store.go
import (
	"io"
	"library-manager/cad-server/internal/database"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

const raftTimeout = 10 * time.Second

// Objeto que implementa a interface do raft (Apply, Snapshot, Restore)
type Store struct {
	raftDir  string
	raftBind string

	mu   sync.Mutex
	raft *raft.Raft

	userRepo database.UserRepo
	bookRepo database.BookRepo
}

func NewStore(raftDir string, user database.UserRepo, book database.BookRepo) *Store {
	return &Store{
		userRepo: user,
		bookRepo: book,

		raftDir: raftDir,
	}
}
func (s *Store) Open(localId string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localId)
	addr, err := net.ResolveTCPAddr("tcp", s.raftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.raftBind, addr, 3, raftTimeout, os.Stderr)
	if err != nil {
		return err
	}
	snapshots, err := raft.NewFileSnapshotStore(s.raftDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	var logStore raft.LogStore
	var stableStore raft.StableStore

}

func (f *Store) Apply(log *raft.Log) interface{} {
	return nil
}

func (f *Store) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *Store) Restore(rc io.ReadCloser) error {
	return nil
}
