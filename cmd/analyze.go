package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"crossplane-ai/pkg/ai"
	"crossplane-ai/pkg/crossplane"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [resource-name]",
	Short: "Analyze Crossplane resources with AI insights",
	Long: `Perform intelligent analysis of your Crossplane resources. Get health checks,
performance insights, security recommendations, and troubleshooting suggestions.`,
	Example: `  # Analyze all resources
  crossplane-ai analyze
  
  # Analyze specific resource
  crossplane-ai analyze my-database
  
  # Analyze by provider
  crossplane-ai analyze --provider aws
  
  # Health check analysis
  crossplane-ai analyze --health-check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Check if running in mock mode
		if IsMockMode() {
			fmt.Println("🔬 Performing AI-powered analysis (MOCK MODE)...")
			fmt.Println()
			return performMockAnalysis(ctx, cmd, args)
		}

		// Initialize clients for real mode
		contextFlag, _ := cmd.Flags().GetString("context")
		kubeconfigFlag, _ := cmd.Flags().GetString("kubeconfig")

		clientOpts := crossplane.ClientOptions{
			Context:    contextFlag,
			Kubeconfig: kubeconfigFlag,
		}

		client, err := crossplane.NewClientWithOptions(ctx, clientOpts)
		if err != nil {
			return fmt.Errorf("failed to initialize Crossplane client: %w", err)
		}

		aiService := ai.NewService()

		// Get flags
		provider, _ := cmd.Flags().GetString("provider")
		namespace, _ := cmd.Flags().GetString("namespace")
		healthCheck, _ := cmd.Flags().GetBool("health-check")
		summary, _ := cmd.Flags().GetBool("summary")

		var resourceName string
		if len(args) > 0 {
			resourceName = args[0]
		}

		return performAnalysis(ctx, client, aiService, resourceName, provider, namespace, healthCheck, summary)
	},
}

func performAnalysis(ctx context.Context, client *crossplane.Client, aiService *ai.Service,
	resourceName, provider, namespace string, healthCheck, summary bool) error {

	// Show AI mode information
	if aiService.IsUsingRealAI() {
		fmt.Println("🔬 Performing AI-powered analysis (OpenAI)...")
	} else {
		fmt.Println("🔬 Performing analysis (template-based)...")
		fmt.Println("💡 Set OPENAI_API_KEY for AI-powered insights")
	}
	fmt.Println()

	// Get resources based on filters
	resources, err := client.GetFilteredResources(ctx, resourceName, provider, namespace)
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	if len(resources) == 0 {
		fmt.Println("No Crossplane resources found matching the criteria.")
		fmt.Println("Please ensure Crossplane is installed and you have created some resources.")
		return nil
	}

	// Perform AI analysis
	analysis, err := aiService.AnalyzeResources(ctx, resources, healthCheck)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if summary {
		printSummary(analysis)
	} else {
		printDetailedAnalysis(analysis)
	}

	return nil
}

func printSummary(analysis *ai.Analysis) {
	fmt.Println("📊 Analysis Summary")
	fmt.Println("==================")
	fmt.Printf("Total Resources: %d\n", analysis.TotalResources)
	fmt.Printf("Healthy: %d\n", analysis.HealthyResources)
	fmt.Printf("Issues Found: %d\n", analysis.IssuesFound)
	fmt.Printf("Recommendations: %d\n", len(analysis.Recommendations))
	fmt.Println()

	if len(analysis.Recommendations) > 0 {
		fmt.Println("🎯 Top Recommendations:")
		for i, rec := range analysis.Recommendations {
			if i >= 3 { // Show only top 3 in summary
				break
			}
			fmt.Printf("• %s\n", rec.Title)
		}
	}
}

func printDetailedAnalysis(analysis *ai.Analysis) {
	// Resource Status Table
	fmt.Println("📋 Resource Status")
	fmt.Println("==================")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tTYPE\tSTATUS\tPROVIDER\tAGE")

	for _, resource := range analysis.Resources {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			resource.Name, resource.Type, resource.Status, resource.Provider, resource.Age)
	}
	_ = w.Flush()
	fmt.Println()

	// Issues
	if len(analysis.Issues) > 0 {
		fmt.Println("⚠️  Issues Detected")
		fmt.Println("==================")
		for _, issue := range analysis.Issues {
			fmt.Printf("• %s: %s\n", issue.Severity, issue.Description)
			if issue.Resolution != "" {
				fmt.Printf("  Resolution: %s\n", issue.Resolution)
			}
		}
		fmt.Println()
	}

	// Recommendations
	if len(analysis.Recommendations) > 0 {
		fmt.Println("💡 AI Recommendations")
		fmt.Println("====================")
		for i, rec := range analysis.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec.Title)
			fmt.Printf("   %s\n", rec.Description)
			if rec.Impact != "" {
				fmt.Printf("   Impact: %s\n", rec.Impact)
			}
		}
		fmt.Println()
	}

	// Health Score
	if analysis.HealthScore > 0 {
		fmt.Printf("🏥 Overall Health Score: %d/100\n", analysis.HealthScore)
		if analysis.HealthScore >= 80 {
			fmt.Println("✅ Your Crossplane setup looks healthy!")
		} else if analysis.HealthScore >= 60 {
			fmt.Println("⚠️  Some areas need attention.")
		} else {
			fmt.Println("🚨 Critical issues detected. Please review recommendations.")
		}
	}
}

// performMockAnalysis performs analysis using mock data
func performMockAnalysis(ctx context.Context, cmd *cobra.Command, args []string) error {
	// Get flags for mock analysis
	provider, _ := cmd.Flags().GetString("provider")
	_, _ = cmd.Flags().GetString("namespace") // namespace not used in mock mode
	healthCheck, _ := cmd.Flags().GetBool("health-check")
	summary, _ := cmd.Flags().GetBool("summary")

	var resourceName string
	if len(args) > 0 {
		resourceName = args[0]
	}

	// Use embedded mock data that works without external files
	fmt.Println("📁 Using embedded mock data (no external files required)")
	fmt.Println()

	// Create mock AI service and analysis
	aiService := ai.NewService()

	// Get embedded mock resources
	mockResources := ai.GetEmbeddedMockResources()

	// Apply filters if specified
	var filteredResources []*ai.ResourceInfo
	for _, res := range mockResources {
		if resourceName != "" && res.Name != resourceName {
			continue
		}
		if provider != "" && res.Provider != provider {
			continue
		}
		// Note: namespace filtering not applicable for these mock resources
		filteredResources = append(filteredResources, res)
	}

	if len(filteredResources) == 0 {
		fmt.Println("No mock Crossplane resources found matching the criteria.")
		fmt.Println("Available mock resources:")
		for _, res := range mockResources {
			fmt.Printf("  - %s (%s, provider: %s)\n", res.Name, res.Type, res.Provider)
		}
		return nil
	}

	// Perform analysis with mock data
	analysis, err := aiService.AnalyzeResources(ctx, filteredResources, healthCheck)
	if err != nil {
		return fmt.Errorf("mock analysis failed: %w", err)
	}

	if summary {
		printSummary(analysis)
	} else {
		printDetailedAnalysis(analysis)
	}

	fmt.Println()
	fmt.Println("🧪 This was a mock analysis using embedded sample data.")
	fmt.Println("To analyze real resources, run without the --mock flag")

	return nil
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().String("provider", "", "filter by provider")
	analyzeCmd.Flags().String("namespace", "", "filter by namespace")
	analyzeCmd.Flags().BoolP("health-check", "H", false, "perform health check analysis")
	analyzeCmd.Flags().BoolP("summary", "s", false, "show summary instead of detailed output")
	analyzeCmd.Flags().String("output", "table", "output format (table, json, yaml)")
}
