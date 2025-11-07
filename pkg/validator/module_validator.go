package validator

import (
	"strings"

	"github.com/terratags/terratags/pkg/config"
	"github.com/terratags/terratags/pkg/parser"
)

// ValidationSummary provides basic validation summary
type ValidationSummary struct {
	TotalResources       int
	CompliantResources   int
	CompliancePercentage float64
}

// ModuleResourceValidation represents validation result for a module resource
type ModuleResourceValidation struct {
	ResourceValidation
	ModulePath   string
	ModuleName   string
	ModuleSource string
	IsExternal   bool
}

// ValidationResultWithModules includes both direct and module resource validation
type ValidationResultWithModules struct {
	DirectResources []ResourceValidation
	ModuleResources []ModuleResourceValidation
	Summary         ValidationSummaryWithModules
}

// ValidationSummaryWithModules provides summary including module resources
type ValidationSummaryWithModules struct {
	ValidationSummary
	ModuleCompliant   int
	ModuleTotal       int
	TotalCompliant    int
	TotalResources    int
	CompliancePercent float64
}

// ValidateWithModules validates both direct and module resources
func ValidateWithModules(directResources []parser.Resource, moduleResources []parser.ModuleResource,
	cfg *config.Config, providerTags map[string]map[string]string) ValidationResultWithModules {

	var result ValidationResultWithModules

	// Validate direct resources (existing logic)
	for _, resource := range directResources {
		validation := validateResource(resource, cfg, providerTags)
		result.DirectResources = append(result.DirectResources, validation)
	}

	// Validate module resources
	for _, moduleResource := range moduleResources {
		validation := validateModuleResource(moduleResource, cfg, providerTags)
		result.ModuleResources = append(result.ModuleResources, validation)
	}

	// Calculate summary
	result.Summary = calculateSummaryWithModules(result.DirectResources, result.ModuleResources)

	return result
}

// validateModuleResource validates a single module resource
func validateModuleResource(moduleResource parser.ModuleResource, cfg *config.Config,
	providerTags map[string]map[string]string) ModuleResourceValidation {

	baseValidation := validateResource(moduleResource.Resource, cfg, providerTags)

	return ModuleResourceValidation{
		ResourceValidation: baseValidation,
		ModulePath:         moduleResource.ModulePath,
		ModuleName:         moduleResource.ModuleName,
		ModuleSource:       moduleResource.ModuleSource,
		IsExternal:         isExternalModule(moduleResource.ModuleSource),
	}
}

// calculateSummary calculates validation summary for direct resources
func calculateSummary(results []ResourceValidation) ValidationSummary {
	compliant := 0
	for _, result := range results {
		if result.IsCompliant {
			compliant++
		}
	}

	total := len(results)
	percentage := 0.0
	if total > 0 {
		percentage = float64(compliant) / float64(total) * 100
	}

	return ValidationSummary{
		TotalResources:       total,
		CompliantResources:   compliant,
		CompliancePercentage: percentage,
	}
}

// calculateSummaryWithModules calculates validation summary including module resources
func calculateSummaryWithModules(directResults []ResourceValidation, moduleResults []ModuleResourceValidation) ValidationSummaryWithModules {
	directSummary := calculateSummary(directResults)

	moduleCompliant := 0
	for _, result := range moduleResults {
		if result.IsCompliant {
			moduleCompliant++
		}
	}

	totalCompliant := directSummary.CompliantResources + moduleCompliant
	totalResources := directSummary.TotalResources + len(moduleResults)
	compliancePercent := 0.0
	if totalResources > 0 {
		compliancePercent = float64(totalCompliant) / float64(totalResources) * 100
	}

	return ValidationSummaryWithModules{
		ValidationSummary: directSummary,
		ModuleCompliant:   moduleCompliant,
		ModuleTotal:       len(moduleResults),
		TotalCompliant:    totalCompliant,
		TotalResources:    totalResources,
		CompliancePercent: compliancePercent,
	}
}

// isExternalModule checks if a module is external (not local)
func isExternalModule(source string) bool {
	// External modules typically start with registry paths or git URLs
	return !strings.HasPrefix(source, "./") &&
		!strings.HasPrefix(source, "../") &&
		!strings.HasPrefix(source, "/")
}
