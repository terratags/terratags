# Terratags

Terratags is a tool for validating tags on resources in Terraform configurations.

## Overview

Consistent tagging is crucial for:

- Cost allocation and tracking
- Resource ownership identification
- Environment classification
- Security and compliance requirements
- Automation and resource management

Terratags helps enforce these tagging standards across your infrastructure, ensuring that all resources are properly tagged according to your organization's policies.

## Key Features

- **Tag Validation**: Validates required tags on AWS and Azure resources
- **Default Tags Support**: Supports AWS provider default_tags
- **AWSCC Support**: Supports AWSCC provider tag format ( Refer [exclusion list](https://github.com/terratags/terratags/blob/main/scripts/update_resources.go#L15) for resources with non compliant tag schema)
- **Azure Support**: Supports Azure providers (azurerm and azapi)
- **Module-Level Tags**: Supports module-level tags
- **Exemption Support**: Supports exemptions for specific resources
- **HTML Reports**: Generates HTML reports of tag compliance
- **Auto-Remediation**: Provides auto-remediation suggestions
- **Plan Integration**: Integrates with Terraform plan output
- **Tag Inheritance**: Tracks tag inheritance from provider default_tags
- **Exemption Tracking**: Tracks and reports on exemptions
- **Excluded Resources**: Tracks AWSCC resources with non-compliant tag schemas

Open issues for other providers:
- [Google provider](https://github.com/terratags/terratags/issues/8)
- [Azure providers](https://github.com/terratags/terratags/issues/7) : Keeping this open as there are additional Azure providers.

## Not validated

- The behavior with provider aliases is not tested and so the evaluation cannot be guaranteed.

## Quick Start

### Installation

```bash
# Using Homebrew
brew install terratags/tap/terratags

# Using Go
go install github.com/terratags/terratags@latest
```

### Basic Usage

```bash
terratags -config config.yaml -dir ./infra
```

Check out the [documentation](configuration.md) for more detailed information on configuration and usage.

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

The HTML report provides a visual representation of tag compliance across your Terraform resources, making it easy to identify which resources need attention and track compliance metrics.
