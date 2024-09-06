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

# Check if the bib-client executable exists and is executable
if [ -x "./bin/bib-client" ]; then
    if [ -n "$PORT" ]; then
        ./bin/bib-client --port "$PORT"
    else
	./bin/bib-client 
    fi
else
    echo "Error: ./bin/bib-client not found or not executable."
    exit 1
fi

