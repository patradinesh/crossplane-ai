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

		// Initialize Crossplane client
		client, err := crossplane.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Crossplane client: %w", err)
		}

		// Initialize AI service
		aiService := ai.NewService()

		var question string
		if len(args) > 0 {
			question = strings.Join(args, " ")
		} else {
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

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().String("provider", "", "filter by specific provider (aws, gcp, azure)")
	askCmd.Flags().String("namespace", "", "filter by namespace")
	askCmd.Flags().BoolP("interactive", "i", false, "start interactive mode")
}
