<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Crossplane AI CLI Project

This is a Go-based command-line tool that brings AI capabilities to Crossplane resource management in Kubernetes clusters.

## Project Guidelines

### Architecture
- This is a CLI application built with Cobra framework
- Uses Kubernetes client-go for interacting with Crossplane CRDs
- Modular design with separate packages for AI, Crossplane client, and CLI commands
- Follows Go best practices and idiomatic patterns

### Key Components
- `cmd/` - Contains all CLI commands (ask, suggest, analyze)
- `pkg/ai/` - AI service for natural language processing and intelligent suggestions  
- `pkg/crossplane/` - Kubernetes client wrapper for Crossplane resources
- `pkg/cli/` - Common CLI utilities and helpers
- `internal/config/` - Configuration management

### Crossplane Integration
- Works with Crossplane CRDs and custom resources
- Supports multiple cloud providers (AWS, GCP, Azure)
- Handles dynamic resource discovery and analysis
- Integrates with Kubernetes RBAC and authentication

### AI Features
- Natural language query processing for resource discovery
- Intelligent resource analysis and health checking
- AI-powered suggestions for optimization and troubleshooting
- Interactive chat-like interface for resource management

### Code Style
- Use descriptive variable and function names
- Include comprehensive error handling
- Add helpful comments for complex logic
- Follow Go naming conventions
- Use structured logging where appropriate

### Dependencies
- github.com/spf13/cobra - CLI framework
- github.com/spf13/viper - Configuration management
- k8s.io/client-go - Kubernetes client library
- k8s.io/apimachinery - Kubernetes API machinery

When implementing new features:
1. Consider how they integrate with existing Crossplane workflows
2. Ensure proper error handling for Kubernetes API calls
3. Make AI responses helpful and actionable
4. Support both interactive and non-interactive modes
5. Include examples in help text and documentation
