GO := go
BUILD_DIR := ../bin
BIN_NAME := crud-client

.PHONY: all build clean

all: build

build:
	@echo "Compilando o programa..."
	@$(GO) build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/main.go

clean:
	@echo "Limpando os arquivos binários..."
	@rm -rf $(BUILD_DIR)/$(BIN_NAME)