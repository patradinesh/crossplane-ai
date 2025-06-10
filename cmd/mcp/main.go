package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"crossplane-ai/pkg/ai"
	"crossplane-ai/pkg/crossplane"
)

// MCP Server for Crossplane AI
// Implements the Model Context Protocol to expose Crossplane AI capabilities
// to Claude Desktop and other MCP-compatible applications

type MCPServer struct {
	aiService        *ai.Service
	crossplaneClient *crossplane.Client
}

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents an MCP response
type MCPResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool definitions for MCP
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Resource definitions for MCP
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType,omitempty"`
}

func NewMCPServer() *MCPServer {
	// Initialize AI service
	aiService := ai.NewService()

	// Initialize Crossplane client
	ctx := context.Background()
	crossplaneClient, err := crossplane.NewClient(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize Crossplane client: %v", err)
		// Continue without Crossplane client for demo purposes
	}

	return &MCPServer{
		aiService:        aiService,
		crossplaneClient: crossplaneClient,
	}
}

func (s *MCPServer) handleRequest(request MCPRequest) MCPResponse {
	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "tools/list":
		return s.handleToolsList(request)
	case "tools/call":
		return s.handleToolsCall(request)
	case "resources/list":
		return s.handleResourcesList(request)
	case "resources/read":
		return s.handleResourcesRead(request)
	default:
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func (s *MCPServer) handleInitialize(request MCPRequest) MCPResponse {
	capabilities := map[string]interface{}{
		"tools": map[string]interface{}{
			"listChanged": false,
		},
		"resources": map[string]interface{}{
			"subscribe":   false,
			"listChanged": false,
		},
	}

	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    capabilities,
		"serverInfo": map[string]interface{}{
			"name":    "crossplane-ai-mcp",
			"version": "1.0.0",
		},
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

func (s *MCPServer) handleToolsList(request MCPRequest) MCPResponse {
	tools := []MCPTool{
		{
			Name:        "crossplane_ask",
			Description: "Ask questions about Crossplane resources using natural language",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"question": map[string]interface{}{
						"type":        "string",
						"description": "Natural language question about Crossplane resources",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"description": "Optional: Filter by specific provider (aws, gcp, azure)",
					},
					"namespace": map[string]interface{}{
						"type":        "string",
						"description": "Optional: Filter by Kubernetes namespace",
					},
				},
				"required": []string{"question"},
			},
		},
		{
			Name:        "crossplane_analyze",
			Description: "Perform AI-powered analysis of Crossplane resources",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"health_check": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to perform health check analysis",
						"default":     true,
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"description": "Optional: Filter by specific provider (aws, gcp, azure)",
					},
				},
			},
		},
		{
			Name:        "crossplane_suggest",
			Description: "Get AI-powered suggestions for Crossplane resources",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"suggestion_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of suggestions (optimization, security, database, etc.)",
						"default":     "optimization",
					},
				},
			},
		},
		{
			Name:        "crossplane_generate",
			Description: "Generate Crossplane manifests from natural language descriptions",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Natural language description of the infrastructure to create",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"description": "Target cloud provider (aws, gcp, azure)",
						"default":     "aws",
					},
				},
				"required": []string{"description"},
			},
		},
		{
			Name:        "crossplane_list_resources",
			Description: "List all Crossplane resources in the cluster",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"resource_type": map[string]interface{}{
						"type":        "string",
						"description": "Optional: Filter by resource type (compositions, providers, etc.)",
					},
					"provider": map[string]interface{}{
						"type":        "string",
						"description": "Optional: Filter by provider (aws, gcp, azure)",
					},
				},
			},
		},
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func (s *MCPServer) handleToolsCall(request MCPRequest) MCPResponse {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Tool name is required",
			},
		}
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	ctx := context.Background()

	switch toolName {
	case "crossplane_ask":
		return s.handleCrossplaneAsk(request, ctx, arguments)
	case "crossplane_analyze":
		return s.handleCrossplaneAnalyze(request, ctx, arguments)
	case "crossplane_suggest":
		return s.handleCrossplaneSuggest(request, ctx, arguments)
	case "crossplane_generate":
		return s.handleCrossplaneGenerate(request, ctx, arguments)
	case "crossplane_list_resources":
		return s.handleCrossplaneListResources(request, ctx, arguments)
	default:
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Unknown tool",
			},
		}
	}
}

