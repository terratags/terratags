# Usage

Terratags can be used in various ways to validate tags on AWS resources in your Terraform configurations.

## Basic Usage

The basic usage of Terratags is:

```bash
terratags -config config.yaml -dir ./infra
```

This command will analyze all Terraform files in the specified directory and validate that AWS resources have the required tags as defined in your configuration file.

## Command Examples

### Generate HTML Report

Generate a detailed HTML report of tag compliance:

```bash
terratags -config config.yaml -dir ./infra -report report.html
```

### Validate Terraform Plan

Validate tags in a Terraform plan output:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

### Show Auto-remediation Suggestions

Get suggestions for fixing non-compliant resources:

```bash
terratags -config config.yaml -dir ./infra -remediate
```

### Use Exemptions

Apply exemptions to specific resources:

```bash
terratags -config config.yaml -dir ./infra -exemptions exemptions.yaml
```

## Verbose Output

For more detailed output, use the `-verbose` flag:

```bash
terratags -config config.yaml -dir ./infra -verbose
```

This will show additional information about the validation process, including:

- Files being analyzed
- Resources being checked
- Tag inheritance from default_tags
- Exemptions being applied

## Exit Codes

Terratags uses the following exit codes:

- `0`: All resources are compliant with tagging requirements
- `1`: One or more resources are missing required tags
- `2`: Error in configuration or execution

This makes it easy to integrate Terratags into CI/CD pipelines and fail builds when tag requirements are not met.

## Working with Large Codebases

For large Terraform codebases, you can:

1. Run Terratags on specific directories:
   ```bash
   terratags -config config.yaml -dir ./infra/modules/networking
   ```

2. Use the plan-based approach to only validate resources that are changing:
   ```bash
   terraform plan -out=tfplan -target=module.networking
   terraform show -json tfplan > plan.json
   terratags -config config.yaml -plan plan.json
   ```

## HTML Reports

The HTML report provides a visual representation of tag compliance across your Terraform resources, making it easy to identify which resources need attention and track compliance metrics.

To generate a report:

```bash
terratags -config config.yaml -dir ./infra -report report.html
```

The report includes:

- Overall compliance percentage
- List of compliant and non-compliant resources
- Missing tags for each non-compliant resource
- Summary statistics

## Log Levels

Terratags supports different log levels to control the verbosity of output:

```bash
terratags -config config.yaml -dir ./infra -log-level INFO
```

Available log levels:

- `DEBUG`: Shows all debug information, including detailed tag discovery
- `INFO`: Shows informational messages (same as using the `-verbose` flag)
- `WARN`: Shows only warnings and errors
- `ERROR`: Shows only errors (default)

For backward compatibility, the `-verbose` flag is equivalent to `-log-level INFO`.

## Case-Insensitive Tag Matching

By default, Terratags performs case-sensitive matching for tag keys. To enable case-insensitive matching, use the `-ignore-case` flag:

```bash
terratags -config config.yaml -dir ./infra -ignore-case
```

With this option enabled, tag keys like "Environment", "ENVIRONMENT", and "environment" will all match a required tag key "Environment".