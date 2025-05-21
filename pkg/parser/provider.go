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

	// Find all provider blocks - improved pattern to better match provider blocks with default_tags
	providerPattern := `provider\s+"([^"]+)"\s*{([\s\S]*?default_tags[\s\S]*?{[\s\S]*?}[\s\S]*?)}`
	providerRegex := regexp.MustCompile(`(?s)` + providerPattern)
	providerMatches := providerRegex.FindAllStringSubmatch(fileContent, -1)

	for _, match := range providerMatches {
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

	return providers, nil
}

// extractDefaultTagsFromProviderBody extracts default_tags from a provider block body
func extractDefaultTagsFromProviderBody(providerBody string) map[string]string {
	defaultTags := make(map[string]string)

	// Find the default_tags block within the provider - improved pattern
	defaultTagsPattern := `default_tags\s*{[\s\S]*?tags\s*=\s*{([\s\S]*?)}`
	defaultTagsRegex := regexp.MustCompile(`(?s)` + defaultTagsPattern)
	defaultTagsMatch := defaultTagsRegex.FindStringSubmatch(providerBody)

	if len(defaultTagsMatch) > 1 {
		// Extract key-value pairs from the tags block
		tagContent := defaultTagsMatch[1]
		keyValuePattern := `["']?([A-Za-z0-9_-]+)["']?\s*=\s*["']?([^,"'}\s]*)["']?`
		keyValueRegex := regexp.MustCompile(keyValuePattern)
		keyValueMatches := keyValueRegex.FindAllStringSubmatch(tagContent, -1)

		for _, match := range keyValueMatches {
			if len(match) > 2 {
				key := match[1]
				value := match[2]
				defaultTags[key] = value
			}
		}
	}

	return defaultTags
}
