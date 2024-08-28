#!/bin/bash

if [ -x "./bin/bib-server" ]; then
    ./bin/bib-server --port "$1"
fi
