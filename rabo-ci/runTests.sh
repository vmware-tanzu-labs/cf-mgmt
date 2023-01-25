#!/bin/bash

go version

echo "Running go vet"
go vet ./...
echo "Running go test"
go test ./...

