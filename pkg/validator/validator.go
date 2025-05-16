package validator

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/terratags/terratags/pkg/config"
	"github.com/terratags/terratags/pkg/parser"
)

// TagViolation represents a tag validation violation
type TagViolation struct {
	ResourceType string
	ResourceName string
	ResourcePath string
	MissingTags  []string
	IsExempt     bool
	ExemptReason string
}

// TagComplianceStats represents statistics about tag compliance
type TagComplianceStats struct {
	TotalResources     int
	CompliantResources int
	ExemptResources    int
	ViolationsByTag    map[string]int
}

// ValidateResources validates that all resources have the required tags
func ValidateResources(resources []parser.Resource, providers []parser.ProviderConfig, cfg *config.Config) (bool, []TagViolation, TagComplianceStats, []parser.Resource) {
	var violations []TagViolation
	stats := TagComplianceStats{
		TotalResources:  len(resources),
		ViolationsByTag: make(map[string]int),
	}
	valid := true

	// Create a map of provider default tags by path
	defaultTagsByPath := make(map[string]map[string]string)
	for _, provider := range providers {
		defaultTagsByPath[provider.Path] = provider.DefaultTags
	}

	for _, resource := range resources {
		// Get default tags for this resource's path
		defaultTags := defaultTagsByPath[resource.Path]

		// Track tag sources
		for k, v := range resource.Tags {
			resource.TagSources[k] = parser.TagSource{
				Source: "resource",
				Value:  v,
			}
		}

		// Add default tags to tag sources
		for k, v := range defaultTags {
			if _, exists := resource.TagSources[k]; !exists {
				resource.TagSources[k] = parser.TagSource{
					Source: "provider_default",
					Value:  v,
				}
			}
		}

		// Check for missing required tags
		var missingTags []string
		isExempt := false
		var exemptReason string

		for _, requiredTag := range cfg.Required {
			// Check if the tag is in the resource's tags
			if _, exists := resource.Tags[requiredTag]; !exists {
				// If not in resource tags, check if it's in default tags
				if _, existsInDefault := defaultTags[requiredTag]; !existsInDefault {
					// Check if this resource is exempt from this tag requirement
					exempt, reason := cfg.IsExemptFromTag(resource.Type, resource.Name, requiredTag)
					if exempt {
						isExempt = true
						exemptReason = reason
					} else {
						missingTags = append(missingTags, requiredTag)
						stats.ViolationsByTag[requiredTag]++
					}
				} else {
					// Tag exists in default_tags, so it's valid
					if len(defaultTags) > 0 {
						fmt.Printf("Resource %s '%s' inherits tag '%s' from provider default_tags\n",
							resource.Type, resource.Name, requiredTag)
					}
				}
			}
		}

		if len(missingTags) > 0 {
			valid = false
			violations = append(violations, TagViolation{
				ResourceType: resource.Type,
				ResourceName: resource.Name,
				ResourcePath: resource.Path,
				MissingTags:  missingTags,
				IsExempt:     isExempt,
				ExemptReason: exemptReason,
			})
		} else if isExempt {
			stats.ExemptResources++
		} else {
			stats.CompliantResources++
		}
	}

	return valid, violations, stats, resources
}

