#!/bin/bash

# Default port value
PORT=""

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --port) PORT="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

# Check if the cad-server executable exists and is executable
if [ -x "./bin/cad-server" ]; then
    if [ -n "$PORT" ]; then
        ./bin/cad-server --port "$PORT"
    else
	./bin/cad-server 
    fi
else
    echo "Error: ./bin/cad-server not found or not executable."
    exit 1
fi

