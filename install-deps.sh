#!/bin/bash

echo "Instalando dependências...\n\n\n"
sudo yay -S --needed go mosquitto protoc

echo "Gerando os arquivos de código a partir do arquivo proto\n\n\n"
# Gerar os arquivos de código a partir do arquivo proto no cliente do crud
protoc --go_out=./crud-terminal-client --go_opt=paths=source_relative --go-grpc_out=./crud-terminal-client --go-grpc_opt=paths=source_relative api/portal-administrativo.proto

protoc --go_out=./crud-terminal-server --go_opt=paths=source_relative --go-grpc_out=./crud-terminal-server --go-grpc_opt=paths=source_relative api/portal-administrativo.proto

echo "Gerando os binários"

cd crud-terminal-client

make build

cd ..
