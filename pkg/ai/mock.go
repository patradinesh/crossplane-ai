package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// MockAIService implements a mock AI service for testing
type MockAIService struct {
	mockDataDir string
}

// NewMockService creates a new mock AI service
func NewMockService() *MockAIService {
	mockDataDir := os.Getenv("CROSSPLANE_AI_MOCK_DATA_DIR")
	if mockDataDir == "" {
		mockDataDir = "./examples"
	}
	return &MockAIService{
		mockDataDir: mockDataDir,
	}
}

// Ask responds to a question with mock data
func (s *MockAIService) Ask(ctx context.Context, question string, clusterInfo map[string]interface{}) (string, error) {
	// For mock purposes, just list the files in the mock data directory
	files, err := os.ReadDir(s.mockDataDir)
	if err != nil {
		return "", fmt.Errorf("failed to read mock data directory: %w", err)
	}

	response := fmt.Sprintf("Mock AI response to: %s\n\n", question)
	response += "Found the following Crossplane resources:\n"

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			response += fmt.Sprintf("- %s\n", file.Name())
		}
	}

	return response, nil
}

// GenerateResources generates mock resources
func (s *MockAIService) GenerateResources(ctx context.Context, prompt string, clusterInfo map[string]interface{}) (string, error) {
	return fmt.Sprintf("Mock resource generation for: %s\n\nHere's a sample YAML manifest for a Crossplane resource.", prompt), nil
}

// AnalyzeResources analyzes mock resources
func (s *MockAIService) AnalyzeResources(ctx context.Context, resources []map[string]interface{}) (string, error) {
	return "Mock resource analysis. Found 3 resources, all appear to be healthy.", nil
}
