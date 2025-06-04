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

		// Initialize clients
		client, err := crossplane.NewClient(ctx)
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

	fmt.Println("🔬 Performing AI-powered analysis...")
	fmt.Println()

	// Get resources based on filters
	resources, err := client.GetFilteredResources(ctx, resourceName, provider, namespace)
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	if len(resources) == 0 {
		fmt.Println("No resources found matching the criteria.")
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
	fmt.Fprintln(w, "NAME\tTYPE\tSTATUS\tPROVIDER\tAGE")

	for _, resource := range analysis.Resources {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			resource.Name, resource.Type, resource.Status, resource.Provider, resource.Age)
	}
	w.Flush()
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

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().String("provider", "", "filter by provider")
	analyzeCmd.Flags().String("namespace", "", "filter by namespace")
	analyzeCmd.Flags().BoolP("health-check", "H", false, "perform health check analysis")
	analyzeCmd.Flags().BoolP("summary", "s", false, "show summary instead of detailed output")
	analyzeCmd.Flags().String("output", "table", "output format (table, json, yaml)")
}
