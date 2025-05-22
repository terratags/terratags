package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Schema represents the structure of the Terraform provider schema
type Schema struct {
	ProviderSchemas map[string]ProviderSchema `json:"provider_schemas"`
}

// ProviderSchema represents a provider's schema
type ProviderSchema struct {
	ResourceSchemas map[string]ResourceSchema `json:"resource_schemas"`
}

// ResourceSchema represents a resource's schema
type ResourceSchema struct {
	Block Block `json:"block"`
}

// Block represents a block in the schema
type Block struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// hasLabelsAttribute checks if a resource schema has a 'labels' attribute
func hasLabelsAttribute(schema ResourceSchema) bool {
	_, hasLabels := schema.Block.Attributes["labels"]
	return hasLabels
}

// createTerraformConfig creates a temporary Terraform configuration with Google provider
func createTerraformConfig(tempDir string) error {
	config := `terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }
}

provider "google" {
  project = "my-project-id"
  region  = "us-central1"
}
`
	return os.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(config), 0644)
}

// generateGoFile generates a Go file with the list of taggable Google resources
func generateGoFile(googleResources []string, outputFile string) error {
	// Sort resources alphabetically
	sort.Strings(googleResources)

	var content strings.Builder
	content.WriteString("package parser\n\n")
	content.WriteString("// Google resources that do not properly support labeling\n")
	content.WriteString("// These resources are excluded from the taggable resources list\n")
	content.WriteString("var GoogleExcludedResources = map[string]bool{\n")
	content.WriteString("\t// Add excluded resources here as they are identified\n")
	content.WriteString("}\n\n")
	content.WriteString("// Google taggable resources\n")
	content.WriteString("// This list is automatically generated from the provider schemas\n")
	content.WriteString("// and represents resources that support the 'labels' attribute\n")
	content.WriteString("var googleTaggableResources = map[string]bool{\n")

	// Google resources
	content.WriteString("\t// Google Provider resources\n")
	for _, resource := range googleResources {
		content.WriteString(fmt.Sprintf("\t\"%s\": true,\n", resource))
	}

	content.WriteString("}")

	return os.WriteFile(outputFile, []byte(content.String()), 0644)
}

func main() {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "terraform-providers")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Create Terraform configuration
	err = createTerraformConfig(tempDir)
	if err != nil {
		fmt.Printf("Error creating Terraform configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize Terraform
	fmt.Println("Initializing Terraform providers...")
	cmd := exec.Command("terraform", "init")
	cmd.Dir = tempDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error initializing Terraform: %v\n", err)
		os.Exit(1)
	}

	// Extract provider schemas
	fmt.Println("Extracting provider schemas...")
	schemaFile := filepath.Join(tempDir, "schema.json")
	cmd = exec.Command("terraform", "providers", "schema", "-json")
	cmd.Dir = tempDir
	schemaOutput, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error extracting provider schemas: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(schemaFile, schemaOutput, 0644)
	if err != nil {
		fmt.Printf("Error writing schema file: %v\n", err)
		os.Exit(1)
	}

	// Parse schemas
	fmt.Println("Parsing schemas to find taggable resources...")
	schemaData, err := os.ReadFile(schemaFile)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	var schema Schema
	err = json.Unmarshal(schemaData, &schema)
	if err != nil {
		fmt.Printf("Error parsing schema JSON: %v\n", err)
		os.Exit(1)
	}

	var googleResources []string

	// Process Google provider
	googleSchema, ok := schema.ProviderSchemas["registry.terraform.io/hashicorp/google"]
	if ok {
		for resourceName, resourceSchema := range googleSchema.ResourceSchemas {
			if hasLabelsAttribute(resourceSchema) {
				googleResources = append(googleResources, resourceName)
			}
		}
	}

	// Determine output file path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	scriptDir := filepath.Dir(execPath)
	repoRoot := filepath.Dir(scriptDir)
	outputFile := filepath.Join(repoRoot, "pkg", "parser", "google_taggable_resources.go")

	// If running from the scripts directory directly
	if filepath.Base(scriptDir) == "scripts" {
		repoRoot = filepath.Dir(scriptDir)
		outputFile = filepath.Join(repoRoot, "pkg", "parser", "google_taggable_resources.go")
	} else {
		// If running from the repo root with go run
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		outputFile = filepath.Join(currentDir, "pkg", "parser", "google_taggable_resources.go")
	}

	// Generate Go file
	err = generateGoFile(googleResources, outputFile)
	if err != nil {
		fmt.Printf("Error generating Go file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated taggable resources list with %d Google resources\n",
		len(googleResources))
}
