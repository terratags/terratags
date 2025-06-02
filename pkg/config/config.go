package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for tag validation
type Config struct {
	Required      []string            `json:"required_tags" yaml:"required_tags"`
	Exemptions    []ResourceExemption `json:"exemptions" yaml:"exemptions"`
	ReportPath    string              `json:"report_path" yaml:"report_path"`
	IgnoreTagCase bool                `json:"-" yaml:"-"` // Runtime option, not from config file
}

// ResourceExemption represents a resource that is exempt from tag requirements
type ResourceExemption struct {
	ResourceType string   `json:"resource_type" yaml:"resource_type"`
	ResourceName string   `json:"resource_name" yaml:"resource_name"`
	ExemptTags   []string `json:"exempt_tags" yaml:"exempt_tags"`
	Reason       string   `json:"reason" yaml:"reason"`
}

// LoadConfig loads the configuration from a JSON or YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	ext := strings.ToLower(filepath.Ext(path))

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

	return &config, nil
}

// LoadExemptions loads exemptions from a JSON or YAML file
func LoadExemptions(path string) ([]ResourceExemption, error) {
	data, err := os.ReadFile(path)
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
