package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/terratags/terratags/pkg/logging"
)

// Resource represents a Terraform resource with its tags
type Resource struct {
	Type string
	Name string
	Tags map[string]string
	Path string
	// New field to track tag sources
	TagSources map[string]TagSource
}

// TagSource represents the source of a tag
type TagSource struct {
	Source string // "provider_default", "resource", "module"
	Value  string
}

// ParseFile parses a Terraform file and extracts resources with their tags
func ParseFile(path string, logLevel string) ([]Resource, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	var resources []Resource

	// Create a more permissive schema that allows all block types
	// but we'll only process resource and module blocks
	content2, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
			{
				Type:       "module",
				LabelNames: []string{"name"},
			},
			// Allow other block types but don't process them
			{
				Type:       "provider",
				LabelNames: []string{"name"},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
			},
			{
				Type: "locals",
			},
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
			{
				Type: "terraform",
			},
		},
	})

	if diags.HasErrors() {
		// If we still have errors, try a more permissive approach
		logging.Warn("Warning: Some blocks in %s couldn't be parsed, but we'll continue with what we can parse", path)
	}

	for _, block := range content2.Blocks {
		switch block.Type {
		case "resource":
			resourceType := block.Labels[0]
			resourceName := block.Labels[1]

			// Check if this resource type supports tagging
			if isTaggableResource(resourceType) {
				tags := extractTagsFromContent(content, resourceType, resourceName)
				resources = append(resources, Resource{
					Type:       resourceType,
					Name:       resourceName,
					Tags:       tags,
					Path:       path,
					TagSources: make(map[string]TagSource),
				})
			}
		case "module":
			moduleName := block.Labels[0]
			// Extract module resources and their tags
			moduleTags := extractModuleTagsFromContent(content, moduleName)
			if len(moduleTags) > 0 {
				resources = append(resources, Resource{
					Type:       "module",
					Name:       moduleName,
					Tags:       moduleTags,
					Path:       path,
					TagSources: make(map[string]TagSource),
				})
			}
		// Ignore other block types (provider, data, locals, etc.)
		default:
			// Skip processing for non-resource, non-module blocks
		}
	}

	return resources, nil
}

// isTaggableResource checks if a resource type supports tagging
func isTaggableResource(resourceType string) bool {
	// First check if it's in the excluded list
	if AwsccExcludedResources[resourceType] {
		logging.Debug("%s is in the excluded awscc resources list", resourceType)
		return true // Still return true so we process it, but we'll skip validation later
	}
	// Check if it's in the Google excluded list
	if GoogleExcludedResources[resourceType] {
		logging.Debug("%s is in the excluded google resources list", resourceType)
		return true // Still return true so we process it, but we'll skip validation later
	}
	// Check if it's a Google resource
	if strings.HasPrefix(resourceType, "google_") {
		return googleTaggableResources[resourceType]
	}
	// Use the comprehensive list of AWS and AWSCC taggable resources
	return awsTaggableResources[resourceType]
}

