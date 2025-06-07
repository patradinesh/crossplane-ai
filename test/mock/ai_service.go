package mock

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// AIService implements a mock AI service for testing
type AIService struct {
	mockDataDir string
}

// NewAIService creates a new mock AI service
func NewAIService() *AIService {
	mockDataDir := os.Getenv("CROSSPLANE_AI_MOCK_DATA_DIR")
	if mockDataDir == "" {
		mockDataDir = "./examples"
	}
	return &AIService{
		mockDataDir: mockDataDir,
	}
}

// Ask responds to a question with mock data
func (s *AIService) Ask(ctx context.Context, question string, clusterInfo map[string]interface{}) (string, error) {
	// For mock purposes, list the files in the mock data directory
	files, err := os.ReadDir(s.mockDataDir)
	if err != nil {
		return "", fmt.Errorf("failed to read mock data directory: %w", err)
	}

	response := fmt.Sprintf("Mock AI response to: %s\n\n", question)
	response += "Found the following Crossplane resources:\n"

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			response += fmt.Sprintf("- %s\n", file.Name())

			// Optionally load and parse the content of the file
			if strings.Contains(question, "detail") || strings.Contains(question, "content") {
				content, err := s.loadResourceContent(filepath.Join(s.mockDataDir, file.Name()))
				if err == nil {
					response += fmt.Sprintf("  Summary: %s\n", content)
				}
			}
		}
	}

	return response, nil
}

// GenerateResources generates mock resources
func (s *AIService) GenerateResources(ctx context.Context, prompt string, clusterInfo map[string]interface{}) (string, error) {
	resourceType := "generic"

	// Determine the resource type based on the prompt
	if strings.Contains(strings.ToLower(prompt), "database") {
		resourceType = "database"
	} else if strings.Contains(strings.ToLower(prompt), "bucket") ||
		strings.Contains(strings.ToLower(prompt), "storage") {
		resourceType = "storage"
	} else if strings.Contains(strings.ToLower(prompt), "cluster") ||
		strings.Contains(strings.ToLower(prompt), "kubernetes") {
		resourceType = "cluster"
	}

	// Generate mock YAML based on the resource type
	return s.generateMockYaml(resourceType, prompt), nil
}

// AnalyzeResources analyzes mock resources
func (s *AIService) AnalyzeResources(ctx context.Context, resources []map[string]interface{}) (string, error) {
	// Count how many resources we have
	resourceCount := len(resources)
	if resourceCount == 0 {
		// Try to load resources from the mock directory
		files, err := os.ReadDir(s.mockDataDir)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
					resourceCount++
				}
			}
		}
	}

	// Generate a mock analysis
	return fmt.Sprintf("Mock resource analysis. Found %d resources.\n\n"+
		"Health Summary:\n"+
		"- Healthy: %d\n"+
		"- Degraded: %d\n"+
		"- Failed: %d\n\n"+
		"Performance: Good\n"+
		"Security: No issues detected\n"+
		"Recommendations: Consider enabling monitoring for all resources",
		resourceCount, resourceCount-1, 0, 1), nil
}

// Helper functions

// loadResourceContent loads and summarizes the content of a YAML resource file
func (s *AIService) loadResourceContent(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Parse YAML to extract basic information
	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return "YAML file (could not parse details)", nil
	}

	kind := "Unknown"
	if k, ok := data["kind"].(string); ok {
		kind = k
	}

	name := "unnamed"
	if metadata, ok := data["metadata"].(map[interface{}]interface{}); ok {
		if n, ok := metadata["name"].(string); ok {
			name = n
		}
	}

	return fmt.Sprintf("%s '%s'", kind, name), nil
}

// generateMockYaml generates a mock YAML manifest for a resource
func (s *AIService) generateMockYaml(resourceType, prompt string) string {
	timestamp := time.Now().Format(time.RFC3339)

	switch resourceType {
	case "database":
		return fmt.Sprintf(`apiVersion: database.example.org/v1alpha1
kind: PostgreSQLInstance
metadata:
  name: db-from-prompt
  annotations:
    crossplane.io/external-name: db-%d
    timestamp: "%s"
spec:
  parameters:
    storageGB: 20
    version: "14"
    replicas: 1
    backup:
      enabled: true
      retentionDays: 7
  writeConnectionSecretToRef:
    name: db-conn
    namespace: default
`, time.Now().Unix()%1000, timestamp)

	case "storage":
		return fmt.Sprintf(`apiVersion: storage.example.org/v1alpha1
kind: Bucket
metadata:
  name: bucket-from-prompt
  annotations:
    crossplane.io/external-name: bucket-%d
    timestamp: "%s"
spec:
  parameters:
    region: us-west-2
    acl: private
    versioning: true
`, time.Now().Unix()%1000, timestamp)

	case "cluster":
		return fmt.Sprintf(`apiVersion: compute.example.org/v1alpha1
kind: KubernetesCluster
metadata:
  name: cluster-from-prompt
  annotations:
    crossplane.io/external-name: cluster-%d
    timestamp: "%s"
spec:
  parameters:
    version: "1.25"
    nodeCount: 3
    nodeSize: medium
    region: us-central1
`, time.Now().Unix()%1000, timestamp)

	default:
		return fmt.Sprintf(`apiVersion: example.org/v1alpha1
kind: Resource
metadata:
  name: resource-from-prompt
  annotations:
    crossplane.io/external-name: res-%d
    timestamp: "%s"
spec:
  parameters:
    key1: value1
    key2: value2
`, time.Now().Unix()%1000, timestamp)
	}
}
