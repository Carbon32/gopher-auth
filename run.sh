#!/bin/bash

if [ "$1" == "-t" ]; then
    go test ./...
else
    go run cmd/main.go
fi