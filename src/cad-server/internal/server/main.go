package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	book_database "library-manager/book-database/database"
	api_cad "library-manager/shared/api/cad"
	"library-manager/shared/utils"
	user_database "library-manager/user-database/database"
)

type Server struct {
	api_cad.UnimplementedPortalCadastroServer
	userDatabaseAddr string // Endereço http do banco, ex. http://localhost:21000
	bookDatabaseAddr string
}

func NewServer(userDatabaseAddr, bookDatabaseAddr string) *Server {
	return &Server{
		userDatabaseAddr: userDatabaseAddr,
		bookDatabaseAddr: bookDatabaseAddr,
	}
}

func (s *Server) NovoUsuario(ctx context.Context, usuario *api_cad.Usuario) (*api_cad.Status, error) {
	cpf := utils.CPF(usuario.Cpf)
	if !cpf.Validate() {
		log.Printf("CPF inválido")
		return &api_cad.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	jsonData, err := json.Marshal(usuario)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	req, err := http.Post(s.userDatabaseAddr+"/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao criar usuário: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar usuário"}, nil
	}
	defer req.Body.Close()

	return &api_cad.Status{Status: 0}, nil
}

func (s *Server) ObtemUsuario(ctx context.Context, request *api_cad.Identificador) (*api_cad.Usuario, error) {

	log.Default().Printf("Obtendo usuário com CPF: %s", request.Id)

	cpf := utils.CPF(request.Id)
	if !cpf.Validate() {
		log.Printf("CPF inválido")
		return nil, errors.New("CPF inválido")
	}

	req, err := http.NewRequest("GET", s.userDatabaseAddr+"/user/"+request.Id, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return nil, errors.New("Erro ao criar requisição")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao obter usuário: %v", err)
		return nil, errors.New("Erro ao obter usuário")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro ao obter usuário: código de status %d", resp.StatusCode)
		return nil, errors.New("Usuário não encontrado")
	}

	var usuario user_database.User
	jsonData := json.NewDecoder(resp.Body)
	if err := jsonData.Decode(&usuario); err != nil {
		log.Printf("Erro ao decodificar dados do usuário: %v", err)
		return nil, errors.New("Erro ao decodificar dados do usuário")

	}
	usuario_proto := user_database.UserToProto(usuario)

	return &usuario_proto, nil
}

