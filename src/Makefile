GO := go
BUILD_DIR := ../bin

# Inputs
BIB_SERVER := ./bib-server/cmd/main.go
CAD_SERVER := ./cad-server/cmd/main.go
BIB_CLIENT := ./bib-client/cmd/main.go
CAD_CLIENT := ./cad-client/cmd/main.go

BOOK_DATABASE := ./book-database/cmd/main.go  
USER_DATABASE := ./user-database/cmd/main.go   


# Binaries
BIB_SERVER_BIN := bib-server
CAD_SERVER_BIN := cad-server
BIB_CLIENT_BIN := bib-client
CAD_CLIENT_BIN := cad-client

BOOK_DATABASE_BIN := book-database
USER_DATABASE_BIN := user-database


.PHONY: all install-deps bib-server cad-server bib-client cad-client clean

all: clean bib-server cad-server bib-client cad-client

install-deps:
	@echo "Instalando dependências..."
	@$(GO) mod download

bib-server: install-deps
	@echo "Compilando o servidor do Portal Biblioteca..."
	@$(GO) build -o $(BUILD_DIR)/$(BIB_SERVER_BIN) $(BIB_SERVER)

cad-server: install-deps
	@echo "Compilando o servidor do Portal Cadastro..."
	@$(GO) build -o $(BUILD_DIR)/$(CAD_SERVER_BIN) $(CAD_SERVER)

bib-client: install-deps
	@echo "Compilando o cliente do Portal Biblioteca..."
	@$(GO) build -o $(BUILD_DIR)/$(BIB_CLIENT_BIN) $(BIB_CLIENT)

cad-client: install-deps
	@echo "Compilando o cliente do Portal Cadastro..."
	@$(GO) build -o $(BUILD_DIR)/$(CAD_CLIENT_BIN) $(CAD_CLIENT)

book-database: install-deps
	@echo "Compilando o banco de dados de livros..."
	@$(GO) build -o $(BUILD_DIR)/$(BOOK_DATABASE_BIN) $(BOOK_DATABASE)


user-database: install-deps
	@echo "Compilando o banco de dados de usuários..."
	@$(GO) build -o $(BUILD_DIR)/$(USER_DATABASE_BIN) $(USER_DATABASE)


clean:
	@echo "Limpando os arquivos binários..."
	@rm -rf $(BUILD_DIR)
