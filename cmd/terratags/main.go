package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"terratags/pkg/config"
	"terratags/pkg/parser"
	"terratags/pkg/validator"
)

func main() {
	var (
		configFile   string
		terraformDir string
		verbose      bool
	)

	flag.StringVar(&configFile, "config", "", "Path to the config file (JSON/YAML) containing required tag keys")
	flag.StringVar(&terraformDir, "dir", ".", "Path to the Terraform directory to analyze")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	if configFile == "" {
		fmt.Println("Error: Config file is required")
		fmt.Println("Usage: terratags -config <config_file.json|yaml> [-dir <terraform_directory>] [-verbose]")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Loaded configuration with %d required tags\n", len(cfg.Required))
	}

	// Find all Terraform files
	var terraformFiles []string
	err = filepath.Walk(terraformDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".tf") {
			terraformFiles = append(terraformFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Found %d Terraform files to analyze\n", len(terraformFiles))
	}

	// Parse Terraform files to extract resources and provider configurations
	var allResources []parser.Resource
	var allProviders []parser.ProviderConfig
	
	for _, file := range terraformFiles {
		if verbose {
			fmt.Printf("Analyzing file: %s\n", file)
		}

		// Extract resources
		resources, err := parser.ParseFile(file)
		if err != nil {
			fmt.Printf("Error parsing file %s for resources: %v\n", file, err)
		} else {
			allResources = append(allResources, resources...)
		}

		// Extract provider configurations
		providers, err := parser.ParseProviderBlocks(file)
		if err != nil {
			fmt.Printf("Error parsing file %s for providers: %v\n", file, err)
		} else {
			allProviders = append(allProviders, providers...)
		}
	}

	if verbose {
		fmt.Printf("Found %d taggable resources\n", len(allResources))
		fmt.Printf("Found %d provider configurations with default tags\n", len(allProviders))
	}

	// Validate required tags, considering provider default tags
	results := validator.ValidateRequiredTags(allResources, allProviders, cfg.Required)

	// Print results
	if len(results) > 0 {
		fmt.Println("\nTag validation issues found:")
		for _, result := range results {
			fmt.Printf("  - %s '%s': %s\n", result.ResourceType, result.ResourceName, result.Message)
		}
		fmt.Println("\nTag validation failed. Please fix the issues above.")
		os.Exit(1)
	} else {
		fmt.Println("All resources have the required tags!")
	}
}
