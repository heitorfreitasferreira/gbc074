#!/bin/bash

if [ -x "./bin/crud-client" ]; then
    ./bin/crud-client --port "$1"
fi