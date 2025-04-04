package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// hasTagsAttribute checks if a resource schema has a 'tags' attribute
func hasTagsAttribute(schema ResourceSchema) bool {
	_, hasTags := schema.Block.Attributes["tags"]
	_, hasTagsAll := schema.Block.Attributes["tags_all"]
	return hasTags || hasTagsAll
}

// createTerraformConfig creates a temporary Terraform configuration with AWS and AWSCC providers
func createTerraformConfig(tempDir string) error {
	config := `terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    awscc = {
      source = "hashicorp/awscc"
    }
  }
}

provider "aws" {
  region = "us-west-2"
}

provider "awscc" {
  region = "us-west-2"
}
`
	return ioutil.WriteFile(filepath.Join(tempDir, "main.tf"), []byte(config), 0644)
}

// generateGoFile generates a Go file with the list of taggable resources
func generateGoFile(awsResources, awsccResources []string, outputFile string) error {
	// Sort resources alphabetically
	sort.Strings(awsResources)
	sort.Strings(awsccResources)

	var content strings.Builder
	content.WriteString("package parser\n\n")
	content.WriteString("// AWS and AWSCC taggable resources\n")
	content.WriteString("// This list is automatically generated from the provider schemas\n")
	content.WriteString("// and represents resources that support the 'tags' attribute\n")
	content.WriteString("var awsTaggableResources = map[string]bool{\n")

	// AWS resources
	content.WriteString("\t// AWS Provider resources\n")
	for _, resource := range awsResources {
		content.WriteString(fmt.Sprintf("\t\"%s\": true,\n", resource))
	}

	content.WriteString("\n\t// AWSCC Provider resources\n")
	for _, resource := range awsccResources {
		content.WriteString(fmt.Sprintf("\t\"%s\": true,\n", resource))
	}

	content.WriteString("}")

	return ioutil.WriteFile(outputFile, []byte(content.String()), 0644)
}

func main() {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "terraform-providers")
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

	err = ioutil.WriteFile(schemaFile, schemaOutput, 0644)
	if err != nil {
		fmt.Printf("Error writing schema file: %v\n", err)
		os.Exit(1)
	}

	// Parse schemas
	fmt.Println("Parsing schemas to find taggable resources...")
	schemaData, err := ioutil.ReadFile(schemaFile)
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

	var awsResources []string
	var awsccResources []string

	// Process AWS provider
	awsSchema, ok := schema.ProviderSchemas["registry.terraform.io/hashicorp/aws"]
	if ok {
		for resourceName, resourceSchema := range awsSchema.ResourceSchemas {
			if hasTagsAttribute(resourceSchema) {
				awsResources = append(awsResources, resourceName)
			}
		}
	}

	// Process AWSCC provider
	awsccSchema, ok := schema.ProviderSchemas["registry.terraform.io/hashicorp/awscc"]
	if ok {
		for resourceName, resourceSchema := range awsccSchema.ResourceSchemas {
			if hasTagsAttribute(resourceSchema) {
				awsccResources = append(awsccResources, resourceName)
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
	outputFile := filepath.Join(repoRoot, "pkg", "parser", "aws_taggable_resources.go")

	// If running from the scripts directory directly
	if filepath.Base(scriptDir) == "scripts" {
		repoRoot = filepath.Dir(scriptDir)
		outputFile = filepath.Join(repoRoot, "pkg", "parser", "aws_taggable_resources.go")
	} else {
		// If running from the repo root with go run
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		outputFile = filepath.Join(currentDir, "pkg", "parser", "aws_taggable_resources.go")
	}

	// Generate Go file
	err = generateGoFile(awsResources, awsccResources, outputFile)
	if err != nil {
		fmt.Printf("Error generating Go file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated taggable resources list with %d AWS resources and %d AWSCC resources\n", 
		len(awsResources), len(awsccResources))
}
