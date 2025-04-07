package validator

import (
	"fmt"
	"path/filepath"

	"terratags/pkg/config"
	"terratags/pkg/parser"
)

// ValidateResources validates that all resources have the required tags
func ValidateResources(resources []parser.Resource, providers []parser.ProviderConfig, cfg *config.Config) (bool, []string) {
	var issues []string
	valid := true

	// Create a map of provider default tags by path
	defaultTagsByPath := make(map[string]map[string]string)
	for _, provider := range providers {
		defaultTagsByPath[provider.Path] = provider.DefaultTags
	}

	for _, resource := range resources {
		// Get default tags for this resource's path
		defaultTags := defaultTagsByPath[resource.Path]
		
		// Check if all required tags are present
		for _, requiredTag := range cfg.Required {
			// Check if the tag is in the resource's tags
			if _, exists := resource.Tags[requiredTag]; !exists {
				// If not in resource tags, check if it's in default tags
				if _, existsInDefault := defaultTags[requiredTag]; !existsInDefault {
					valid = false
					issues = append(issues, fmt.Sprintf("  - %s '%s': Missing required tag: %s", resource.Type, resource.Name, requiredTag))
				} else {
					// Tag exists in default_tags, so it's valid
					if len(defaultTags) > 0 {
						fmt.Printf("Resource %s '%s' inherits tag '%s' from provider default_tags\n", 
							resource.Type, resource.Name, requiredTag)
					}
				}
			}
		}
	}

	return valid, issues
}

// ValidateDirectory validates all Terraform files in a directory
func ValidateDirectory(dir string, cfg *config.Config, verbose bool) (bool, []string) {
	// Find all Terraform files in the directory
	files, err := filepath.Glob(filepath.Join(dir, "*.tf"))
	if err != nil {
		return false, []string{fmt.Sprintf("Error finding Terraform files: %s", err)}
	}

	if verbose {
		fmt.Printf("Found %d Terraform files to analyze\n", len(files))
	}

	var allResources []parser.Resource
	var allProviders []parser.ProviderConfig

	// Parse each file
	for _, file := range files {
		if verbose {
			fmt.Printf("Analyzing file: %s\n", file)
		}

		// Parse resources
		resources, err := parser.ParseFile(file)
		if err != nil {
			if verbose {
				fmt.Printf("Error parsing file %s: %s\n", file, err)
			}
			continue
		}
		allResources = append(allResources, resources...)

		// Parse provider blocks
		providers, err := parser.ParseProviderBlocks(file)
		if err != nil {
			if verbose {
				fmt.Printf("Error parsing provider blocks in %s: %s\n", file, err)
			}
			continue
		}
		allProviders = append(allProviders, providers...)
	}

	if verbose {
		fmt.Printf("Found %d taggable resources\n", len(allResources))
		fmt.Printf("Found %d provider configurations with default tags\n", len(allProviders))
	}

	// Validate resources
	valid, issues := ValidateResources(allResources, allProviders, cfg)
	return valid, issues
}
