package parser

import (
	"path/filepath"

	"github.com/terratags/terratags/pkg/logging"
)

// ModuleTagInheritance handles tag inheritance from module calls to resources
type ModuleTagInheritance struct {
	moduleTags map[string]map[string]string // module name -> tags
}

// NewModuleTagInheritance creates a new tag inheritance handler
func NewModuleTagInheritance() *ModuleTagInheritance {
	return &ModuleTagInheritance{
		moduleTags: make(map[string]map[string]string),
	}
}

// LoadModuleTags loads tags from module calls in Terraform files
func (m *ModuleTagInheritance) LoadModuleTags(terraformDir string) error {
	files, err := filepath.Glob(filepath.Join(terraformDir, "*.tf"))
	if err != nil {
		return err
	}

	for _, file := range files {
		resources, err := ParseFile(file, "ERROR")
		if err != nil {
			logging.Debug("Skipping file %s due to parse error: %v", file, err)
			continue // Skip files with parse errors
		}

		for _, resource := range resources {
			if resource.Type == "module" {
				m.moduleTags[resource.Name] = resource.Tags
				logging.Debug("Loaded tags for module %s: %v", resource.Name, resource.Tags)
			}
		}
	}

	return nil
}

// InheritTags applies module tags to module resources
func (m *ModuleTagInheritance) InheritTags(moduleResource *ModuleResource) {
	if moduleResource == nil {
		return
	}
	
	if moduleTags, exists := m.moduleTags[moduleResource.ModuleName]; exists {
		// Merge module tags with resource tags (resource tags take precedence)
		for key, value := range moduleTags {
			if _, exists := moduleResource.Tags[key]; !exists {
				moduleResource.Tags[key] = value
				moduleResource.TagSources[key] = TagSource{
					Source: "module_call",
					Value:  value,
				}
				logging.Debug("Inherited tag %s=%s for resource %s from module %s", 
					key, value, moduleResource.Name, moduleResource.ModuleName)
			}
		}
	}
}
