package validator

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/terratags/terratags/pkg/config"
	"github.com/terratags/terratags/pkg/logging"
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
	TotalResources           int
	CompliantResources       int
	FullyExemptResources     int
	PartiallyExemptResources int
	ExcludedAWSCCResources   []string
	ExcludedResourcesCount   int
	ViolationsByTag          map[string]int
}

// ValidateResources validates that all resources have the required tags
func ValidateResources(resources []parser.Resource, providers []parser.ProviderConfig, cfg *config.Config) (bool, []TagViolation, TagComplianceStats, []parser.Resource) {
	var violations []TagViolation
	stats := TagComplianceStats{
		ViolationsByTag: make(map[string]int),
	}
	valid := true

	// Count resources that are not excluded
	var nonExcludedResources int

	// We'll track excluded resources later during resource processing

	// Create a map of provider default tags by path
	defaultTagsByPath := make(map[string]map[string]string)
	for _, provider := range providers {
		defaultTagsByPath[provider.Path] = provider.DefaultTags
	}

	for _, resource := range resources {
		// Check if this is an excluded AWSCC resource
		if parser.AwsccExcludedResources[resource.Type] {
			// Add to excluded resources list if not already there
			found := false
			for _, excludedType := range stats.ExcludedAWSCCResources {
				if excludedType == resource.Type {
					found = true
					break
				}
			}
			if !found {
				stats.ExcludedAWSCCResources = append(stats.ExcludedAWSCCResources, resource.Type)
			}
			// Increment the excluded resources count
			stats.ExcludedResourcesCount++
			// Skip validation for this resource
			continue
		}

		// Count this as a non-excluded resource
		nonExcludedResources++

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
		var exemptTags []string
		var nonExemptMissingTags []string
		var exemptReason string

		for _, requiredTag := range cfg.Required {
			// Check if the tag is in the resource's tags
			tagExists := false
			if cfg.IgnoreTagCase {
				// Case-insensitive comparison
				requiredTagLower := strings.ToLower(requiredTag)
				for tagKey := range resource.Tags {
					if strings.ToLower(tagKey) == requiredTagLower {
						tagExists = true
						break
					}
				}
			} else {
				// Case-sensitive comparison (original behavior)
				_, tagExists = resource.Tags[requiredTag]
			}

			if !tagExists {
				// If not in resource tags, check if it's in default tags
				defaultTagExists := false
				if cfg.IgnoreTagCase {
					// Case-insensitive comparison for default tags
					requiredTagLower := strings.ToLower(requiredTag)
					for tagKey := range defaultTags {
						if strings.ToLower(tagKey) == requiredTagLower {
							defaultTagExists = true
							break
						}
					}
				} else {
					// Case-sensitive comparison (original behavior)
					_, defaultTagExists = defaultTags[requiredTag]
				}

				if !defaultTagExists {
					// Check if this resource is exempt from this tag requirement
					exempt, reason := cfg.IsExemptFromTag(resource.Type, resource.Name, requiredTag)
					if exempt {
						exemptTags = append(exemptTags, requiredTag)
						if exemptReason == "" {
							exemptReason = reason
						}
						// Add to missingTags so it shows up in the report
						missingTags = append(missingTags, requiredTag)
					} else {
						missingTags = append(missingTags, requiredTag)
						nonExemptMissingTags = append(nonExemptMissingTags, requiredTag)
						stats.ViolationsByTag[requiredTag]++
					}
				} else {
					// Tag exists in default_tags, so it's valid
					if len(defaultTags) > 0 {
						logging.Debug("Resource %s '%s' inherits tag '%s' from provider default_tags",
							resource.Type, resource.Name, requiredTag)
					}
				}
			}
		}

		// Determine if the resource has any exemptions
		isExempt := len(exemptTags) > 0

		// Determine if the resource is fully exempt (all missing tags are exempt)
		isFullyExempt := isExempt && len(nonExemptMissingTags) == 0 && len(missingTags) > 0

		// Determine if the resource is partially exempt (some missing tags are exempt, but others aren't)
		isPartiallyExempt := isExempt && len(nonExemptMissingTags) > 0

		// If the resource has any missing tags (exempt or not), add it to violations
		if len(missingTags) > 0 {
			// If there are any non-exempt missing tags, the resource is not fully compliant
			if len(nonExemptMissingTags) > 0 {
				valid = false
			}

			violations = append(violations, TagViolation{
				ResourceType: resource.Type,
				ResourceName: resource.Name,
				ResourcePath: resource.Path,
				MissingTags:  missingTags,
				IsExempt:     isExempt,
				ExemptReason: exemptReason,
			})

			// Update statistics based on exemption status
			if isFullyExempt {
				stats.FullyExemptResources++
			} else if isPartiallyExempt {
				stats.PartiallyExemptResources++
			}
		} else {
			// No missing tags, resource is compliant
			stats.CompliantResources++
		}
	}

	// Set the total resources to only count non-excluded resources
	stats.TotalResources = nonExcludedResources

	return valid, violations, stats, resources
}

