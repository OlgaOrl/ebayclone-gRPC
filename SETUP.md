# eBayClone gRPC Setup Guide

This guide provides step-by-step instructions to set up and run the eBayClone gRPC service.

## Step 1: Install Prerequisites

### Install Go (Required)
1. Visit https://golang.org/doc/install
2. Download Go 1.21 or later for your operating system
3. Follow installation instructions
4. Verify: `go version`

### Install Protocol Buffers Compiler (Required)
1. Visit https://grpc.io/docs/protoc-installation/
2. Download protoc for your operating system
3. Add protoc to your PATH
4. Verify: `protoc --version`

### Install Git (Required)
1. Visit https://git-scm.com/downloads
2. Download and install Git
3. Verify: `git --version`

## Step 2: Build and Run

### Option A: Using Build Scripts

**On Unix/Linux/macOS:**
```bash
# Make scripts executable
chmod +x scripts/run.sh scripts/verify.sh tests/test.sh

# Verify setup
./scripts/verify.sh

# Build and run server
./scripts/run.sh
```

**On Windows:**
```cmd
# Verify setup (if using WSL/Git Bash)
scripts\verify.sh

# Build and run server
scripts\run.bat
```

### Option B: Using Makefile (Unix/Linux/macOS)
```bash
# Install dependencies
make deps

# Build everything
make all

# Run server
make run
```

### Option C: Manual Steps
```bash
# 1. Install Go dependencies
go mod tidy

# 2. Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 3. Generate protobuf code
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/ebayclone.proto

# 4. Build server
mkdir -p bin
go build -o bin/server src/main.go

# 5. Run server
./bin/server  # Unix/Linux/macOS
# OR
bin\server.exe  # Windows
```

## Step 3: Test the Implementation

### Run Client Example
In a new terminal:
```bash
go run client/example.go
```

### Run Automated Tests
```bash
# Install grpcurl for testing
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Run tests
./tests/test.sh  # Unix/Linux/macOS
```

### Manual Testing with grpcurl
```bash
# List services
grpcurl -plaintext localhost:50051 list

# Test user creation
grpcurl -plaintext -d '{"username":"test","email":"test@example.com","password":"pass"}' \
localhost:50051 ebayclone.UserService/CreateUser
```

## Step 4: Python Client (Optional)

```bash
# Install Python dependencies
pip install -r client/requirements.txt

# Generate Python protobuf files
./scripts/generate_python_proto.sh

# Run Python client
python3 client/example.py
```

## Troubleshooting

### Common Issues

1. **"protoc: command not found"**
   - Install Protocol Buffers compiler from https://grpc.io/docs/protoc-installation/

2. **"go: command not found"**
   - Install Go from https://golang.org/doc/install
   - Make sure Go is in your PATH

3. **"Failed to listen: address already in use"**
   - Kill any existing server: `pkill -f server` or `taskkill /f /im server.exe`
   - Or change port in src/main.go

4. **Import errors in Go**
   - Run `go mod tidy` to resolve dependencies
   - Make sure protobuf files are generated

5. **Permission denied on scripts**
   - On Unix: `chmod +x scripts/*.sh tests/*.sh`
   - On Windows: Use Git Bash or WSL

### Verification Checklist

- [ ] Go 1.21+ installed and in PATH
- [ ] protoc installed and in PATH
- [ ] All .proto files compile without errors
- [ ] Server builds successfully
- [ ] Server starts on port 50051
- [ ] Client example runs without errors
- [ ] All tests pass

## Success Indicators

When everything is working correctly, you should see:

1. **Server startup:**
   ```
   gRPC server starting on port 50051...
   ```

2. **Client example output:**
   ```
   === eBayClone gRPC Client Example ===
   1. Creating user...
   Created user: ID=1, Username=testuser, Email=test@example.com
   ...
   === All operations completed successfully! ===
   ```

3. **Test results:**
   ```
   === Test Results ===
   Tests Passed: 12
   Tests Failed: 0
   All tests passed! ✓
   ```

## Next Steps

Once the basic setup is working:
1. Explore the API using grpcurl
2. Implement additional client applications
3. Add persistent storage (PostgreSQL, MongoDB)
4. Implement proper authentication middleware
5. Add TLS/SSL for production deployment
