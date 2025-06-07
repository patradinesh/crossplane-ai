#!/bin/bash
# Setup script for Crossplane AI mock testing

# Ensure we're in the project root directory
cd "$(dirname "$0")/../.." || exit 1

echo "Setting up Crossplane AI mock testing environment..."

# Install required Go packages for mock testing
echo "Installing dependencies..."
go get gopkg.in/yaml.v2

# Build the Crossplane AI binary if it doesn't exist
if [ ! -f "./crossplane-ai" ]; then
  echo "Building Crossplane AI binary..."
  go build -o crossplane-ai .
fi

# Create mock data directory if it doesn't exist
if [ ! -d "./examples" ]; then
  echo "Creating examples directory..."
  mkdir -p ./examples
fi

# Copy example YAML files to the examples directory if they don't exist
if [ ! -f "./examples/xdatabase-definition.yaml" ] && [ -f "./test/mock/examples/xdatabase-definition.yaml" ]; then
  echo "Copying example YAML files..."
  cp ./test/mock/examples/*.yaml ./examples/
fi

echo "Setup complete!"
echo
echo "To run Crossplane AI with mock mode:"
echo "  ./test/mock/run-mock.sh ask \"what resources do I have?\""
echo "  ./test/mock/run-mock.sh suggest database"
echo "  ./test/mock/run-mock.sh --help"
echo
