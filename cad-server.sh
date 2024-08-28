#!/bin/bash

if [ -x "./bin/cad-server" ]; then
    ./bin/cad-server --port "$1"
fi
