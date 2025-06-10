package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"crossplane-ai/pkg/ai"
	"crossplane-ai/pkg/crossplane"

	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask AI about your Crossplane resources",
	Long: `Ask natural language questions about your Crossplane resources and get intelligent answers.
The AI will analyze your cluster's Crossplane resources and provide helpful insights.`,
	Example: `  # Ask about available resources
  crossplane-ai ask "what resources do I have?"
  
  # Ask about specific providers
  crossplane-ai ask "show me all AWS resources"
  
  # Ask for troubleshooting help
  crossplane-ai ask "why is my database not ready?"
  
  # Interactive mode (no question provided)
  crossplane-ai ask`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		var question string
		if len(args) > 0 {
			question = strings.Join(args, " ")
		}

		// Check if running in mock mode
		if IsMockMode() {
			fmt.Println("🤖 AI Assistant (MOCK MODE)")
			fmt.Println("===========================")
			if question == "" {
				fmt.Println("Interactive mode not supported in mock mode.")
				fmt.Println("Please provide a question as an argument.")
				return nil
			}
			return handleMockAsk(ctx, question)
		}

		// Initialize Crossplane client for real mode
		client, err := crossplane.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Crossplane client: %w", err)
		}

		// Initialize AI service
		aiService := ai.NewService()

		// Show AI mode information
		if aiService.IsUsingRealAI() {
			fmt.Println("🤖 AI Assistant (POWERED BY OPENAI)")
		} else {
			fmt.Println("🤖 AI Assistant (TEMPLATE MODE)")
			fmt.Println("💡 Set OPENAI_API_KEY for real AI capabilities")
		}
		fmt.Println("===========================")
		fmt.Println()

		if question == "" {
			// Interactive mode
			return runInteractiveMode(ctx, client, aiService)
		}

		return processQuestion(ctx, client, aiService, question)
	},
}

func runInteractiveMode(ctx context.Context, client *crossplane.Client, aiService *ai.Service) error {
	fmt.Println("🤖 Crossplane AI Interactive Mode")
	fmt.Println("Ask me anything about your Crossplane resources! Type 'exit' to quit.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("💬 You: ")
		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}

		if strings.ToLower(question) == "exit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		fmt.Print("🤖 AI: ")
		if err := processQuestion(ctx, client, aiService, question); err != nil {
			fmt.Printf("Sorry, I encountered an error: %v\n", err)
		}
		fmt.Println()
	}

	return scanner.Err()
}

func processQuestion(ctx context.Context, client *crossplane.Client, aiService *ai.Service, question string) error {
	// Get current cluster state
	resources, err := client.GetAllResources(ctx)
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	// Process with AI
	response, err := aiService.ProcessQuery(ctx, question, resources)
	if err != nil {
		return fmt.Errorf("AI processing failed: %w", err)
	}

	fmt.Println(response)
	return nil
}