func (s *MCPServer) handleCrossplaneAsk(request MCPRequest, ctx context.Context, args map[string]interface{}) MCPResponse {
	question, ok := args["question"].(string)
	if !ok {
		return s.errorResponse(request.ID, -32602, "Question is required")
	}

	// Get resources for context
	resources, err := s.getResources(args)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Failed to get resources: %v", err))
	}

	// Process query with AI
	response, err := s.aiService.ProcessQuery(ctx, question, resources)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("AI processing failed: %v", err))
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": response,
				},
			},
		},
	}
}

func (s *MCPServer) handleCrossplaneAnalyze(request MCPRequest, ctx context.Context, args map[string]interface{}) MCPResponse {
	healthCheck := true
	if hc, ok := args["health_check"].(bool); ok {
		healthCheck = hc
	}

	// Get resources for analysis
	resources, err := s.getResources(args)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Failed to get resources: %v", err))
	}

	// Perform analysis
	analysis, err := s.aiService.AnalyzeResources(ctx, resources, healthCheck)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Analysis failed: %v", err))
	}

	// Format analysis results
	result := fmt.Sprintf(`📊 Crossplane Analysis Results

🔍 Overview:
• Total Resources: %d
• Healthy Resources: %d
• Issues Found: %d
• Health Score: %d/100

`, analysis.TotalResources, analysis.HealthyResources, analysis.IssuesFound, analysis.HealthScore)

	if len(analysis.Issues) > 0 {
		result += "⚠️ Issues Detected:\n"
		for _, issue := range analysis.Issues {
			result += fmt.Sprintf("• %s: %s\n", issue.Severity, issue.Description)
		}
		result += "\n"
	}

	if len(analysis.Recommendations) > 0 {
		result += "💡 Recommendations:\n"
		for _, rec := range analysis.Recommendations {
			result += fmt.Sprintf("• %s: %s\n", rec.Title, rec.Description)
		}
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

func (s *MCPServer) handleCrossplaneSuggest(request MCPRequest, ctx context.Context, args map[string]interface{}) MCPResponse {
	suggestionType := "optimization"
	if st, ok := args["suggestion_type"].(string); ok {
		suggestionType = st
	}

	// Get resources for context
	resources, err := s.getResources(args)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Failed to get resources: %v", err))
	}

	// Generate suggestions
	suggestions, err := s.aiService.GenerateSuggestions(ctx, suggestionType, resources)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Suggestion generation failed: %v", err))
	}

	// Format suggestions
	result := fmt.Sprintf("💡 Crossplane %s Suggestions:\n\n", suggestionType)
	for i, suggestion := range suggestions {
		result += fmt.Sprintf("%d. %s\n", i+1, suggestion.Title)
		result += fmt.Sprintf("   %s\n", suggestion.Description)
		if suggestion.Priority != "" {
			result += fmt.Sprintf("   Priority: %s\n", suggestion.Priority)
		}
		result += "\n"
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

func (s *MCPServer) handleCrossplaneGenerate(request MCPRequest, ctx context.Context, args map[string]interface{}) MCPResponse {
	description, ok := args["description"].(string)
	if !ok {
		return s.errorResponse(request.ID, -32602, "Description is required")
	}

	provider := "aws"
	if p, ok := args["provider"].(string); ok {
		provider = p
	}

	// Generate manifest
	manifest, err := s.aiService.GenerateManifest(ctx, description, provider)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Manifest generation failed: %v", err))
	}

	result := fmt.Sprintf("📝 Generated Crossplane Manifest:\n\n```yaml\n%s\n```", manifest)

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

func (s *MCPServer) handleCrossplaneListResources(request MCPRequest, ctx context.Context, args map[string]interface{}) MCPResponse {
	// Get resources
	resources, err := s.getResources(args)
	if err != nil {
		return s.errorResponse(request.ID, -32603, fmt.Sprintf("Failed to get resources: %v", err))
	}

	// Format resource list
	result := "📋 Crossplane Resources:\n\n"
	if resourceList, ok := resources.([]interface{}); ok {
		for i, resource := range resourceList {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				if name, ok := resourceMap["name"].(string); ok {
					result += fmt.Sprintf("%d. %s", i+1, name)
					if kind, ok := resourceMap["kind"].(string); ok {
						result += fmt.Sprintf(" (%s)", kind)
					}
					if status, ok := resourceMap["status"].(string); ok {
						result += fmt.Sprintf(" - %s", status)
					}
					result += "\n"
				}
			}
		}
	} else {
		result += "No resources found or mock data being used.\n"
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
	}
}

