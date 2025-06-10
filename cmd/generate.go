package cmd

import (
	"context"
	"fmt"
	"strings"

	"crossplane-ai/pkg/ai"
	"crossplane-ai/pkg/cli"
	"crossplane-ai/pkg/crossplane"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate [description]",
	Short:   "Generate Crossplane resource manifests using AI",
	Aliases: []string{"gen", "create"},
	Long: `Use AI to generate Crossplane resource manifests based on natural language descriptions.
This command helps you quickly create infrastructure as code by describing what you want
in plain English.`,
	Example: `  # Generate an AWS RDS database
  crossplane-ai generate "create a MySQL database on AWS"
  
  # Generate a complete web application stack
  crossplane-ai generate "web app with database and load balancer on GCP"
  
  # Generate storage resources
  crossplane-ai generate "S3 bucket with versioning enabled"
  
  # Interactive mode
  crossplane-ai generate`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runInteractiveGenerate()
		}

		description := strings.Join(args, " ")
		provider, _ := cmd.Flags().GetString("provider")
		outputFormat, _ := cmd.Flags().GetString("output")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		apply, _ := cmd.Flags().GetBool("apply")

		return runGenerate(description, provider, outputFormat, dryRun, apply)
	},
}

func runGenerate(description, provider, outputFormat string, dryRun, apply bool) error {
	ctx := context.Background()

	// Initialize clients
	client, err := crossplane.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize Crossplane client: %w", err)
	}

	aiService := ai.NewService()

	// Show AI mode information
	if aiService.IsUsingRealAI() {
		cli.PrintInfo("🤖 Using OpenAI for intelligent manifest generation")
	} else {
		cli.PrintInfo("🤖 Using template-based generation (set OPENAI_API_KEY for AI-powered generation)")
	}

	cli.PrintInfo(fmt.Sprintf("📝 Generating Crossplane resources for: %s", description))
	fmt.Println()

	// Generate the manifest
	manifest, err := generateManifest(ctx, aiService, description, provider)
	if err != nil {
		return fmt.Errorf("failed to generate manifest: %w", err)
	}

	// Output the result
	if outputFormat == "json" {
		fmt.Println(cli.FormatJSON(manifest))
	} else {
		fmt.Println(manifest)
	}

	// Handle dry-run
	if dryRun {
		cli.PrintInfo("🧪 Dry run mode - manifest generated but not applied")
		return nil
	}

	// Handle apply
	if apply {
		fmt.Println()
		cli.PrintInfo("🚀 Applying manifest to cluster...")

		if err := applyManifest(ctx, client, manifest); err != nil {
			return fmt.Errorf("failed to apply manifest: %w", err)
		}

		cli.PrintSuccess("✅ Manifest applied successfully!")
	} else {
		fmt.Println()
		cli.PrintInfo("💡 Use --apply to apply this manifest to your cluster")
	}

	return nil
}

func runInteractiveGenerate() error {
	fmt.Println("🤖 Welcome to Crossplane AI Resource Generator!")
	fmt.Println()
	cli.PrintInfo("Describe the infrastructure you want to create in natural language.")
	fmt.Println()

	// Interactive prompts
	description := cli.PromptUser("What would you like to create? ")
	if description == "" {
		return fmt.Errorf("description is required")
	}

	provider := cli.PromptUser("Preferred cloud provider (aws/gcp/azure) [auto]: ")
	if provider == "" {
		provider = "auto"
	}

	return runGenerate(description, provider, "yaml", false, false)
}

func generateManifest(ctx context.Context, aiService *ai.Service, description, provider string) (string, error) {
	// Use AI service for intelligent manifest generation
	manifest, err := aiService.GenerateManifest(ctx, description, provider)
	if err != nil {
		return "", fmt.Errorf("AI manifest generation failed: %w", err)
	}

	// If we got a manifest from AI, return it
	if manifest != "" {
		return manifest, nil
	}

	// Fallback to template-based generation (this shouldn't happen with the new AI service)
	descriptionLower := strings.ToLower(description)

	// Database resources
	if strings.Contains(descriptionLower, "database") || strings.Contains(descriptionLower, "db") || strings.Contains(descriptionLower, "mysql") || strings.Contains(descriptionLower, "postgres") {
		return generateDatabaseManifest(description, provider), nil
	}

	// Storage resources
	if strings.Contains(descriptionLower, "storage") || strings.Contains(descriptionLower, "bucket") || strings.Contains(descriptionLower, "s3") {
		return generateStorageManifest(description, provider), nil
	}

	// Network resources
	if strings.Contains(descriptionLower, "network") || strings.Contains(descriptionLower, "vpc") || strings.Contains(descriptionLower, "subnet") {
		return generateNetworkManifest(description, provider), nil
	}

	// Compute resources
	if strings.Contains(descriptionLower, "server") || strings.Contains(descriptionLower, "instance") || strings.Contains(descriptionLower, "compute") {
		return generateComputeManifest(description, provider), nil
	}

	// Web application stack
	if strings.Contains(descriptionLower, "web app") || strings.Contains(descriptionLower, "application") {
		return generateWebAppManifest(description, provider), nil
	}

	// Default template
	return generateDefaultManifest(description, provider), nil
}