// extractTagsFromContent extracts tags directly from the file content
func extractTagsFromContent(content []byte, resourceType, resourceName string) map[string]string {
	tags := make(map[string]string)

	// Convert content to string
	fileContent := string(content)

	// Find the resource block with proper handling of nested blocks
	// This pattern matches the entire resource block including nested blocks
	resourcePattern := fmt.Sprintf(`resource\s+"%s"\s+"%s"\s*{[\s\S]*?(?:^}|\n}|\r\n})`,
		regexp.QuoteMeta(resourceType), regexp.QuoteMeta(resourceName))
	resourceRegex := regexp.MustCompile(`(?sm)` + resourcePattern)
	resourceMatch := resourceRegex.FindString(fileContent)

	if resourceMatch != "" {
		// Check if this is a Google resource
		isGoogle := strings.HasPrefix(resourceType, "google_")
		// Check if this is an AWSCC resource
		isAWSCC := strings.HasPrefix(resourceType, "awscc_")

		if isGoogle {
			// For Google resources, labels are in the format: labels = {Key = "Value"}
			labelsPattern := `labels\s*=\s*{([\s\S]*?)}`
			labelsRegex := regexp.MustCompile(`(?s)` + labelsPattern)
			labelsMatch := labelsRegex.FindStringSubmatch(resourceMatch)

			if len(labelsMatch) > 1 {
				logging.Debug("Found labels attribute in %s %s", resourceType, resourceName)

				// Extract key-value pairs
				labelContent := labelsMatch[1]
				keyValuePattern := `["']?([A-Za-z0-9_-]+)["']?\s*=\s*["']?([^,"'}\s]*)["']?`
				keyValueRegex := regexp.MustCompile(keyValuePattern)
				keyValueMatches := keyValueRegex.FindAllStringSubmatch(labelContent, -1)

				for _, match := range keyValueMatches {
					if len(match) > 2 {
						key := match[1]
						value := match[2]
						logging.Debug("Found label key: %s", key)
						tags[key] = value
					}
				}
			} else {
				logging.Debug("No labels attribute found in %s %s", resourceType, resourceName)
			}
		} else if isAWSCC {
			// For AWSCC resources, tags are in the format: tags = [{key = "Key", value = "Value"}]
			tagsPattern := `tags\s*=\s*\[([\s\S]*?)\]`
			tagsRegex := regexp.MustCompile(`(?s)` + tagsPattern)
			tagsMatch := tagsRegex.FindStringSubmatch(resourceMatch)

			if len(tagsMatch) > 1 {
				logging.Debug("Found AWSCC tags attribute in %s %s", resourceType, resourceName)

				// Extract key-value pairs from the list of maps format
				tagContent := tagsMatch[1]
				// Match each {key = "...", value = "..."} block
				tagBlockPattern := `{[\s\S]*?key\s*=\s*["']?([^"'\s,}]*)["']?[\s\S]*?value\s*=\s*["']?([^"'\s,}]*)["']?[\s\S]*?}`
				tagBlockRegex := regexp.MustCompile(`(?s)` + tagBlockPattern)
				tagBlockMatches := tagBlockRegex.FindAllStringSubmatch(tagContent, -1)

				for _, match := range tagBlockMatches {
					if len(match) > 2 {
						key := match[1]
						value := match[2]
						logging.Debug("Found AWSCC tag key: %s with value: %s", key, value)
						tags[key] = value
					}
				}
			} else {
				logging.Debug("No tags attribute found in AWSCC resource %s %s", resourceType, resourceName)
			}
		} else {
			// For AWS resources, tags are in the format: tags = {Key = "Value"}
			tagsPattern := `tags\s*=\s*{([\s\S]*?)}`
			tagsRegex := regexp.MustCompile(`(?s)` + tagsPattern)
			tagsMatch := tagsRegex.FindStringSubmatch(resourceMatch)

			if len(tagsMatch) > 1 {
				logging.Debug("Found tags attribute in %s %s", resourceType, resourceName)

				// Extract key-value pairs
				tagContent := tagsMatch[1]
				keyValuePattern := `["']?([A-Za-z0-9_-]+)["']?\s*=\s*["']?([^,"'}\s]*)["']?`
				keyValueRegex := regexp.MustCompile(keyValuePattern)
				keyValueMatches := keyValueRegex.FindAllStringSubmatch(tagContent, -1)

				for _, match := range keyValueMatches {
					if len(match) > 2 {
						key := match[1]
						value := match[2]
						logging.Debug("Found tag key: %s", key)
						tags[key] = value
					}
				}
			} else {
				logging.Debug("No tags attribute found in %s %s", resourceType, resourceName)
			}
		}
	}

	return tags
}