// ValidateDirectory validates all Terraform files in a directory
func ValidateDirectory(dir string, cfg *config.Config, logLevel string) (bool, []TagViolation, TagComplianceStats, []parser.Resource) {
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

	logging.Info("Found %d Terraform files to analyze", len(files))

	var allResources []parser.Resource
	var allProviders []parser.ProviderConfig

	// Parse each file
	for _, file := range files {
		logging.Info("Analyzing file: %s", file)

		// Parse resources
		resources, err := parser.ParseFile(file, logLevel)
		if err != nil {
			logging.Warn("Error parsing file %s: %s", file, err)
			continue
		}
		allResources = append(allResources, resources...)

		// Parse provider blocks
		providers, err := parser.ParseProviderBlocks(file)
		if err != nil {
			logging.Warn("Error parsing provider blocks in %s: %s", file, err)
			continue
		}
		allProviders = append(allProviders, providers...)
	}

	logging.Info("Found %d taggable resources", len(allResources))
	logging.Info("Found %d provider configurations with default tags", len(allProviders))

	// Validate resources
	valid, violations, stats, _ := ValidateResources(allResources, allProviders, cfg)
	return valid, violations, stats, allResources
}

// ValidateTerraformPlan validates a Terraform plan file
func ValidateTerraformPlan(planPath string, cfg *config.Config, logLevel string) (bool, []TagViolation, TagComplianceStats, []parser.Resource) {
	logging.Info("Analyzing Terraform plan: %s", planPath)

	// Parse the plan
	resources, err := parser.ParseTerraformPlan(planPath, logLevel)
	if err != nil {
		return false, []TagViolation{{
			ResourceType: "error",
			ResourceName: "error",
			ResourcePath: planPath,
			MissingTags:  []string{fmt.Sprintf("Error parsing Terraform plan: %s", err)},
		}}, TagComplianceStats{}, nil
	}

	logging.Info("Found %d taggable resources in plan", len(resources))

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

	sb.WriteString("# AWS Provider\n")
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
	sb.WriteString("}\n\n")

	// Add Azure API provider suggestion
	sb.WriteString("# Azure API Provider\n")
	sb.WriteString("provider \"azapi\" {\n")
	sb.WriteString("  # Existing provider configuration preserved\n\n")
	sb.WriteString("  default_tags = {\n")

	// Add missing tags with placeholder values
	for _, tag := range missingTags {
		sb.WriteString(fmt.Sprintf("    %s = \"CHANGE_ME\"  # Added missing required tag\n", tag))
	}

	sb.WriteString("  }\n")
	sb.WriteString("}")

	return sb.String()
}

