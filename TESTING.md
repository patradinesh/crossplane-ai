# Testing Crossplane AI

This document provides instructions for testing the Crossplane AI CLI tool in different environments.

## Testing with a Mock Environment

If you don't have a Crossplane installation available, you can use the mock mode to test the functionality of the tool:

```bash
# Run with the mock testing script
./test/mock/run-mock.sh ask "what resources do I have?"
./test/mock/run-mock.sh suggest database
./test/mock/run-mock.sh generate "postgresql database with 20GB storage"
./test/mock/run-mock.sh interactive

# For more options
./test/mock/run-mock.sh --help

# Or set environment variables directly
export CROSSPLANE_AI_MODE=mock
export CROSSPLANE_AI_MOCK_DATA_DIR="./examples"
./crossplane-ai ask "what resources do I have?"
```

See the [mock testing README](./test/mock/README.md) for more details on the mock testing framework.

## Testing with a Real Crossplane Installation

When you have a Crossplane installation in a Kubernetes cluster:

1. Make sure your `kubeconfig` is properly configured to connect to your cluster
2. Run the Crossplane AI commands directly:

```bash
./crossplane-ai ask "what resources do I have?"
./crossplane-ai suggest database
./crossplane-ai analyze
./crossplane-ai generate "postgresql database with 20GB storage"
```

## Building the Tool

```bash
go build -o crossplane-ai .
```

## Configuration

The tool can be configured using:

1. Environment variables
2. Command-line flags
3. Configuration file (default: `$HOME/.crossplane-ai.yaml`)

Example configuration:

```yaml
# AI Service Configuration
ai:
  provider: "openai"  # or "mock" for testing
  api_key: "${OPENAI_API_KEY}"
  model: "gpt-4"

# Kubernetes Configuration  
kubernetes:
  kubeconfig: "/path/to/kubeconfig"
  context: "my-context"
  namespace: "default"

# Crossplane Configuration
crossplane:
  providers:
    - aws
    - gcp
    - azure
```

## Usage in CI/CD Pipelines

For using Crossplane AI in CI/CD pipelines:

```bash
export CROSSPLANE_AI_MODE=non-interactive
./crossplane-ai analyze --output json > analysis.json
./crossplane-ai suggest optimize --output yaml > suggestions.yaml
```
