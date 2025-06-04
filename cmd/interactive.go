package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"crossplane-ai/pkg/ai"
	"crossplane-ai/pkg/cli"
	"crossplane-ai/pkg/crossplane"

	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Short:   "Start interactive AI-powered Crossplane management session",
	Aliases: []string{"i", "chat"},
	Long: `Launch an interactive session where you can chat with AI about your Crossplane resources.
This mode provides a conversational interface for managing, analyzing, and troubleshooting
your Crossplane infrastructure.`,
	Example: `  # Start interactive mode
  crossplane-ai interactive
  
  # Start with banner
  crossplane-ai interactive --banner
  
  # Start with initial analysis
  crossplane-ai interactive --analyze`,
	RunE: func(cmd *cobra.Command, args []string) error {
		showBanner, _ := cmd.Flags().GetBool("banner")
		initialAnalyze, _ := cmd.Flags().GetBool("analyze")

		return runInteractiveSession(showBanner, initialAnalyze)
	},
}

func runInteractiveSession(showBanner, initialAnalyze bool) error {
	ctx := context.Background()

	// Initialize clients
	client, err := crossplane.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize Crossplane client: %w", err)
	}

	aiService := ai.NewService()

	// Show banner if requested
	if showBanner {
		cli.PrintBanner()
	}

	// Welcome message
	fmt.Println("🤖 Welcome to Crossplane AI Interactive Mode!")
	fmt.Println()
	cli.PrintInfo("I can help you manage and analyze your Crossplane resources using natural language.")
	fmt.Println()

	// Show available commands
	printInteractiveHelp()

	// Perform initial analysis if requested
	if initialAnalyze {
		fmt.Println("🔍 Performing initial analysis of your Crossplane resources...")
		if err := performQuickAnalysis(ctx, client, aiService); err != nil {
			cli.PrintWarning(fmt.Sprintf("Initial analysis failed: %v", err))
		}
		fmt.Println()
	}

	// Start interactive loop
	return startInteractiveLoop(ctx, client, aiService)
}

func printInteractiveHelp() {
	cli.PrintSubHeader("Available Commands")
	fmt.Println("💬 Ask anything about your resources: 'what AWS resources do I have?'")
	fmt.Println("🔍 analyze - Perform detailed resource analysis")
	fmt.Println("💡 suggest [type] - Get AI suggestions (e.g., 'suggest database')")
	fmt.Println("📊 status - Show resource status overview")
	fmt.Println("🏥 health - Perform health check")
	fmt.Println("❓ help - Show this help message")
	fmt.Println("👋 exit/quit - Exit interactive mode")
	fmt.Println()
}

func startInteractiveLoop(ctx context.Context, client *crossplane.Client, aiService *ai.Service) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("🤖 crossplane-ai> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle special commands
		if handled, exit := handleSpecialCommands(ctx, client, aiService, input); handled {
			if exit {
				break
			}
			continue
		}

		// Process as natural language query
		if err := processInteractiveQuery(ctx, client, aiService, input); err != nil {
			cli.PrintError(fmt.Sprintf("Error: %v", err))
		}

		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println("👋 Thanks for using Crossplane AI! Goodbye!")
	return nil
}

func handleSpecialCommands(ctx context.Context, client *crossplane.Client, aiService *ai.Service, input string) (bool, bool) {
	command := strings.ToLower(strings.TrimSpace(input))

	switch {
	case command == "exit" || command == "quit" || command == "q":
		return true, true

	case command == "help" || command == "?":
		printInteractiveHelp()
		return true, false

	case command == "analyze":
		performDetailedAnalysis(ctx, client, aiService)
		return true, false

	case command == "status":
		showResourceStatus(ctx, client)
		return true, false

	case command == "health":
		performHealthCheck(ctx, client, aiService)
		return true, false

	case strings.HasPrefix(command, "suggest"):
		parts := strings.Fields(input)
		suggestionType := "general"
		if len(parts) > 1 {
			suggestionType = parts[1]
		}
		showSuggestions(ctx, client, aiService, suggestionType)
		return true, false
	}

	return false, false
}

func processInteractiveQuery(ctx context.Context, client *crossplane.Client, aiService *ai.Service, query string) error {
	// Get resources for context
	resources, err := client.GetAllResources(ctx)
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	// Process with AI
	response, err := aiService.ProcessQuery(ctx, query, resources)
	if err != nil {
		return fmt.Errorf("AI processing failed: %w", err)
	}

	fmt.Println(response)
	return nil
}

