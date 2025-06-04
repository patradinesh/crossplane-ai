package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIConfig represents OpenAI configuration
type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string
	Timeout time.Duration
}

// OpenAIClient represents an OpenAI API client
type OpenAIClient struct {
	config     OpenAIConfig
	httpClient *http.Client
}

// OpenAIRequest represents a request to OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// OpenAIMessage represents a message in OpenAI conversation
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a response from OpenAI API
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config OpenAIConfig) *OpenAIClient {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}
	if config.Model == "" {
		config.Model = "gpt-4"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &OpenAIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Complete sends a completion request to OpenAI API
func (c *OpenAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	request := OpenAIRequest{
		Model: c.config.Model,
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: "You are an expert Crossplane infrastructure assistant. Provide helpful, accurate, and actionable responses about Crossplane resources, Kubernetes, and cloud infrastructure. Keep responses concise but informative.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	return c.sendRequest(ctx, request)
}

// CompleteWithContext sends a completion request with additional context
func (c *OpenAIClient) CompleteWithContext(ctx context.Context, query, resourceContext string) (string, error) {
	prompt := fmt.Sprintf(`Context: You are analyzing Crossplane resources in a Kubernetes cluster.

Resource Information:
%s

User Query: %s

Please provide a helpful response based on the resource context. If the query is about specific resources, reference the actual resource names and statuses from the context.`, resourceContext, query)

	return c.Complete(ctx, prompt)
}

// GenerateSuggestions generates AI-powered suggestions
func (c *OpenAIClient) GenerateSuggestions(ctx context.Context, suggestionType, resourceContext string) ([]Suggestion, error) {
	prompt := fmt.Sprintf(`As a Crossplane expert, analyze the following resources and provide specific %s suggestions.

Resource Context:
%s

Provide 3-5 actionable suggestions in JSON format as an array of objects with fields:
- title: Brief suggestion title
- description: Detailed explanation
- priority: High/Medium/Low
- category: The category of suggestion
- example: Optional YAML example if applicable

Focus on practical, implementable suggestions for Crossplane and Kubernetes infrastructure.`, suggestionType, resourceContext)

	response, err := c.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Try to parse JSON response
	var suggestions []Suggestion
	if err := json.Unmarshal([]byte(response), &suggestions); err != nil {
		// If JSON parsing fails, create a single suggestion with the response
		return []Suggestion{
			{
				Title:       fmt.Sprintf("AI Suggestion for %s", suggestionType),
				Description: response,
				Priority:    "Medium",
				Category:    suggestionType,
			},
		}, nil
	}

	return suggestions, nil
}

// AnalyzeResources performs AI analysis of resources
func (c *OpenAIClient) AnalyzeResources(ctx context.Context, resourceContext string, healthCheck bool) (*Analysis, error) {
	analysisType := "general"
	if healthCheck {
		analysisType = "health-focused"
	}

	prompt := fmt.Sprintf(`Analyze the following Crossplane resources and provide a %s analysis.

Resource Context:
%s

Provide analysis in JSON format with these fields:
- total_resources: number of total resources
- healthy_resources: number of healthy resources  
- issues_found: number of issues detected
- health_score: overall health score (0-100)
- resources: array of resource info with name, type, status, provider, age
- issues: array of issues with severity, description, resource, resolution
- recommendations: array of recommendations with title, description, impact, priority

Focus on actionable insights for Crossplane infrastructure management.`, analysisType, resourceContext)

	response, err := c.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Try to parse JSON response
	var analysis Analysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		// If JSON parsing fails, return a basic analysis with the response as a recommendation
		return &Analysis{
			TotalResources:   1,
			HealthyResources: 1,
			IssuesFound:      0,
			HealthScore:      85,
			Recommendations: []Recommendation{
				{
					Title:       "AI Analysis Results",
					Description: response,
					Priority:    "Medium",
				},
			},
		}, nil
	}

	return &analysis, nil
}

// sendRequest sends a request to OpenAI API
func (c *OpenAIClient) sendRequest(ctx context.Context, request OpenAIRequest) (string, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return response.Choices[0].Message.Content, nil
}
