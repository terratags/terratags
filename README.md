<p align="center">
  <img src="docs/assets/terratags-logo.svg" alt="Terratags Logo" height="80" style="vertical-align:middle">
  <span style="font-size:48px; font-weight:bold; vertical-align:middle">Terratags</span>
</p>

<p align="center">Terratags is a tool for validating tags on AWS, Azure, Google Cloud, and Alibaba Cloud resources in Terraform configurations.</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/terratags/terratags)](https://goreportcard.com/report/github.com/terratags/terratags)

## Features

- Validates required tags on AWS, Azure, Google Cloud, and Alibaba Cloud resources
- **Advanced pattern matching** with regex validation for tag values
- **Module resource validation** - validates resources created by external modules via Terraform plan analysis
- **Remote config files** - load config from HTTP/HTTPS URLs or Git repositories
- Supports AWS provider default_tags
- Supports AWSCC provider tag format (see [AWSCC exclusion list](https://github.com/terratags/terratags/blob/main/scripts/update_resources.go#L15) for resources with non-compliant tag schemas)
- Supports Azure providers (azurerm and azapi)
- Supports azapi provider default_tags
- Supports Google Cloud provider with labels (GCP uses 'labels' instead of 'tags')
- Supports Google provider default_labels
- Supports Google Cloud Beta provider (google-beta) with labels and default_labels
- Supports Alibaba Cloud provider with tags (uses same format as AWS)
- Supports module-level tags with tag inheritance
- Supports exemptions for specific resources
- Generates HTML reports of tag compliance
- Provides auto-remediation suggestions
- Integrates with Terraform plan output
- Tracks tag inheritance from provider default_tags
- Exemption tracking and reporting
- Excluded resources tracking for AWSCC resources with non-compliant tag schemas

Open issues for other providers:
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

### Remote Config Files

Terratags supports loading config files from remote locations:

```bash
# HTTP/HTTPS URL
terratags -config https://example.com/configs/terratags.yaml -dir ./infra

# Git repository (HTTPS)
terratags -config https://github.com/org/configs.git//terratags.yaml?ref=main -dir ./infra

# Git repository (SSH)
terratags -config git@github.com:org/configs.git//path/to/config.yaml?ref=v1.0.0 -dir ./infra
```

See [Remote Config Examples](examples/remote_config/README.md) for more details.

### Options

- `-config`, `-c`: Path or URL to the config file (JSON/YAML) containing required tag keys (required)
  - Supports local paths: `./config.yaml`, `/path/to/config.json`
  - Supports HTTP/HTTPS URLs: `https://example.com/config.yaml`
  - Supports Git URLs: `https://github.com/org/repo.git//path/to/config.yaml?ref=main`
- `-dir`, `-d`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`, `-v`: Enable verbose output
- `-log-level`, `-l`: Set logging level: DEBUG, INFO, WARN, ERROR (default: ERROR)
- `-plan`, `-p`: Path to Terraform plan JSON file to analyze (includes module resource validation)
- `-report`, `-r`: Path to output HTML report file
- `-remediate`, `-re`: Show auto-remediation suggestions for non-compliant resources
- `-exemptions`, `-e`: Path to exemptions file (JSON/YAML)
- `-ignore-case`, `-i`: Ignore case when comparing required tag keys
- `-help`, `-h`: Show help message
- `-version`, `-V`: Show version information

## Pattern Matching

> **Note**: Pattern validation was introduced in version 0.3.0 and provides advanced regex-based tag value validation.

Terratags supports advanced pattern validation using regular expressions to validate tag values. This allows you to enforce specific formats, naming conventions, and business rules for your tags.

### Pattern Validation Features

- **Regex Support**: Use any valid Go regex pattern for tag value validation
- **Case Sensitivity**: Patterns are case-sensitive by default (use `--ignore-case` for case-insensitive tag name matching)
- **Flexible Configuration**: Mix pattern validation with simple presence validation
- **Clear Error Messages**: Detailed feedback when patterns don't match
- **Backward Compatibility**: Existing simple configurations continue to work

### Configuration Format

Pattern validation uses an object format instead of the simple array format:

```yaml
required_tags:
  TagName:
    pattern: "regex_pattern_here"
  
  # Tag without pattern (just presence validation)
  SimpleTag: {}
```

### Common Pattern Examples

#### Environment Validation
```yaml
Environment:
  pattern: "^(dev|test|staging|prod)$"
```
- ✅ Matches: `dev`, `test`, `staging`, `prod`
- ❌ Rejects: `development`, `production`, `DEV`, `Test`

#### Email Validation
```yaml
Owner:
  pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```
- ✅ Matches: `devops@company.com`, `team.lead@company.com`
- ❌ Rejects: `username`, `user@domain`, `@company.com`

#### Project Code Format
```yaml
Project:
  pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"
```
- ✅ Matches: `WEB-123456`, `DATA-567890`, `INFRA-890123`
- ❌ Rejects: `web-123`, `PROJECT`, `ABC-12`, `TOOLONG-1234567`

#### Cost Center Format
```yaml
CostCenter:
  pattern: "^CC-[0-9]{4}$"
```
- ✅ Matches: `CC-1234`, `CC-5678`, `CC-9012`
- ❌ Rejects: `CC123`, `CC-12345`, `cc-1234`, `CostCenter-1234`

#### No Whitespace
```yaml
Name:
  pattern: "^\\S+$"
```
- ✅ Matches: `web-server-01`, `data-bucket`, `main-vpc`
- ❌ Rejects: `web server`, `database 01`, `api gateway`

#### Semantic Versioning
```yaml
Version:
  pattern: "^v?[0-9]+\\.[0-9]+\\.[0-9]+$"
```
- ✅ Matches: `1.0.0`, `v2.1.3`, `10.15.2`
- ❌ Rejects: `1.0`, `v1`, `1.0.0-beta`, `latest`

#### Alphanumeric with Dashes
```yaml
ResourceName:
  pattern: "^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$"
```
- ✅ Matches: `web-server`, `api-gateway-v2`, `database01`
- ❌ Rejects: `-web-server`, `api-gateway-`, `web--server`

### Pattern Validation Examples

#### Complete Configuration
```yaml
required_tags:
  # Strict environment values
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  # Valid email for ownership
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  
  # Project code format
  Project:
    pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"
  
  # Cost center format
  CostCenter:
    pattern: "^CC-[0-9]{4}$"
  
  # No whitespace in names
  Name:
    pattern: "^\\S+$"
  
  # Simple presence validation (no pattern)
  Team: {}
```

#### Mixed Validation
```yaml
required_tags:
  # Pattern validation for critical tags
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  
  # Simple validation for others
  Name: {}
  Project: {}
  Team: {}
```

### Error Messages

When pattern validation fails, Terratags provides clear error messages:

```
Resource aws_instance 'web_server' has tag pattern violations:
  - Tag 'Environment': value 'Production' does not match required pattern '^(dev|test|staging|prod)$'
  - Tag 'Owner': value 'DevOps Team' does not match required pattern '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'
  - Tag 'Project': value 'website' does not match required pattern '^[A-Z]{2,4}-[0-9]{3,6}$'
```

### Testing Patterns

Use the provided examples to test pattern validation:

```bash
# Test passing patterns
terratags -config examples/config-patterns.yaml -dir examples/pattern_validation_passing

# Test failing patterns
terratags -config examples/config-patterns.yaml -dir examples/pattern_validation_failing

# Generate detailed report
terratags -config examples/config-patterns.yaml -dir examples/pattern_validation_failing -report pattern-report.html
```

### Pattern Development Tips

1. **Start Simple**: Begin with basic patterns and refine as needed
2. **Test Thoroughly**: Use both passing and failing examples to validate patterns
3. **Escape Special Characters**: Remember to escape backslashes in YAML/JSON (`\\` instead of `\`)
4. **Case Sensitivity**: Patterns are case-sensitive unless using `--ignore-case` flag
5. **Anchor Patterns**: Use `^` and `$` to match the entire string
6. **Test Online**: Use regex testing tools to validate patterns before deployment

### Regex Reference

Common regex elements used in tag patterns:

| Element | Description | Example |
|---------|-------------|---------|
| `^` | Start of string | `^dev` matches strings starting with "dev" |
| `$` | End of string | `prod$` matches strings ending with "prod" |
| `\S` | Non-whitespace character | `^\S+$` matches strings without spaces |
| `[a-z]` | Character class | `[a-zA-Z]` matches any letter |
| `[0-9]` | Digit class | `[0-9]{4}` matches exactly 4 digits |
| `{n,m}` | Quantifier | `{2,4}` matches 2 to 4 occurrences |
| `+` | One or more | `[a-z]+` matches one or more letters |
| `*` | Zero or more | `[a-z]*` matches zero or more letters |
| `\|` | Alternation | `(dev\|test\|prod)` matches any of the three |
| `\.` | Literal dot | `\.com` matches ".com" literally |
| `\\` | Escape character | `\\S` in YAML becomes `\S` in regex |

### Migration from Simple Format

Existing simple configurations work unchanged:

```yaml
# This continues to work
required_tags:
  - Name
  - Environment
  - Owner
```

To add pattern validation, convert to object format:

```yaml
# Enhanced with patterns
required_tags:
  Name: {}  # Just presence validation
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

### Quick Pattern Reference

| Use Case | Pattern | Example Values |
|----------|---------|----------------|
| Environment | `^(dev\|test\|staging\|prod)$` | `dev`, `test`, `staging`, `prod` |
| Email | `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$` | `devops@company.com`, `team.lead@company.com` |
| Project Code | `^[A-Z]{2,4}-[0-9]{3,6}$` | `WEB-123456`, `DATA-567890`, `SEC-123456` |
| Cost Center | `^CC-[0-9]{4}$` | `CC-1234`, `CC-5678`, `CC-9012` |
| No Whitespace | `^\\S+$` | `web-server-01`, `data-bucket`, `main-vpc` |
| Version | `^v?[0-9]+\\.[0-9]+\\.[0-9]+$` | `1.0.0`, `v2.1.3`, `10.15.2` |

### Advanced Pattern Matching

For comprehensive pattern matching documentation, including advanced techniques, common patterns library, and troubleshooting guide, see [Pattern Matching Guide](docs/pattern-matching.md).

## Configuration

### Required Tags Configuration

Terratags supports two configuration formats for specifying required tags: a simple format for basic tag presence validation, and an advanced format with regex pattern validation for tag values.

#### Simple Format (Legacy - Fully Supported)

The simple format validates that required tags are present but doesn't validate their values:

**YAML Format:**
```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

**JSON Format:**
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

#### Advanced Format with Pattern Validation

The advanced format allows you to specify regex patterns to validate tag values:

**YAML Format:**
```yaml
required_tags:
  Name:
    pattern: "^\\S+$"
  
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  
  Project:
    pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"
  
  # Tag required but no pattern validation
  Team: {}
```

**JSON Format:**
```json
{
  "required_tags": {
    "Name": {
      "pattern": "^\\S+$"
    },
    "Environment": {
      "pattern": "^(dev|test|staging|prod)$"
    },
    "Owner": {
      "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    },
    "Project": {
      "pattern": "^[A-Z]{2,4}-[0-9]{3,6}$"
    },
    "Team": {}
  }
}
```

#### Mixed Format

You can mix both simple and advanced formats in the same configuration:

```yaml
required_tags:
  # Simple tags (just check presence)
  Name: {}
  Project: {}
  
  # Advanced tags with pattern validation
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

#### Pattern Validation Features

- **Regex Support**: Use any valid Go regex pattern
- **Case Sensitivity**: Respects the `--ignore-case` flag for tag name matching
- **Backward Compatibility**: Existing simple configurations continue to work unchanged

#### Common Patterns

Here are some commonly used patterns:

```yaml
required_tags:
  # Environment validation
  Environment:
    pattern: "^(dev|test|staging|prod|production)$"
  
  # Email validation
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  
  # Cost center format
  CostCenter:
    pattern: "^CC-[0-9]{4}$"
  
  # Semantic version
  Version:
    pattern: "^v?[0-9]+\\.[0-9]+\\.[0-9]+$"
  
  # No whitespace
  Name:
    pattern: "^\\S+$"
  
  # Alphanumeric with dashes
  Project:
    pattern: "^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$"
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

### Google Cloud Provider Support

Terratags supports the Google Cloud provider for GCP resources. GCP uses 'labels' instead of 'tags' for resource metadata.

#### Google Provider

The Google provider uses a map of key/value pairs for labels:

```terraform
resource "google_compute_instance" "example" {
  name         = "example-instance"
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
  }

  labels = {
    Environment = "Production"
    Project     = "Terratags"
    Name        = "example-instance"
  }
}
```

#### Google Provider with default_labels

The Google provider supports default_labels at the provider level:

```terraform
provider "google" {
  project = "my-project-id"
  region  = "us-central1"

  default_labels = {
    Environment = "Production"
    Owner       = "team-a"
  }
}

resource "google_storage_bucket" "example" {
  name     = "example-bucket"
  location = "US"

  labels = {
    Name    = "example-bucket"
    Project = "demo"
  }
}
```

In this example, the bucket will have all four required labels: `Name` and `Project` from the resource-level labels, and `Environment` and `Owner` from the provider's default_labels.

**Note:** GCP uses 'labels' instead of 'tags', but terratags treats them the same way for validation purposes.

## Examples

### Basic Usage (Direct Resources Only)

```bash
terratags -config config.yaml -dir ./infra
```

### Validate Terraform Plan (Includes Module Resources)

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

### Generate HTML Report

```bash
terratags -config config.yaml -dir ./infra -report report.html
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
       rev: v0.3.0  # Use the latest version (available from v0.3.0+)
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
    rev: v0.3.0
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
    rev: v0.3.0
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
# GitHub Actions example - Directory validation
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
        
      - name: Validate Tags (Direct Resources)
        run: terratags -config config.yaml -dir ./infra
```

For comprehensive validation including module resources, use plan-based validation. See [CI/CD documentation](docs/ci-cd.md) for detailed examples.

## Enhanced Reporting

Terratags provides enhanced HTML reports with detailed information about tag compliance:

- Visual indicators for compliant, non-compliant, and exempt resources
- Detailed breakdown of tag status for each resource
- **Separate sections for direct resources and module-created resources**
- **Module path and source information for module resources**
- Tracking of tag sources (resource-level vs provider default_tags vs module inheritance)
- Exemption details including reasons for exemptions
- Summary statistics including exempt resources
- Tag violation counts by tag name

The HTML report provides a visual representation of tag compliance across your Terraform resources, making it easy to identify which resources need attention and track compliance metrics. When using plan validation, the report clearly distinguishes between direct resources and those created by modules. You can view the generated HTML report in any web browser.

![Sample Terratags Report](docs/assets/sample_report.png)
