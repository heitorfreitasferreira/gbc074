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

    "library-manager/bib-server/internal/database"
    "library-manager/bib-server/internal/queue/handlers"
    api_bib "library-manager/shared/api/bib"
    api_cad "library-manager/shared/api/cad"
    "library-manager/shared/utils"
)

type Server struct {
    api_cad.UnimplementedPortalCadastroServer
    userDatabaseAddr string // Endereço http do banco, ex. http://localhost:21000
    bookDatabaseAddr string
    raftServerAddr   string // Endereço http do servidor Raft, ex. http://localhost:22000
}

func NewServer(userDatabaseAddr, bookDatabaseAddr, raftServerAddr string) *Server {
    return &Server{
        userDatabaseAddr: userDatabaseAddr,
        bookDatabaseAddr: bookDatabaseAddr,
        raftServerAddr:   raftServerAddr,
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

        resp, err := http.Post(fmt.Sprintf("%s/realiza-emprestimo", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
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

        if err := stream.Send(&status); err != nil {
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

        resp, err := http.Post(fmt.Sprintf("%s/realiza-devolucao", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
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

        if err := stream.Send(&status); err != nil {
            return err
        }
    }

    return nil
}

func (s *Server) BloqueiaUsuarios(ctx context.Context, request *api_cad.ListaIdentificadores) (*api_cad.Status, error) {
    jsonData, err := json.Marshal(request)

    if err != nil {
        log.Printf("Erro ao converter dados dos identificadores para JSON: %v", err)
        return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
    }

    resp, err := http.Post(fmt.Sprintf("%s/bloqueia-usuarios", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Erro ao bloquear usuários: %v", err)
        return &api_cad.Status{Status: 1, Msg: "Erro ao bloquear usuários"}, nil
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
        return &api_cad.Status{Status: 1, Msg: "Erro na resposta do servidor Raft"}, errors.New(resp.Status)
    }

    var status api_cad.Status
    if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
        log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
        return &api_cad.Status{Status: 1, Msg: "Erro ao decodificar resposta do servidor Raft"}, err
    }

    return &status, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, request *api_cad.ListaIdentificadores) (*api_cad.Status, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        log.Printf("Erro ao converter dados dos identificadores para JSON: %v", err)
        return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
    }

    resp, err := http.Post(fmt.Sprintf("%s/libera-usuarios", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Erro ao liberar usuários: %v", err)
        return &api_cad.Status{Status: 1, Msg: "Erro ao liberar usuários"}, nil
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
        return &api_cad.Status{Status: 1, Msg: "Erro na resposta do servidor Raft"}, errors.New(resp.Status)
    }

    var status api_cad.Status

    if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
        log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
        return &api_cad.Status{Status: 1, Msg: "Erro ao decodificar resposta do servidor Raft"}, err
    }

    return &status, nil
}

func (s *Server) ListaUsuariosBloqueados(ctx context.Context, request *api_cad.Empty) (*api_cad.ListaUsuarios, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        log.Printf("Erro ao converter dados para JSON: %v", err)
        return nil, err
    }

    resp, err := http.Post(fmt.Sprintf("%s/lista-usuarios-bloqueados", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Erro ao listar usuários bloqueados: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
        return nil, errors.New(resp.Status)
    }

    var listaUsuarios api_cad.ListaUsuarios
    if err := json.NewDecoder(resp.Body).Decode(&listaUsuarios); err != nil {
        log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
        return nil, err
    }

    return &listaUsuarios, nil
}

func (s *Server) ListaLivrosEmprestados(ctx context.Context, request *api_bib.Empty) (*api_bib.ListaLivros, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        log.Printf("Erro ao converter dados para JSON: %v", err)
        return nil, err
    }

    resp, err := http.Post(fmt.Sprintf("%s/lista-livros-emprestados", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Erro ao listar livros emprestados: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
        return nil, errors.New(resp.Status)
    }

    var listaLivros api_bib.ListaLivros
    if err := json.NewDecoder(resp.Body).Decode(&listaLivros); err != nil {
        log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
        return nil, err
    }

    return &listaLivros, nil
}

func (s *Server) ListaLivrosEmFalta(ctx context.Context, request *api_bib.Empty) (*api_bib.ListaLivros, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        log.Printf("Erro ao converter dados para JSON: %v", err)
        return nil, err
    }

    resp, err := http.Post(fmt.Sprintf("%s/lista-livros-em-falta", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Erro ao listar livros em falta: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
        return nil, errors.New(resp.Status)
    }

    var listaLivros api_bib.ListaLivros
    if err := json.NewDecoder(resp.Body).Decode(&listaLivros); err != nil {
        log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
        return nil, err
    }

    return &listaLivros, nil
}

func (s *Server) PesquisaLivro(ctx context.Context, request *api_bib.PesquisaLivroRequest) (*api_bib.ListaLivros, error) {
    jsonData, err := json.Marshal(request)
    if err != nil {
        log.Printf("Erro ao converter dados da pesquisa para JSON: %v", err)
        return nil, err
    }

    resp, err := http.Post(fmt.Sprintf("%s/pesquisa-livro", s.raftServerAddr), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Erro ao pesquisar livro: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Erro na resposta do servidor Raft: %v", resp.Status)
        return nil, errors.New(resp.Status)
    }

    var listaLivros api_bib.ListaLivros
    if err := json.NewDecoder(resp.Body).Decode(&listaLivros); err != nil {
        log.Printf("Erro ao decodificar resposta do servidor Raft: %v", err)
        return nil, err
    }

    return &listaLivros, nil
}
