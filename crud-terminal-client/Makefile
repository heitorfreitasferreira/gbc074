GO := go
BUILD_DIR := ../bin
BIN_NAME := crud-client

.PHONY: all build clean

all: build

install-deps:
	@echo "Instalando as dependências do cliente administrativo..."
	@$(GO) mod download

build:install-deps
	@echo "Compilando o cliente administrativo."
	@$(GO) build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/main.go

clean:
	@echo "Limpando os arquivos binários..."
	@rm -rf $(BUILD_DIR)/$(BIN_NAME)