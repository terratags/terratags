# Pre-commit Hook Integration

Terratags can be integrated with [pre-commit](https://pre-commit.com/) to automatically validate tags before commits are made to your repository.

## Prerequisites

1. Install pre-commit:
   ```bash
   pip install pre-commit
   ```

2. Ensure you have a terratags configuration file in your repository (see [Configuration](../README.md#required-tags-configuration))

## Basic Setup

1. Create or update your `.pre-commit-config.yaml` file in your repository root:

   ```yaml
   repos:
     - repo: https://github.com/terratags/terratags
       rev: v0.3.0  # Use the latest version (available from v0.3.0+)
       hooks:
         - id: terratags
   ```

2. Install the pre-commit hook:
   ```bash
   pre-commit install
   ```

3. Create your `terratags.yaml` configuration file:
   ```yaml
   required_tags:
     - Environment
     - Owner
     - Project
   ```

## Advanced Configuration

### Custom Configuration File

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      - id: terratags
        args: [--config=custom-config.yaml]
```

### Generate HTML Report

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      - id: terratags
        args: [--config=terratags.yaml, --report=tag-report.html]
```



### Show Remediation Suggestions

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      - id: terratags
        args: [--config=terratags.yaml, --remediate]
```

### Use Exemptions

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      - id: terratags
        args: [--config=terratags.yaml, --exemptions=exemptions.yaml]
```

### Custom Directory

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      - id: terratags
        args: [--config=terratags.yaml, --dir=./infrastructure]
```

## Multiple Hook Configurations

You can define multiple hooks for different purposes:

```yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      # Basic validation on every commit
      - id: terratags
        name: terratags-validate
        args: [--config=terratags.yaml]
      
      # Generate report (manual stage)
      - id: terratags
        name: terratags-report
        args: [--config=terratags.yaml, --report=reports/tags.html]
        stages: [manual]
      
      # Show remediation suggestions
      - id: terratags
        name: terratags-remediate
        args: [--config=terratags.yaml, --remediate]
        stages: [manual]
```

## Usage Examples

### Basic Workflow

1. Make changes to your Terraform files
2. Attempt to commit:
   ```bash
   git add .
   git commit -m "Add new resources"
   ```
3. Terratags will automatically run and validate your tags
4. If validation fails, fix the issues and commit again

### Manual Report Generation

```bash
# Run terratags report generation
pre-commit run terratags-report --hook-stage manual

# Run terratags with remediation suggestions
pre-commit run terratags-remediate --hook-stage manual
```

### Skip Hook for Emergency Commits

```bash
# Skip all pre-commit hooks
git commit -m "Emergency fix" --no-verify

# Skip only terratags hook
SKIP=terratags git commit -m "Skip terratags validation"
```

## File Filtering

The terratags pre-commit hook is configured to only run on Terraform configuration files:
- `*.tf` files

This ensures the hook only runs when relevant files are changed, improving performance.

**Note:** Pre-commit hooks validate Terraform source files only. For validating Terraform plan output, use the `--plan` flag in your CI/CD pipeline as described in the main [README](../README.md#integration-with-cicd).

## Troubleshooting

### Hook Not Running

1. Ensure pre-commit is installed: `pre-commit --version`
2. Ensure hooks are installed: `pre-commit install`
3. Check your `.pre-commit-config.yaml` syntax

### Configuration File Not Found

1. Ensure your terratags configuration file exists in the repository root
2. Use the `--config` argument to specify a custom path
3. Check the file name matches what you specified in args

### Validation Failures

1. Use `--remediate` to see suggested fixes
2. Check exemptions if certain resources should be excluded
3. Review your required tags configuration

## Integration with CI/CD

Pre-commit hooks work well with CI/CD pipelines. You can run the same validations in your CI:

```yaml
# GitHub Actions example
- name: Run pre-commit
  uses: pre-commit/action@v3.0.0
```

This ensures that even if developers skip local pre-commit hooks, the validation still runs in CI.
