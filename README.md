<p align="center">
  <img src="docs/assets/terratags-logo.svg" alt="Terratags Logo" height="80" style="vertical-align:middle">
  <span style="font-size:48px; font-weight:bold; vertical-align:middle">Terratags</span>
</p>

<p align="center">Terratags is a tool for validating tags on AWS and Azure resources in Terraform configurations.</p>

## Features

- Validates required tags on AWS and Azure resources
- Supports AWS provider default_tags
- Supports AWSCC provider tag format ( Refer [exclusion list]([https://github.com/terratags/terratags/issues/9](https://github.com/terratags/terratags/blob/main/scripts/update_resources.go#L15)) for resources with non compliant tag schema)
- Supports Azure providers (azurerm and azapi)
- Supports azapi provider default_tags
- Supports module-level tags
- Supports exemptions for specific resources
- Generates HTML reports of tag compliance
- Provides auto-remediation suggestions
- Integrates with Terraform plan output
- Tracks tag inheritance from provider default_tags
- Exemption tracking and reporting
- Excluded resources tracking for AWSCC resources with non-compliant tag schemas

Open issues for other providers:
- [Google provider](https://github.com/terratags/terratags/issues/8)
- [Azure providers](https://github.com/terratags/terratags/issues/7) : Keeping this open as there are additional Azure providers.

## Not validated

- The behavior with provider aliases is not tested and so the evaluation cannot be guaranteed.

## Installation

### Using Homebrew

```bash
brew install terratags/tap/terratags
```

### Using Go

```bash
go install github.com/terratags/terratags@latest
```

### Binary Download

Download the appropriate binary from the [GitHub Releases](https://github.com/terratags/terratags/releases) page.

See [installation docs](docs/installation.md) for more options.

## Usage

```bash
terratags -config config.yaml -dir ./infra
```

### Options

- `-config`, `-c`: Path to the config file (JSON/YAML) containing required tag keys (required)
- `-dir`, `-d`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`, `-v`: Enable verbose output
- `-log-level`, `-l`: Set logging level: DEBUG, INFO, WARN, ERROR (default: ERROR)
- `-plan`, `-p`: Path to Terraform plan JSON file to analyze
- `-report`, `-r`: Path to output HTML report file
- `-remediate`, `-re`: Show auto-remediation suggestions for non-compliant resources
- `-exemptions`, `-e`: Path to exemptions file (JSON/YAML)
- `-ignore-case`, `-i`: Ignore case when comparing required tag keys
- `-help`, `-h`: Show help message
- `-version`, `-V`: Show version information

## Configuration

### Required Tags Configuration

Terratags requires a configuration file that specifies which tags must be present on your AWS resources. This file can be in either YAML or JSON format.

#### YAML Format

```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

#### JSON Format

```json
{
  "required_tags": [
    "Name",
    "Environment",
    "Owner",
    "Project"
  ]
}
```

### Exemptions Configuration

Exemptions allow you to exclude specific resources or resource types from certain tag requirements. Create a YAML or JSON file with your exemptions.

#### Exemption Fields

- `resource_type`: The AWS resource type (e.g., aws_s3_bucket, aws_instance)
- `resource_name`: The name of the specific resource to exempt. Use "*" to exempt all resources of the specified type
- `exempt_tags`: List of tags that are not required for this resource
- `reason`: A description explaining why this exemption exists

#### YAML Example

```yaml
exemptions:
  - resource_type: aws_s3_bucket
    resource_name: logs_bucket
    exempt_tags: [Owner, Project]
    reason: "Legacy bucket used for system logs only"
  
  - resource_type: aws_dynamodb_table
    resource_name: "*"
    exempt_tags: [Environment]
    reason: "DynamoDB tables use environment from provider default_tags"
```

### AWSCC Provider Support

Terratags supports the AWS Cloud Control (AWSCC) provider's tag format, which differs from the standard AWS provider tag format.

#### AWS vs AWSCC Tag Formats

**AWS Provider** uses a map of key/value pairs:
```hcl
tags = {
  Environment = "test"
  Project     = "demo"
}
```

**AWSCC Provider** uses a list of maps with `key` and `value` fields:
```hcl
tags = [{
  key   = "Environment"
  value = "test"
}, {
  key   = "Project"
  value = "demo"
}]
```

**Note:** AWSCC provider does not support the `default_tags` feature. Some AWSCC resources have non-compliant tag schemas and are excluded from validation.

See [AWSCC Support](docs/awscc_support.md) for more details.

### Provider Default Tags Support

Terratags integrates with AWS provider's `default_tags` feature. When you define default tags in your AWS provider configuration, Terratags will recognize these tags and consider them when validating resources.

#### How Default Tags Work

1. Tags defined in the AWS provider's `default_tags` block are automatically applied to all resources created by that provider
2. Terratags tracks tag inheritance from provider default_tags to individual resources
3. Resources only need to specify tags not covered by default_tags
4. Default tags can be overridden at the resource level if needed

#### Example with Default Tags

```terraform
provider "aws" {
  region = "us-west-2"
  
  default_tags {
    tags = {
      Environment = "dev"
      Owner       = "team-a"
      Project     = "demo"
    }
  }
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  # Only need to specify Name tag, as other required tags come from default_tags
  tags = {
    Name = "example-instance"
  }
}
```

In this example, the AWS instance will have all four required tags: `Name` from the resource-level tags, and `Environment`, `Owner`, and `Project` from the provider's default_tags.

### AWSCC Provider Tag Support

Terratags also supports the AWSCC provider's tag format, which differs from the standard AWS provider format.

#### AWSCC Tag Format

The AWSCC provider uses a list of maps with key/value pairs for tags:

```terraform
resource "awscc_apigateway_rest_api" "example" {
  name        = "example-api"
  description = "Example API"
  
  tags = [{
    key   = "Environment"
    value = "Production"
  }, {
    key   = "Project"
    value = "Terratags"
  }]
}
```

**Note:** AWSCC provider does not support `default_tags`, so all required tags must be specified at the resource level.

See [AWSCC Support documentation](docs/awscc_support.md) for more details.

### Azure Providers Support

Terratags supports both the Azurerm and azapi providers for Azure resources.

#### Azurerm Provider

The Azurerm provider uses a map of key/value pairs for tags, similar to the AWS provider:

```terraform
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
  
  tags = {
    Environment = "Production"
    Project     = "Terratags"
  }
}
```

**Note:** Azurerm provider does not support `default_tags`, so all required tags must be specified at the resource level.

#### azapi Provider 

The azapi provider supports tags at both the provider level (via `default_tags`) and at the resource level:

```terraform
provider "azapi" {
  default_tags = {
    Environment = "Production"
    Project     = "Terratags"
  }
}

resource "azapi_resource" "example" {
  type      = "Microsoft.Storage/storageAccounts@2022-05-01"
  name      = "examplestorageaccount"
  parent_id = azurerm_resource_group.example.id
  location  = azurerm_resource_group.example.location
  
  tags = {
    Name = "example-storage"
  }
}
```

See [Azure Support documentation](docs/azure-support.md) for more details.

## Examples

### Basic Usage

```bash
terratags -config config.yaml -dir ./infra
```

### Generate HTML Report

```bash
terratags -config config.yaml -dir ./infra -report report.html
```

### Validate Terraform Plan

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

### Show Auto-remediation Suggestions

```bash
terratags -config config.yaml -dir ./infra -remediate
```

### Use Exemptions

```bash
terratags -config config.yaml -dir ./infra -exemptions exemptions.yaml
```

## Pre-commit Hook Integration

Terratags can be integrated with [pre-commit](https://pre-commit.com/) to automatically validate tags before commits are made to your repository.

### Quick Setup

1. Install pre-commit:
   ```bash
   pip install pre-commit
   ```

2. Add to your `.pre-commit-config.yaml`:
   ```yaml
   repos:
     - repo: https://github.com/terratags/terratags
       rev: v1.x.x  # Use the latest version
       hooks:
         - id: terratags
   ```

3. Install the hook:
   ```bash
   pre-commit install
   ```

4. Ensure you have a `terratags.yaml` configuration file in your repository root

### Advanced Configuration

You can customize the hook with additional arguments:

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v1.x.x
    hooks:
      - id: terratags
        args: [
          --config=custom-config.yaml,
          --exemptions=exemptions.yaml,
          --remediate
        ]
```

### Multiple Hook Configurations

Define different hooks for different purposes:

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v1.x.x
    hooks:
      # Basic validation on every commit
      - id: terratags
        name: terratags-validate
        args: [--config=terratags.yaml]
      
      # Generate report (run manually)
      - id: terratags
        name: terratags-report
        args: [--config=terratags.yaml, --report=tag-report.html]
        stages: [manual]
```

See [Pre-commit Documentation](docs/pre-commit.md) for detailed setup instructions and advanced configurations.

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
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          
      - name: Install Terratags
        run: go install github.com/terratags/terratags@latest
        
      - name: Validate Tags
        run: terratags -config config.yaml -dir ./infra
```

## Enhanced Reporting

Terratags now provides enhanced HTML reports with detailed information about tag compliance:

- Visual indicators for compliant, non-compliant, and exempt resources
- Detailed breakdown of tag status for each resource
- Tracking of tag sources (resource-level vs provider default_tags)
- Exemption details including reasons for exemptions
- Summary statistics including exempt resources
- Tag violation counts by tag name

The HTML report provides a visual representation of tag compliance across your Terraform resources, making it easy to identify which resources need attention and track compliance metrics. You can view the generated HTML report in any web browser.

![Sample Terratags Report](docs/assets/sample_report.png)
