#!/bin/bash

cd ./crud-terminal-server
go get -v -t -d ./...
make build

cd ..

cd ./crud-terminal-client
go get -v -t -d ./...
make build

cd ..
