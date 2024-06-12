GO := go
BUILD_DIR := ../bin
BIN_NAME := cms-server

.PHONY: all build clean

all: build

install-deps:
	@echo "Instalando as dependências do servidor CMS."
	@$(GO) mod download

build:install-deps
	@echo "Compilando o servidor CMS."
	@$(GO) build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/main.go

clean:
	@echo "Limpando os arquivos binários..."
	@rm -rf $(BUILD_DIR)/$(BIN_NAME)
