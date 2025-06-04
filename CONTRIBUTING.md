# Contributing to Crossplane AI

Thank you for your interest in contributing to Crossplane AI! This document provides guidelines for contributing to the project.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Issues

Before creating an issue, please check if it already exists. When creating an issue:

1. Use a clear and descriptive title
2. Describe the exact steps to reproduce the problem
3. Provide specific examples and expected vs actual behavior
4. Include environment details (OS, Go version, Kubernetes version)

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

1. A clear description of the proposed feature
2. Explanation of why this enhancement would be useful
3. Examples of how the feature would work

### Pull Requests

1. Fork the repository and create your branch from `main`
2. If you've added code that should be tested, add tests
3. Ensure the test suite passes
4. Make sure your code follows the existing style
5. Issue the pull request

## Development Setup

### Prerequisites

- Go 1.24 or higher
- Access to a Kubernetes cluster with Crossplane (for testing)
- Make (optional, for convenient commands)

### Local Development

1. Clone your fork:
   ```bash
   git clone https://github.com/your-username/crossplane-ai.git
   cd crossplane-ai
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build -o crossplane-ai .
   ```

4. Run tests:
   ```bash
   go test ./...
   ```

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and small
- Handle errors appropriately

### Testing

- Write unit tests for new functionality
- Include integration tests where appropriate
- Ensure all tests pass before submitting PR
- Aim for good test coverage

### Documentation

- Update README.md if you change functionality
- Add inline code comments for complex logic
- Update help text for CLI commands if needed

## Project Structure

```
crossplane-ai/
├── cmd/                    # CLI commands
├── pkg/                    # Core packages
│   ├── ai/                # AI service integration
│   ├── crossplane/        # Crossplane client
│   └── cli/               # CLI utilities
├── internal/              # Internal packages
└── main.go               # Application entry point
```

## Commit Messages

Use clear and descriptive commit messages:

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

Example:
```
Add OpenAI integration for real AI responses

- Implement OpenAIClient with completion API
- Add configuration support for API keys
- Update service to use real AI when configured
- Fallback to mock responses when API key not provided

Fixes #123
```

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release PR
4. Tag release after merge
5. Build and publish binaries

## Getting Help

- Check existing issues and discussions
- Join our community discussions
- Reach out to maintainers

Thank you for contributing to Crossplane AI! 🚀
