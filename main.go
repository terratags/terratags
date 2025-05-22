package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/terratags/terratags/pkg/config"
	"github.com/terratags/terratags/pkg/logging"
	"github.com/terratags/terratags/pkg/parser"
	"github.com/terratags/terratags/pkg/validator"
)

// version is set during build time using ldflags
// Build with: go build -ldflags "-X main.version=0.1.0" -o terratags main.go
var version = "dev"

// Custom usage function to display both long and short forms of flags
func printUsage() {
	version, _, err := getVersion()
	if err != nil {
		version = "unknown"
	}
	fmt.Fprintf(os.Stderr, "Terratags v%s - AWS Resource Tag Validator for Terraform\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage: terratags [OPTIONS]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  --config, -c <file>       Path to the config file (JSON/YAML) containing required tag keys\n")
	fmt.Fprintf(os.Stderr, "  --dir, -d <directory>     Path to the Terraform directory to analyze (default: \".\")\n")
	fmt.Fprintf(os.Stderr, "  --log-level, -l <level>   Set logging level: DEBUG, INFO, WARN, ERROR (default: ERROR)\n")
	fmt.Fprintf(os.Stderr, "  --verbose, -v             Enable verbose output (same as --log-level=INFO)\n")
	fmt.Fprintf(os.Stderr, "  --plan, -p <file>         Path to Terraform plan JSON file to analyze\n")
	fmt.Fprintf(os.Stderr, "  --report, -r <file>       Path to output HTML report file\n")
	fmt.Fprintf(os.Stderr, "  --remediate, -re          Show auto-remediation suggestions for non-compliant resources\n")
	fmt.Fprintf(os.Stderr, "  --exemptions, -e <file>   Path to exemptions file (JSON/YAML)\n")
	fmt.Fprintf(os.Stderr, "  --help, -h                Show this help message\n")
	fmt.Fprintf(os.Stderr, "  --version, -V             Show version information\n")
}

