package parser

import (
	"encoding/json"
	"testing"
)

func TestExtractTagsFromPlanResource_DefaultTags(t *testing.T) {
	tests := []struct {
		name     string
		resource map[string]any
		expected map[string]string
	}{
		{
			name: "AWS resource with default_tags (tags_all present)",
			resource: map[string]any{
				"tags": map[string]any{
					"Name": "My bucket",
				},
				"tags_all": map[string]any{
					"Environment": "dev",
					"Name":        "My bucket",
				},
			},
			expected: map[string]string{
				"Environment": "dev",
				"Name":        "My bucket",
			},
		},
		{
			name: "AWS resource without default_tags (only tags)",
			resource: map[string]any{
				"tags": map[string]any{
					"Name":        "My bucket",
					"Environment": "prod",
				},
			},
			expected: map[string]string{
				"Name":        "My bucket",
				"Environment": "prod",
			},
		},
		{
			name: "AWSCC resource with tags as list",
			resource: map[string]any{
				"tags": []any{
					map[string]any{
						"key":   "Environment",
						"value": "test",
					},
					map[string]any{
						"key":   "Name",
						"value": "Test Resource",
					},
				},
			},
			expected: map[string]string{
				"Environment": "test",
				"Name":        "Test Resource",
			},
		},
		{
			name: "Google resource with default_labels (effective_labels present)",
			resource: map[string]any{
				"labels": map[string]any{
					"name": "test-instance",
				},
				"effective_labels": map[string]any{
					"environment": "dev",
					"name":        "test-instance",
				},
			},
			expected: map[string]string{
				"environment": "dev",
				"name":        "test-instance",
			},
		},
		{
			name: "Google resource without default_labels (only labels)",
			resource: map[string]any{
				"labels": map[string]any{
					"environment": "dev",
					"name":        "test-instance",
				},
			},
			expected: map[string]string{
				"environment": "dev",
				"name":        "test-instance",
			},
		},
		{
			name: "Resource with no tags",
			resource: map[string]any{
				"bucket": "my-bucket",
			},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTagsFromPlanResource(tt.resource)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d tags, got %d", len(tt.expected), len(result))
			}
			
			for key, expectedValue := range tt.expected {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("Expected tag %s not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("Expected tag %s to have value %s, got %s", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestExtractTagsFromPlanResource_RealPlanData(t *testing.T) {
	// Test with actual Terraform plan JSON structure
	planJSON := `{
		"resource_changes": [
			{
				"address": "aws_s3_bucket.example",
				"type": "aws_s3_bucket",
				"name": "example",
				"change": {
					"actions": ["create"],
					"after": {
						"bucket": "my-tf-test-bucket-example",
						"tags": {
							"Name": "My bucket"
						},
						"tags_all": {
							"Environment": "dev",
							"Name": "My bucket"
						}
					}
				}
			}
		]
	}`
	
	var plan struct {
		ResourceChanges []struct {
			Change struct {
				After map[string]any `json:"after"`
			} `json:"change"`
		} `json:"resource_changes"`
	}
	
	err := json.Unmarshal([]byte(planJSON), &plan)
	if err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}
	
	result := extractTagsFromPlanResource(plan.ResourceChanges[0].Change.After)
	
	expected := map[string]string{
		"Environment": "dev",
		"Name":        "My bucket",
	}
	
	if len(result) != len(expected) {
		t.Errorf("Expected %d tags, got %d", len(expected), len(result))
	}
	
	for key, expectedValue := range expected {
		if actualValue, exists := result[key]; !exists {
			t.Errorf("Expected tag %s not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected tag %s to have value %s, got %s", key, expectedValue, actualValue)
		}
	}
}
