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

- `-config`, `-c`: Path to the config file (JSON/YAML) containing required tag keys (required)
- `-dir`, `-d`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`, `-v`: Enable verbose output
- `-plan`, `-p`: Path to Terraform plan JSON file to analyze
- `-report`, `-r`: Path to output HTML report file
- `-remediate`, `-m`: Show auto-remediation suggestions for non-compliant resources
- `-exemptions`, `-e`: Path to exemptions file (JSON/YAML)
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
          go-version: '1.24'
          
      - name: Install Terratags
        run: go install github.com/terratags/terratags@latest
        
      - name: Validate Tags
        run: terratags -config config.yaml -dir ./terraform
```

## Sample Report

When you generate an HTML report with Terratags, it will look similar to this:

```
┌─────────────────────────────────────────────────────┐
│           Terraform Tag Compliance Report           │
├─────────────────────────────────────────────────────┤
│ Generated on: 2025-05-16                           │
│                                                     │
│ Summary:                                            │
│ ✓ Total Resources: 4                                │
│ ✓ Compliant Resources: 2                            │
│ ✗ Non-compliant Resources: 2                        │
│                                                     │
│ [████████████████████████████████--------] 50.0%     │
│                                                     │
│ Non-compliant Resources:                            │
│ ✗ aws_s3_bucket "data_bucket"                       │
│   Missing Tags: Environment, Owner, Project         │
└─────────────────────────────────────────────────────┘
```

The HTML report provides a visual representation of tag compliance across your Terraform resources, making it easy to identify which resources need attention and track compliance metrics. You can view the generated HTML report in any web browser.