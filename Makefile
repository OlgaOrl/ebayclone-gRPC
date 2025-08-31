.PHONY: proto build run client test clean

# Default target
all: proto build

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@mkdir -p proto
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/ebayclone.proto

# Build the server
build: proto
	@echo "Building server..."
	@mkdir -p bin
	go build -o bin/server src/main.go

# Run the server
run: build
	@echo "Starting gRPC server on port 50051..."
	./bin/server

# Run client example
client:
	@echo "Running client example..."
	go run client/example.go

# Run tests
test:
	@echo "Running tests..."
	chmod +x tests/test.sh
	./tests/test.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f proto/*.pb.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Help
help:
	@echo "Available targets:"
	@echo "  all     - Generate proto and build (default)"
	@echo "  proto   - Generate protobuf code"
	@echo "  build   - Build the server"
	@echo "  run     - Build and run the server"
	@echo "  client  - Run client example"
	@echo "  test    - Run automated tests"
	@echo "  clean   - Clean build artifacts"
	@echo "  deps    - Install dependencies"
	@echo "  help    - Show this help"
