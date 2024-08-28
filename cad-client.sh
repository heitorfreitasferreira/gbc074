#!/bin/bash

if [ -x "./bin/cad-client" ]; then
    ./bin/cad-client --port "$1"
fi