func (s *Server) ObtemTodosUsuarios(request *api_cad.Vazia, stream api_cad.PortalCadastro_ObtemTodosUsuariosServer) error {
	req, err := http.Get(s.userDatabaseAddr + "/user")
	if err != nil {
		log.Printf("Erro ao obter usuários: %v", err)
		return errors.New("Erro ao obter usuários")
	}
	defer req.Body.Close()

	jsonData := json.NewDecoder(req.Body)
	var users []user_database.User
	if err := jsonData.Decode(&users); err != nil {
		log.Printf("Erro ao decodificar dados: %v", err)
		return errors.New("Erro ao decodificar dados")
	}

	for _, user := range users {
		protoUser := user_database.UserToProto(user)
		err := stream.Send(&protoUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) EditaUsuario(ctx context.Context, usuario *api_cad.Usuario) (*api_cad.Status, error) {

	log.Default().Printf("Editando usuário %s", usuario.Cpf)

	user := user_database.ProtoToUser(usuario)

	if !user.Cpf.Validate() {
		log.Printf("CPF inválido")
		return &api_cad.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	req, err := http.NewRequest("PUT", s.userDatabaseAddr+"/user"+usuario.Cpf, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Printf("Erro ao criar requisição HTTP: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar requisição HTTP"}, nil
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao editar usuário: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao editar usuário"}, nil
	}
	defer resp.Body.Close()

	return &api_cad.Status{Status: 0, Msg: "Usuário editado com sucesso"}, nil
}

func (s *Server) RemoveUsuario(ctx context.Context, request *api_cad.Identificador) (*api_cad.Status, error) {

	cpf := utils.CPF(request.Id)
	if !cpf.Validate() {
		log.Printf("CPF inválido")
		return &api_cad.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	req, err := http.NewRequest("DELETE", s.userDatabaseAddr+"/user/"+request.Id, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição HTTP: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar requisição HTTP"}, nil
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao remover usuário: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao remover usuário"}, nil
	}
	defer resp.Body.Close()

	return &api_cad.Status{Status: 0, Msg: "Usuário removido com sucesso"}, nil

}

func (s *Server) NovoLivro(ctx context.Context, livro *api_cad.Livro) (*api_cad.Status, error) {
	isbn := utils.ISBN(livro.Isbn)
	if !isbn.Validate() {
		return &api_cad.Status{Status: 1, Msg: "ISBN inválido"}, nil
	}

	jsonData, err := json.Marshal(livro)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	req, err := http.NewRequest("POST", s.bookDatabaseAddr+"/book", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao criar livro: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar livro"}, nil
	}

	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar requisição para criar livro: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar livro"}, nil
	}
	defer resp.Body.Close()

	return &api_cad.Status{Status: 0, Msg: "Livro criado com sucesso"}, nil
}

func (s *Server) ObtemLivro(ctx context.Context, request *api_cad.Identificador) (*api_cad.Livro, error) {

	isbn := utils.ISBN(request.Id)
	if !isbn.Validate() {
		return nil, errors.New("ISBN inválido")
	}
	req, err := http.NewRequest("GET", s.userDatabaseAddr+"/book/"+request.Id, nil)

	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return nil, errors.New("Erro ao criar requisição")
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Erro ao obter livro: %v", err)
		return nil, errors.New("Erro ao obter livro")
	}
	defer resp.Body.Close()

	var book book_database.Book
	jsonData := json.NewDecoder(resp.Body)
	if err := jsonData.Decode(&book); err != nil {
		log.Printf("Erro ao decodificar dados do usuário: %v", err)
		return nil, errors.New("Erro ao decodificar dados do usuário")

	}
	book_proto := book_database.BookToProto(book)

	return &book_proto, nil
}

func (s *Server) ObtemTodosLivros(request *api_cad.Vazia, stream api_cad.PortalCadastro_ObtemTodosLivrosServer) error {

	req, err := http.Get(s.bookDatabaseAddr + "/books/")
	if err != nil {
		log.Printf("Erro ao fazer requisição para obter todos os livros: %v", err)
		return fmt.Errorf("erro ao fazer requisição para obter todos os livros: %w", err)
	}
	defer req.Body.Close()

	var books []book_database.Book
	jsonData := json.NewDecoder(req.Body)
	if err := jsonData.Decode(&books); err != nil {
		log.Printf("Erro ao decodificar dados dos livros: %v", err)
		return fmt.Errorf("erro ao decodificar dados dos livros: %w", err)
	}

	for _, book := range books {
		protoBook := book_database.BookToProto(book)
		if err := stream.Send(&protoBook); err != nil {
			log.Printf("Erro ao enviar livro na stream: %v", err)
			return fmt.Errorf("erro ao enviar livro na stream: %w", err)
		}
	}

	return nil
}

func (s *Server) EditaLivro(ctx context.Context, livro *api_cad.Livro) (*api_cad.Status, error) {

	isbn := utils.ISBN(livro.Isbn)
	if !isbn.Validate() {
		return &api_cad.Status{Status: 1, Msg: "ISBN inválido"}, nil
	}

	jsonData, err := json.Marshal(livro)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	req, err := http.NewRequest("PUT", s.bookDatabaseAddr+"/book/"+livro.Isbn, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar requisição"}, nil
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao executar requisição para editar livro: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao editar livro"}, nil
	}
	defer resp.Body.Close()

	return &api_cad.Status{Status: 0, Msg: "Livro atualizado com sucesso"}, nil
}

func (s *Server) RemoveLivro(ctx context.Context, request *api_cad.Identificador) (*api_cad.Status, error) {

	isbn := utils.ISBN(request.Id)
	if !isbn.Validate() {
		return &api_cad.Status{Status: 1, Msg: "ISBN inválido"}, nil
	}

	req, err := http.NewRequest("DELETE", s.bookDatabaseAddr+"/book/"+request.Id, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao criar requisição"}, nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar requisição: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao enviar requisição"}, nil
	}
	defer resp.Body.Close()

	return &api_cad.Status{Status: 0, Msg: "Livro removido com sucesso"}, nil
}
