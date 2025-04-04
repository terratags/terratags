package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

// ProviderConfig represents a Terraform provider configuration
type ProviderConfig struct {
	Name       string
	DefaultTags map[string]string
	Path       string
}

// ParseProviderBlocks parses a Terraform file and extracts provider configurations
func ParseProviderBlocks(path string) ([]ProviderConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	var providers []ProviderConfig

	// Create a schema that focuses on provider blocks
	content2, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "provider",
				LabelNames: []string{"name"},
			},
		},
	})

	if diags.HasErrors() {
		fmt.Printf("Warning: Some provider blocks in %s couldn't be parsed, but we'll continue with what we can parse\n", path)
	}

	for _, block := range content2.Blocks {
		if block.Type == "provider" {
			providerName := block.Labels[0]
			
			// We're only interested in AWS providers that might have default_tags
			if strings.HasPrefix(providerName, "aws") {
				defaultTags := extractDefaultTags(block)
				if len(defaultTags) > 0 {
					providers = append(providers, ProviderConfig{
						Name:       providerName,
						DefaultTags: defaultTags,
						Path:       path,
					})
				}
			}
		}
	}

	return providers, nil
}

// extractDefaultTags extracts default_tags from an AWS provider block
func extractDefaultTags(block *hcl.Block) map[string]string {
	defaultTags := make(map[string]string)

	// Look for default_tags block
	content, _ := block.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "default_tags",
			},
		},
	})

	for _, defaultTagsBlock := range content.Blocks {
		if defaultTagsBlock.Type == "default_tags" {
			// Extract tags attribute from default_tags block
			defaultTagsContent, _ := defaultTagsBlock.Body.Content(&hcl.BodySchema{
				Attributes: []hcl.AttributeSchema{
					{Name: "tags", Required: false},
				},
			})

			if attr, exists := defaultTagsContent.Attributes["tags"]; exists {
				// Create an evaluation context
				ctx := &hcl.EvalContext{
					Variables: map[string]cty.Value{},
				}
				
				val, diags := attr.Expr.Value(ctx)
				if diags.HasErrors() {
					fmt.Printf("Error evaluating default_tags: %s\n", diags.Error())
				} else {
					// Try to process as map if possible
					if val.Type().IsMapType() || val.Type().IsObjectType() {
						val.ForEachElement(func(key cty.Value, value cty.Value) bool {
							if key.Type() == cty.String {
								keyStr := key.AsString()
								defaultTags[keyStr] = "dummy-value"
								if value.Type() == cty.String {
									defaultTags[keyStr] = value.AsString()
								}
							}
							return true
						})
					}
				}
			}
		}
	}

	return defaultTags
}
