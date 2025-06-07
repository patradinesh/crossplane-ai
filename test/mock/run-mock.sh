#!/bin/bash
# Mock test script for Crossplane AI CLI tool
# This script allows testing the Crossplane AI tool without a full Crossplane installation

# Default values
MOCK_DATA_DIR="../examples"
CLI_BINARY="../crossplane-ai"
VERBOSE=false

# Parse command line options
while [[ $# -gt 0 ]]; do
  case $1 in
    --mock-dir)
      MOCK_DATA_DIR="$2"
      shift 2
      ;;
    --binary)
      CLI_BINARY="$2"
      shift 2
      ;;
    --verbose)
      VERBOSE=true
      shift
      ;;
    --help)
      echo "Usage: $0 [OPTIONS] COMMAND [ARGS...]"
      echo
      echo "Run Crossplane AI CLI in mock mode for testing"
      echo
      echo "Options:"
      echo "  --mock-dir DIR    Directory containing mock data (default: ../examples)"
      echo "  --binary PATH     Path to crossplane-ai binary (default: ../crossplane-ai)"
      echo "  --verbose         Enable verbose output"
      echo "  --help            Show this help message"
      echo
      echo "Examples:"
      echo "  $0 ask \"what resources do I have?\""
      echo "  $0 suggest database"
      echo "  $0 --mock-dir ./custom-mocks analyze"
      exit 0
      ;;
    *)
      break
      ;;
  esac
done

# Set environment variables for testing
export CROSSPLANE_AI_MODE=mock
export CROSSPLANE_AI_MOCK_DATA_DIR="$MOCK_DATA_DIR"

# Run the Crossplane AI tool with the specified command
echo "Running Crossplane AI in mock mode..."
echo "Mock data directory: $MOCK_DATA_DIR"
echo "Command: $*"
echo "---------------------------------"

# Verbose output if requested
if [ "$VERBOSE" = true ]; then
  echo "Environment variables:"
  echo "  CROSSPLANE_AI_MODE=$CROSSPLANE_AI_MODE"
  echo "  CROSSPLANE_AI_MOCK_DATA_DIR=$CROSSPLANE_AI_MOCK_DATA_DIR"
  echo
  echo "Using binary: $CLI_BINARY"
  echo
fi

# Execute the Crossplane AI tool
"$CLI_BINARY" "$@"

# Capture exit code
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
  echo
  echo "Command failed with exit code: $EXIT_CODE"
fi

exit $EXIT_CODE
