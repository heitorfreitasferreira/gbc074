package handlers

import (
	"context"
	"fmt"
	"time"

	br_ufu_facom_gbc074_projeto_cadastro "github.com/rpc-mqtt-library-manager/crud-terminal-client/api"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/screens"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/utils"
)

func NewCreateUserUseCase() {
	userRead, err := screens.NewReadUserScreen()
	if err != nil {
		panic(err)
	}

	conn := utils.GetConn()
	defer conn.Close()

	client := br_ufu_facom_gbc074_projeto_cadastro.NewPortalCadastroClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	response, err := client.NovoUsuario(ctx, &br_ufu_facom_gbc074_projeto_cadastro.Usuario{
		Nome: userRead.Nome,
		Cpf:  userRead.Cpf,
	})
	defer cancel()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response[%d]: %s\n", response.Status, response.Msg)
}
