package cmd

import (
	"context"
	"fmt"
	"strings"

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

		var suggestionType string
		if len(args) > 0 {
			suggestionType = args[0]
		} else {
			suggestionType = "general"
		}

		// Check if running in mock mode
		if IsMockMode() {
			fmt.Printf("🔍 Analyzing your Crossplane setup for %s suggestions (using embedded sample data)...\n\n", suggestionType)
			return generateMockSuggestions(ctx, suggestionType)
		}

		// Initialize clients for real mode
		client, err := crossplane.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Crossplane client: %w", err)
		}

		aiService := ai.NewService()

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

// generateMockSuggestions generates mock suggestions for testing using embedded data
func generateMockSuggestions(ctx context.Context, suggestionType string) error {
	fmt.Println("💡 AI Suggestions:")
	fmt.Println("==================")

	// Get embedded mock resources to make suggestions contextual
	mockResources := ai.GetEmbeddedMockResources()

	// Count resources by type for contextual suggestions
	providerCount := 0
	dbCount := 0
	healthyResources := 0

	for _, resource := range mockResources {
		if resource.Type == "providers" {
			providerCount++
		}
		if strings.Contains(strings.ToLower(resource.Type), "db") || strings.Contains(strings.ToLower(resource.Name), "database") {
			dbCount++
		}
		if resource.Status == "Ready" {
			healthyResources++
		}
	}

	// Generate contextual suggestions based on suggestion type and mock resources
	switch strings.ToLower(suggestionType) {
	case "database", "db":
		fmt.Printf("1. Enable Automated Backups for %d Database Resources\n", dbCount)
		fmt.Println("   Configure automated backups for your RDS instances to ensure data protection")
		fmt.Println("   Priority: High")
		fmt.Println()
		fmt.Println("2. Implement Read Replicas")
		fmt.Println("   Add read replicas to distribute read traffic and improve performance")
		fmt.Println("   Priority: Medium")
		fmt.Println()
		fmt.Printf("3. Configure Multi-AZ Deployment (Based on %d healthy resources)\n", healthyResources)
		fmt.Println("   Enable Multi-AZ for high availability and automatic failover")
		fmt.Println("   Priority: High")

	case "security":
		fmt.Println("1. Enable Encryption at Rest")
		fmt.Println("   Encrypt your databases and storage resources to enhance security")
		fmt.Println("   Priority: High")
		fmt.Println()
		fmt.Println("2. Review IAM Policies")
		fmt.Println("   Audit and tighten IAM policies to follow principle of least privilege")
		fmt.Println("   Priority: High")
		fmt.Println()
		fmt.Println("3. Enable VPC Security Groups")
		fmt.Println("   Configure proper security groups to restrict network access")
		fmt.Println("   Priority: Medium")

	case "optimize", "optimization":
		fmt.Println("1. Right-size Resources")
		fmt.Println("   Analyze resource utilization and adjust instance sizes accordingly")
		fmt.Println("   Priority: Medium")
		fmt.Println()
		fmt.Println("2. Implement Auto-scaling")
		fmt.Println("   Configure auto-scaling groups to optimize resource usage")
		fmt.Println("   Priority: Medium")
		fmt.Println()
		fmt.Println("3. Use Reserved Instances")
		fmt.Println("   Consider reserved instances for predictable workloads to reduce costs")
		fmt.Println("   Priority: Low")

	case "network", "networking":
		fmt.Println("1. Configure VPC Peering")
		fmt.Println("   Set up VPC peering for secure cross-region communication")
		fmt.Println("   Priority: Medium")
		fmt.Println()
		fmt.Println("2. Implement Load Balancers")
		fmt.Println("   Use Application Load Balancers to distribute traffic efficiently")
		fmt.Println("   Priority: High")
		fmt.Println()
		fmt.Println("3. Set up VPN Gateway")
		fmt.Println("   Configure VPN gateway for secure on-premises connectivity")
		fmt.Println("   Priority: Low")

	default: // general
		fmt.Printf("1. Health Check %d Resources\n", len(mockResources))
		fmt.Println("   Run regular health checks on all Crossplane resources")
		fmt.Println("   Priority: Medium")
		fmt.Println()
		fmt.Printf("2. Update %d Providers\n", providerCount)
		fmt.Println("   Keep your Crossplane providers up to date for latest features and security fixes")
		fmt.Println("   Priority: Low")
		fmt.Println()
		fmt.Println("3. Implement Monitoring")
		fmt.Println("   Set up comprehensive monitoring for all your infrastructure resources")
		fmt.Println("   Priority: High")
		fmt.Println()
		fmt.Println("4. Create Backup Strategy")
		fmt.Println("   Develop and implement a comprehensive backup and disaster recovery plan")
		fmt.Println("   Priority: High")
	}

	fmt.Println()
	fmt.Println("🧪 These are mock suggestions using embedded sample data.")
	fmt.Println("To get real suggestions, run without the --mock flag")

	return nil
}

func init() {
	rootCmd.AddCommand(suggestCmd)

	suggestCmd.Flags().String("provider", "", "focus suggestions on specific provider")
	suggestCmd.Flags().String("category", "", "suggestion category (security, performance, cost, reliability)")
	suggestCmd.Flags().BoolP("detailed", "d", false, "show detailed suggestions with examples")
	suggestCmd.Flags().IntP("limit", "l", 5, "maximum number of suggestions to show")
}
