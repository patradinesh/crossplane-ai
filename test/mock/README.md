# Mock Testing for Crossplane AI

This directory contains tools and resources for testing the Crossplane AI CLI tool without requiring a full Crossplane installation or Kubernetes cluster.

## Overview

The mock testing framework provides:

1. **Mock AI Service** - A simulated AI service that returns predefined responses based on the query
2. **Mock Crossplane Client** - A simulated Crossplane client that loads resources from YAML files
3. **Test Script** - A convenient script for running Crossplane AI with mock mode enabled

## Usage

### Basic Usage

```bash
# From the root directory
./test/mock/run-mock.sh ask "what resources do I have?"
./test/mock/run-mock.sh suggest database
./test/mock/run-mock.sh analyze
```

### Advanced Options

```bash
# Use custom mock data directory
./test/mock/run-mock.sh --mock-dir ./path/to/mocks ask "what resources do I have?"

# Enable verbose output
./test/mock/run-mock.sh --verbose generate "postgresql database with 20GB storage"

# Show help
./test/mock/run-mock.sh --help
```

## Mock Data

The mock system looks for YAML files in the specified mock data directory (defaults to `./examples`). These files should contain valid Crossplane resource definitions that will be used by the mock system.

Example mock data files:

- `xdatabase-definition.yaml` - CompositeResourceDefinition for databases
- `xdatabase-composition.yaml` - Composition for databases  
- `database-claim.yaml` - Database claim

## Creating Custom Mocks

You can create your own mock data by adding YAML files to the mock data directory. The mock system will automatically discover and use these files.

## Implementation Details

The mock implementation consists of:

- `ai_service.go` - Implements a mock AI service for testing
- `crossplane_client.go` - Implements a mock Crossplane client
- `run-mock.sh` - Script for running the Crossplane AI tool in mock mode
