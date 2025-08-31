#!/bin/bash

set -e

echo "=== eBayClone gRPC Functional Tests ==="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name="$1"
    local test_command="$2"
    
    echo -e "${YELLOW}Running: $test_name${NC}"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ PASSED: $test_name${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ FAILED: $test_name${NC}"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# Check if server is running
check_server() {
    echo "Checking if gRPC server is running on localhost:50051..."
    if ! nc -z localhost 50051 2>/dev/null; then
        echo "Error: gRPC server is not running on port 50051"
        echo "Please start the server first with: ./scripts/run.sh"
        exit 1
    fi
    echo "✓ gRPC server is running"
    echo ""
}

# Check if grpcurl is available for testing
check_grpcurl() {
    if ! command -v grpcurl &> /dev/null; then
        echo "Installing grpcurl for testing..."
        go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
    fi
}

echo "=== Pre-test Setup ==="
check_server
check_grpcurl

echo "=== Running gRPC Service Tests ==="

# Test 1: Create User
run_test "Create User" '
grpcurl -plaintext -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}" \
localhost:50051 ebayclone.UserService/CreateUser | grep -q "testuser"
'

# Test 2: Login
run_test "User Login" '
grpcurl -plaintext -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}" \
localhost:50051 ebayclone.SessionService/Login | grep -q "token"
'

# Test 3: Get User
run_test "Get User" '
grpcurl -plaintext -d "{\"id\":1}" \
localhost:50051 ebayclone.UserService/GetUser | grep -q "testuser"
'

# Test 4: Create Listing
run_test "Create Listing" '
grpcurl -plaintext -d "{\"title\":\"iPhone 13\",\"description\":\"Great phone\",\"price\":999.99,\"category\":\"electronics\",\"condition\":\"new\"}" \
localhost:50051 ebayclone.ListingService/CreateListing | grep -q "iPhone 13"
'

# Test 5: Get Listings
run_test "Get Listings" '
grpcurl -plaintext -d "{\"search\":\"iPhone\"}" \
localhost:50051 ebayclone.ListingService/GetListings | grep -q "iPhone 13"
'

# Test 6: Get Listing by ID
run_test "Get Listing by ID" '
grpcurl -plaintext -d "{\"id\":1}" \
localhost:50051 ebayclone.ListingService/GetListing | grep -q "iPhone 13"
'

# Test 7: Create Order
run_test "Create Order" '
grpcurl -plaintext -d "{\"listingId\":1,\"quantity\":1,\"shippingAddress\":{\"street\":\"123 Main St\",\"city\":\"New York\",\"country\":\"USA\"}}" \
localhost:50051 ebayclone.OrderService/CreateOrder | grep -q "pending"
'

# Test 8: Get Orders
run_test "Get Orders" '
grpcurl -plaintext -d "{\"page\":1,\"limit\":10}" \
localhost:50051 ebayclone.OrderService/GetOrders | grep -q "orders"
'

# Test 9: Update Order Status
run_test "Update Order Status" '
grpcurl -plaintext -d "{\"id\":1,\"status\":\"confirmed\"}" \
localhost:50051 ebayclone.OrderService/UpdateOrderStatus | grep -q "confirmed"
'

# Test 10: Cancel Order
run_test "Cancel Order" '
grpcurl -plaintext -d "{\"id\":1,\"cancelReason\":\"Changed my mind\"}" \
localhost:50051 ebayclone.OrderService/CancelOrder | grep -q "cancelled"
'

# Test 11: Error Handling - Invalid User ID
run_test "Error Handling - Invalid User ID" '
grpcurl -plaintext -d "{\"id\":999}" \
localhost:50051 ebayclone.UserService/GetUser 2>&1 | grep -q "NotFound"
'

# Test 12: Error Handling - Invalid Input
run_test "Error Handling - Invalid Input" '
grpcurl -plaintext -d "{\"username\":\"\",\"email\":\"\",\"password\":\"\"}" \
localhost:50051 ebayclone.UserService/CreateUser 2>&1 | grep -q "InvalidArgument"
'

echo "=== Test Results ==="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo -e "Total Tests: $((TESTS_PASSED + TESTS_FAILED))"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi
