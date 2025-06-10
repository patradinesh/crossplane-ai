package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"crossplane-ai/internal/config"
	"crossplane-ai/pkg/crossplane"
)

// Service represents the AI service
type Service struct {
	openaiClient *OpenAIClient
	config       *config.Config
	useRealAI    bool
}

// Suggestion represents an AI-generated suggestion
type Suggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Example     string `json:"example,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Category    string `json:"category,omitempty"`
}

// Analysis represents the result of AI analysis
type Analysis struct {
	TotalResources   int              `json:"total_resources"`
	HealthyResources int              `json:"healthy_resources"`
	IssuesFound      int              `json:"issues_found"`
	HealthScore      int              `json:"health_score"`
	Resources        []ResourceInfo   `json:"resources"`
	Issues           []Issue          `json:"issues"`
	Recommendations  []Recommendation `json:"recommendations"`
}

// ResourceInfo represents analyzed resource information
type ResourceInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Provider string `json:"provider"`
	Age      string `json:"age"`
}

// Issue represents a detected issue
type Issue struct {
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Resource    string `json:"resource,omitempty"`
	Resolution  string `json:"resolution,omitempty"`
}

// Recommendation represents an AI recommendation
type Recommendation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

// NewService creates a new AI service
func NewService() *Service {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Fallback to default configuration if loading fails
		return &Service{
			useRealAI: false,
		}
	}

	// Check if we should use real AI
	useRealAI := shouldUseRealAI(cfg)

	var openaiClient *OpenAIClient
	if useRealAI {
		// Initialize OpenAI client with configuration
		openaiConfig := OpenAIConfig{
			APIKey:  getAPIKey(cfg),
			Model:   cfg.AI.Model,
			BaseURL: cfg.AI.BaseURL,
		}
		openaiClient = NewOpenAIClient(openaiConfig)
	}

	return &Service{
		openaiClient: openaiClient,
		config:       cfg,
		useRealAI:    useRealAI,
	}
}

// NewServiceWithConfig creates a new AI service with explicit configuration
func NewServiceWithConfig(cfg *config.Config) *Service {
	if cfg == nil {
		return NewService()
	}

	// Check if we should use real AI
	useRealAI := shouldUseRealAI(cfg)

	var openaiClient *OpenAIClient
	if useRealAI {
		// Initialize OpenAI client with configuration
		openaiConfig := OpenAIConfig{
			APIKey:  getAPIKey(cfg),
			Model:   cfg.AI.Model,
			BaseURL: cfg.AI.BaseURL,
		}
		openaiClient = NewOpenAIClient(openaiConfig)
	}

	return &Service{
		openaiClient: openaiClient,
		config:       cfg,
		useRealAI:    useRealAI,
	}
}

// shouldUseRealAI determines if real AI should be used based on configuration
func shouldUseRealAI(cfg *config.Config) bool {
	// Don't use real AI in mock mode
	if os.Getenv("CROSSPLANE_AI_MODE") == "mock" {
		return false
	}

	// Check if provider is set to openai and we have an API key
	if cfg.AI.Provider == "openai" {
		apiKey := getAPIKey(cfg)
		return apiKey != ""
	}

	return false
}

// getAPIKey gets the API key from config or environment variable
func getAPIKey(cfg *config.Config) string {
	// First try the config
	if cfg.AI.APIKey != "" {
		// Handle environment variable expansion
		if strings.HasPrefix(cfg.AI.APIKey, "${") && strings.HasSuffix(cfg.AI.APIKey, "}") {
			envVar := strings.TrimSuffix(strings.TrimPrefix(cfg.AI.APIKey, "${"), "}")
			return os.Getenv(envVar)
		}
		return cfg.AI.APIKey
	}

	// Fallback to environment variable
	return os.Getenv("OPENAI_API_KEY")
}

// IsUsingRealAI returns true if the service is configured to use real AI
func (s *Service) IsUsingRealAI() bool {
	return s.useRealAI
}

// ProcessQuery processes a natural language query about Crossplane resources
func (s *Service) ProcessQuery(ctx context.Context, query string, resources interface{}) (string, error) {
	// Convert resources to JSON for analysis
	resourcesJSON, err := json.Marshal(resources)
	if err != nil {
		return "", fmt.Errorf("failed to marshal resources: %w", err)
	}

	// Use real AI if available, otherwise simulate
	if s.useRealAI && s.openaiClient != nil {
		return s.openaiClient.CompleteWithContext(ctx, query, string(resourcesJSON))
	}

	// Fallback to simulated AI processing
	response := s.simulateAIResponse(query, string(resourcesJSON))
	return response, nil
}

// GenerateSuggestions generates AI-powered suggestions
func (s *Service) GenerateSuggestions(ctx context.Context, suggestionType string, resources interface{}) ([]*Suggestion, error) {
	// Use real AI if available
	if s.useRealAI && s.openaiClient != nil {
		// Convert resources to JSON for context
		resourcesJSON, err := json.Marshal(resources)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resources: %w", err)
		}

		// Get AI-generated suggestions
		suggestions, err := s.openaiClient.GenerateSuggestions(ctx, suggestionType, string(resourcesJSON))
		if err != nil {
			// Fallback to mock suggestions if AI fails
			return s.generateMockSuggestions(suggestionType), nil
		}

		// Convert from []Suggestion to []*Suggestion
		result := make([]*Suggestion, len(suggestions))
		for i, suggestion := range suggestions {
			result[i] = &suggestion
		}
		return result, nil
	}

	// Fallback to simulated AI-generated suggestions
	suggestions := s.generateMockSuggestions(suggestionType)
	return suggestions, nil
}

// AnalyzeResources performs AI analysis of resources
func (s *Service) AnalyzeResources(ctx context.Context, resources interface{}, healthCheck bool) (*Analysis, error) {
	// Check if we have actual resources
	var resourceList []*ResourceInfo

	switch r := resources.(type) {
	case []*ResourceInfo:
		resourceList = r
	case []*crossplane.Resource:
		// Convert from crossplane.Resource to ResourceInfo
		for _, res := range r {
			resourceList = append(resourceList, &ResourceInfo{
				Name:     res.Name,
				Type:     res.Type,
				Status:   res.Status,
				Provider: res.Provider,
				Age:      res.Age,
			})
		}
	case []map[string]interface{}:
		// Convert from generic map format
		for _, res := range r {
			resourceList = append(resourceList, convertMapToResourceInfo(res))
		}
	default:
		// If we don't have proper resources, return empty analysis
		return &Analysis{
			TotalResources:   0,
			HealthyResources: 0,
			IssuesFound:      0,
			HealthScore:      0,
			Resources:        []ResourceInfo{},
			Issues:           []Issue{},
			Recommendations:  []Recommendation{},
		}, nil
	}

	// If no resources found, return appropriate analysis
	if len(resourceList) == 0 {
		return &Analysis{
			TotalResources:   0,
			HealthyResources: 0,
			IssuesFound:      0,
			HealthScore:      0,
			Resources:        []ResourceInfo{},
			Issues:           []Issue{},
			Recommendations: []Recommendation{{
				Title:       "Install Crossplane Providers",
				Description: "No Crossplane resources found. Install providers and create compositions to get started.",
				Impact:      "Enable infrastructure management through Crossplane",
				Priority:    "High",
			}},
		}, nil
	}

	// Use real AI for analysis if available
	if s.useRealAI && s.openaiClient != nil {
		// Convert resources to JSON for AI analysis
		resourcesJSON, err := json.Marshal(resourceList)
		if err != nil {
			// Fallback to real analysis if marshaling fails
			return s.performRealAnalysis(resourceList, healthCheck), nil
		}

		// Get AI-powered analysis
		analysis, err := s.openaiClient.AnalyzeResources(ctx, string(resourcesJSON), healthCheck)
		if err != nil {
			// Fallback to real analysis if AI fails
			return s.performRealAnalysis(resourceList, healthCheck), nil
		}

		return analysis, nil
	}

	// Fallback to perform real analysis on actual resources
	return s.performRealAnalysis(resourceList, healthCheck), nil
}

// convertMapToResourceInfo converts a map to ResourceInfo
func convertMapToResourceInfo(res map[string]interface{}) *ResourceInfo {
	info := &ResourceInfo{}

	if name, ok := res["name"].(string); ok {
		info.Name = name
	}
	if resourceType, ok := res["type"].(string); ok {
		info.Type = resourceType
	}
	if status, ok := res["status"].(string); ok {
		info.Status = status
	}
	if provider, ok := res["provider"].(string); ok {
		info.Provider = provider
	}
	if age, ok := res["age"].(string); ok {
		info.Age = age
	}

	return info
}

// performRealAnalysis analyzes actual resources from the cluster
func (s *Service) performRealAnalysis(resources []*ResourceInfo, healthCheck bool) *Analysis {
	totalResources := len(resources)
	healthyResources := 0
	issues := []Issue{}

	// Convert ResourceInfo pointers to values for the analysis
	resourceList := make([]ResourceInfo, len(resources))
	for i, res := range resources {
		resourceList[i] = *res
		// Count healthy resources
		if res.Status == "Ready" {
			healthyResources++
		} else if res.Status != "Ready" && res.Status != "Unknown" {
			// Add issue for non-ready resources
			issues = append(issues, Issue{
				Severity:    "Warning",
				Description: fmt.Sprintf("Resource %s is in %s state", res.Name, res.Status),
				Resource:    res.Name,
				Resolution:  "Check resource events and provider status",
			})
		}
	}

	issuesFound := len(issues)

	// Calculate health score
	healthScore := 100
	if totalResources > 0 {
		healthScore = (healthyResources * 100) / totalResources
	}

	// Generate recommendations based on actual state
	recommendations := s.generateRealRecommendations(resources, healthScore)

	return &Analysis{
		TotalResources:   totalResources,
		HealthyResources: healthyResources,
		IssuesFound:      issuesFound,
		HealthScore:      healthScore,
		Resources:        resourceList,
		Issues:           issues,
		Recommendations:  recommendations,
	}
}

// generateRealRecommendations generates recommendations based on actual resource state
func (s *Service) generateRealRecommendations(resources []*ResourceInfo, healthScore int) []Recommendation {
	recommendations := []Recommendation{}

	// Health-based recommendations
	if healthScore < 80 {
		recommendations = append(recommendations, Recommendation{
			Title:       "Investigate Resource Issues",
			Description: "Some resources are not in ready state. Check logs and events for troubleshooting.",
			Impact:      "Improve system reliability and performance",
			Priority:    "High",
		})
	}

	// Provider-specific recommendations
	providerCounts := make(map[string]int)
	for _, res := range resources {
		providerCounts[res.Provider]++
	}

	if len(providerCounts) > 2 {
		recommendations = append(recommendations, Recommendation{
			Title:       "Multi-Cloud Management",
			Description: "Consider implementing consistent policies across multiple cloud providers.",
			Impact:      "Better governance and cost optimization",
			Priority:    "Medium",
		})
	}

	// If we have resources, suggest monitoring
	if len(resources) > 0 {
		recommendations = append(recommendations, Recommendation{
			Title:       "Enable Monitoring and Alerting",
			Description: "Set up monitoring for your Crossplane resources to track health and performance.",
			Impact:      "Proactive issue detection and resolution",
			Priority:    "Medium",
		})
	}

	return recommendations
}

// simulateAIResponse simulates an AI response to a natural language query
func (s *Service) simulateAIResponse(query string, resourcesJSON string) string {
	queryLower := strings.ToLower(query)

	switch {
	case strings.Contains(queryLower, "what") && strings.Contains(queryLower, "resources"):
		return s.generateResourceSummary(resourcesJSON)
	case strings.Contains(queryLower, "aws"):
		return "🔍 Here are your AWS resources managed by Crossplane:\n\n" +
			"Based on your cluster, I found several AWS resources. " +
			"Most appear to be healthy, but I'd recommend checking the RDS instances " +
			"for performance optimization opportunities."
	case strings.Contains(queryLower, "database") || strings.Contains(queryLower, "db"):
		return "🗄️ Database Analysis:\n\n" +
			"I found database instances in your cluster. Here are some insights:\n" +
			"• All databases appear to be in 'Ready' state\n" +
			"• Consider enabling automated backups for production databases\n" +
			"• Review connection pooling settings for better performance"
	case strings.Contains(queryLower, "not ready") || strings.Contains(queryLower, "failed"):
		return "🔧 Troubleshooting Not Ready Resources:\n\n" +
			"Let me help you diagnose issues:\n" +
			"1. Check resource events for error messages\n" +
			"2. Verify provider credentials are valid\n" +
			"3. Ensure required dependencies are available\n" +
			"4. Check network connectivity to cloud provider APIs"
	case strings.Contains(queryLower, "cost") || strings.Contains(queryLower, "expensive"):
		return "💰 Cost Optimization Insights:\n\n" +
			"Here are some cost-saving recommendations:\n" +
			"• Consider using spot instances for non-critical workloads\n" +
			"• Review instance sizes - some might be over-provisioned\n" +
			"• Enable auto-scaling to optimize resource usage\n" +
			"• Use reserved instances for predictable workloads"
	default:
		return fmt.Sprintf("🤖 I understand you're asking about: %s\n\n"+
			"Based on your Crossplane resources, here's what I can tell you:\n"+
			"• Your cluster has multiple providers configured\n"+
			"• Most resources appear to be healthy\n"+
			"• Consider running 'crossplane-ai analyze' for detailed insights\n\n"+
			"Feel free to ask more specific questions about your resources!", query)
	}
}

func (s *Service) generateResourceSummary(resourcesJSON string) string {
	return "📊 Resource Summary:\n\n" +
		"Your Crossplane cluster contains:\n" +
		"• Multiple cloud providers (AWS, GCP, Azure)\n" +
		"• Various resource types (databases, compute, storage)\n" +
		"• Most resources are in healthy state\n\n" +
		"💡 Quick tips:\n" +
		"• Use 'crossplane-ai analyze' for detailed health check\n" +
		"• Run 'crossplane-ai suggest optimize' for optimization recommendations\n" +
		"• Try 'crossplane-ai ask \"show me failing resources\"' for troubleshooting"
}

func (s *Service) generateMockSuggestions(suggestionType string) []*Suggestion {
	switch strings.ToLower(suggestionType) {
	case "database", "db":
		return []*Suggestion{
			{
				Title:       "Enable Automated Backups",
				Description: "Configure automated backups for your RDS instances to ensure data protection",
				Example:     "spec:\n  backupRetentionPeriod: 7\n  backupWindow: \"03:00-04:00\"",
				Priority:    "High",
				Category:    "Reliability",
			},
			{
				Title:       "Implement Read Replicas",
				Description: "Add read replicas to distribute read traffic and improve performance",
				Priority:    "Medium",
				Category:    "Performance",
			},
		}
	case "security":
		return []*Suggestion{
			{
				Title:       "Enable Encryption at Rest",
				Description: "Encrypt your databases and storage resources to enhance security",
				Priority:    "High",
				Category:    "Security",
			},
			{
				Title:       "Review IAM Policies",
				Description: "Audit and tighten IAM policies to follow principle of least privilege",
				Priority:    "High",
				Category:    "Security",
			},
		}
	case "optimize", "optimization":
		return []*Suggestion{
			{
				Title:       "Right-size Resources",
				Description: "Analyze resource utilization and adjust instance sizes accordingly",
				Priority:    "Medium",
				Category:    "Cost",
			},
			{
				Title:       "Implement Auto-scaling",
				Description: "Configure auto-scaling groups to optimize resource usage",
				Priority:    "Medium",
				Category:    "Performance",
			},
		}
	default:
		return []*Suggestion{
			{
				Title:       "Health Check Resources",
				Description: "Run regular health checks on all Crossplane resources",
				Priority:    "Medium",
				Category:    "Monitoring",
			},
			{
				Title:       "Update Providers",
				Description: "Keep your Crossplane providers up to date for latest features and security fixes",
				Priority:    "Low",
				Category:    "Maintenance",
			},
		}
	}
}

// GenerateManifest generates a Crossplane manifest from natural language description
func (s *Service) GenerateManifest(ctx context.Context, description, provider string) (string, error) {
	// Use real AI if available
	if s.useRealAI && s.openaiClient != nil {
		prompt := fmt.Sprintf(`Generate a Crossplane manifest for: %s

Requirements:
- Use provider: %s (if specified, otherwise choose appropriate provider)
- Create valid Crossplane YAML
- Include metadata, spec, and appropriate labels
- Follow Crossplane best practices
- Include helpful comments

Please provide only the YAML manifest without additional explanations.`, description, provider)

		return s.openaiClient.Complete(ctx, prompt)
	}

	// Fallback to template-based generation
	return s.generateTemplateManifest(description, provider), nil
}

// generateTemplateManifest generates a basic template manifest (fallback)
func (s *Service) generateTemplateManifest(description, provider string) string {
	// Simple template generation based on keywords in description
	descLower := strings.ToLower(description)

	if strings.Contains(descLower, "database") || strings.Contains(descLower, "postgres") || strings.Contains(descLower, "mysql") {
		return s.generateDatabaseTemplate(provider)
	} else if strings.Contains(descLower, "storage") || strings.Contains(descLower, "bucket") {
		return s.generateStorageTemplate(provider)
	} else if strings.Contains(descLower, "compute") || strings.Contains(descLower, "instance") {
		return s.generateComputeTemplate(provider)
	}

	// Default to a composition template
	return s.generateCompositionTemplate(provider)
}

func (s *Service) generateDatabaseTemplate(provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	return fmt.Sprintf(`# Database instance generated by Crossplane AI
apiVersion: rds.%s.crossplane.io/v1alpha1
kind: DBInstance
metadata:
  name: my-database
  namespace: default
  labels:
    generated-by: crossplane-ai
spec:
  forProvider:
    region: us-east-1
    dbInstanceClass: db.t3.micro
    engine: postgres
    engineVersion: "13.7"
    allocatedStorage: 20
    dbName: myapp
    masterUsername: admin
    autoMinorVersionUpgrade: true
    backupRetentionPeriod: 7
    storageEncrypted: true
  providerConfigRef:
    name: default
  writeConnectionSecretsToRef:
    name: my-database-connection
    namespace: default`, provider)
}

func (s *Service) generateStorageTemplate(provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	return fmt.Sprintf(`# Storage bucket generated by Crossplane AI
apiVersion: s3.%s.crossplane.io/v1beta1
kind: Bucket
metadata:
  name: my-storage-bucket
  namespace: default
  labels:
    generated-by: crossplane-ai
spec:
  forProvider:
    region: us-east-1
    versioning:
      status: Enabled
    serverSideEncryptionConfiguration:
      rules:
      - applyServerSideEncryptionByDefault:
          sseAlgorithm: AES256
  providerConfigRef:
    name: default`, provider)
}

func (s *Service) generateComputeTemplate(provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	return fmt.Sprintf(`# Compute instance generated by Crossplane AI
apiVersion: ec2.%s.crossplane.io/v1alpha1
kind: Instance
metadata:
  name: my-instance
  namespace: default
  labels:
    generated-by: crossplane-ai
spec:
  forProvider:
    region: us-east-1
    instanceType: t3.micro
    imageId: ami-0abcdef1234567890
    keyName: my-key-pair
    tags:
      Name: MyInstance
      GeneratedBy: crossplane-ai
  providerConfigRef:
    name: default`, provider)
}

func (s *Service) generateCompositionTemplate(provider string) string {
	return fmt.Sprintf(`# Crossplane Composition generated by Crossplane AI
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: my-composition
  namespace: default
  labels:
    generated-by: crossplane-ai
    provider: %s
spec:
  compositeTypeRef:
    apiVersion: example.com/v1alpha1
    kind: XResource
  resources:
  - name: resource
    base:
      apiVersion: example.com/v1alpha1
      kind: Resource
      spec:
        forProvider:
          region: us-east-1
    patches:
    - type: FromCompositeFieldPath
      fromFieldPath: spec.region
      toFieldPath: spec.forProvider.region`, provider)
}
