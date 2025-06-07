package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"crossplane-ai/pkg/ai"

	"github.com/spf13/cobra"
)

var generateExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Generate example Crossplane YAML files",
	Long: `Generate example Crossplane YAML files for testing and learning.
These files can be used as templates or with the mock mode.`,
	Example: `  # Generate examples in current directory
  crossplane-ai generate examples
  
  # Generate examples in specific directory
  crossplane-ai generate examples --output ./examples
  
  # List available example types
  crossplane-ai generate examples --list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir, _ := cmd.Flags().GetString("output")
		listOnly, _ := cmd.Flags().GetBool("list")

		if listOnly {
			return listExampleTypes()
		}

		return generateExampleFiles(outputDir)
	},
}

func listExampleTypes() error {
	fmt.Println("📋 Available Example Types:")
	fmt.Println("==========================")

	examples := ai.GetEmbeddedMockYAMLExamples()
	for exampleType := range examples {
		fmt.Printf("• %s\n", exampleType)
	}

	fmt.Println()
	fmt.Println("These examples demonstrate common Crossplane patterns and can be used")
	fmt.Println("as starting points for your own infrastructure definitions.")

	return nil
}

func generateExampleFiles(outputDir string) error {
	if outputDir == "" {
		outputDir = "./examples"
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("📁 Generating example files in: %s\n", outputDir)
	fmt.Println()

	examples := ai.GetEmbeddedMockYAMLExamples()
	fileMap := map[string]string{
		"composition": "xdatabase-composition.yaml",
		"xrd":         "xdatabase-definition.yaml",
		"claim":       "database-claim.yaml",
		"provider":    "provider-aws.yaml",
	}

	for exampleType, content := range examples {
		filename := fileMap[exampleType]
		if filename == "" {
			filename = fmt.Sprintf("%s.yaml", exampleType)
		}

		filePath := filepath.Join(outputDir, filename)

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		fmt.Printf("✅ Created: %s\n", filename)
	}

	fmt.Println()
	fmt.Println("🎉 Example files generated successfully!")
	fmt.Println()
	fmt.Println("You can now:")
	fmt.Printf("• Apply them to your cluster: kubectl apply -f %s/\n", outputDir)
	fmt.Printf("• Use them for mock testing: crossplane-ai --mock --mock-data-dir %s analyze\n", outputDir)
	fmt.Println("• Modify them as templates for your own resources")

	return nil
}

func init() {
	generateCmd.AddCommand(generateExamplesCmd)

	generateExamplesCmd.Flags().String("output", "./examples", "output directory for generated files")
	generateExamplesCmd.Flags().Bool("list", false, "list available example types without generating files")
}
