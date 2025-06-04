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
	Short: "AI-powered Crossplane management CLI",
	Long: `Crossplane AI is a command-line tool that brings artificial intelligence
to Crossplane resource management. It allows you to interact with your Kubernetes
cluster's Crossplane resources using natural language, get intelligent suggestions,
and automate complex infrastructure operations.

Features:
- Natural language queries for Crossplane resources
- AI-powered resource discovery and analysis
- Intelligent resource recommendations
- Automated troubleshooting and optimization
- Interactive resource management`,
	Example: `  # Ask AI about your resources
  crossplane-ai ask "what AWS resources do I have?"
  
  # Get resource recommendations
  crossplane-ai suggest database
  
  # Interactive mode
  crossplane-ai interactive`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.crossplane-ai.yaml)")
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file")
	rootCmd.PersistentFlags().String("context", "", "kubernetes context to use")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
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
