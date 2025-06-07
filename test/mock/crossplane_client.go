package mock

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// CrossplaneClient implements a mock Crossplane client for testing
type CrossplaneClient struct {
	mockDataDir string
	resources   map[string][]map[string]interface{}
}

// NewCrossplaneClient creates a new mock Crossplane client
func NewCrossplaneClient() *CrossplaneClient {
	mockDataDir := os.Getenv("CROSSPLANE_AI_MOCK_DATA_DIR")
	if mockDataDir == "" {
		mockDataDir = "./examples"
	}

	client := &CrossplaneClient{
		mockDataDir: mockDataDir,
		resources:   make(map[string][]map[string]interface{}),
	}

	// Load mock resources
	client.loadMockResources()

	return client
}

// GetResources retrieves mock resources
func (c *CrossplaneClient) GetResources(ctx context.Context, resourceType string) ([]map[string]interface{}, error) {
	if resources, ok := c.resources[resourceType]; ok {
		return resources, nil
	}

	// If not found, return empty slice for this type
	return []map[string]interface{}{}, nil
}

// GetAllResources retrieves all mock resources
func (c *CrossplaneClient) GetAllResources(ctx context.Context) (map[string][]map[string]interface{}, error) {
	return c.resources, nil
}

// ApplyManifest simulates applying a manifest
func (c *CrossplaneClient) ApplyManifest(ctx context.Context, manifest string) error {
	// For mock, we simply log the manifest
	fmt.Println("Mock: Applied manifest:")
	fmt.Println("---")
	fmt.Println(manifest)
	fmt.Println("---")
	return nil
}

// Helper methods

// loadMockResources loads resources from YAML files in the mock directory
func (c *CrossplaneClient) loadMockResources() {
	files, err := os.ReadDir(c.mockDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to read mock data directory: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			filePath := filepath.Join(c.mockDataDir, file.Name())

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to read file %s: %v\n", filePath, err)
				continue
			}

			// Parse YAML
			var resource map[string]interface{}
			if err := yaml.Unmarshal(content, &resource); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse YAML in %s: %v\n", filePath, err)
				continue
			}

			// Determine resource type
			kind := "unknown"
			if k, ok := resource["kind"].(string); ok {
				kind = k
			}

			// Add resource to the map
			if _, ok := c.resources[kind]; !ok {
				c.resources[kind] = []map[string]interface{}{}
			}

			// Add mock status if not present
			c.addMockStatus(resource)

			c.resources[kind] = append(c.resources[kind], resource)
		}
	}
}

// addMockStatus adds a mock status to a resource if it doesn't have one
func (c *CrossplaneClient) addMockStatus(resource map[string]interface{}) {
	if _, ok := resource["status"]; !ok {
		// Add a mock status
		resource["status"] = map[string]interface{}{
			"conditions": []map[string]interface{}{
				{
					"type":               "Ready",
					"status":             "True",
					"lastTransitionTime": time.Now().Format(time.RFC3339),
					"reason":             "ResourceReady",
					"message":            "Resource is ready",
				},
			},
			"phase": "Bound",
		}
	}
}
