package crossplane

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client represents a Crossplane client
type Client struct {
	kubeClient    kubernetes.Interface
	dynamicClient dynamic.Interface
	restConfig    *rest.Config
}

// Resource represents a Crossplane resource
type Resource struct {
	Name      string                     `json:"name"`
	Namespace string                     `json:"namespace"`
	Type      string                     `json:"type"`
	Provider  string                     `json:"provider"`
	Status    string                     `json:"status"`
	Age       string                     `json:"age"`
	Labels    map[string]string          `json:"labels,omitempty"`
	Spec      interface{}                `json:"spec,omitempty"`
	Raw       *unstructured.Unstructured `json:"-"`
}

// NewClient creates a new Crossplane client
func NewClient(ctx context.Context) (*Client, error) {
	// Try to get kubeconfig from various sources
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Build config from kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// Try in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
		}
	}

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Create dynamic client for CRDs
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Client{
		kubeClient:    kubeClient,
		dynamicClient: dynamicClient,
		restConfig:    config,
	}, nil
}

// GetAllResources retrieves all Crossplane resources
func (c *Client) GetAllResources(ctx context.Context) ([]*Resource, error) {
	var allResources []*Resource

	// Common Crossplane resource types to check
	resourceTypes := []schema.GroupVersionResource{
		// Core Crossplane resources
		{Group: "apiextensions.crossplane.io", Version: "v1", Resource: "compositions"},
		{Group: "apiextensions.crossplane.io", Version: "v1", Resource: "compositeresourcedefinitions"},
		{Group: "pkg.crossplane.io", Version: "v1", Resource: "providers"},
		{Group: "pkg.crossplane.io", Version: "v1", Resource: "configurations"},

		// AWS Provider resources (common ones)
		{Group: "rds.aws.crossplane.io", Version: "v1alpha1", Resource: "dbinstances"},
		{Group: "ec2.aws.crossplane.io", Version: "v1alpha1", Resource: "instances"},
		{Group: "s3.aws.crossplane.io", Version: "v1alpha1", Resource: "buckets"},
		{Group: "eks.aws.crossplane.io", Version: "v1alpha1", Resource: "clusters"},

		// GCP Provider resources (common ones)
		{Group: "sql.gcp.crossplane.io", Version: "v1alpha1", Resource: "databaseinstances"},
		{Group: "compute.gcp.crossplane.io", Version: "v1alpha1", Resource: "instances"},
		{Group: "storage.gcp.crossplane.io", Version: "v1alpha1", Resource: "buckets"},

		// Azure Provider resources (common ones)
		{Group: "sql.azure.crossplane.io", Version: "v1alpha1", Resource: "servers"},
		{Group: "compute.azure.crossplane.io", Version: "v1alpha1", Resource: "virtualmachines"},
		{Group: "storage.azure.crossplane.io", Version: "v1alpha1", Resource: "accounts"},
	}

	for _, gvr := range resourceTypes {
		resources, err := c.getResourcesOfType(ctx, gvr)
		if err != nil {
			// Continue with other resources even if one type fails
			continue
		}
		allResources = append(allResources, resources...)
	}

	return allResources, nil
}

// GetFilteredResources retrieves filtered Crossplane resources
func (c *Client) GetFilteredResources(ctx context.Context, name, provider, namespace string) ([]*Resource, error) {
	allResources, err := c.GetAllResources(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*Resource
	for _, resource := range allResources {
		// Apply filters
		if name != "" && resource.Name != name {
			continue
		}
		if provider != "" && resource.Provider != provider {
			continue
		}
		if namespace != "" && resource.Namespace != namespace {
			continue
		}
		filtered = append(filtered, resource)
	}

	return filtered, nil
}

func (c *Client) getResourcesOfType(ctx context.Context, gvr schema.GroupVersionResource) ([]*Resource, error) {
	list, err := c.dynamicClient.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var resources []*Resource
	for _, item := range list.Items {
		resource := c.convertToResource(&item, gvr)
		resources = append(resources, resource)
	}

	return resources, nil
}

func (c *Client) convertToResource(obj *unstructured.Unstructured, gvr schema.GroupVersionResource) *Resource {
	// Extract basic information
	name := obj.GetName()
	namespace := obj.GetNamespace()
	labels := obj.GetLabels()

	// Determine provider from group
	provider := extractProviderFromGroup(gvr.Group)

	// Get status
	status := "Unknown"
	if statusObj, found, _ := unstructured.NestedMap(obj.Object, "status"); found {
		if ready, found, _ := unstructured.NestedBool(statusObj, "ready"); found {
			if ready {
				status = "Ready"
			} else {
				status = "Not Ready"
			}
		}

		// Check for conditions
		if conditions, found, _ := unstructured.NestedSlice(statusObj, "conditions"); found {
			status = extractStatusFromConditions(conditions)
		}
	}

	// Calculate age
	age := "Unknown"
	if creationTime := obj.GetCreationTimestamp(); !creationTime.IsZero() {
		age = time.Since(creationTime.Time).Round(time.Minute).String()
	}

	// Get spec
	spec, _, _ := unstructured.NestedMap(obj.Object, "spec")

	return &Resource{
		Name:      name,
		Namespace: namespace,
		Type:      gvr.Resource,
		Provider:  provider,
		Status:    status,
		Age:       age,
		Labels:    labels,
		Spec:      spec,
		Raw:       obj,
	}
}

func extractProviderFromGroup(group string) string {
	if group == "pkg.crossplane.io" || group == "apiextensions.crossplane.io" {
		return "crossplane"
	}

	// Extract provider from group name (e.g., "rds.aws.crossplane.io" -> "aws")
	parts := strings.Split(group, ".")
	if len(parts) >= 2 {
		return parts[1]
	}

	return "unknown"
}

func extractStatusFromConditions(conditions []interface{}) string {
	for _, conditionRaw := range conditions {
		if condition, ok := conditionRaw.(map[string]interface{}); ok {
			if condType, found, _ := unstructured.NestedString(condition, "type"); found {
				if condStatus, found, _ := unstructured.NestedString(condition, "status"); found {
					if condType == "Ready" {
						if condStatus == "True" {
							return "Ready"
						} else {
							return "Not Ready"
						}
					}
				}
			}
		}
	}
	return "Unknown"
}

// GetProviders returns all installed Crossplane providers
func (c *Client) GetProviders(ctx context.Context) ([]*Resource, error) {
	gvr := schema.GroupVersionResource{Group: "pkg.crossplane.io", Version: "v1", Resource: "providers"}
	return c.getResourcesOfType(ctx, gvr)
}

// GetCompositions returns all Crossplane compositions
func (c *Client) GetCompositions(ctx context.Context) ([]*Resource, error) {
	gvr := schema.GroupVersionResource{Group: "apiextensions.crossplane.io", Version: "v1", Resource: "compositions"}
	return c.getResourcesOfType(ctx, gvr)
}