func (s *MCPServer) handleResourcesList(request MCPRequest) MCPResponse {
	resources := []MCPResource{
		{
			URI:         "crossplane://cluster/resources",
			Name:        "Crossplane Cluster Resources",
			Description: "All Crossplane resources in the cluster",
			MimeType:    "application/json",
		},
		{
			URI:         "crossplane://cluster/providers",
			Name:        "Crossplane Providers",
			Description: "Installed Crossplane providers",
			MimeType:    "application/json",
		},
		{
			URI:         "crossplane://cluster/compositions",
			Name:        "Crossplane Compositions",
			Description: "Available Crossplane compositions",
			MimeType:    "application/json",
		},
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"resources": resources,
		},
	}
}

func (s *MCPServer) handleResourcesRead(request MCPRequest) MCPResponse {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		return s.errorResponse(request.ID, -32602, "Invalid params")
	}

	uri, ok := params["uri"].(string)
	if !ok {
		return s.errorResponse(request.ID, -32602, "URI is required")
	}

	// Get resource data based on URI
	var content string
	switch uri {
	case "crossplane://cluster/resources":
		resources, err := s.getResources(make(map[string]interface{}))
		if err != nil {
			return s.errorResponse(request.ID, -32603, fmt.Sprintf("Failed to get resources: %v", err))
		}
		resourcesJSON, _ := json.MarshalIndent(resources, "", "  ")
		content = string(resourcesJSON)
	case "crossplane://cluster/providers":
		content = `{
  "providers": [
    {"name": "provider-aws", "status": "Ready"},
    {"name": "provider-gcp", "status": "Ready"},
    {"name": "provider-azure", "status": "Ready"}
  ]
}`
	case "crossplane://cluster/compositions":
		content = `{
  "compositions": [
    {"name": "sample-database-composition", "status": "Ready"}
  ]
}`
	default:
		return s.errorResponse(request.ID, -32602, "Unknown resource URI")
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"uri":      uri,
					"mimeType": "application/json",
					"text":     content,
				},
			},
		},
	}
}

func (s *MCPServer) getResources(args map[string]interface{}) (interface{}, error) {
	// Try to get real resources from Crossplane client
	if s.crossplaneClient != nil {
		provider := ""
		if p, ok := args["provider"].(string); ok {
			provider = p
		}

		namespace := ""
		if n, ok := args["namespace"].(string); ok {
			namespace = n
		}

		resources, err := s.crossplaneClient.GetFilteredResources(context.Background(), "", provider, namespace)
		if err == nil && len(resources) > 0 {
			return resources, nil
		}
	}

	// Fallback to mock data
	return ai.GetEmbeddedMockResources(), nil
}

func (s *MCPServer) errorResponse(id interface{}, code int, message string) MCPResponse {
	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
}

func main() {
	server := NewMCPServer()

	log.Println("Starting Crossplane AI MCP Server...")
	log.Println("Reading JSON-RPC requests from stdin...")

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request MCPRequest
		if err := decoder.Decode(&request); err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Error decoding request: %v", err)
			continue
		}

		response := server.handleRequest(request)

		if err := encoder.Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}
