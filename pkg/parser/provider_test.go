package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseProviderBlocks_GoogleBeta(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []ProviderConfig
	}{
		{
			name: "Google Beta provider with default_labels",
			content: `
provider "google-beta" {
  project = "my-project"
  region  = "us-central1"
  
  default_labels = {
    environment = "test"
    team        = "platform"
    project     = "terratags"
  }
}
`,
			expected: []ProviderConfig{
				{
					Name: "google-beta",
					DefaultTags: map[string]string{
						"environment": "test",
						"team":        "platform",
						"project":     "terratags",
					},
				},
			},
		},
		{
			name: "Both Google and Google Beta providers",
			content: `
provider "google" {
  project = "my-project"
  region  = "us-central1"
  
  default_labels = {
    environment = "prod"
    team        = "backend"
  }
}

provider "google-beta" {
  project = "my-project"
  region  = "us-central1"
  
  default_labels = {
    environment = "beta"
    team        = "frontend"
    experimental = "true"
  }
}
`,
			expected: []ProviderConfig{
				{
					Name: "google",
					DefaultTags: map[string]string{
						"environment": "prod",
						"team":        "backend",
					},
				},
				{
					Name: "google-beta",
					DefaultTags: map[string]string{
						"environment": "beta",
						"team":        "frontend",
						"experimental": "true",
					},
				},
			},
		},
		{
			name: "Google Beta provider without default_labels",
			content: `
provider "google-beta" {
  project = "my-project"
  region  = "us-central1"
}
`,
			expected: []ProviderConfig{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "main.tf")
			
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Parse the provider blocks
			providers, err := ParseProviderBlocks(tmpFile)
			if err != nil {
				t.Fatalf("ParseProviderBlocks failed: %v", err)
			}

			// Check the number of providers
			if len(providers) != len(tt.expected) {
				t.Errorf("Expected %d providers, got %d", len(tt.expected), len(providers))
				return
			}

			// Check each provider
			for i, expected := range tt.expected {
				if i >= len(providers) {
					t.Errorf("Missing provider at index %d", i)
					continue
				}

				provider := providers[i]
				
				if provider.Name != expected.Name {
					t.Errorf("Expected provider name %s, got %s", expected.Name, provider.Name)
				}

				if len(provider.DefaultTags) != len(expected.DefaultTags) {
					t.Errorf("Expected %d default tags, got %d", len(expected.DefaultTags), len(provider.DefaultTags))
				}

				for key, expectedValue := range expected.DefaultTags {
					if actualValue, exists := provider.DefaultTags[key]; !exists {
						t.Errorf("Missing default tag %s", key)
					} else if actualValue != expectedValue {
						t.Errorf("Expected default tag %s=%s, got %s=%s", key, expectedValue, key, actualValue)
					}
				}
			}
		})
	}
}