// ValidateDirectory validates all Terraform files in a directory
func ValidateDirectory(dir string, cfg *config.Config, verbose bool) (bool, []TagViolation, TagComplianceStats, []parser.Resource) {
	// Find all Terraform files in the directory
	files, err := filepath.Glob(filepath.Join(dir, "*.tf"))
	if err != nil {
		return false, []TagViolation{{
			ResourceType: "error",
			ResourceName: "error",
			ResourcePath: dir,
			MissingTags:  []string{fmt.Sprintf("Error finding Terraform files: %s", err)},
		}}, TagComplianceStats{}, nil
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
	valid, violations, stats, _ := ValidateResources(allResources, allProviders, cfg)
	return valid, violations, stats, allResources
}

// ValidateTerraformPlan validates a Terraform plan file
func ValidateTerraformPlan(planPath string, cfg *config.Config, verbose bool) (bool, []TagViolation, TagComplianceStats, []parser.Resource) {
	if verbose {
		fmt.Printf("Analyzing Terraform plan: %s\n", planPath)
	}

	// Parse the plan
	resources, err := parser.ParseTerraformPlan(planPath)
	if err != nil {
		return false, []TagViolation{{
			ResourceType: "error",
			ResourceName: "error",
			ResourcePath: planPath,
			MissingTags:  []string{fmt.Sprintf("Error parsing Terraform plan: %s", err)},
		}}, TagComplianceStats{}, nil
	}

	if verbose {
		fmt.Printf("Found %d taggable resources in plan\n", len(resources))
	}

	// Validate resources (no provider default tags in plan output)
	valid, violations, stats, _ := ValidateResources(resources, []parser.ProviderConfig{}, cfg)
	return valid, violations, stats, resources
}

// GenerateRemediationCode generates HCL code to fix missing tags
func GenerateRemediationCode(resourceType, resourceName, resourcePath string, missingTags []string, existingTags map[string]string) string {
	var sb strings.Builder

	// Start with the resource block
	sb.WriteString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", resourceType, resourceName))
	sb.WriteString("  # Existing attributes preserved\n\n")

	// Generate the tags block with existing and suggested tags
	sb.WriteString("  tags = {\n")

	// Add existing tags
	for k, v := range existingTags {
		sb.WriteString(fmt.Sprintf("    %s = \"%s\"\n", k, v))
	}

	// Add missing tags with placeholder values
	for _, tag := range missingTags {
		sb.WriteString(fmt.Sprintf("    %s = \"CHANGE_ME\"  # Added missing required tag\n", tag))
	}

	sb.WriteString("  }\n")
	sb.WriteString("}")

	return sb.String()
}

// SuggestProviderDefaultTagsUpdate suggests an update to provider default_tags
func SuggestProviderDefaultTagsUpdate(missingTags []string) string {
	var sb strings.Builder

	sb.WriteString("provider \"aws\" {\n")
	sb.WriteString("  # Existing provider configuration preserved\n\n")
	sb.WriteString("  default_tags {\n")
	sb.WriteString("    tags = {\n")

	// Add missing tags with placeholder values
	for _, tag := range missingTags {
		sb.WriteString(fmt.Sprintf("      %s = \"CHANGE_ME\"  # Added missing required tag\n", tag))
	}

	sb.WriteString("    }\n")
	sb.WriteString("  }\n")
	sb.WriteString("}")

	return sb.String()
}

// GenerateHTMLReport generates an HTML report of tag compliance
func GenerateHTMLReport(violations []TagViolation, stats TagComplianceStats, cfg *config.Config) string {
	var sb strings.Builder

	// HTML header
	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Terraform Tag Compliance Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .summary { margin-bottom: 20px; }
        .progress-bar { 
            width: 100%; 
            background-color: #f3f3f3; 
            border-radius: 5px; 
            margin-bottom: 20px;
        }
        .progress { 
            height: 30px; 
            background-color: #4CAF50; 
            border-radius: 5px; 
            text-align: center;
            line-height: 30px;
            color: white;
        }
        .resource { margin-bottom: 10px; padding: 10px; border: 1px solid #ddd; border-radius: 5px; }
        .compliant { background-color: #dff0d8; }
        .non-compliant { background-color: #f2dede; }
        .exempt { background-color: #fcf8e3; }
        .tag-table { width: 100%; border-collapse: collapse; }
        .tag-table th, .tag-table td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        .tag-table th { background-color: #f2f2f2; }
        .missing { color: red; }
        .present { color: green; }
        .exempt-tag { color: orange; }
    </style>
</head>
<body>
    <h1>Terraform Tag Compliance Report</h1>
    <p>Generated on: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>`)

	// Summary statistics
	compliancePercentage := 0.0
	if stats.TotalResources > 0 {
		compliancePercentage = float64(stats.CompliantResources) / float64(stats.TotalResources) * 100
	}

	sb.WriteString(fmt.Sprintf(`
    <div class="summary">
        <h2>Summary</h2>
        <p>Total Resources: %d</p>
        <p>Compliant Resources: %d</p>
        <p>Non-compliant Resources: %d</p>
        <p>Exempt Resources: %d</p>
        <div class="progress-bar">
            <div class="progress" style="width: %.1f%%;">%.1f%% Compliant</div>
        </div>
    </div>`,
		stats.TotalResources, stats.CompliantResources,
		stats.TotalResources-stats.CompliantResources-stats.ExemptResources,
		stats.ExemptResources,
		compliancePercentage, compliancePercentage))

	// Required tags section
	sb.WriteString(`<h2>Required Tags</h2>
    <ul>`)
	for _, tag := range cfg.Required {
		sb.WriteString(fmt.Sprintf(`<li>%s</li>`, tag))
	}
	sb.WriteString(`</ul>`)

	// Violations by tag
	sb.WriteString(`<h2>Violations by Tag</h2>
    <table class="tag-table">
        <tr>
            <th>Tag</th>
            <th>Violations</th>
        </tr>`)

	for tag, count := range stats.ViolationsByTag {
		sb.WriteString(fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%d</td>
        </tr>`, tag, count))
	}
	sb.WriteString(`</table>`)

	// Non-compliant resources
	sb.WriteString(`<h2>Non-compliant Resources</h2>`)

	if len(violations) == 0 {
		sb.WriteString(`<p>All resources are compliant!</p>`)
	} else {
		for _, v := range violations {
			if v.IsExempt {
				sb.WriteString(fmt.Sprintf(`
                <div class="resource exempt">
                    <h3>%s "%s"</h3>
                    <p>Path: %s</p>
                    <p>Status: <span class="exempt-tag">EXEMPT</span> - %s</p>
                </div>`, v.ResourceType, v.ResourceName, v.ResourcePath, v.ExemptReason))
			} else {
				sb.WriteString(fmt.Sprintf(`
                <div class="resource non-compliant">
                    <h3>%s "%s"</h3>
                    <p>Path: %s</p>
                    <p>Missing Tags: %s</p>
                </div>`, v.ResourceType, v.ResourceName, v.ResourcePath, strings.Join(v.MissingTags, ", ")))
			}
		}
	}

	sb.WriteString(`</body></html>`)

	return sb.String()
}
