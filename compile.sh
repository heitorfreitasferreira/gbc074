#!/bin/bash

directories=(
    "./crud-terminal-server"
    "./crud-terminal-client"
    "./cms-server"
)

for dir in "${directories[@]}"
do
    cd "$dir"
    make build
    cd ..
done
