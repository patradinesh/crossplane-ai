#!/bin/bash
# Test script for Crossplane AI without a full Crossplane installation

# Set environment variables for testing
export CROSSPLANE_AI_MODE=mock
export CROSSPLANE_AI_MOCK_DATA_DIR="./examples"

# Run the Crossplane AI tool with the specified command
echo "Running Crossplane AI in mock mode..."
echo "Command: $@"
echo "---------------------------------"

# Execute the Crossplane AI tool
./crossplane-ai "$@"
