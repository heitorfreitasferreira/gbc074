package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	book_database "library-manager/book-database/database"
	api_bib "library-manager/shared/api/bib"
	"library-manager/user-database/database"
)

// TODO: Usar endpoints validos para obter os dados necessários.
// TODO: Se basear em como o arquivo era antes (https://github.com/heitorfreitasferreira/gbc074/blob/897c656f9661339bc2329405fcfc1fbb886c3c0b/src/bib-server/internal/server/main.go) para fazer as requisições ao raft.
// Implementar regra de negócio que tiver faltando nos endpoints raft (user-database e book-database).

type Server struct {
	api_bib.UnsafePortalBibliotecaServer
	userDatabaseAddr string // Endereço http do banco, ex. http://localhost:21000
	bookDatabaseAddr string
}

func NewServer(userDatabaseAddr, bookDatabaseAddr string) *Server {
	return &Server{
		userDatabaseAddr: userDatabaseAddr,
		bookDatabaseAddr: bookDatabaseAddr,
	}
}

var qos byte = 2

func (s *Server) RealizaEmprestimo(stream api_bib.PortalBiblioteca_RealizaEmprestimoServer) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Erro ao converter dados do empréstimo para JSON: %v", err)
			return err
		}

		resp, err := http.Post(fmt.Sprintf("%s/realiza-emprestimo", s.bookDatabaseAddr), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Erro ao realizar empréstimo: %v", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
			return errors.New(resp.Status)
		}

		var status api_bib.Status
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
			return err
		}
	}

	return nil
}

func (s *Server) RealizaDevolucao(stream api_bib.PortalBiblioteca_RealizaDevolucaoServer) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Erro ao converter dados da devolução para JSON: %v", err)
			return err
		}

		resp, err := http.Post(fmt.Sprintf("%s/realiza-devolucao", s.bookDatabaseAddr), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Erro ao realizar devolução: %v", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
			return errors.New(resp.Status)
		}

		var status api_bib.Status
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
			return err
		}
	}

	return nil
}

func (s *Server) BloqueiaUsuarios(ctx context.Context, request *api_bib.Vazia) (*api_bib.Status, error) {
	jsonData, err := json.Marshal(request)

	if err != nil {
		log.Printf("Erro ao converter dados dos identificadores para JSON: %v", err)
		return &api_bib.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	resp, err := http.Post(fmt.Sprintf("%s/bloqueia-usuarios", s.userDatabaseAddr), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao bloquear usuários: %v", err)
		return &api_bib.Status{Status: 1, Msg: "Erro ao bloquear usuários"}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
		return &api_bib.Status{Status: 1, Msg: "Erro na resposta do servidor Raft"}, errors.New(resp.Status)
	}

	var status api_bib.Status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
		return &api_bib.Status{Status: 1, Msg: "Erro ao decodificar resposta do servidor Raft"}, err
	}

	return &status, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, request *api_bib.Vazia) (*api_bib.Status, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao converter dados dos identificadores para JSON: %v", err)
		return &api_bib.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	resp, err := http.Post(fmt.Sprintf("%s/libera-usuarios", s.userDatabaseAddr), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao liberar usuários: %v", err)
		return &api_bib.Status{Status: 1, Msg: "Erro ao liberar usuários"}, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
		return &api_bib.Status{Status: 1, Msg: "Erro na resposta do servidor Raft"}, errors.New(resp.Status)
	}

	var status api_bib.Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
		return &api_bib.Status{Status: 1, Msg: "Erro ao decodificar resposta do servidor Raft"}, err
	}

	return &status, nil
}

func (s *Server) ListaUsuariosBloqueados(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	resp, err := http.Get(fmt.Sprintf("%s/lista-usuarios-bloqueados", s.userDatabaseAddr))
	if err != nil {
		log.Printf("Erro ao listar usuários bloqueados: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
		return errors.New(resp.Status)
	}

	var listaUsuarios []database.User
	if err := json.NewDecoder(resp.Body).Decode(&listaUsuarios); err != nil {
		log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
		return err
	}

	// TODO: Resgatar livros que fizeram usuário ficar bloqueado.
	for _, usuario := range listaUsuarios {
		protoUser := database.UserToProto(usuario)

		blockedUser := &api_bib.UsuarioBloqueado{
			Usuario: &api_bib.Usuario{
				Cpf:  protoUser.Cpf,
				Nome: protoUser.Nome,
			},
			Livros: make([]*api_bib.Livro, 0),
		}

		stream.Send(blockedUser)
	}

	return nil
}

func (s *Server) ListaLivrosEmprestados(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	resp, err := http.Get(fmt.Sprintf("%s/lista-livros-emprestados", s.bookDatabaseAddr))
	if err != nil {
		log.Printf("Erro ao listar livros emprestados: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
		return errors.New(resp.Status)
	}

	// TODO: corrigir listagem, o endpoint do raft não retorna nesse formato.
	var listaLivros []book_database.Book
	if err := json.NewDecoder(resp.Body).Decode(&listaLivros); err != nil {
		log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
		return err
	}

	for _, livro := range listaLivros {
		stream.Send(&api_bib.Livro{
			Isbn:   string(livro.Isbn),
			Titulo: livro.Titulo,
			Autor:  livro.Autor,
			Total:  livro.Total,
		})
	}

	return nil
}

func (s *Server) ListaLivrosEmFalta(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	resp, err := http.Get(fmt.Sprintf("%s/lista-livros-em-falta", s.bookDatabaseAddr))
	if err != nil {
		log.Printf("Erro ao listar livros em falta: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
		return errors.New(resp.Status)
	}

	// TODO: corrigir listagem, o endpoint do raft não retorna nesse formato.
	var listaLivros []book_database.Book
	if err := json.NewDecoder(resp.Body).Decode(&listaLivros); err != nil {
		log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
		return err
	}

	for _, livro := range listaLivros {
		stream.Send(&api_bib.Livro{
			Isbn:   string(livro.Isbn),
			Titulo: livro.Titulo,
			Autor:  livro.Autor,
			Total:  livro.Total,
		})
	}

	return nil
}

func (s *Server) PesquisaLivro(req *api_bib.Criterio, stream api_bib.PortalBiblioteca_PesquisaLivroServer) error {

	// TODO: parse do critério de busca. Enviar critério de busca como parâmetro da requisição.
	// Tratar endpoint do raft para usar critério de busca.
	resp, err := http.Get(fmt.Sprintf("%s/pesquisa-livro", s.bookDatabaseAddr))
	if err != nil {
		log.Printf("Erro ao pesquisar livro: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
		return errors.New(resp.Status)
	}

	var listaLivros []book_database.Book
	if err := json.NewDecoder(resp.Body).Decode(&listaLivros); err != nil {
		log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
		return err
	}

	for _, livro := range listaLivros {
		stream.Send(&api_bib.Livro{
			Isbn:   string(livro.Isbn),
			Titulo: livro.Titulo,
			Autor:  livro.Autor,
			Total:  livro.Total,
		})
	}

	return nil
}