// extractModuleTagsFromContent extracts module tags directly from the file content
func extractModuleTagsFromContent(content []byte, moduleName string) map[string]string {
	tags := make(map[string]string)

	// Convert content to string
	fileContent := string(content)

	// Find the module block with proper handling of nested blocks
	modulePattern := fmt.Sprintf(`module\s+"%s"\s*{[\s\S]*?(?:^}|\n}|\r\n})`, regexp.QuoteMeta(moduleName))
	moduleRegex := regexp.MustCompile(`(?sm)` + modulePattern)
	moduleMatch := moduleRegex.FindString(fileContent)

	if moduleMatch != "" {
		// Find the tags block within the module
		// Improved pattern to handle tags that might appear after nested blocks
		tagsPattern := `tags\s*=\s*{([\s\S]*?)}`
		tagsRegex := regexp.MustCompile(`(?s)` + tagsPattern)
		tagsMatch := tagsRegex.FindStringSubmatch(moduleMatch)

		if len(tagsMatch) > 1 {
			logging.Debug("Found tags attribute in module %s", moduleName)

			// Extract key-value pairs
			tagContent := tagsMatch[1]
			keyValuePattern := `["']?([A-Za-z0-9_-]+)["']?\s*=\s*["']?([^,"'}\s]*)["']?`
			keyValueRegex := regexp.MustCompile(keyValuePattern)
			keyValueMatches := keyValueRegex.FindAllStringSubmatch(tagContent, -1)

			for _, match := range keyValueMatches {
				if len(match) > 2 {
					key := match[1]
					value := match[2]
					logging.Debug("Found module tag key: %s", key)
					tags[key] = value
				}
			}
		} else {
			logging.Debug("No tags attribute found in module %s", moduleName)
		}
	}

	return tags
}

// ParseTerraformPlan parses a Terraform plan JSON file and extracts resources with their tags
func ParseTerraformPlan(planPath string, logLevel string) ([]Resource, error) {
	// Read the plan file
	planData, err := os.ReadFile(planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	// Parse the plan JSON
	var plan struct {
		ResourceChanges []struct {
			Address string `json:"address"`
			Type    string `json:"type"`
			Change  struct {
				After map[string]any `json:"after"`
			} `json:"change"`
		} `json:"resource_changes"`
	}

	if err := json.Unmarshal(planData, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w", err)
	}

	var resources []Resource

	// Process each resource change
	for _, rc := range plan.ResourceChanges {
		// Check if this is a taggable resource
		if isTaggableResource(rc.Type) {
			// Extract resource name from address
			nameParts := strings.Split(rc.Address, ".")
			resourceName := nameParts[len(nameParts)-1]

			// Extract tags from the "after" state
			tags := extractTagsFromPlanResource(rc.Change.After)

			resources = append(resources, Resource{
				Type:       rc.Type,
				Name:       resourceName,
				Tags:       tags,
				Path:       planPath,
				TagSources: make(map[string]TagSource),
			})
		}
	}

	return resources, nil
}

// extractTagsFromPlanResource extracts tags from a resource in the plan
func extractTagsFromPlanResource(resource map[string]any) map[string]string {
	tags := make(map[string]string)

	// Check if the resource has tags or labels (for Google resources)
	if labelsInterface, ok := resource["labels"]; ok {
		// This is likely a Google resource with labels
		logging.Debug("Found Google labels in plan resource")
		if labelsMap, ok := labelsInterface.(map[string]any); ok {
			for k, v := range labelsMap {
				if strValue, ok := v.(string); ok {
					logging.Debug("Found label key: %s", k)
					tags[k] = strValue
				}
			}
		}
	} else if tagsInterface, ok := resource["tags"]; ok {
		// Check if this is an AWSCC resource (tags will be a list of maps with key/value pairs)
		if tagsList, ok := tagsInterface.([]any); ok {
			// This is likely an AWSCC resource with tags as a list of maps
			logging.Debug("Found AWSCC tags in plan resource")
			for _, tagItem := range tagsList {
				if tagMap, ok := tagItem.(map[string]any); ok {
					// Extract key and value from the map
					if keyInterface, keyOk := tagMap["key"]; keyOk {
						if valueInterface, valueOk := tagMap["value"]; valueOk {
							if key, keyOk := keyInterface.(string); keyOk {
								if value, valueOk := valueInterface.(string); valueOk {
									logging.Debug("Found AWSCC tag key: %s with value: %s", key, value)
									tags[key] = value
								}
							}
						}
					}
				}
			}
		} else if tagsMap, ok := tagsInterface.(map[string]any); ok {
			// This is the standard AWS provider tags format (map of key/value)
			logging.Debug("Found AWS tags in plan resource")
			for k, v := range tagsMap {
				if strValue, ok := v.(string); ok {
					logging.Debug("Found tag key: %s", k)
					tags[k] = strValue
				}
			}
		}
	}

	return tags
}