// handleMockAsk handles ask command in mock mode
func handleMockAsk(ctx context.Context, question string) error {
	fmt.Printf("Question: %s\n\n", question)

	// Get embedded mock resources for more dynamic responses
	mockResources := ai.GetEmbeddedMockResources()

	// Mock responses based on question patterns
	questionLower := strings.ToLower(question)

	var response string
	switch {
	case strings.Contains(questionLower, "what") && strings.Contains(questionLower, "resources"):
		// Count resources by type and provider
		compositions := 0
		providers := 0
		infraResources := 0
		awsResources := 0

		for _, res := range mockResources {
			switch res.Type {
			case "compositions", "compositeresourcedefinitions":
				compositions++
			case "providers":
				providers++
			default:
				infraResources++
				if res.Provider == "aws" {
					awsResources++
				}
			}
		}

		response = fmt.Sprintf(`📊 Mock Resource Summary:

Based on your mock Crossplane cluster, you have:
• %d Composition(s) and CompositeResourceDefinition(s)
• %d Provider(s) installed
• %d Infrastructure resource(s) total
• %d AWS-managed resources

Resource Status: %d Ready, %d Not Ready

All core resources are operational! 🎉`,
			compositions, providers, infraResources, awsResources,
			len(mockResources)-1, 1) // -1 ready, 1 not ready (failing-test-resource)

	case strings.Contains(questionLower, "aws"):
		var awsResources []string
		for _, res := range mockResources {
			if res.Provider == "aws" {
				status := "✅"
				if res.Status != "Ready" {
					status = "⚠️"
				}
				awsResources = append(awsResources, fmt.Sprintf("• %s (%s) - %s %s", res.Name, res.Type, res.Status, status))
			}
		}

		response = fmt.Sprintf(`🔍 AWS Resources (Mock):

Your AWS resources managed by Crossplane:
%s

Most resources are healthy and operational.`, strings.Join(awsResources, "\n"))

	case strings.Contains(questionLower, "database") || strings.Contains(questionLower, "db"):
		var dbResources []string
		for _, res := range mockResources {
			if strings.Contains(res.Type, "db") || strings.Contains(res.Name, "database") {
				status := "✅"
				if res.Status != "Ready" {
					status = "⚠️"
				}
				dbResources = append(dbResources, fmt.Sprintf("• %s (%s) - %s %s, Age: %s", res.Name, res.Type, res.Status, status, res.Age))
			}
		}

		if len(dbResources) == 0 {
			response = `🗄️ Database Resources (Mock):

No database resources found in the current mock data.
Try asking about other resource types or run 'crossplane-ai --mock analyze' to see all resources.`
		} else {
			response = fmt.Sprintf(`🗄️ Database Resources (Mock):

Found database resources:
%s

The databases appear to be healthy and operational.`, strings.Join(dbResources, "\n"))
		}

	case strings.Contains(questionLower, "not ready") || strings.Contains(questionLower, "failed") || strings.Contains(questionLower, "problem"):
		var failingResources []string
		for _, res := range mockResources {
			if res.Status != "Ready" {
				failingResources = append(failingResources, fmt.Sprintf("• %s (%s) - %s", res.Name, res.Type, res.Status))
			}
		}

		if len(failingResources) == 0 {
			response = `🔧 Troubleshooting (Mock):

Good news! In this mock environment, all resources are healthy.

In a real environment, here's how to troubleshoot:
1. Check resource events: kubectl describe <resource>
2. Verify provider credentials
3. Check network connectivity
4. Review Crossplane logs`
		} else {
			response = fmt.Sprintf(`🔧 Troubleshooting (Mock):

Found resources with issues:
%s

In a real environment, here's how to troubleshoot:
1. Check resource events: kubectl describe <resource>
2. Verify provider credentials  
3. Check network connectivity
4. Review Crossplane logs`, strings.Join(failingResources, "\n"))
		}

	case strings.Contains(questionLower, "providers") || strings.Contains(questionLower, "provider"):
		var providers []string
		for _, res := range mockResources {
			if res.Type == "providers" {
				status := "✅"
				if res.Status != "Ready" {
					status = "⚠️"
				}
				providers = append(providers, fmt.Sprintf("• %s - %s %s, Age: %s", res.Name, res.Status, status, res.Age))
			}
		}

		response = fmt.Sprintf(`🏗️ Providers (Mock):

Installed providers:
%s

These providers enable management of cloud resources through Crossplane.`, strings.Join(providers, "\n"))

	default:
		response = fmt.Sprintf(`🤖 Mock AI Response:

I understand you're asking: "%s"

In this mock environment, I can tell you:
• Your Crossplane setup includes %d total resources
• Multiple providers are configured (AWS, GCP, Azure)  
• Most resources are healthy
• You have compositions and resource definitions ready

Try asking about:
• "what resources do I have?" - for a full summary
• "show me AWS resources" - for provider-specific info
• "tell me about databases" - for resource type info
• "what providers are installed?" - for provider status

🧪 This is mock mode - responses use embedded sample data.`, question, len(mockResources))
	}

	fmt.Println(response)
	fmt.Println()
	fmt.Println("🧪 This was a mock response using embedded sample data.")

	return nil
}

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().String("provider", "", "filter by specific provider (aws, gcp, azure)")
	askCmd.Flags().String("namespace", "", "filter by namespace")
	askCmd.Flags().BoolP("interactive", "i", false, "start interactive mode")
}
