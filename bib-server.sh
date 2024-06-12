#!/bin/bash

if [ -x "./bin/cms-server" ]; then
    ./bin/cms-server --port "$1"
fi