func performQuickAnalysis(ctx context.Context, client *crossplane.Client, aiService *ai.Service) error {
	resources, err := client.GetAllResources(ctx)
	if err != nil {
		return err
	}

	if len(resources) == 0 {
		cli.PrintWarning("No Crossplane resources found in the cluster")
		return nil
	}

	cli.PrintSuccess(fmt.Sprintf("Found %d Crossplane resources", len(resources)))

	// Quick status summary
	readyCount := 0
	for _, r := range resources {
		if strings.ToLower(r.Status) == "ready" {
			readyCount++
		}
	}

	fmt.Printf("📊 Status: %d/%d resources are ready\n", readyCount, len(resources))

	if readyCount < len(resources) {
		cli.PrintWarning("Some resources are not ready - type 'health' for details")
	}

	return nil
}

func performDetailedAnalysis(ctx context.Context, client *crossplane.Client, aiService *ai.Service) {
	fmt.Println("🔬 Performing detailed analysis...")

	resources, err := client.GetAllResources(ctx)
	if err != nil {
		cli.PrintError(fmt.Sprintf("Failed to get resources: %v", err))
		return
	}

	analysis, err := aiService.AnalyzeResources(ctx, resources, true)
	if err != nil {
		cli.PrintError(fmt.Sprintf("Analysis failed: %v", err))
		return
	}

	// Print analysis results
	cli.PrintHeader("Analysis Results")
	fmt.Printf("Total Resources: %d\n", analysis.TotalResources)
	fmt.Printf("Healthy Resources: %d\n", analysis.HealthyResources)
	fmt.Printf("Issues Found: %d\n", analysis.IssuesFound)
	fmt.Printf("Health Score: %d/100\n", analysis.HealthScore)

	if len(analysis.Issues) > 0 {
		cli.PrintSubHeader("Issues Detected")
		for _, issue := range analysis.Issues {
			cli.PrintWarning(fmt.Sprintf("%s: %s", issue.Severity, issue.Description))
		}
	}

	if len(analysis.Recommendations) > 0 {
		cli.PrintSubHeader("Recommendations")
		for i, rec := range analysis.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec.Title)
		}
	}
}

func showResourceStatus(ctx context.Context, client *crossplane.Client) {
	fmt.Println("📋 Resource Status Overview")

	resources, err := client.GetAllResources(ctx)
	if err != nil {
		cli.PrintError(fmt.Sprintf("Failed to get resources: %v", err))
		return
	}

	if len(resources) == 0 {
		cli.PrintWarning("No resources found")
		return
	}

	// Create table
	headers := []string{"NAME", "TYPE", "STATUS", "PROVIDER", "AGE"}
	var rows [][]string

	for _, r := range resources {
		rows = append(rows, []string{
			cli.TruncateString(r.Name, 30),
			r.Type,
			cli.FormatStatus(r.Status),
			r.Provider,
			cli.FormatAge(r.Age),
		})
	}

	cli.PrintTable(headers, rows)
}

func performHealthCheck(ctx context.Context, client *crossplane.Client, aiService *ai.Service) {
	fmt.Println("🏥 Performing health check...")

	resources, err := client.GetAllResources(ctx)
	if err != nil {
		cli.PrintError(fmt.Sprintf("Failed to get resources: %v", err))
		return
	}

	// Simple health check
	totalResources := len(resources)
	healthyResources := 0

	for _, r := range resources {
		if strings.ToLower(r.Status) == "ready" {
			healthyResources++
		}
	}

	healthPercentage := 0
	if totalResources > 0 {
		healthPercentage = (healthyResources * 100) / totalResources
	}

	fmt.Printf("🏥 Health Score: %d/100\n", healthPercentage)
	fmt.Printf("📊 %d/%d resources are healthy\n", healthyResources, totalResources)

	if healthPercentage >= 80 {
		cli.PrintSuccess("Your Crossplane setup is healthy!")
	} else if healthPercentage >= 60 {
		cli.PrintWarning("Some resources need attention")
	} else {
		cli.PrintError("Critical issues detected - run 'analyze' for details")
	}
}

func showSuggestions(ctx context.Context, client *crossplane.Client, aiService *ai.Service, suggestionType string) {
	fmt.Printf("💡 Generating %s suggestions...\n", suggestionType)

	resources, err := client.GetAllResources(ctx)
	if err != nil {
		cli.PrintError(fmt.Sprintf("Failed to get resources: %v", err))
		return
	}

	suggestions, err := aiService.GenerateSuggestions(ctx, suggestionType, resources)
	if err != nil {
		cli.PrintError(fmt.Sprintf("Failed to generate suggestions: %v", err))
		return
	}

	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion.Title)
		fmt.Printf("   %s\n", suggestion.Description)
		if suggestion.Priority != "" {
			fmt.Printf("   Priority: %s\n", suggestion.Priority)
		}
		fmt.Println()
	}
}

func init() {
	rootCmd.AddCommand(interactiveCmd)

	interactiveCmd.Flags().BoolP("banner", "b", false, "show banner on startup")
	interactiveCmd.Flags().BoolP("analyze", "a", false, "perform initial analysis")
}
