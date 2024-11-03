// Package httpd provides the HTTP server for accessing the distributed key-value store.
// It also provides the endpoint for other nodes to join an existing cluster.
package httpd

import (
	"encoding/json"
	"io"
	"library-manager/shared/utils"
	"library-manager/user-database/database"
	"log"
	"net"
	"net/http"
	"strings"
)

// Store is the interface Raft-backed stores must implement
type Store interface {
	GetUser(cpf utils.CPF) (database.User, error)
	GetAllUsers() ([]database.User, error)
	CreateUser(value database.User) error
	EditUser(cpf utils.CPF, value database.User) error
	DeleteUser(utils.CPF) error
	CreateUserBook(value database.UserBook) error
	DeleteUserBook(value database.UserBook) error
	GetAllUserBooks() []database.UserBook
	GetUserLoans(userId string) []database.LoanBookAndTime
	RemoveUserLoans(value database.UserBook) error
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
}

// ServeHTTP allows Service to serve HTTP requests.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/user") {
		s.handleUserRequest(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/user/loan") {
		s.handleUserLoanRequest(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/user-book") {
		s.handleUserBookRequest(w, r)
	} else if r.URL.Path == "/join" {
		s.handleJoin(w, r)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
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

func (s *Service) handleUserRequest(w http.ResponseWriter, r *http.Request) {
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
			v, err := s.store.GetAllUsers()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b, err := json.Marshal(v)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			io.Writer.Write(w, b)
			return
		}
		// If there is a key get the item.
		v, err := s.store.GetUser(utils.CPF(k))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Writer.Write(w, b)

	case "POST":
		// Read the value from the POST body.
		v := database.User{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.store.CreateUser(v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "PUT":
		k := getKey()
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		v := database.User{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.EditUser(utils.CPF(k), v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "DELETE":
		k := getKey()
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.DeleteUser(utils.CPF(k)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) handleUserBookRequest(w http.ResponseWriter, r *http.Request) {
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
		// If there is a key, return bad request.
		if k != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		v := s.store.GetAllUserBooks()

		b, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.Writer.Write(w, b)
		return

	case "POST":
		// Read the value from the POST body.
		v := database.UserBook{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.store.CreateUserBook(v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "DELETE":
		// Read the value from the DELETE body.
		v := database.UserBook{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.DeleteUserBook(v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) handleUserLoanRequest(w http.ResponseWriter, r *http.Request) {
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
		if k == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		userLoans := s.store.GetUserLoans(k)
		b, err := json.Marshal(userLoans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Writer.Write(w, b)
		break
	case "DELETE":
		// Read the value from the DELETE body.
		v := database.UserBook{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := s.store.RemoveUserLoans(v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

}

// Addr returns the address on which the Service is listening
func (s *Service) Addr() net.Addr {
	return s.ln.Addr()
}
