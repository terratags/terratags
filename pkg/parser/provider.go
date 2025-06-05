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

	// Find AWS-style provider blocks with default_tags as a nested block
	awsProviderPattern := `provider\s+"([^"]+)"\s*{([\s\S]*?default_tags[\s\S]*?{[\s\S]*?}[\s\S]*?)}`
	awsProviderRegex := regexp.MustCompile(`(?s)` + awsProviderPattern)
	awsProviderMatches := awsProviderRegex.FindAllStringSubmatch(fileContent, -1)

	for _, match := range awsProviderMatches {
		if len(match) > 2 {
			providerName := match[1]
			providerBody := match[2]

			// We're only interested in AWS providers that might have default_tags
			if strings.HasPrefix(providerName, "aws") && !strings.HasPrefix(providerName, "awscc") {
				defaultTags := extractDefaultTagsFromProviderBody(providerBody)
				if len(defaultTags) > 0 {
					logging.Debug("Found AWS provider %s with default_tags", providerName)
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

	// Find Azure API-style provider blocks with default_tags as a direct attribute
	azapiProviderPattern := `provider\s+"azapi"\s*{([\s\S]*?default_tags\s*=\s*{([\s\S]*?)}[\s\S]*?)}`
	azapiProviderRegex := regexp.MustCompile(`(?s)` + azapiProviderPattern)
	azapiProviderMatches := azapiProviderRegex.FindAllStringSubmatch(fileContent, -1)

	for _, match := range azapiProviderMatches {
		if len(match) > 2 {
			providerName := "azapi"
			tagContent := match[2]
			
			// Extract key-value pairs from the tags block
			defaultTags := make(map[string]string)
			keyValuePattern := `["']?([A-Za-z0-9_-]+)["']?\s*=\s*["']?([^,"'}\s]*)["']?`
			keyValueRegex := regexp.MustCompile(keyValuePattern)
			keyValueMatches := keyValueRegex.FindAllStringSubmatch(tagContent, -1)

			for _, kvMatch := range keyValueMatches {
				if len(kvMatch) > 2 {
					key := kvMatch[1]
					value := kvMatch[2]
					defaultTags[key] = value
				}
			}
			
			if len(defaultTags) > 0 {
				logging.Debug("Found Azure API provider with default_tags")
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
