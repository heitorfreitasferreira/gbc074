package server

import (
	"context"
	"fmt"

	br_ufu_facom_gbc074_projeto_cadastro "github.com/rpc-mqtt-library-manager/crud-terminal-server/api"
)

type Server struct {
	br_ufu_facom_gbc074_projeto_cadastro.UnimplementedPortalCadastroServer
}

func (s *Server) NovoUsuario(ctx context.Context, usuario *br_ufu_facom_gbc074_projeto_cadastro.Usuario) (*br_ufu_facom_gbc074_projeto_cadastro.Status, error) {
	fmt.Printf("Name: %s\nCPF: %s", usuario.Nome, usuario.Cpf)
	return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 0}, nil
}
