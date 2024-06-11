#!/bin/bash

directories=(
    "./crud-terminal-server"
    "./crud-terminal-client"
)

for dir in "${directories[@]}"
do
    cd "$dir"
    echo "Getting dependencies for $dir"

    go get -v -t -d ./...
    make build

    cd ..
done
