GO := go
BUILD_DIR := ../bin

# Inputs
BIB_SERVER := ./bib-server/cmd/main.go
CAD_SERVER := ./cad-server/cmd/main.go
BIB_CLIENT := ./bib-client/cmd/main.go
CAD_CLIENT := ./cad-client/cmd/main.go


# Binaries
BIB_SERVER_BIN := bib-server
CAD_SERVER_BIN := cad-server
BIB_CLIENT_BIN := bib-client
CAD_CLIENT_BIN := cad-client


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

clean:
	@echo "Limpando os arquivos binários..."
	@rm -rf $(BUILD_DIR)
