package parser

import (
	"fmt"
	"strings"
)

// ModuleResource represents a resource created by a module
type ModuleResource struct {
	Resource
	ModulePath   string // e.g., "module.vpc", "module.vpc.module.subnets"
	ModuleName   string // e.g., "vpc"
	ModuleSource string // e.g., "terraform-aws-modules/vpc/aws"
}

// extractModuleName extracts the top-level module name from module address
func extractModuleName(moduleAddress string) string {
	// moduleAddress format: "module.vpc" or "module.vpc.module.subnets"
	parts := strings.Split(moduleAddress, ".")
	if len(parts) >= 2 && parts[0] == "module" {
		return parts[1]
	}
	return moduleAddress
}

// getModuleSource retrieves the source of a module from configuration
func getModuleSource(moduleName string, moduleCalls map[string]struct {
	Source  string `json:"source"`
	Version string `json:"version,omitempty"`
}) string {
	if moduleCall, exists := moduleCalls[moduleName]; exists {
		if moduleCall.Version != "" {
			return fmt.Sprintf("%s@%s", moduleCall.Source, moduleCall.Version)
		}
		return moduleCall.Source
	}
	return "unknown"
}

// isExternalModule checks if a module is external (not local)
func isExternalModule(source string) bool {
	// External modules typically start with registry paths or git URLs
	return !strings.HasPrefix(source, "./") &&
		!strings.HasPrefix(source, "../") &&
		!strings.HasPrefix(source, "/")
}