func main() {
	var (
		configFile     string
		terraformDir   string
		logLevel       string
		planFile       string
		reportFile     string
		autoRemediate  bool
		exemptionsFile string
		showHelp       bool
		showVersion    bool
	)

	// Define flags with both long and short forms
	flag.StringVar(&configFile, "config", "", "Path to the config file (JSON/YAML) containing required tag keys")
	flag.StringVar(&configFile, "c", "", "Path to the config file (JSON/YAML) containing required tag keys")

	flag.StringVar(&terraformDir, "dir", ".", "Path to the Terraform directory to analyze")
	flag.StringVar(&terraformDir, "d", ".", "Path to the Terraform directory to analyze")

	flag.StringVar(&logLevel, "log-level", "ERROR", fmt.Sprintf("Log level (options: %s)", strings.Join(logging.ValidLogLevels, ", ")))
	flag.StringVar(&logLevel, "l", "ERROR", "Log level")

	// Keep verbose flag for backward compatibility
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output (same as --log-level=INFO)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (same as --log-level=INFO)")

	flag.StringVar(&planFile, "plan", "", "Path to Terraform plan JSON file to analyze")
	flag.StringVar(&planFile, "p", "", "Path to Terraform plan JSON file to analyze")

	flag.StringVar(&reportFile, "report", "", "Path to output HTML report file")
	flag.StringVar(&reportFile, "r", "", "Path to output HTML report file")

	flag.BoolVar(&autoRemediate, "remediate", false, "Show auto-remediation suggestions for non-compliant resources")
	flag.BoolVar(&autoRemediate, "re", false, "Show auto-remediation suggestions for non-compliant resources")

	flag.StringVar(&exemptionsFile, "exemptions", "", "Path to exemptions file (JSON/YAML)")
	flag.StringVar(&exemptionsFile, "e", "", "Path to exemptions file (JSON/YAML)")

	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message")

	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "V", false, "Show version information")

	// Override default usage function
	flag.Usage = printUsage

	flag.Parse()

	// Handle verbose flag for backward compatibility
	if verbose && logLevel == "ERROR" {
		logLevel = "INFO"
	}

	// Initialize logging
	if err := logging.Initialize(logLevel); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Show version if requested
	if showVersion {
		version, platform, err := getVersion()
		if err != nil {
			logging.Error("Error reading version: %v", err)
			os.Exit(1)
		}
		fmt.Printf("Terratags v%s (%s)\n", version, platform)
		os.Exit(0)
	}

	// Show help if requested or if no arguments were provided
	if showHelp || len(os.Args) <= 1 {
		printUsage()
		os.Exit(0)
	}

	if configFile == "" {
		logging.Error("Error: Config file is required")
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Error("Error loading config: %v", err)
		os.Exit(1)
	}

	// Load exemptions if provided
	if exemptionsFile != "" {
		exemptions, err := config.LoadExemptions(exemptionsFile)
		if err != nil {
			logging.Error("Error loading exemptions: %v", err)
			os.Exit(1)
		}
		cfg.Exemptions = exemptions
		logging.Info("Loaded %d exemptions", len(exemptions))
	}

	logging.Info("Loaded configuration with %d required tags", len(cfg.Required))

	// Determine which validation to run
	var valid bool
	var violations []validator.TagViolation
	var stats validator.TagComplianceStats
	var resources []parser.Resource

	if planFile != "" {
		// Validate the Terraform plan
		logging.Info("Validating Terraform plan: %s", planFile)
		valid, violations, stats, resources = validator.ValidateTerraformPlan(planFile, cfg, logLevel)
	} else {
		// Validate the directory
		logging.Info("Validating Terraform directory: %s", terraformDir)
		valid, violations, stats, resources = validator.ValidateDirectory(terraformDir, cfg, logLevel)
	}

	// Generate HTML report if requested
	if reportFile != "" {
		reportContent := validator.GenerateHTMLReport(violations, stats, cfg)
		reportDir := filepath.Dir(reportFile)
		if reportDir != "." {
			if err := os.MkdirAll(reportDir, 0755); err != nil {
				logging.Error("Error creating report directory: %v", err)
			}
		}
		if err := os.WriteFile(reportFile, []byte(reportContent), 0644); err != nil {
			logging.Error("Error writing report file: %v", err)
		} else {
			logging.Print("Report written to %s", reportFile)
		}
	}

	// Print results
	if !valid {
		logging.Print("\nTag validation issues found:")
		for _, violation := range violations {
			logging.Print("Resource %s '%s' is missing required tags: %s",
				violation.ResourceType, violation.ResourceName, strings.Join(violation.MissingTags, ", "))

			// Show auto-remediation suggestions if requested
			if autoRemediate {
				logging.Print("\nSuggested remediation:")

				// Get existing tags for this resource
				existingTags := make(map[string]string)
				for _, resource := range resources {
					if resource.Type == violation.ResourceType && resource.Name == violation.ResourceName {
						existingTags = resource.Tags
						break
					}
				}

				// Generate remediation code
				remediation := validator.GenerateRemediationCode(
					violation.ResourceType,
					violation.ResourceName,
					violation.ResourcePath,
					violation.MissingTags,
					existingTags)
				logging.Print("%s", remediation)

				// Suggest provider default_tags/default_labels update if appropriate
				if strings.HasPrefix(violation.ResourceType, "aws_") {
					logging.Print("\nAlternatively, consider using provider default_tags:")
					logging.Print("%s", validator.SuggestProviderDefaultTagsUpdate(violation.MissingTags))
				} else if strings.HasPrefix(violation.ResourceType, "google_") {
					logging.Print("\nAlternatively, consider using provider default_labels:")
					logging.Print("%s", validator.SuggestProviderDefaultLabelsUpdate(violation.MissingTags))
				}
			}
		}

		// Print summary statistics
		logging.Print("\nSummary: %d/%d resources compliant (%.1f%%)",
			stats.CompliantResources,
			stats.TotalResources,
			float64(stats.CompliantResources)/float64(stats.TotalResources)*100)

		totalExemptResources := stats.FullyExemptResources + stats.PartiallyExemptResources
		if totalExemptResources > 0 {
			logging.Print("%d resources exempt from validation (%d fully exempt, %d partially exempt)",
				totalExemptResources, stats.FullyExemptResources, stats.PartiallyExemptResources)
		}

		logging.Print("\nTag validation failed. Please fix the issues above.")
		os.Exit(1)
	} else {
		logging.Print("All resources have the required tags!")
	}
}

// getVersion returns the version and platform information of the application
// The version is set at build time using ldflags
// Example: go build -ldflags "-X main.version=0.1.0" -o terratags main.go
func getVersion() (string, string, error) {
	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	if version != "" {
		return version, platform, nil
	}
	return "dev", platform, nil
}
