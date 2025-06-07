package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crossplane-ai",
	Short: "AI-powered Crossplane resource management",
	Long: `Crossplane AI is a command-line tool that brings artificial intelligence
capabilities to Crossplane resource management in Kubernetes clusters.

Ask questions, get suggestions, and perform intelligent analysis of your
cloud infrastructure resources managed by Crossplane.`,
	Example: `  # Ask about your resources
  crossplane-ai ask "What databases do I have?"
  
  # Get optimization suggestions
  crossplane-ai suggest optimization
  
  # Analyze resource health
  crossplane-ai analyze --health-check
  
  # Use specific cluster context
  crossplane-ai --context eks-cluster analyze
  
  # Run in mock mode for testing/demos (uses embedded data)
  crossplane-ai --mock analyze
  
  # Generate example files for learning or custom mock data
  crossplane-ai generate examples
  
  # Use custom mock data directory
  crossplane-ai --mock --mock-data-dir ./my-examples analyze`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.crossplane-ai.yaml)")
	rootCmd.PersistentFlags().String("context", "", "kubectl context to use (overrides current context)")
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().Bool("mock", false, "run in mock mode with embedded sample data (for testing and demos)")
	rootCmd.PersistentFlags().String("mock-data-dir", "", "directory containing mock data files (optional, uses embedded data if not specified)")

	// Bind flags to viper
	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	_ = viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("mock", rootCmd.PersistentFlags().Lookup("mock"))
	_ = viper.BindPFlag("mock-data-dir", rootCmd.PersistentFlags().Lookup("mock-data-dir"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".crossplane-ai")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// IsMockMode checks if the tool should run in mock mode
func IsMockMode() bool {
	// Check command line flag first
	if viper.GetBool("mock") {
		return true
	}
	// Fall back to environment variable for backward compatibility
	return os.Getenv("CROSSPLANE_AI_MODE") == "mock"
}

// GetMockDataDir returns the mock data directory
func GetMockDataDir() string {
	// Check command line flag first
	if mockDir := viper.GetString("mock-data-dir"); mockDir != "" {
		return mockDir
	}
	// Fall back to environment variable for backward compatibility
	mockDir := os.Getenv("CROSSPLANE_AI_MOCK_DATA_DIR")
	if mockDir == "" {
		// Return empty string to indicate we should use embedded data
		return ""
	}
	return mockDir
}
