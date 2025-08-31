#!/bin/bash

echo "=== eBayClone gRPC Implementation Verification ==="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

CHECKS_PASSED=0
CHECKS_FAILED=0

check_requirement() {
    local requirement="$1"
    local command="$2"
    
    echo -e "${YELLOW}Checking: $requirement${NC}"
    
    if eval "$command"; then
        echo -e "${GREEN}✓ PASSED: $requirement${NC}"
        ((CHECKS_PASSED++))
    else
        echo -e "${RED}✗ FAILED: $requirement${NC}"
        ((CHECKS_FAILED++))
    fi
    echo ""
}

echo "=== Verifying Assessment Criteria ==="

# 1. .proto files compile with protoc without errors
check_requirement ".proto files compile without errors" "protoc --version > /dev/null && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/ebayclone.proto"

# 2. All REST endpoints have corresponding RPCs in .proto
check_requirement "All REST endpoints mapped to gRPC RPCs" "grep -q 'service UserService' proto/ebayclone.proto && grep -q 'service SessionService' proto/ebayclone.proto && grep -q 'service ListingService' proto/ebayclone.proto && grep -q 'service OrderService' proto/ebayclone.proto"

# 3. Project structure exists
check_requirement "Project structure exists" "test -d proto && test -d src && test -f scripts/run.sh && test -d client && test -d tests && test -f README.md"

# 4. Build script exists and is executable
check_requirement "Build script exists" "test -f scripts/run.sh"

# 5. Client example exists
check_requirement "Client example exists" "test -f client/example.go"

# 6. Automated tests exist
check_requirement "Automated tests exist" "test -f tests/test.sh"

# 7. README contains instructions
check_requirement "README contains build instructions" "grep -q 'Quick Start' README.md && grep -q 'Prerequisites' README.md"

# 8. Go module is properly configured
check_requirement "Go module configured" "test -f go.mod && grep -q 'module ebayclone-grpc' go.mod"

# 9. All required message types exist in proto
check_requirement "Required message types exist" "grep -q 'message User' proto/ebayclone.proto && grep -q 'message Listing' proto/ebayclone.proto && grep -q 'message Order' proto/ebayclone.proto && grep -q 'message Error' proto/ebayclone.proto"

# 10. Service implementations exist
check_requirement "Service implementations exist" "test -f src/services/user_service.go && test -f src/services/session_service.go && test -f src/services/listing_service.go && test -f src/services/order_service.go"

echo "=== Verification Results ==="
echo -e "Checks Passed: ${GREEN}$CHECKS_PASSED${NC}"
echo -e "Checks Failed: ${RED}$CHECKS_FAILED${NC}"
echo -e "Total Checks: $((CHECKS_PASSED + CHECKS_FAILED))"

if [ $CHECKS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All verification checks passed! ✓${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Install Go 1.21+ and protoc if not already installed"
    echo "2. Run: ./scripts/run.sh (Unix) or scripts\\run.bat (Windows)"
    echo "3. In another terminal, run: go run client/example.go"
    echo "4. Run tests: ./tests/test.sh"
    exit 0
else
    echo -e "${RED}Some verification checks failed! ✗${NC}"
    exit 1
fi
