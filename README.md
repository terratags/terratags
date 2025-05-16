# Terratags

Terratags is a tool for validating tags on AWS resources in Terraform configurations.

## Features

- Validates required tags on AWS resources
- Supports AWS provider default_tags
- Supports module-level tags
- Supports exemptions for specific resources
- Generates HTML reports of tag compliance
- Provides auto-remediation suggestions
- Integrates with Terraform plan output
- Tracks tag inheritance from provider default_tags

## Installation

```bash
go install github.com/terratags/terratags@latest
```

## Usage

```bash
terratags -config config.yaml -dir ./terraform
```

### Options

- `-config`: Path to the config file (JSON/YAML) containing required tag keys (required)
- `-dir`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`: Enable verbose output
- `-plan`: Path to Terraform plan JSON file to analyze
- `-report`: Path to output HTML report file
- `-remediate`: Show auto-remediation suggestions for non-compliant resources
- `-exemptions`: Path to exemptions file (JSON/YAML)

## Configuration

### Required Tags

Create a YAML or JSON file with the required tags:

```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

### Exemptions

Create a YAML or JSON file with exemptions:

```yaml
exemptions:
  - resource_type: aws_s3_bucket
    resource_name: logs_bucket
    exempt_tags: [Owner, Project]
    reason: "Legacy bucket used for system logs only"
    expires_at: "2025-12-31"
```

## Examples

### Basic Usage

```bash
terratags -config config.yaml -dir ./terraform
```

### Generate HTML Report

```bash
terratags -config config.yaml -dir ./terraform -report report.html
```

### Validate Terraform Plan

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

### Show Auto-remediation Suggestions

```bash
terratags -config config.yaml -dir ./terraform -remediate
```

### Use Exemptions

```bash
terratags -config config.yaml -dir ./terraform -exemptions exemptions.yaml
```

## Integration with CI/CD

Add Terratags to your CI/CD pipeline to enforce tag compliance:

```yaml
# GitHub Actions example
name: Validate Tags

on:
  pull_request:
    paths:
      - '**.tf'

jobs:
  validate-tags:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
          
      - name: Install Terratags
        run: go install github.com/terratags/terratags@latest
        
      - name: Validate Tags
        run: terratags -config config.yaml -dir ./terraform
```

## License

MIT