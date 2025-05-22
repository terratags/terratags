
package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/terratags/terratags/pkg/logging"
)

// ProviderConfig represents a Terraform provider configuration
type ProviderConfig struct {
	Name        string
	DefaultTags map[string]string
	Path        string
}

// ParseProviderBlocks parses a Terraform file and extracts provider configurations
func ParseProviderBlocks(path string) ([]ProviderConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Convert content to string
	fileContent := string(content)

	var providers []ProviderConfig

	// Find all provider blocks with default_tags (AWS) or default_labels (Google)
	// AWS provider pattern
	awsProviderPattern := `provider\s+"([^"]+)"\s*{([\s\S]*?default_tags[\s\S]*?{[\s\S]*?}[\s\S]*?)}`
	awsProviderRegex := regexp.MustCompile(`(?s)` + awsProviderPattern)
	awsProviderMatches := awsProviderRegex.FindAllStringSubmatch(fileContent, -1)

	// Process AWS providers
	for _, match := range awsProviderMatches {
		if len(match) > 2 {
			providerName := match[1]
			providerBody := match[2]

			// We're only interested in AWS providers that might have default_tags
			// Note: AWSCC provider doesn't support default_tags
			if strings.HasPrefix(providerName, "aws") && !strings.HasPrefix(providerName, "awscc") {
				defaultTags := extractDefaultTagsFromProviderBody(providerBody)
				if len(defaultTags) > 0 {
					logging.Debug("Found provider %s with default_tags", providerName)
					for tag, value := range defaultTags {
						logging.Debug("Found default tag key: %s with value: %s", tag, value)
					}
					providers = append(providers, ProviderConfig{
						Name:        providerName,
						DefaultTags: defaultTags,
						Path:        path,
					})
				}
			}
		}
	}

	// Google provider pattern
	googleProviderPattern := `provider\s+"([^"]+)"\s*{([\s\S]*?default_labels[\s\S]*?{[\s\S]*?}[\s\S]*?)}`
	googleProviderRegex := regexp.MustCompile(`(?s)` + googleProviderPattern)
	googleProviderMatches := googleProviderRegex.FindAllStringSubmatch(fileContent, -1)

	// Process Google providers
	for _, match := range googleProviderMatches {
		if len(match) > 2 {
			providerName := match[1]
			providerBody := match[2]

			// We're only interested in Google providers that might have default_labels
			if strings.HasPrefix(providerName, "google") {
				defaultLabels := extractDefaultLabelsFromProviderBody(providerBody)
				if len(defaultLabels) > 0 {
					logging.Debug("Found provider %s with default_labels", providerName)
					for label, value := range defaultLabels {
						logging.Debug("Found default label key: %s with value: %s", label, value)
					}
					providers = append(providers, ProviderConfig{
						Name:        providerName,
						DefaultTags: defaultLabels, // We use the same field for both tags and labels
						Path:        path,
					})
				}
			}
		}
	}

	return providers, nil
}

// extractDefaultAttributesFromProviderBody extracts default_tags or default_labels from a provider block body
func extractDefaultAttributesFromProviderBody(providerBody string, blockName string, attributeName string) map[string]string {
	defaultAttributes := make(map[string]string)

	// Find the default_tags or default_labels block within the provider
	pattern := fmt.Sprintf(`%s\s*{[\s\S]*?%s\s*=\s*{([\s\S]*?)}`, blockName, attributeName)
	regex := regexp.MustCompile(`(?s)` + pattern)
	match := regex.FindStringSubmatch(providerBody)

	if len(match) > 1 {
		// Extract key-value pairs from the tags/labels block
		content := match[1]
		keyValuePattern := `["']?([A-Za-z0-9_-]+)["']?\s*=\s*["']?([^,"'}\s]*)["']?`
		keyValueRegex := regexp.MustCompile(keyValuePattern)
		keyValueMatches := keyValueRegex.FindAllStringSubmatch(content, -1)

		for _, match := range keyValueMatches {
			if len(match) > 2 {
				key := match[1]
				value := match[2]
				defaultAttributes[key] = value
			}
		}
	}

	return defaultAttributes
}

// extractDefaultTagsFromProviderBody extracts default_tags from a provider block body
func extractDefaultTagsFromProviderBody(providerBody string) map[string]string {
	return extractDefaultAttributesFromProviderBody(providerBody, "default_tags", "tags")
}

// extractDefaultLabelsFromProviderBody extracts default_labels from a provider block body
func extractDefaultLabelsFromProviderBody(providerBody string) map[string]string {
	return extractDefaultAttributesFromProviderBody(providerBody, "default_labels", "labels")
}
