#!/bin/bash

set -e

echo "=== eBayClone gRPC Server Build and Run Script ==="

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install Protocol Buffers compiler."
    echo "Visit: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

echo "1. Installing Go dependencies..."
go mod tidy

echo "2. Installing protoc-gen-go and protoc-gen-go-grpc..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "3. Creating proto output directory..."
mkdir -p proto

echo "4. Compiling proto files..."
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/ebayclone.proto

echo "5. Building gRPC server..."
go build -o bin/server src/main.go

echo "6. Starting gRPC server on port 50051..."
echo "Server will be available at localhost:50051"
echo "Press Ctrl+C to stop the server"
echo ""

./bin/server
