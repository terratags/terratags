package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// TagRequirement represents a tag requirement with optional pattern validation
type TagRequirement struct {
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	// Internal field to store compiled regex (not serialized)
	compiledPattern *regexp.Regexp `json:"-" yaml:"-"`
}

// Config represents the configuration for tag validation
type Config struct {
	RequiredTags  map[string]TagRequirement `json:"required_tags" yaml:"required_tags"`
	Exemptions    []ResourceExemption       `json:"exemptions" yaml:"exemptions"`
	ReportPath    string                    `json:"report_path" yaml:"report_path"`
	IgnoreTagCase bool                      `json:"-" yaml:"-"` // Runtime option, not from config file
	
	// Legacy support - will be populated from RequiredTags for backward compatibility
	Required []string `json:"-" yaml:"-"`
}

// ResourceExemption represents a resource that is exempt from tag requirements
type ResourceExemption struct {
	ResourceType string   `json:"resource_type" yaml:"resource_type"`
	ResourceName string   `json:"resource_name" yaml:"resource_name"`
	ExemptTags   []string `json:"exempt_tags" yaml:"exempt_tags"`
	Reason       string   `json:"reason" yaml:"reason"`
}

// LoadConfig loads the configuration from a JSON or YAML file (local or remote)
// Supports:
// - Local file paths: ./config.yaml, /path/to/config.json
// - HTTP/HTTPS URLs: https://example.com/config.yaml
// - Git HTTPS: https://github.com/org/repo.git//path/to/config.yaml?ref=main
// - Git SSH: git@github.com:org/repo.git//path/to/config.yaml?ref=main
func LoadConfig(path string) (*Config, error) {
	var data []byte
	var err error
	var ext string
	
	// Check if it's a remote URL
	if IsRemoteURL(path) {
		data, err = FetchRemoteConfig(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote config: %w", err)
		}
		// Extract extension from URL
		ext = strings.ToLower(filepath.Ext(path))
		if idx := strings.Index(ext, "?"); idx != -1 {
			ext = ext[:idx]
		}
	} else {
		// Local file
		data, err = os.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		ext = strings.ToLower(filepath.Ext(path))
	}

	var config Config

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	// Compile regex patterns and validate
	if err := config.compilePatterns(); err != nil {
		return nil, fmt.Errorf("failed to compile regex patterns: %w", err)
	}

	// Populate legacy Required field for backward compatibility
	config.populateLegacyRequired()

	return &config, nil
}

// LoadExemptions loads exemptions from a JSON or YAML file
func LoadExemptions(path string) ([]ResourceExemption, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read exemptions file: %w", err)
	}

	var exemptions struct {
		Exemptions []ResourceExemption `json:"exemptions" yaml:"exemptions"`
	}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &exemptions); err != nil {
			return nil, fmt.Errorf("failed to parse JSON exemptions: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &exemptions); err != nil {
			return nil, fmt.Errorf("failed to parse YAML exemptions: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported exemptions file format: %s", ext)
	}

	return exemptions.Exemptions, nil
}

// IsExemptFromTag checks if a resource is exempt from a specific tag requirement
func (c *Config) IsExemptFromTag(resourceType, resourceName, tagName string) (bool, string) {
	for _, exemption := range c.Exemptions {
		if (exemption.ResourceType == resourceType || exemption.ResourceType == "*") &&
			(exemption.ResourceName == resourceName || exemption.ResourceName == "*") {

			for _, exemptTag := range exemption.ExemptTags {
				if c.IgnoreTagCase {
					// Case-insensitive comparison
					if strings.EqualFold(exemptTag, tagName) || exemptTag == "*" {
						return true, exemption.Reason
					}
				} else {
					// Case-sensitive comparison (original behavior)
					if exemptTag == tagName || exemptTag == "*" {
						return true, exemption.Reason
					}
				}
			}
		}
	}
	return false, ""
}

