package validator

import (
	"fmt"
	"strings"

	"terratags/pkg/parser"
)

// ValidationResult represents the result of a tag validation
type ValidationResult struct {
	ResourceType string
	ResourceName string
	Message      string
}

// ValidateRequiredTags validates that resources have all required tags
// considering both resource-specific tags and provider default tags
func ValidateRequiredTags(resources []parser.Resource, providers []parser.ProviderConfig, requiredTags []string) []ValidationResult {
	var results []ValidationResult

	// Extract default tags from AWS providers
	defaultTags := make(map[string]string)
	for _, provider := range providers {
		if strings.HasPrefix(provider.Name, "aws") {
			for k, v := range provider.DefaultTags {
				defaultTags[k] = v
			}
		}
	}

	for _, resource := range resources {
		// For AWS resources, consider both resource tags and default tags
		effectiveTags := make(map[string]string)
		
		// Add provider default tags first (if this is an AWS resource)
		if strings.HasPrefix(resource.Type, "aws") {
			for k, v := range defaultTags {
				effectiveTags[k] = v
			}
		}
		
		// Add resource-specific tags (these override default tags)
		for k, v := range resource.Tags {
			effectiveTags[k] = v
		}
		
		// Now validate against the effective tags
		var missingTags []string
		for _, required := range requiredTags {
			if _, exists := effectiveTags[required]; !exists {
				missingTags = append(missingTags, required)
			}
		}
		
		if len(missingTags) > 0 {
			results = append(results, ValidationResult{
				ResourceType: resource.Type,
				ResourceName: resource.Name,
				Message:      fmt.Sprintf("Missing required tags: %s", strings.Join(missingTags, ", ")),
			})
		}
	}

	return results
}
