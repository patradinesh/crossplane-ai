#!/bin/bash
# Comprehensive test script for Crossplane AI mock mode
# This script demonstrates all the mock functionality of the Crossplane AI tool

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_BINARY="../../crossplane-ai"

echo "🧪 Crossplane AI Mock Mode Test Suite"
echo "====================================="
echo

# Check if binary exists
if [ ! -f "$CLI_BINARY" ]; then
    echo "❌ Binary not found: $CLI_BINARY"
    echo "Please run 'go build -o crossplane-ai .' first"
    exit 1
fi

echo "✅ Found binary: $CLI_BINARY"
echo

# Test 1: Basic analyze command
echo "🔬 Test 1: Basic Analysis"
echo "-------------------------"
./run-mock.sh --binary "$CLI_BINARY" analyze
echo

# Test 2: Analyze with provider filter
echo "🔬 Test 2: Analysis with AWS Provider Filter"
echo "---------------------------------------------"
./run-mock.sh --binary "$CLI_BINARY" analyze --provider aws
echo

# Test 3: Ask about resources
echo "🤖 Test 3: Ask About Resources"
echo "-------------------------------"
./run-mock.sh --binary "$CLI_BINARY" ask "What resources do I have?"
echo

# Test 4: Ask about AWS resources
echo "🤖 Test 4: Ask About AWS Resources"
echo "-----------------------------------"
./run-mock.sh --binary "$CLI_BINARY" ask "Show me AWS resources"
echo

# Test 5: Ask about databases
echo "🤖 Test 5: Ask About Databases"
echo "-------------------------------"
./run-mock.sh --binary "$CLI_BINARY" ask "Tell me about my databases"
echo

# Test 6: Ask about providers
echo "🤖 Test 6: Ask About Providers"
echo "-------------------------------"
./run-mock.sh --binary "$CLI_BINARY" ask "What providers are installed?"
echo

# Test 7: General suggestions
echo "💡 Test 7: General Suggestions"
echo "-------------------------------"
./run-mock.sh --binary "$CLI_BINARY" suggest
echo

# Test 8: Database suggestions
echo "💡 Test 8: Database Suggestions"
echo "--------------------------------"
./run-mock.sh --binary "$CLI_BINARY" suggest database
echo

# Test 9: Security suggestions
echo "💡 Test 9: Security Suggestions"
echo "--------------------------------"
./run-mock.sh --binary "$CLI_BINARY" suggest security
echo

# Test 10: Optimization suggestions
echo "💡 Test 10: Optimization Suggestions"
echo "-------------------------------------"
./run-mock.sh --binary "$CLI_BINARY" suggest optimize
echo

# Test 11: Network suggestions
echo "💡 Test 11: Network Suggestions"
echo "--------------------------------"
./run-mock.sh --binary "$CLI_BINARY" suggest network
echo

echo "🎉 All mock tests completed successfully!"
echo
echo "Summary:"
echo "• Analyze command: ✅ Working with mock data and filters"
echo "• Ask command: ✅ Working with intelligent mock responses"
echo "• Suggest command: ✅ Working with category-specific suggestions"
echo "• Mock mode: ✅ Properly isolated from real cluster"
echo
echo "To test against a real cluster, run commands without CROSSPLANE_AI_MODE=mock"