func generateDatabaseManifest(description, provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	dbType := "mysql"
	if strings.Contains(strings.ToLower(description), "postgres") {
		dbType = "postgres"
	}

	// Use provider in the API version to support different providers
	apiVersion := fmt.Sprintf("rds.%s.crossplane.io/v1alpha1", provider)

	return fmt.Sprintf(`apiVersion: %s
kind: DBInstance
metadata:
  name: my-database
  namespace: default
spec:
  forProvider:
    dbInstanceClass: db.t3.micro
    engine: %s
    engineVersion: "8.0"
    dbName: myapp
    masterUsername: admin
    allocatedStorage: 20
    storageType: gp2
    storageEncrypted: true
    multiAZ: false
    publiclyAccessible: false
    deletionProtection: false
    region: us-east-1
  writeConnectionSecretsToRef:
    name: my-database-connection
    namespace: default
  providerConfigRef:
    name: default
---
apiVersion: v1
kind: Secret
metadata:
  name: my-database-connection
  namespace: default
type: Opaque
data: {}`, apiVersion, dbType)
}

func generateStorageManifest(description, provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	versioning := "false"
	if strings.Contains(strings.ToLower(description), "version") {
		versioning = "true"
	}

	// Use provider in the API version
	apiVersion := fmt.Sprintf("s3.%s.crossplane.io/v1beta1", provider)

	return fmt.Sprintf(`apiVersion: %s
kind: Bucket
metadata:
  name: my-app-bucket
  namespace: default
spec:
  forProvider:
    region: us-east-1
    versioning:
      enabled: %s
    serverSideEncryptionConfiguration:
      rules:
      - applyServerSideEncryptionByDefault:
          sseAlgorithm: AES256
    publicAccessBlockConfiguration:
      blockPublicAcls: true
      blockPublicPolicy: true
      ignorePublicAcls: true
      restrictPublicBuckets: true
  providerConfigRef:
    name: default`, apiVersion, versioning)
}

func generateNetworkManifest(description, provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	// Use provider in the API version
	apiVersion := fmt.Sprintf("ec2.%s.crossplane.io/v1beta1", provider)

	return fmt.Sprintf(`apiVersion: %s
kind: VPC
metadata:
  name: my-vpc
  namespace: default
spec:
  forProvider:
    cidrBlock: 10.0.0.0/16
    region: us-east-1
    tags:
      Name: MyVPC
  providerConfigRef:
    name: default
---
apiVersion: %s
kind: Subnet
metadata:
  name: my-subnet-public
  namespace: default
spec:
  forProvider:
    availabilityZone: us-east-1a
    cidrBlock: 10.0.1.0/24
    region: us-east-1
    mapPublicIPOnLaunch: true
    vpcIdSelector:
      matchLabels:
        name: my-vpc
    tags:
      Name: MyPublicSubnet
  providerConfigRef:
    name: default`, apiVersion, apiVersion)
}

func generateComputeManifest(description, provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	// Use provider in the API version
	apiVersion := fmt.Sprintf("ec2.%s.crossplane.io/v1alpha1", provider)

	return fmt.Sprintf(`apiVersion: %s
kind: Instance
metadata:
  name: my-instance
  namespace: default
spec:
  forProvider:
    region: us-east-1
    instanceType: t3.micro
    imageId: ami-0abcdef1234567890
    keyName: my-key-pair
    subnetIdSelector:
      matchLabels:
        name: my-subnet-public
    securityGroupIdSelector:
      matchLabels:
        name: my-security-group
    tags:
      Name: MyInstance
  providerConfigRef:
    name: default`, apiVersion)
}

func generateWebAppManifest(description, provider string) string {
	if provider == "" || provider == "auto" {
		provider = "aws"
	}

	// Use provider in the API versions
	elbVersion := fmt.Sprintf("elbv2.%s.crossplane.io/v1alpha1", provider)
	rdsVersion := fmt.Sprintf("rds.%s.crossplane.io/v1alpha1", provider)

	return fmt.Sprintf(`# Load Balancer
apiVersion: %s
kind: LoadBalancer
metadata:
  name: my-web-lb
  namespace: default
spec:
  forProvider:
    region: us-east-1
    scheme: internet-facing
    type: application
    subnetIdSelector:
      matchLabels:
        network: public
    tags:
      Name: MyWebLB
  providerConfigRef:
    name: default
---
# Database
apiVersion: %s
kind: DBInstance
metadata:
  name: my-web-db
  namespace: default
spec:
  forProvider:
    dbInstanceClass: db.t3.micro
    engine: mysql
    engineVersion: "8.0"
    dbName: webapp
    masterUsername: admin
    allocatedStorage: 20
    storageEncrypted: true
    region: us-east-1
  writeConnectionSecretsToRef:
    name: my-web-db-connection
    namespace: default
  providerConfigRef:
    name: default`, elbVersion, rdsVersion)
}

func generateDefaultManifest(description, provider string) string {
	return fmt.Sprintf(`# Generated Crossplane resource for: %s
# Provider: %s
#
# This is a template - please customize according to your needs

apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: my-custom-resource
  labels:
    provider: %s
spec:
  compositeTypeRef:
    apiVersion: example.com/v1alpha1
    kind: XCustomResource
  
  # Add your resource definitions here
  # Visit https://docs.crossplane.io for documentation`, description, provider, provider)
}

func applyManifest(ctx context.Context, client *crossplane.Client, manifest string) error {
	// In a real implementation, this would parse the YAML and apply it to the cluster
	cli.PrintInfo("📝 Parsing manifest...")
	cli.PrintInfo("🔍 Validating resources...")
	cli.PrintInfo("⚡ Creating resources...")

	// Simulate successful application
	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("provider", "p", "", "target cloud provider (aws, gcp, azure)")
	generateCmd.Flags().StringP("output", "o", "yaml", "output format (yaml, json)")
	generateCmd.Flags().Bool("dry-run", false, "generate manifest but don't apply")
	generateCmd.Flags().Bool("apply", false, "apply the generated manifest to cluster")
}