// GenerateHTMLReport generates an enhanced HTML report of tag compliance using html/template with Bootstrap styling
func GenerateHTMLReport(violations []TagViolation, stats TagComplianceStats, cfg *config.Config) string {
	// Calculate compliance percentage - only considering non-excluded resources
	compliancePercentage := 0.0
	if stats.TotalResources > 0 {
		compliancePercentage = float64(stats.CompliantResources) / float64(stats.TotalResources) * 100
	}

	// Define the template inline to avoid file I/O
	const tmplStr = `<!DOCTYPE html>
<html>
<head>
    <title>Terraform Tag Compliance Report</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
    <style>
        .header-logo { 
            max-height: 60px; 
            margin-right: 15px;
        }
        .github-link {
            margin-left: auto;
            text-decoration: none;
        }
        .progress { height: 30px; }
        .exempt-tag { color: #fd7e14; }
        .logo-svg {
            height: 60px;
            width: 60px;
            margin-right: 15px;
        }
    </style>
</head>
<body class="bg-light">
    <div class="container py-4">
        <!-- Header with Logo and GitHub Link -->
        <div class="d-flex align-items-center mb-4">
            <svg class="logo-svg" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
              <!-- Cloud shape representing infrastructure -->
              <path d="M120,90 C130,75 150,75 160,85 C170,75 190,80 190,100 C190,115 175,120 160,120 
                       C160,125 155,130 145,130 C135,130 90,130 80,130 C65,130 50,120 50,105 
                       C50,90 65,80 80,85 C85,70 105,70 120,90 Z" 
                    fill="#FFE0B2" stroke="#FF8C00" stroke-width="3"/>
              
              <!-- Tag symbol on the cloud -->
              <path d="M100,100 L115,100 L115,115 L107.5,122 L100,115 Z" 
                    fill="#FF8C00" stroke="#FF8C00" stroke-width="1"/>
              <circle cx="107.5" cy="105" r="2" fill="#FFE0B2"/>
              
              <!-- Magnifying glass rim -->
              <circle cx="70" cy="110" r="35" fill="none" stroke="#FF8C00" stroke-width="6"/>
              
              <!-- Magnifying glass handle -->
              <line x1="95" y1="135" x2="120" y2="160" stroke="#FF8C00" stroke-width="10" stroke-linecap="round"/>
              
              <!-- Magnifying glass lens highlight -->
              <circle cx="70" cy="110" r="28" fill="none" stroke="#FFE0B2" stroke-width="2" stroke-opacity="0.7"/>
              
              <!-- Scan lines in magnifying glass -->
              <line x1="50" y1="110" x2="90" y2="110" stroke="#FF8C00" stroke-width="2" stroke-opacity="0.5"/>
              <line x1="70" y1="90" x2="70" y2="130" stroke="#FF8C00" stroke-width="2" stroke-opacity="0.5"/>
            </svg>
            <h1 class="mb-0">Terraform Tag Compliance Report</h1>
            <a href="https://github.com/terratags/terratags" class="github-link" target="_blank">
                <img src="https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png" alt="GitHub" width="32">
            </a>
        </div>
        
        <p class="text-muted">Generated on: {{.GeneratedTime}}</p>
        
        <!-- Summary Card -->
        <div class="card mb-4">
            <div class="card-header bg-primary text-white">
                <h2 class="card-title h5 mb-0">Summary</h2>
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-md-3">
                        <div class="card text-center mb-3">
                            <div class="card-body">
                                <h3 class="h2">{{.Stats.TotalResources}}</h3>
                                <p class="mb-0">Total Resources</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2">
                        <div class="card text-center mb-3 bg-success text-white">
                            <div class="card-body">
                                <h3 class="h2">{{.Stats.CompliantResources}}</h3>
                                <p class="mb-0">Compliant</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2">
                        <div class="card text-center mb-3 bg-danger text-white">
                            <div class="card-body">
                                <h3 class="h2">{{.NonCompliantCount}}</h3>
                                <p class="mb-0">Non-compliant</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-2">
                        <div class="card text-center mb-3 bg-warning">
                            <div class="card-body">
                                <h3 class="h2">{{.TotalExemptResources}}</h3>
                                <p class="mb-0">Exempt</p>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-3">
                        <div class="card text-center mb-3 bg-info text-white">
                            <div class="card-body">
                                <h3 class="h2">{{.Stats.ExcludedResourcesCount}}</h3>
                                <p class="mb-0">Excluded</p>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- Exemption Details -->
                {{if gt .TotalExemptResources 0}}
                <div class="row mt-3">
                    <div class="col-12">
                        <div class="card">
                            <div class="card-header bg-warning">
                                <h3 class="card-title h6 mb-0">Exemption Details</h3>
                            </div>
                            <div class="card-body p-2">
                                <div class="row">
                                    <div class="col-md-6">
                                        <div class="d-flex justify-content-between align-items-center">
                                            <span>Fully Exempt Resources:</span>
                                            <span class="badge bg-warning">{{.Stats.FullyExemptResources}}</span>
                                        </div>
                                    </div>
                                    <div class="col-md-6">
                                        <div class="d-flex justify-content-between align-items-center">
                                            <span>Partially Exempt Resources:</span>
                                            <span class="badge bg-warning">{{.Stats.PartiallyExemptResources}}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                {{end}}
                
                <div class="progress mt-3">
                    <div class="progress-bar bg-success" role="progressbar" 
                         style="width: {{printf "%.1f" .CompliancePercentage}}%;" 
                         aria-valuenow="{{printf "%.1f" .CompliancePercentage}}" 
                         aria-valuemin="0" aria-valuemax="100">
                        {{printf "%.1f" .CompliancePercentage}}% Compliant
                    </div>
                </div>
            </div>
        </div>
        
        <!-- Required Tags Card -->
        <div class="card mb-4">
            <div class="card-header bg-info text-white">
                <h2 class="card-title h5 mb-0">Required Tags</h2>
            </div>
            <div class="card-body">
                <div class="row">
                    {{range .RequiredTags}}
                    <div class="col-md-3 mb-2">
                        <span class="badge bg-primary">{{.}}</span>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
        
        <!-- Violations by Tag -->
        <div class="card mb-4">
            <div class="card-header bg-danger text-white">
                <h2 class="card-title h5 mb-0">Violations by Tag</h2>
            </div>
            <div class="card-body">
                <table class="table table-striped">
                    <thead>
                        <tr>
                            <th>Tag</th>
                            <th>Violations</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range $tag, $count := .Stats.ViolationsByTag}}
                        <tr>
                            <td><code>{{$tag}}</code></td>
                            <td>{{$count}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <!-- Non-compliant Resources with Collapsible Sections -->
        <div class="card">
            <div class="card-header bg-secondary text-white">
                <h2 class="card-title h5 mb-0">Non-compliant Resources</h2>
            </div>
            <div class="card-body">
                {{if eq (len .Violations) 0}}
                <div class="alert alert-success">All resources are compliant!</div>
                {{else}}
                <div class="accordion" id="resourceAccordion">
                    {{range $index, $v := .Violations}}
                    <div class="accordion-item">
                        <h2 class="accordion-header" id="heading{{$index}}">
                            <button class="accordion-button {{if $v.IsExempt}}bg-warning{{else}}bg-danger text-white{{end}} collapsed" type="button" 
                                    data-bs-toggle="collapse" data-bs-target="#collapse{{$index}}" 
                                    aria-expanded="false" aria-controls="collapse{{$index}}">
                                {{$v.ResourceType}} "{{$v.ResourceName}}"
                                {{if $v.IsExempt}}
                                <span class="badge bg-warning ms-2">EXEMPT</span>
                                {{else}}
                                <span class="badge bg-danger ms-2">{{len $v.MissingTags}} missing tags</span>
                                {{end}}
                            </button>
                        </h2>
                        <div id="collapse{{$index}}" class="accordion-collapse collapse" 
                             aria-labelledby="heading{{$index}}" data-bs-parent="#resourceAccordion">
                            <div class="accordion-body">
                                <p><strong>Path:</strong> {{$v.ResourcePath}}</p>
                                {{if $v.IsExempt}}
                                <p><strong>Status:</strong> <span class="exempt-tag">EXEMPT</span> - {{$v.ExemptReason}}</p>
                                {{else}}
                                <p><strong>Missing Tags:</strong></p>
                                <ul>
                                    {{range $v.MissingTags}}
                                    <li><code>{{.}}</code></li>
                                    {{end}}
                                </ul>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>
        
        <!-- AWSCC Resources Excluded from Evaluation -->
        {{if .HasExcludedResources}}
        <div class="card mt-4">
            <div class="card-header bg-info text-white">
                <h2 class="card-title h5 mb-0">AWSCC Resources Excluded from Evaluation</h2>
            </div>
            <div class="card-body">
                <p>The following AWSCC resources are excluded from tag validation due to tag schema mismatches:</p>
                <div class="row">
                    {{range .Stats.ExcludedAWSCCResources}}
                    <div class="col-md-4 mb-2">
                        <code>{{.}}</code>
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
        {{end}}
        
        <footer class="mt-4 text-center text-muted">
            <p>Generated by <a href="https://github.com/terratags/terratags" target="_blank">Terratags</a></p>
        </footer>
    </div>
</body>
</html>`

	// Create template with a custom function for joining strings
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(tmplStr)

	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	// Prepare data for the template
	data := struct {
		GeneratedTime        string
		Stats                TagComplianceStats
		NonCompliantCount    int
		TotalExemptResources int
		CompliancePercentage float64
		RequiredTags         []string
		Violations           []TagViolation
		HasExcludedResources bool // Add this
	}{
		GeneratedTime:        time.Now().Format("2006-01-02 15:04:05"),
		Stats:                stats,
		NonCompliantCount:    stats.TotalResources - stats.CompliantResources - stats.FullyExemptResources - stats.PartiallyExemptResources,
		TotalExemptResources: stats.FullyExemptResources + stats.PartiallyExemptResources,
		CompliancePercentage: compliancePercentage,
		RequiredTags:         cfg.Required,
		Violations:           violations,
		HasExcludedResources: len(stats.ExcludedAWSCCResources) > 0,
	}

	// Create a buffer to store the rendered template
	var buf bytes.Buffer

	// Execute the template
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}
