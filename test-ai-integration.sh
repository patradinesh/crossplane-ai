#!/bin/bash

# Test script for Crossplane AI integration
# Tests mock mode, template mode, and AI mode detection

set -e

echo "🧪 Testing Crossplane AI Integration"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Build the tool
echo -e "${BLUE}Building crossplane-ai...${NC}"
go build -o crossplane-ai . || { echo -e "${RED}❌ Build failed${NC}"; exit 1; }
echo -e "${GREEN}✅ Build successful${NC}"

# Test 1: Mock Mode
echo -e "\n${BLUE}Test 1: Mock Mode${NC}"
echo "Running: ./crossplane-ai --mock ask 'test'"
OUTPUT=$(./crossplane-ai --mock ask "test" 2>&1)
if echo "$OUTPUT" | grep -q "MOCK MODE"; then
    echo -e "${GREEN}✅ Mock mode working correctly${NC}"
else
    echo -e "${RED}❌ Mock mode test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 2: Template Mode (no API key)
echo -e "\n${BLUE}Test 2: Template Mode (no API key)${NC}"
echo "Running: ./crossplane-ai ask 'test'"
OUTPUT=$(./crossplane-ai ask "test" 2>&1)
if echo "$OUTPUT" | grep -q "TEMPLATE MODE"; then
    echo -e "${GREEN}✅ Template mode working correctly${NC}"
else
    echo -e "${RED}❌ Template mode test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 3: AI Mode Detection (with fake API key)
echo -e "\n${BLUE}Test 3: AI Mode Detection (with test API key)${NC}"
echo "Running: OPENAI_API_KEY=test-key ./crossplane-ai --config ./config.yaml ask 'test'"
OUTPUT=$(OPENAI_API_KEY=test-key ./crossplane-ai --config ./config.yaml ask "test" 2>&1)
if echo "$OUTPUT" | grep -q "POWERED BY OPENAI"; then
    echo -e "${GREEN}✅ AI mode detection working correctly${NC}"
    if echo "$OUTPUT" | grep -q "401"; then
        echo -e "${GREEN}✅ API call attempted (expected 401 with fake key)${NC}"
    fi
else
    echo -e "${RED}❌ AI mode detection test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 4: Mock Analyze Command
echo -e "\n${BLUE}Test 4: Mock Analyze Command${NC}"
echo "Running: ./crossplane-ai --mock analyze"
OUTPUT=$(./crossplane-ai --mock analyze 2>&1)
if echo "$OUTPUT" | grep -q "embedded mock data"; then
    echo -e "${GREEN}✅ Mock analyze working correctly${NC}"
else
    echo -e "${RED}❌ Mock analyze test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 5: Template Generate Command
echo -e "\n${BLUE}Test 5: Template Generate Command${NC}"
echo "Running: ./crossplane-ai generate 'test database'"
OUTPUT=$(./crossplane-ai generate "test database" 2>&1)
if echo "$OUTPUT" | grep -q "template-based generation"; then
    echo -e "${GREEN}✅ Template generate working correctly${NC}"
else
    echo -e "${RED}❌ Template generate test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 6: AI Generate Detection (with fake API key)
echo -e "\n${BLUE}Test 6: AI Generate Detection${NC}"
echo "Running: OPENAI_API_KEY=test-key ./crossplane-ai --config ./config.yaml generate 'test database'"
OUTPUT=$(OPENAI_API_KEY=test-key ./crossplane-ai --config ./config.yaml generate "test database" 2>&1)
if echo "$OUTPUT" | grep -q "Using OpenAI for intelligent"; then
    echo -e "${GREEN}✅ AI generate detection working correctly${NC}"
else
    echo -e "${RED}❌ AI generate detection test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 7: Generate Examples Command
echo -e "\n${BLUE}Test 7: Generate Examples Command${NC}"
echo "Running: ./crossplane-ai generate examples"
OUTPUT=$(./crossplane-ai generate examples 2>&1)
if echo "$OUTPUT" | grep -q "Generated example files"; then
    echo -e "${GREEN}✅ Generate examples working correctly${NC}"
else
    echo -e "${RED}❌ Generate examples test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Test 8: Help Command
echo -e "\n${BLUE}Test 8: Help Command${NC}"
echo "Running: ./crossplane-ai --help"
OUTPUT=$(./crossplane-ai --help 2>&1)
if echo "$OUTPUT" | grep -q "Crossplane AI is a command-line tool"; then
    echo -e "${GREEN}✅ Help command working correctly${NC}"
else
    echo -e "${RED}❌ Help command test failed${NC}"
    echo "$OUTPUT"
    exit 1
fi

# Summary
echo -e "\n${GREEN}🎉 All tests passed!${NC}"
echo -e "\n${YELLOW}Summary of modes tested:${NC}"
echo -e "✅ Mock Mode: Uses embedded sample data"
echo -e "✅ Template Mode: Smart templates without AI"
echo -e "✅ AI Mode Detection: Properly detects OpenAI API key"
echo -e "✅ All commands working correctly"

echo -e "\n${BLUE}To test with real OpenAI:${NC}"
echo -e "export OPENAI_API_KEY=your-real-api-key"
echo -e "./crossplane-ai --config ./config.yaml ask 'What can you help me with?'"

echo -e "\n${GREEN}Integration testing complete! 🚀${NC}"
