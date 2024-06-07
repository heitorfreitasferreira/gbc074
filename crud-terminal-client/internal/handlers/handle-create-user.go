package handlers

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/screens"
)

func NewCreateUserUseCase() {
	userRead, err := screens.NewReadUserScreen()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("CPF: %s\n", userRead.Cpf)
	}
}
