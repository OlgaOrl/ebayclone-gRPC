# eBayClone gRPC API

A complete gRPC implementation of the eBayClone marketplace API, converted from the original OpenAPI REST specification. This implementation provides functionally equivalent services for user management, authentication, listings, and orders.

## Features

- **Complete gRPC Services**: User management, authentication, listings, and orders
- **JWT Authentication**: Secure token-based authentication
- **Search & Filtering**: Advanced search capabilities for listings and orders
- **Pagination**: Efficient pagination for large datasets
- **Error Handling**: Proper gRPC status codes with detailed error messages
- **File Upload Support**: Image upload for listings
- **In-Memory Storage**: Ready-to-run implementation with in-memory data storage

## Project Structure

```
/
├── proto/                  # Protocol Buffer definitions
│   └── ebayclone.proto    # Main service definitions
├── src/                   # Go source code
│   ├── main.go           # Server entry point
│   ├── services/         # gRPC service implementations
│   │   ├── user_service.go
│   │   ├── session_service.go
│   │   ├── listing_service.go
│   │   └── order_service.go
│   └── storage/          # Data storage layer
│       └── storage.go
├── client/               # Client examples
│   └── example.go       # Demonstration client
├── scripts/             # Build and run scripts
│   ├── run.sh          # Unix/Linux/macOS build script
│   └── run.bat         # Windows build script
├── tests/              # Automated tests
│   └── test.sh        # Functional test suite
├── go.mod             # Go module definition
└── README.md          # This file
```

## Prerequisites

### Required Software

1. **Go 1.21 or later**
   - Download: https://golang.org/doc/install
   - Verify: `go version`

2. **Protocol Buffers Compiler (protoc)**
   - Download: https://grpc.io/docs/protoc-installation/
   - Verify: `protoc --version`

3. **Git** (for cloning and version control)
   - Download: https://git-scm.com/downloads
   - Verify: `git --version`

### Optional (for testing)

4. **grpcurl** (for command-line testing)
   - Install: `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`
   - Verify: `grpcurl --version`

5. **netcat/nc** (for port checking in tests)
   - Usually pre-installed on Unix systems
   - Windows: Install via WSL or use telnet

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd ebayclone-gRPC
```

### 2. Build and Run Server

**On Unix/Linux/macOS:**
```bash
chmod +x scripts/run.sh
./scripts/run.sh
```

**On Windows:**
```cmd
scripts\run.bat
```

The server will start on `localhost:50051` and display:
```
gRPC server starting on port 50051...
```

### 3. Run Client Example

In a new terminal:
```bash
go run client/example.go
```

### 4. Run Tests

In a new terminal:
```bash
chmod +x tests/test.sh
./tests/test.sh
```

## API Services

### UserService

- `CreateUser(UserCreate) → User` - Create a new user
- `GetUser(GetUserRequest) → User` - Get user by ID
- `UpdateUser(UpdateUserRequest) → User` - Partially update user
- `ReplaceUser(UpdateUserRequest) → User` - Replace user data
- `DeleteUser(DeleteUserRequest) → Empty` - Delete user

### SessionService

- `Login(UserLogin) → LoginResponse` - Authenticate and get JWT token
- `Logout(Empty) → Empty` - Logout (invalidate session)

### ListingService

- `GetListings(ListingsRequest) → ListingsResponse` - Search listings with filters
- `CreateListing(ListingCreate) → Listing` - Create new listing
- `GetListing(GetListingRequest) → Listing` - Get listing by ID
- `UpdateListing(UpdateListingRequest) → Listing` - Update listing
- `DeleteListing(DeleteListingRequest) → Success` - Delete listing

### OrderService

- `GetOrders(OrdersRequest) → OrdersResponse` - Get orders with pagination
- `CreateOrder(OrderCreate) → Order` - Create new order
- `GetOrder(GetOrderRequest) → Order` - Get order by ID
- `UpdateOrder(UpdateOrderRequest) → Order` - Update order
- `DeleteOrder(DeleteOrderRequest) → Success` - Delete order
- `CancelOrder(CancelOrderRequest) → CancelOrderResponse` - Cancel order
- `UpdateOrderStatus(UpdateOrderStatusRequest) → Order` - Update order status

## Manual Testing with grpcurl

Once the server is running, you can test individual endpoints:

### User Operations
```bash
# Create user
grpcurl -plaintext -d '{"username":"john","email":"john@example.com","password":"secret"}' \
localhost:50051 ebayclone.UserService/CreateUser

# Get user
grpcurl -plaintext -d '{"id":1}' \
localhost:50051 ebayclone.UserService/GetUser

# Update user
grpcurl -plaintext -d '{"id":1,"user":{"username":"johnupdated"}}' \
localhost:50051 ebayclone.UserService/UpdateUser
```

### Authentication
```bash
# Login
grpcurl -plaintext -d '{"email":"john@example.com","password":"secret"}' \
localhost:50051 ebayclone.SessionService/Login

# Logout
grpcurl -plaintext -d '{}' \
localhost:50051 ebayclone.SessionService/Logout
```

### Listing Operations
```bash
# Create listing
grpcurl -plaintext -d '{"title":"iPhone 13","description":"Great phone","price":999.99,"category":"electronics","condition":"new"}' \
localhost:50051 ebayclone.ListingService/CreateListing

# Search listings
grpcurl -plaintext -d '{"search":"iPhone","priceMin":500,"priceMax":1500}' \
localhost:50051 ebayclone.ListingService/GetListings

