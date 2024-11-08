// Package httpd provides the HTTP server for accessing the distributed key-value store.
// It also provides the endpoint for other nodes to join an existing cluster.
package httpd

import (
	"encoding/json"
	"io"
	"library-manager/book-database/database"
	"log"
	"net"
	"net/http"
	"strings"
)

// Store is the interface Raft-backed stores must implement
type Store interface {
	Get(isbn database.ISBN) (database.Book, error)
	GetAll() ([]database.Book, error)
	Create(value database.Book) error
	Edit(isbn database.ISBN, value database.Book) error
	Delete(database.ISBN) error
	// Join joins the node, identitifed by nodeID and reachable at addr, to the cluster.
	Join(nodeID string, addr string) error
}

// Service provides HTTP service.
type Service struct {
	addr string
	ln   net.Listener

	store Store
}

// New returns an uninitialized HTTP service.
func New(addr string, store Store) *Service {
	return &Service{
		addr:  addr,
		store: store,
	}
}

// Start starts the service.
func (s *Service) Start() error {
	server := http.Server{
		Handler: s,
	}

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.ln = ln

	http.Handle("/", s)

	go func() {
		err := server.Serve(s.ln)
		if err != nil {
			log.Fatalf("HTTP serve: %s", err)
		}
	}()

	return nil
}

// Close closes the service.
func (s *Service) Close() {
	s.ln.Close()
	return
}

// ServeHTTP allows Service to serve HTTP requests.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/book") {
		s.handleBookRequest(w, r)
	} else if r.URL.Path == "/join" {
		s.handleJoin(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Service) handleJoin(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(m) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	remoteAddr, ok := m["addr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nodeID, ok := m["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.store.Join(nodeID, remoteAddr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) handleBookRequest(w http.ResponseWriter, r *http.Request) {
	getKey := func() string {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			return ""
		}
		return parts[2]
	}

	switch r.Method {
	case "GET":
		k := getKey()
		// If there is no key, get all items.
		if k == "" {
			v, err := s.store.GetAll()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b, err := json.Marshal(v)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			io.WriteString(w, string(b))
			return
		}
		// If there is a key get the item.
		v, err := s.store.Get(database.ISBN(k))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(b))

	case "POST":
		// Read the value from the POST body.
		v := database.Book{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.store.Create(v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "PUT":
		k := getKey()
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		v := database.Book{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.Edit(database.ISBN(k), v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "DELETE":
		k := getKey()
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.Delete(database.ISBN(k)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

// Addr returns the address on which the Service is listening
func (s *Service) Addr() net.Addr {
	return s.ln.Addr()
}
