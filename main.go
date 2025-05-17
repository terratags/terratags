package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/terratags/terratags/pkg/config"
	"github.com/terratags/terratags/pkg/parser"
	"github.com/terratags/terratags/pkg/validator"
)

// version is set during build time using ldflags
// Build with: go build -ldflags "-X main.version=0.1.0" -o terratags main.go
var version = "dev"

// Custom usage function to display both long and short forms of flags
func printUsage() {
	version, err := getVersion()
	if err != nil {
		version = "unknown"
	}
	fmt.Fprintf(os.Stderr, "Terratags v%s - AWS Resource Tag Validator for Terraform\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage: terratags [OPTIONS]\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  --config, -c <file>       Path to the config file (JSON/YAML) containing required tag keys\n")
	fmt.Fprintf(os.Stderr, "  --dir, -d <directory>     Path to the Terraform directory to analyze (default: \".\")\n")
	fmt.Fprintf(os.Stderr, "  --verbose, -v             Enable verbose output\n")
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
		verbose        bool
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

	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")

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

	// Show version if requested
	if showVersion {
		version, err := getVersion()
		if err != nil {
			fmt.Printf("Error reading version: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Terratags v%s\n", version)
		os.Exit(0)
	}

	// Show help if requested or if no arguments were provided
	if showHelp || len(os.Args) <= 1 {
		printUsage()
		os.Exit(0)
	}

	if configFile == "" {
		fmt.Println("Error: Config file is required")
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Load exemptions if provided
	if exemptionsFile != "" {
		exemptions, err := config.LoadExemptions(exemptionsFile)
		if err != nil {
			fmt.Printf("Error loading exemptions: %v\n", err)
			os.Exit(1)
		}
		cfg.Exemptions = exemptions
		if verbose {
			fmt.Printf("Loaded %d exemptions\n", len(exemptions))
		}
	}

	if verbose {
		fmt.Printf("Loaded configuration with %d required tags\n", len(cfg.Required))
	}

	// Determine which validation to run
	var valid bool
	var violations []validator.TagViolation
	var stats validator.TagComplianceStats
	var resources []parser.Resource

	if planFile != "" {
		// Validate the Terraform plan
		if verbose {
			fmt.Printf("Validating Terraform plan: %s\n", planFile)
		}
		valid, violations, stats, resources = validator.ValidateTerraformPlan(planFile, cfg, verbose)
	} else {
		// Validate the directory
		if verbose {
			fmt.Printf("Validating Terraform directory: %s\n", terraformDir)
		}
		valid, violations, stats, resources = validator.ValidateDirectory(terraformDir, cfg, verbose)
	}

	// Generate HTML report if requested
	if reportFile != "" {
		reportContent := validator.GenerateHTMLReport(violations, stats, cfg)
		reportDir := filepath.Dir(reportFile)
		if reportDir != "." {
			if err := os.MkdirAll(reportDir, 0755); err != nil {
				fmt.Printf("Error creating report directory: %v\n", err)
			}
		}
		if err := os.WriteFile(reportFile, []byte(reportContent), 0644); err != nil {
			fmt.Printf("Error writing report file: %v\n", err)
		} else {
			fmt.Printf("Report written to %s\n", reportFile)
		}
	}

	// Print results
	if !valid {
		fmt.Println("\nTag validation issues found:")
		for _, violation := range violations {
			fmt.Printf("Resource %s '%s' is missing required tags: %s\n",
				violation.ResourceType, violation.ResourceName, strings.Join(violation.MissingTags, ", "))

			// Show auto-remediation suggestions if requested
			if autoRemediate {
				fmt.Println("\nSuggested remediation:")

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
				fmt.Println(remediation)

				// Suggest provider default_tags update if appropriate
				if strings.HasPrefix(violation.ResourceType, "aws_") {
					fmt.Println("\nAlternatively, consider using provider default_tags:")
					fmt.Println(validator.SuggestProviderDefaultTagsUpdate(violation.MissingTags))
				}
			}
		}

		// Print summary statistics
		fmt.Printf("\nSummary: %d/%d resources compliant (%.1f%%)\n",
			stats.CompliantResources,
			stats.TotalResources,
			float64(stats.CompliantResources)/float64(stats.TotalResources)*100)

		if stats.ExemptResources > 0 {
			fmt.Printf("%d resources exempt from validation\n", stats.ExemptResources)
		}

		fmt.Println("\nTag validation failed. Please fix the issues above.")
		os.Exit(1)
	} else {
		fmt.Println("All resources have the required tags!")
	}
}

// getVersion returns the version of the application
// The version is set at build time using ldflags
// Example: go build -ldflags "-X main.version=0.1.0" -o terratags main.go
func getVersion() (string, error) {
	if version != "" {
		return version, nil
	}
	return "dev", nil
}
