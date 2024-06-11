CIENT_BINS := bin/crud-client

SERVER_BINS := bin/crud-server

$(CLIENT_BINS): $(wildcard crud-terminal-client/**/*.go)
	@echo "Compilando o cliente..."
	@go build -o bin/crud-client ./crud-terminal-client/cmd/main.go

$(SERVER_BINS): $(wildcard crud-terminal-server/**/*.go)
	@echo "Compilando o servidor..."
	@go build -o bin/crud-server ./crud-terminal-server/cmd/main.go