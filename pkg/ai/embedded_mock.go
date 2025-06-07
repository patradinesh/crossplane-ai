package ai

// EmbeddedMockData contains hardcoded mock data that doesn't require external files
// This ensures mock mode works even when users download just the binary

// GetEmbeddedMockResources returns a consistent set of mock resources for testing
func GetEmbeddedMockResources() []*ResourceInfo {
	return []*ResourceInfo{
		{
			Name:     "sample-database-composition",
			Type:     "compositions",
			Status:   "Ready",
			Provider: "crossplane",
			Age:      "2h",
		},
		{
			Name:     "xdatabases.example.org",
			Type:     "compositeresourcedefinitions",
			Status:   "Ready",
			Provider: "crossplane",
			Age:      "2h",
		},
		{
			Name:     "provider-aws",
			Type:     "providers",
			Status:   "Ready",
			Provider: "crossplane",
			Age:      "1h",
		},
		{
			Name:     "provider-gcp",
			Type:     "providers",
			Status:   "Ready",
			Provider: "crossplane",
			Age:      "1h",
		},
		{
			Name:     "provider-azure",
			Type:     "providers",
			Status:   "Ready",
			Provider: "crossplane",
			Age:      "1h",
		},
		{
			Name:     "sample-database-instance",
			Type:     "dbinstances",
			Status:   "Ready",
			Provider: "aws",
			Age:      "30m",
		},
		{
			Name:     "web-server-instance",
			Type:     "instances",
			Status:   "Ready",
			Provider: "aws",
			Age:      "45m",
		},
		{
			Name:     "data-storage-bucket",
			Type:     "buckets",
			Status:   "Ready",
			Provider: "aws",
			Age:      "1h",
		},
		{
			Name:     "gcp-database-instance",
			Type:     "databaseinstances",
			Status:   "Ready",
			Provider: "gcp",
			Age:      "20m",
		},
		{
			Name:     "azure-storage-account",
			Type:     "accounts",
			Status:   "Ready",
			Provider: "azure",
			Age:      "35m",
		},
		{
			Name:     "failing-test-resource",
			Type:     "instances",
			Status:   "Not Ready",
			Provider: "aws",
			Age:      "5m",
		},
	}
}

// GetEmbeddedMockYAMLExamples returns example YAML manifests
func GetEmbeddedMockYAMLExamples() map[string]string {
	return map[string]string{
		"composition": `apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: xdatabases.example.org
  labels:
    provider: aws
    service: rds
spec:
  compositeTypeRef:
    apiVersion: example.org/v1alpha1
    kind: XDatabase
  resources:
    - name: rds-instance
      base:
        apiVersion: rds.aws.crossplane.io/v1alpha1
        kind: DBInstance
        spec:
          forProvider:
            dbInstanceClass: db.t3.micro
            engine: postgres
            engineVersion: "13.7"
            allocatedStorage: 20
            storageType: gp2
      patches:
        - type: FromCompositeFieldPath
          fromFieldPath: spec.parameters.storageGB
          toFieldPath: spec.forProvider.allocatedStorage`,

		"xrd": `apiVersion: apiextensions.crossplane.io/v1
kind: CompositeResourceDefinition
metadata:
  name: xdatabases.example.org
spec:
  group: example.org
  names:
    kind: XDatabase
    plural: xdatabases
  versions:
  - name: v1alpha1
    served: true
    referenceable: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              parameters:
                type: object
                properties:
                  storageGB:
                    type: integer
                    default: 20
                required:
                - storageGB
            required:
            - parameters`,

		"claim": `apiVersion: example.org/v1alpha1
kind: XDatabase
metadata:
  name: my-database
spec:
  parameters:
    storageGB: 50
  compositionRef:
    name: xdatabases.example.org`,

		"provider": `apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws
spec:
  package: xpkg.upbound.io/crossplane-contrib/provider-aws:v0.44.0`,
	}
}

// MockScenarios provides different mock scenarios for demonstrations
var MockScenarios = map[string][]*ResourceInfo{
	"healthy": {
		{Name: "web-app-db", Type: "dbinstances", Status: "Ready", Provider: "aws", Age: "2h"},
		{Name: "web-app-server", Type: "instances", Status: "Ready", Provider: "aws", Age: "2h"},
		{Name: "web-app-bucket", Type: "buckets", Status: "Ready", Provider: "aws", Age: "2h"},
	},
	"mixed-health": {
		{Name: "healthy-db", Type: "dbinstances", Status: "Ready", Provider: "aws", Age: "1h"},
		{Name: "failing-server", Type: "instances", Status: "Not Ready", Provider: "aws", Age: "30m"},
		{Name: "pending-bucket", Type: "buckets", Status: "Creating", Provider: "aws", Age: "5m"},
	},
	"multi-cloud": {
		{Name: "aws-database", Type: "dbinstances", Status: "Ready", Provider: "aws", Age: "3h"},
		{Name: "gcp-compute", Type: "instances", Status: "Ready", Provider: "gcp", Age: "2h"},
		{Name: "azure-storage", Type: "accounts", Status: "Ready", Provider: "azure", Age: "1h"},
	},
}
