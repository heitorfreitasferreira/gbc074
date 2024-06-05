package utils

import (
	"log"

	"google.golang.org/grpc"
)

// NÃO SE ESQUEÇA DE FECHAR A CONEXÃO DEPOIS DE USAR!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func GetConn(port int) *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:"+string(rune(port)), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
		panic(err)
	}
	return conn
}
