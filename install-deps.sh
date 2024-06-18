#!/bin/bash

echo "Instalando dependências...\n\n\n"
sudo yay -S --needed go mosquitto protoc

# echo "Gerando os arquivos de código a partir do arquivo proto\n\n\n"
# # Gerar os arquivos de código a partir do arquivo proto no cliente do crud
# protoc --go_out=./crud-terminal-client --go_opt=paths=source_relative --go-grpc_out=./crud-terminal-client --go-grpc_opt=paths=source_relative api/portal-administrativo.proto

# protoc --go_out=./crud-terminal-server --go_opt=paths=source_relative --go-grpc_out=./crud-terminal-server --go-grpc_opt=paths=source_relative api/portal-administrativo.proto

# protoc --go_out=./cms-client --go_opt=paths=source_relative --go-grpc_out=./cms-client --go-grpc_opt=paths=source_relative api/portal-biblioteca.proto

# protoc --go_out=./cms-server --go_opt=paths=source_relative --go-grpc_out=./cms-server --go-grpc_opt=paths=source_relative api/portal-biblioteca.proto

