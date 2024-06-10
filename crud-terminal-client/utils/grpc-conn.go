package utils

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NÃO SE ESQUEÇA DE FECHAR A CONEXÃO DEPOIS DE USAR!!!!!!!!!!!!!!!!!!!!!!!!!!!!
var Port string

func GetConn() *grpc.ClientConn {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
		panic(err)
	}
	return conn
}
