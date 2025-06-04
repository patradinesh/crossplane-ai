package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	AI struct {
		Provider string `yaml:"provider" mapstructure:"provider"`
		APIKey   string `yaml:"api_key" mapstructure:"api_key"`
		Model    string `yaml:"model" mapstructure:"model"`
		BaseURL  string `yaml:"base_url" mapstructure:"base_url"`
	} `yaml:"ai" mapstructure:"ai"`

	Kubernetes struct {
		Kubeconfig string `yaml:"kubeconfig" mapstructure:"kubeconfig"`
		Context    string `yaml:"context" mapstructure:"context"`
		Namespace  string `yaml:"namespace" mapstructure:"namespace"`
	} `yaml:"kubernetes" mapstructure:"kubernetes"`

	Crossplane struct {
		Providers     []string `yaml:"providers" mapstructure:"providers"`
		ResourceTypes []string `yaml:"resource_types" mapstructure:"resource_types"`
	} `yaml:"crossplane" mapstructure:"crossplane"`

	CLI struct {
		OutputFormat string `yaml:"output_format" mapstructure:"output_format"`
		Verbose      bool   `yaml:"verbose" mapstructure:"verbose"`
		Color        bool   `yaml:"color" mapstructure:"color"`
	} `yaml:"cli" mapstructure:"cli"`

	Analysis struct {
		Timeout        int  `yaml:"timeout" mapstructure:"timeout"`
		MaxSuggestions int  `yaml:"max_suggestions" mapstructure:"max_suggestions"`
		Detailed       bool `yaml:"detailed" mapstructure:"detailed"`
	} `yaml:"analysis" mapstructure:"analysis"`
}

var globalConfig *Config

// Load loads the configuration from file and environment variables
func Load() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// Set default values
	setDefaults()

	// Set config name and paths
	viper.SetConfigName(".crossplane-ai")
	viper.SetConfigType("yaml")

	// Add config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Try to find and read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, use defaults
	}

	// Read environment variables
	viper.AutomaticEnv()

	// Unmarshal config
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set global config
	globalConfig = config

	return config, nil
}

// Get returns the global configuration, loading it if necessary
func Get() *Config {
	if globalConfig == nil {
		config, err := Load()
		if err != nil {
			// Return default config if loading fails
			return getDefaultConfig()
		}
		return config
	}
	return globalConfig
}

// setDefaults sets default configuration values
func setDefaults() {
	// AI defaults
	viper.SetDefault("ai.provider", "mock")
	viper.SetDefault("ai.model", "gpt-4")

	// Kubernetes defaults
	if home, err := os.UserHomeDir(); err == nil {
		viper.SetDefault("kubernetes.kubeconfig", filepath.Join(home, ".kube", "config"))
	}

	// Crossplane defaults
	viper.SetDefault("crossplane.providers", []string{"aws", "gcp", "azure", "kubernetes"})
	viper.SetDefault("crossplane.resource_types", []string{
		"compositions", "providers", "configurations",
		"dbinstances", "instances", "buckets", "clusters",
	})

	// CLI defaults
	viper.SetDefault("cli.output_format", "table")
	viper.SetDefault("cli.verbose", false)
	viper.SetDefault("cli.color", true)

	// Analysis defaults
	viper.SetDefault("analysis.timeout", 30)
	viper.SetDefault("analysis.max_suggestions", 10)
	viper.SetDefault("analysis.detailed", true)
}

// getDefaultConfig returns a default configuration
func getDefaultConfig() *Config {
	config := &Config{}
	config.AI.Provider = "mock"
	config.AI.Model = "gpt-4"

	if home, err := os.UserHomeDir(); err == nil {
		config.Kubernetes.Kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config.Crossplane.Providers = []string{"aws", "gcp", "azure", "kubernetes"}
	config.Crossplane.ResourceTypes = []string{
		"compositions", "providers", "configurations",
		"dbinstances", "instances", "buckets", "clusters",
	}

	config.CLI.OutputFormat = "table"
	config.CLI.Verbose = false
	config.CLI.Color = true

	config.Analysis.Timeout = 30
	config.Analysis.MaxSuggestions = 10
	config.Analysis.Detailed = true

	return config
}

// Save saves the current configuration to file
func Save() error {
	if globalConfig == nil {
		return fmt.Errorf("no configuration to save")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".crossplane-ai.yaml")
	return viper.WriteConfigAs(configPath)
}

// IsVerbose returns whether verbose output is enabled
func IsVerbose() bool {
	return Get().CLI.Verbose || viper.GetBool("verbose")
}

// GetOutputFormat returns the configured output format
func GetOutputFormat() string {
	format := Get().CLI.OutputFormat
	if viperFormat := viper.GetString("output"); viperFormat != "" {
		format = viperFormat
	}
	return format
}

// GetKubeconfig returns the kubeconfig path
func GetKubeconfig() string {
	kubeconfig := Get().Kubernetes.Kubeconfig
	if viperKubeconfig := viper.GetString("kubeconfig"); viperKubeconfig != "" {
		kubeconfig = viperKubeconfig
	}
	return kubeconfig
}

// GetContext returns the Kubernetes context
func GetContext() string {
	context := Get().Kubernetes.Context
	if viperContext := viper.GetString("context"); viperContext != "" {
		context = viperContext
	}
	return context
}
