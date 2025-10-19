#!/bin/bash

echo "🚀 GoScraper Feature Validation"
echo "================================"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

success_count=0
total_tests=0

check_test() {
    local test_name="$1"
    local command="$2"
    
    echo -n "Testing $test_name... "
    total_tests=$((total_tests + 1))
    
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
        success_count=$((success_count + 1))
    else
        echo -e "${RED}✗${NC}"
    fi
}

echo -e "\n${YELLOW}1. CLI Tools${NC}"
check_test "Config initialization" "go run ./cmd/cli init"
check_test "Config validation" "go run ./cmd/cli validate"
check_test "Config display" "go run ./cmd/cli config"

echo -e "\n${YELLOW}2. Configuration System${NC}"
check_test "Environment variables" "GOSCRAPER_AI_ENABLED=true go run ./cmd/cli config | grep 'AI Enabled: true'"
check_test "Config file creation" "test -f ~/.goscraper/config.json"

echo -e "\n${YELLOW}3. Build System${NC}"
check_test "API server build" "go build -o /tmp/test-api ./cmd/api"
check_test "CLI tool build" "go build -o /tmp/test-cli ./cmd/cli"

echo -e "\n${YELLOW}4. Dependencies${NC}"
check_test "Go mod tidy" "go mod tidy"
check_test "Go mod verify" "go mod verify"

echo -e "\n${YELLOW}5. Tests${NC}"
check_test "Integration tests" "go test ./tests/ -v"

echo -e "\n${YELLOW}6. API Server (requires manual start)${NC}"
echo "To test API server:"
echo "1. Start: go run ./cmd/api"
echo "2. Test: curl http://localhost:8080/health"
echo "3. Test: curl http://localhost:8080/config"

# Cleanup
rm -f /tmp/test-api /tmp/test-cli

echo -e "\n${YELLOW}Summary${NC}"
echo "================================"
echo -e "Tests passed: ${GREEN}$success_count${NC}/$total_tests"

if [ $success_count -eq $total_tests ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi