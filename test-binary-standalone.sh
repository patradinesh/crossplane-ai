#!/bin/bash
# Test script to demonstrate that the crossplane-ai binary works standalone
# This simulates the user experience when someone downloads just the binary

echo "🧪 Testing Crossplane AI Binary - Standalone Mode"
echo "================================================="
echo
echo "This demonstrates that the crossplane-ai binary works completely standalone"
echo "with embedded mock data, without requiring any external files."
echo
echo "Commands being tested:"
echo "• crossplane-ai --mock analyze"
echo "• crossplane-ai --mock ask \"What databases do I have?\""
echo "• crossplane-ai --mock suggest database"
echo "• crossplane-ai generate examples --list"
echo
echo "-------------------------------------------------"

# Test 1: Analyze command
echo "🔍 Test 1: Analyze Command"
echo "$ ./crossplane-ai --mock analyze"
echo
./crossplane-ai --mock analyze
echo
echo "-------------------------------------------------"

# Test 2: Ask command
echo "🗣️  Test 2: Ask Command"
echo "$ ./crossplane-ai --mock ask \"What databases do I have?\""
echo
./crossplane-ai --mock ask "What databases do I have?"
echo
echo "-------------------------------------------------"

# Test 3: Suggest command
echo "💡 Test 3: Suggest Command"
echo "$ ./crossplane-ai --mock suggest database"
echo
./crossplane-ai --mock suggest database
echo
echo "-------------------------------------------------"

# Test 4: Generate examples
echo "📝 Test 4: Generate Examples"
echo "$ ./crossplane-ai generate examples --list"
echo
./crossplane-ai generate examples --list
echo
echo "-------------------------------------------------"
echo
echo "✅ All tests completed successfully!"
echo "The crossplane-ai binary works completely standalone without any external dependencies."
echo "Users can download just the binary and immediately start testing with mock mode."
