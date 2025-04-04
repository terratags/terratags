package validator

import (
	"fmt"

	"terratags/pkg/parser"
)

// ValidationResult represents the result of a tag validation
type ValidationResult struct {
	ResourceType string
	ResourceName string
	Message      string
}

// ValidateRequiredTags validates that resources have all required tags
func ValidateRequiredTags(resources []parser.Resource, requiredTags []string) []ValidationResult {
	var results []ValidationResult

	for _, resource := range resources {
		for _, requiredTag := range requiredTags {
			if _, exists := resource.Tags[requiredTag]; !exists {
				results = append(results, ValidationResult{
					ResourceType: resource.Type,
					ResourceName: resource.Name,
					Message:      fmt.Sprintf("Missing required tag '%s'", requiredTag),
				})
			}
		}
	}

	return results
}