# Get specific listing
grpcurl -plaintext -d '{"id":1}' \
localhost:50051 ebayclone.ListingService/GetListing
```

### Order Operations
```bash
# Create order
grpcurl -plaintext -d '{"listingId":1,"quantity":1,"shippingAddress":{"street":"123 Main St","city":"New York","country":"USA"}}' \
localhost:50051 ebayclone.OrderService/CreateOrder

# Get orders with pagination
grpcurl -plaintext -d '{"page":1,"limit":10}' \
localhost:50051 ebayclone.OrderService/GetOrders

# Update order status
grpcurl -plaintext -d '{"id":1,"status":"shipped"}' \
localhost:50051 ebayclone.OrderService/UpdateOrderStatus

# Cancel order
grpcurl -plaintext -d '{"id":1,"cancelReason":"Changed my mind"}' \
localhost:50051 ebayclone.OrderService/CancelOrder
```

## Development

### Building from Source

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Generate protobuf code:**
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          proto/ebayclone.proto
   ```

3. **Build server:**
   ```bash
   go build -o bin/server src/main.go
   ```

4. **Run server:**
   ```bash
   ./bin/server
   ```

### Testing

The test suite validates that all gRPC endpoints work correctly and return appropriate responses:

```bash
# Run all tests
./tests/test.sh

# Expected output:
# === eBayClone gRPC Functional Tests ===
# ✓ PASSED: Create User
# ✓ PASSED: User Login
# ✓ PASSED: Get User
# ... (more tests)
# === Test Results ===
# Tests Passed: 12
# Tests Failed: 0
# All tests passed! ✓
```

## Error Handling

The gRPC implementation uses standard gRPC status codes:

- `INVALID_ARGUMENT` (400) - Invalid input data
- `UNAUTHENTICATED` (401) - Authentication required
- `NOT_FOUND` (404) - Resource not found
- `ALREADY_EXISTS` (409) - Resource already exists
- `INTERNAL` (500) - Server error

## REST to gRPC Mapping

| REST Endpoint | gRPC Method | Notes |
|---------------|-------------|-------|
| `POST /users` | `UserService.CreateUser` | Creates new user |
| `GET /users/{id}` | `UserService.GetUser` | Gets user by ID |
| `PUT /users/{id}` | `UserService.ReplaceUser` | Replaces user data |
| `PATCH /users/{id}` | `UserService.UpdateUser` | Partial update |
| `DELETE /users/{id}` | `UserService.DeleteUser` | Deletes user |
| `POST /sessions` | `SessionService.Login` | User authentication |
| `DELETE /sessions` | `SessionService.Logout` | User logout |
| `GET /listings` | `ListingService.GetListings` | Search with filters |
| `POST /listings` | `ListingService.CreateListing` | Create listing |
| `GET /listings/{id}` | `ListingService.GetListing` | Get by ID |
| `PATCH /listings/{id}` | `ListingService.UpdateListing` | Update listing |
| `DELETE /listings/{id}` | `ListingService.DeleteListing` | Delete listing |
| `GET /orders` | `OrderService.GetOrders` | Get with pagination |
| `POST /orders` | `OrderService.CreateOrder` | Create order |
| `GET /orders/{id}` | `OrderService.GetOrder` | Get by ID |
| `PATCH /orders/{id}` | `OrderService.UpdateOrder` | Update order |
| `DELETE /orders/{id}` | `OrderService.DeleteOrder` | Delete order |
| `PATCH /orders/{id}/cancel` | `OrderService.CancelOrder` | Cancel order |
| `PATCH /orders/{id}/status` | `OrderService.UpdateOrderStatus` | Update status |

## Troubleshooting

### Common Issues

1. **"protoc: command not found"**
   - Install Protocol Buffers compiler: https://grpc.io/docs/protoc-installation/

2. **"go: command not found"**
   - Install Go: https://golang.org/doc/install

3. **"Failed to listen: address already in use"**
   - Port 50051 is already in use. Kill existing processes or change port in `src/main.go`

4. **"connection refused"**
   - Make sure the server is running before running tests or client examples

### Logs and Debugging

- Server logs are printed to stdout
- Enable gRPC reflection for debugging: `grpcurl -plaintext localhost:50051 list`
- Use `grpcurl -plaintext localhost:50051 describe ebayclone.UserService` to inspect service definitions

## Docker Support

### Using Docker

1. **Build and run with Docker:**
   ```bash
   docker build -t ebayclone-grpc .
   docker run -p 50051:50051 ebayclone-grpc
   ```

2. **Using Docker Compose:**
   ```bash
   docker-compose up --build
   ```

### Language-Agnostic Clients

The gRPC service can be used from any language that supports gRPC:

#### Python Client
```bash
# Install dependencies
pip install -r client/requirements.txt

# Generate Python protobuf files
./scripts/generate_python_proto.sh

# Run Python client
python3 client/example.py
```

#### Go Client
```bash
go run client/example.go
```

## Testing

### Unit Tests
```bash
go test ./src/services/...
```

### Integration Tests
```bash
./tests/test.sh
```

### Manual Testing
Use grpcurl for manual testing:
```bash
# List available services
grpcurl -plaintext localhost:50051 list

# Describe a service
grpcurl -plaintext localhost:50051 describe ebayclone.UserService
```

## Performance Notes

- In-memory storage is used for simplicity
- For production, replace with persistent storage (PostgreSQL, MongoDB, etc.)
- Add connection pooling and caching for better performance
- Implement proper JWT token validation and refresh

## Security Considerations

- JWT secret should be stored in environment variables
- Implement proper password hashing (bcrypt recommended)
- Add rate limiting and input validation
- Use TLS in production environments

## License

This project is provided as-is for educational and demonstration purposes.