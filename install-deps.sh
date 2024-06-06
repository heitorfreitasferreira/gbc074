#!/bin/bash

sudo pacman -S --needed go mosquitto protoc

# Gerar os arquivos de código a partir do arquivo proto no cliente do crud
protoc --go_out=./crud-terminal-client --go_opt=paths=source_relative --go-grpc_out=./crud-terminal-client --go-grpc_opt=paths=source_relative api/portal-administrativo.proto
