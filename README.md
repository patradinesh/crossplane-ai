# Crossplane AI

An AI-powered command-line tool that brings intelligent capabilities to Crossplane resource management. Similar to kubectl-ai, but specifically designed for Crossplane infrastructure as code workflows.

![Crossplane AI Demo](https://img.shields.io/badge/status-ready-green) ![Go Version](https://img.shields.io/badge/go-1.24+-blue) ![License](https://img.shields.io/badge/license-MIT-blue)

## 🚀 Features

- **Natural Language Queries**: Ask questions about your Crossplane resources in plain English
- **AI-Powered Analysis**: Intelligent analysis of resource health, configuration, and optimization opportunities  
- **Resource Generation**: Generate Crossplane manifests from natural language descriptions
- **Interactive Mode**: Chat-like interface for ongoing infrastructure management
- **Multi-Cloud Support**: Works with AWS, GCP, Azure, and other Crossplane providers
- **Intelligent Suggestions**: Get AI-powered recommendations for optimization, security, and best practices

## 📦 Installation

### From Source

```bash
git clone https://github.com/your-org/crossplane-ai.git
cd crossplane-ai
go build -o crossplane-ai .
sudo mv crossplane-ai /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/your-org/crossplane-ai@latest
```

## 🎯 Quick Start

### Basic Commands

```bash
# Ask about your resources
crossplane-ai ask "what AWS resources do I have?"

# Get intelligent suggestions
crossplane-ai suggest database

# Analyze resource health
crossplane-ai analyze

# Generate new resources  
crossplane-ai generate "create a MySQL database on AWS"

# Interactive mode
crossplane-ai interactive
```

## 🤖 Commands

### `ask` - Natural Language Queries

Ask questions about your Crossplane resources using natural language.

```bash
# General resource queries
crossplane-ai ask "what resources do I have?"
crossplane-ai ask "how many databases are running?"

# Provider-specific queries
crossplane-ai ask "show me all AWS resources"
crossplane-ai ask "what GCP resources are not ready?"

# Troubleshooting queries
crossplane-ai ask "why is my database failing?"
crossplane-ai ask "what resources have errors?"

# Cost and optimization queries
crossplane-ai ask "which resources are most expensive?"
crossplane-ai ask "how can I optimize costs?"
```

### `generate` - AI-Powered Resource Creation

Generate Crossplane manifests from natural language descriptions.

```bash
# Database resources
crossplane-ai generate "create a PostgreSQL database on AWS"
crossplane-ai generate "MySQL cluster with read replicas"

# Storage resources
crossplane-ai generate "S3 bucket with versioning and encryption"
crossplane-ai generate "GCS bucket for backup storage"

# Network resources
crossplane-ai generate "VPC with public and private subnets"
crossplane-ai generate "load balancer for web application"

# Complete stacks
crossplane-ai generate "web application with database and load balancer"

# Options
crossplane-ai generate "database" --provider aws --dry-run
crossplane-ai generate "storage" --apply
crossplane-ai generate "network" --output json
```

### `suggest` - Intelligent Recommendations

Get AI-powered suggestions for optimization, security, and best practices.

```bash
# General suggestions
crossplane-ai suggest

# Category-specific suggestions
crossplane-ai suggest database
crossplane-ai suggest security
crossplane-ai suggest cost-optimization
crossplane-ai suggest performance

# Provider-specific suggestions
crossplane-ai suggest --provider aws
crossplane-ai suggest --provider gcp
```

### `analyze` - Resource Analysis

Perform comprehensive analysis of your Crossplane resources.

```bash
# Basic analysis
crossplane-ai analyze

# Detailed analysis with health checks
crossplane-ai analyze --detailed

# Analysis with specific focus
crossplane-ai analyze --focus security
crossplane-ai analyze --focus cost
crossplane-ai analyze --focus performance

# Provider-specific analysis
crossplane-ai analyze --provider aws
```

### `interactive` - Chat Mode

Start an interactive session for ongoing resource management.

```bash
# Basic interactive mode
crossplane-ai interactive

# With banner and initial analysis
crossplane-ai interactive --banner --analyze

# Aliases
crossplane-ai i
crossplane-ai chat
```

In interactive mode, you can use commands like:
- `analyze` - Run detailed analysis
- `status` - Show resource overview
- `health` - Perform health check
- `suggest [type]` - Get suggestions
- `help` - Show available commands
- `exit` - Exit interactive mode

## ⚙️ Configuration

### Configuration File

Create `~/.crossplane-ai.yaml`:

```yaml
# AI Configuration
ai:
  provider: "openai"  # openai, anthropic, local
  model: "gpt-4"
  api_key: "${OPENAI_API_KEY}"
  
# Output Configuration
output:
  format: "yaml"  # yaml, json, table
  color: true
  verbose: false

# Kubernetes Configuration
kubernetes:
  config_path: "~/.kube/config"
  context: ""
  namespace: "default"

# Provider Preferences
providers:
  default: "aws"
  aws:
    region: "us-east-1"
  gcp:
    project: "my-project"
    region: "us-central1"
  azure:
    subscription: "my-subscription"
    location: "eastus"
```

### Environment Variables

```bash
export CROSSPLANE_AI_CONFIG=~/.crossplane-ai.yaml
export OPENAI_API_KEY=your-api-key
export KUBECONFIG=~/.kube/config
```

### Command Line Flags

```bash
crossplane-ai --config /path/to/config.yaml \
              --kubeconfig /path/to/kubeconfig \
              --context my-context \
              --verbose \
              [command]
```

## 🎨 Examples

### Example Workflow

```bash
# 1. Start with analysis
crossplane-ai analyze

# 2. Ask specific questions
crossplane-ai ask "which resources are not ready?"

# 3. Get recommendations
crossplane-ai suggest optimization

# 4. Generate new resources
crossplane-ai generate "backup storage for my database" --apply

# 5. Interactive session for ongoing management
crossplane-ai interactive
```

### Advanced Usage

```bash
# Complex resource generation
crossplane-ai generate "3-tier web application with:
- Application Load Balancer
- Auto-scaling ECS service
- RDS MySQL database with read replica
- ElastiCache Redis cluster
- All with proper security groups and encryption"

# Troubleshooting workflow
crossplane-ai ask "what's wrong with my infrastructure?"
crossplane-ai analyze --focus errors
crossplane-ai suggest troubleshooting

# Cost optimization workflow
crossplane-ai ask "how much am I spending on compute?"
crossplane-ai suggest cost-optimization
crossplane-ai analyze --focus cost
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Commands  │    │   AI Service    │    │ Crossplane API  │
│                 │    │                 │    │                 │
│ • ask           │───▶│ • OpenAI/Claude │───▶│ • Custom        │
│ • generate      │    │ • Local Models  │    │   Resources     │
│ • suggest       │    │ • Context Mgmt  │    │ • Compositions  │
│ • analyze       │    │                 │    │ • Providers     │
│ • interactive   │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        │                       │                       │
        ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Configuration   │    │ Resource Cache  │    │ Kubernetes API  │
│ Management      │    │ & Context       │    │ Server          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
git clone https://github.com/your-org/crossplane-ai.git
cd crossplane-ai
go mod download
go build -o crossplane-ai .
```

### Running Tests

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [kubectl-ai](https://github.com/GoogleCloudPlatform/kubectl-ai)
- Built on [Crossplane](https://crossplane.io/)
- Powered by [Cobra CLI](https://github.com/spf13/cobra)

## 🔗 Related Projects

- [kubectl-ai](https://github.com/GoogleCloudPlatform/kubectl-ai) - AI for Kubernetes
- [Crossplane](https://crossplane.io/) - Infrastructure as Code
- [Terramate](https://terramate.io/) - Terraform automation
- [k9s](https://k9scli.io/) - Kubernetes CLI dashboard

## Overview

Crossplane AI brings artificial intelligence to your Crossplane infrastructure management workflow. Instead of manually writing complex kubectl commands or YAML files, you can now interact with your Crossplane resources using natural language queries and get intelligent insights, suggestions, and automated analysis.

## Features

🎯 **Natural Language Queries**: Ask questions about your resources in plain English
- "What AWS resources do I have?"
- "Show me all failing databases"
- "Which resources are using the most resources?"

🔍 **Intelligent Analysis**: AI-powered resource analysis and health checking
- Automatic health scoring
- Issue detection and resolution suggestions
- Performance and cost optimization recommendations

💡 **Smart Suggestions**: Get AI-generated recommendations for:
- Resource optimization
- Security improvements
- Cost reduction
- Performance tuning
- Best practices implementation

🤖 **Interactive Mode**: Chat-like interface for continuous resource management
- Real-time resource monitoring
- Interactive troubleshooting
- Guided resource creation and modification

📊 **Multiple Output Formats**: Support for table, JSON, and YAML output formats

## Installation

### Prerequisites

- Go 1.24 or later
- Access to a Kubernetes cluster with Crossplane installed
- kubectl configured with appropriate permissions

### Build from Source

```bash
git clone https://github.com/your-org/crossplane-ai.git
cd crossplane-ai
go build -o crossplane-ai .
```

### Install Binary

```bash
# Move to your PATH
sudo mv crossplane-ai /usr/local/bin/
```

## Quick Start

### 1. Basic Usage

```bash
# Ask about your resources
crossplane-ai ask "what resources do I have?"

# Get suggestions for optimization
crossplane-ai suggest optimize

# Analyze resource health
crossplane-ai analyze
```

### 2. Interactive Mode

```bash
# Start interactive session
crossplane-ai interactive

# Or with initial analysis
crossplane-ai interactive --analyze --banner
```

### 3. Specific Resource Analysis

```bash
# Analyze specific resource
crossplane-ai analyze my-database

# Filter by provider
crossplane-ai analyze --provider aws

# Health check only
crossplane-ai analyze --health-check
```

## Commands

### `ask` - Natural Language Queries
Ask questions about your Crossplane resources using natural language.

```bash
crossplane-ai ask "what AWS RDS instances do I have?"
crossplane-ai ask "show me resources that are not ready"
crossplane-ai ask "help me optimize my database setup"
```

**Options:**
- `--provider` - Filter by specific provider (aws, gcp, azure)
- `--namespace` - Filter by namespace
- `--interactive` - Start interactive mode

### `analyze` - Resource Analysis
Perform detailed AI-powered analysis of your resources.

```bash
crossplane-ai analyze                    # Analyze all resources
crossplane-ai analyze my-resource        # Analyze specific resource
crossplane-ai analyze --provider aws     # Analyze AWS resources only
crossplane-ai analyze --health-check     # Focus on health checking
```

**Options:**
- `--provider` - Filter by provider
- `--namespace` - Filter by namespace  
- `--health-check` - Perform health check analysis
- `--summary` - Show summary instead of detailed output
- `--output` - Output format (table, json, yaml)

### `suggest` - AI Suggestions
Get intelligent recommendations for your infrastructure.

```bash
crossplane-ai suggest                # General suggestions
crossplane-ai suggest database      # Database-specific suggestions
crossplane-ai suggest security      # Security recommendations
crossplane-ai suggest optimize      # Optimization suggestions
```

**Options:**
- `--provider` - Focus on specific provider
- `--category` - Suggestion category (security, performance, cost, reliability)
- `--detailed` - Show detailed suggestions with examples
- `--limit` - Maximum number of suggestions

### `interactive` - Interactive Mode
Start an interactive AI-powered session.

```bash
crossplane-ai interactive           # Basic interactive mode
crossplane-ai i --banner           # With banner
crossplane-ai chat --analyze       # With initial analysis
```

**Options:**
- `--banner` - Show banner on startup
- `--analyze` - Perform initial analysis

## Configuration

### Configuration File

Create `~/.crossplane-ai.yaml`:

```yaml
# AI Service Configuration
ai:
  provider: "mock"  # or "openai", "google", "azure"
  # api_key: "${OPENAI_API_KEY}"
  # model: "gpt-4"

# Kubernetes Configuration  
kubernetes:
  kubeconfig: "~/.kube/config"
  context: ""
  namespace: ""

# Crossplane Configuration
crossplane:
  providers:
    - aws
    - gcp
    - azure
  resource_types:
    - compositions
    - providers
    - dbinstances
    - instances

# CLI Configuration
cli:
  output_format: "table"
  verbose: false
  color: true

# Analysis Configuration
analysis:
  timeout: 30
  max_suggestions: 10
  detailed: true
```

### Environment Variables

```bash
export KUBECONFIG=/path/to/kubeconfig
export CROSSPLANE_AI_VERBOSE=true
export OPENAI_API_KEY=your-api-key  # For real AI integration
```

### Global Flags

- `--config` - Path to config file
- `--kubeconfig` - Path to kubeconfig file
- `--context` - Kubernetes context to use
- `--verbose` - Enable verbose output

## Examples

### Resource Discovery

```bash
# Find all AWS resources
crossplane-ai ask "show me all AWS resources"

# Find failing resources
crossplane-ai ask "what resources are failing?"

# Database-specific queries
crossplane-ai ask "show me database performance metrics"
```

### Health Monitoring

```bash
# Overall health check
crossplane-ai analyze --health-check

# Provider-specific health
crossplane-ai analyze --provider aws --health-check

# Interactive monitoring
crossplane-ai interactive
> health
> analyze
> suggest optimize
```

### Optimization

```bash
# Get cost optimization suggestions
crossplane-ai suggest optimize --category cost

# Security recommendations
crossplane-ai suggest security --detailed

# Performance tuning
crossplane-ai suggest performance --provider aws
```

## Architecture

### Project Structure

```
crossplane-ai/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and global flags
│   ├── ask.go             # Natural language queries
│   ├── analyze.go         # Resource analysis
│   ├── suggest.go         # AI suggestions
│   └── interactive.go     # Interactive mode
├── pkg/
│   ├── ai/                # AI service integration
│   ├── crossplane/        # Crossplane client wrapper
│   └── cli/               # CLI utilities
├── internal/
│   └── config/            # Configuration management
└── main.go                # Application entry point
```

### Key Components

- **AI Service**: Handles natural language processing and intelligent analysis
- **Crossplane Client**: Kubernetes client wrapper for Crossplane CRDs
- **CLI Framework**: Built on Cobra for robust command-line interface
- **Configuration**: Viper-based configuration management

## Development

### Prerequisites

- Go 1.22+
- Access to Kubernetes cluster with Crossplane
- Make (optional)

### Building

```bash
# Build binary
go build -o crossplane-ai .

# Run tests
go test ./...

# Run linter
golangci-lint run
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [ ] **Real AI Integration**: OpenAI, Google AI, Azure AI integration
- [ ] **Advanced Analytics**: Resource usage trends and predictions
- [ ] **Automated Actions**: AI-driven resource provisioning and scaling
- [ ] **Policy Engine**: AI-generated governance policies
- [ ] **Web Dashboard**: Web-based interface for visual management
- [ ] **Slack/Teams Integration**: ChatOps support
- [ ] **Custom Plugins**: Extensible plugin system

## Troubleshooting

### Common Issues

**"Failed to initialize Crossplane client"**
- Ensure kubectl is configured and you have access to the cluster
- Check if Crossplane is installed in the cluster
- Verify RBAC permissions

**"No resources found"**
- Confirm Crossplane providers are installed and configured
- Check if resources exist in the specified namespace
- Verify resource types in configuration

**"AI processing failed"**
- Currently using mock AI - this is expected
- For real AI integration, configure API keys in config file
- Check network connectivity for API calls

### Debug Mode

```bash
# Enable verbose logging
crossplane-ai --verbose analyze

# Show configuration
crossplane-ai --verbose ask "test" 2>&1 | grep -i config
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Crossplane](https://crossplane.io/) - The amazing cloud-native control plane framework
- [Kubernetes](https://kubernetes.io/) - Container orchestration platform
- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [Viper](https://github.com/spf13/viper) - Configuration management

---

**Made with ❤️ for the Crossplane community**
