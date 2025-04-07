package main

import (
	"flag"
	"fmt"
	"os"

	"terratags/pkg/config"
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

	// Validate the directory
	valid, issues := validator.ValidateDirectory(terraformDir, cfg, verbose)

	// Print results
	if !valid {
		fmt.Println("\nTag validation issues found:")
		for _, issue := range issues {
			fmt.Println(issue)
		}
		fmt.Println("\nTag validation failed. Please fix the issues above.")
		os.Exit(1)
	} else {
		fmt.Println("All resources have the required tags!")
	}
}
