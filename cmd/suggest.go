package cmd

import (
	"context"
	"fmt"

	"crossplane-ai/pkg/ai"
	"crossplane-ai/pkg/crossplane"

	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest [resource-type]",
	Short: "Get AI-powered suggestions for Crossplane resources",
	Long: `Get intelligent suggestions for creating, optimizing, or troubleshooting 
Crossplane resources. The AI analyzes your current setup and provides tailored recommendations.`,
	Example: `  # Get database suggestions
  crossplane-ai suggest database
  
  # Get networking suggestions
  crossplane-ai suggest network
  
  # Get general optimization suggestions
  crossplane-ai suggest optimize
  
  # Get security recommendations
  crossplane-ai suggest security`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Initialize clients
		client, err := crossplane.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Crossplane client: %w", err)
		}

		aiService := ai.NewService()

		var suggestionType string
		if len(args) > 0 {
			suggestionType = args[0]
		} else {
			suggestionType = "general"
		}

		return generateSuggestions(ctx, client, aiService, suggestionType)
	},
}

func generateSuggestions(ctx context.Context, client *crossplane.Client, aiService *ai.Service, suggestionType string) error {
	fmt.Printf("🔍 Analyzing your Crossplane setup for %s suggestions...\n\n", suggestionType)

	// Get current resources
	resources, err := client.GetAllResources(ctx)
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	// Get AI suggestions
	suggestions, err := aiService.GenerateSuggestions(ctx, suggestionType, resources)
	if err != nil {
		return fmt.Errorf("failed to generate suggestions: %w", err)
	}

	fmt.Println("💡 AI Suggestions:")
	fmt.Println("==================")

	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion.Title)
		fmt.Printf("   %s\n", suggestion.Description)

		if suggestion.Example != "" {
			fmt.Printf("   Example:\n   %s\n", suggestion.Example)
		}

		if suggestion.Priority != "" {
			fmt.Printf("   Priority: %s\n", suggestion.Priority)
		}

		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(suggestCmd)

	suggestCmd.Flags().String("provider", "", "focus suggestions on specific provider")
	suggestCmd.Flags().String("category", "", "suggestion category (security, performance, cost, reliability)")
	suggestCmd.Flags().BoolP("detailed", "d", false, "show detailed suggestions with examples")
	suggestCmd.Flags().IntP("limit", "l", 5, "maximum number of suggestions to show")
}
