GO := go
BUILD_DIR := ../bin
BIN_NAME := crud-server

.PHONY: all build clean

all: build

build:
	@echo "Compilando o servidor administrativo..."
	@$(GO) build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/main.go

clean:
	@echo "Limpando os arquivos binários..."
	@rm -rf $(BUILD_DIR)/$(BIN_NAME)