// UnmarshalJSON implements custom JSON unmarshaling to support both array and object formats
func (c *Config) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a struct with the new format
	type configAlias Config
	var temp struct {
		RequiredTags interface{}         `json:"required_tags"`
		Exemptions   []ResourceExemption `json:"exemptions"`
		ReportPath   string              `json:"report_path"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Copy non-required_tags fields
	c.Exemptions = temp.Exemptions
	c.ReportPath = temp.ReportPath
	c.RequiredTags = make(map[string]TagRequirement)

	// Handle required_tags field which can be array or object
	if temp.RequiredTags != nil {
		switch v := temp.RequiredTags.(type) {
		case []interface{}:
			// Array format (legacy)
			for _, item := range v {
				if tagName, ok := item.(string); ok {
					c.RequiredTags[tagName] = TagRequirement{}
				}
			}
		case map[string]interface{}:
			// Object format (new)
			for tagName, tagConfig := range v {
				var req TagRequirement
				if configMap, ok := tagConfig.(map[string]interface{}); ok {
					if pattern, exists := configMap["pattern"]; exists {
						if patternStr, ok := pattern.(string); ok {
							req.Pattern = patternStr
						}
					}
				}
				c.RequiredTags[tagName] = req
			}
		default:
			return fmt.Errorf("required_tags must be an array of strings or an object")
		}
	}

	return nil
}

// UnmarshalYAML implements custom YAML unmarshaling to support both array and object formats
func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	// First unmarshal the basic structure
	type configAlias struct {
		RequiredTags interface{}         `yaml:"required_tags"`
		Exemptions   []ResourceExemption `yaml:"exemptions"`
		ReportPath   string              `yaml:"report_path"`
	}
	
	var temp configAlias
	if err := value.Decode(&temp); err != nil {
		return err
	}

	// Copy non-required_tags fields
	c.Exemptions = temp.Exemptions
	c.ReportPath = temp.ReportPath
	c.RequiredTags = make(map[string]TagRequirement)

	// Handle required_tags field which can be array or object
	if temp.RequiredTags != nil {
		switch v := temp.RequiredTags.(type) {
		case []interface{}:
			// Array format (legacy)
			for _, item := range v {
				if tagName, ok := item.(string); ok {
					c.RequiredTags[tagName] = TagRequirement{}
				}
			}
		case map[string]interface{}:
			// Object format (new)
			for tagName, tagConfig := range v {
				var req TagRequirement
				if tagConfig == nil {
					// Empty object like "Name: {}"
					req = TagRequirement{}
				} else if configMap, ok := tagConfig.(map[string]interface{}); ok {
					if pattern, exists := configMap["pattern"]; exists {
						if patternStr, ok := pattern.(string); ok {
							req.Pattern = patternStr
						}
					}
				}
				c.RequiredTags[tagName] = req
			}
		default:
			return fmt.Errorf("required_tags must be an array of strings or an object")
		}
	}

	return nil
}

// compilePatterns compiles all regex patterns in the configuration
func (c *Config) compilePatterns() error {
	for tagName, req := range c.RequiredTags {
		if req.Pattern != "" {
			compiled, err := regexp.Compile(req.Pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern for tag '%s': %w", tagName, err)
			}
			// Update the requirement with compiled pattern
			req.compiledPattern = compiled
			c.RequiredTags[tagName] = req
		}
	}
	return nil
}

// populateLegacyRequired populates the legacy Required field for backward compatibility
func (c *Config) populateLegacyRequired() {
	c.Required = make([]string, 0, len(c.RequiredTags))
	for tagName := range c.RequiredTags {
		c.Required = append(c.Required, tagName)
	}
}

// ValidateTagValue validates a tag value against its pattern if one is defined
func (c *Config) ValidateTagValue(tagName, tagValue string) (bool, string) {
	// Find the tag requirement (case-sensitive or case-insensitive)
	var req TagRequirement
	var found bool
	
	if c.IgnoreTagCase {
		for name, requirement := range c.RequiredTags {
			if strings.EqualFold(name, tagName) {
				req = requirement
				found = true
				break
			}
		}
	} else {
		req, found = c.RequiredTags[tagName]
	}
	
	if !found || req.compiledPattern == nil {
		// No pattern defined, so value is valid
		return true, ""
	}
	
	if req.compiledPattern.MatchString(tagValue) {
		return true, ""
	}
	
	// Pattern validation failed
	errorMsg := fmt.Sprintf("value '%s' does not match required pattern '%s'", tagValue, req.Pattern)
	
	return false, errorMsg
}
