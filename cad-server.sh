#!/bin/bash

if [ -x "./bin/crud-server" ]; then
    ./bin/crud-server --port "$1"
fi