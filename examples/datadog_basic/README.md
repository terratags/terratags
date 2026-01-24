# Datadog Examples

This directory contains examples demonstrating how to use Terratags with Datadog resources.

## Examples

### 1. Basic Datadog Resources (`datadog_basic/`)

Demonstrates basic tag validation for common Datadog resources:
- `datadog_monitor` - Monitoring alerts
- `datadog_dashboard` - Dashboards  
- `datadog_synthetics_test` - Synthetic tests

**Usage:**
```bash
terratags -config examples/datadog_basic/config.yaml -dir examples/datadog_basic/
```

### 2. Provider Default Tags (`datadog_default_tags/`)

Shows how to use Datadog provider's `default_tags` feature to automatically apply tags to all resources:

```hcl
provider "datadog" {
  default_tags {
    tags = {
      Environment = "production"
      Team        = "platform"
    }
  }
}
```

**Usage:**
```bash
terratags -config examples/datadog_default_tags/config.yaml -dir examples/datadog_default_tags/
```

### 3. Pattern Validation (`datadog_patterns/`)

Demonstrates advanced pattern validation using regex to enforce tag value formats:

```yaml
required_tags:
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  Team:
    pattern: "^[a-z][a-z0-9-]*[a-z0-9]$"
  Service:
    pattern: "^[a-z][a-z0-9-]*[a-z0-9]$"
  Version:
    pattern: "^v?[0-9]+\\.[0-9]+\\.[0-9]+$"
```

**Usage:**
```bash
# This will show validation failures for pattern violations
terratags -config examples/datadog_patterns/config.yaml -dir examples/datadog_patterns/
```

## Datadog Tag Format

Datadog uses a different tag format compared to AWS/Azure:

```hcl
# Datadog format: list of "key:value" strings
resource "datadog_monitor" "example" {
  tags = [
    "Environment:production",
    "Team:platform",
    "Service:web-api"
  ]
}

# Compare to AWS format: map of key = "value"
resource "aws_instance" "example" {
  tags = {
    Environment = "production"
    Team        = "platform"
    Service     = "web-api"
  }
}
```

## Supported Datadog Resources

Terratags supports 21+ Datadog resources that have a `tags` attribute, including:

- **Monitoring**: `datadog_monitor`, `datadog_dashboard`, `datadog_service_level_objective`
- **Synthetics**: `datadog_synthetics_test`, `datadog_synthetics_global_variable`, `datadog_synthetics_private_location`
- **Security**: `datadog_security_monitoring_rule`, `datadog_appsec_waf_custom_rule`, `datadog_sensitive_data_scanner_rule`
- **Logs**: `datadog_logs_custom_pipeline`
- **Integrations**: `datadog_integration_confluent_account`, `datadog_integration_fastly_service`
- **And more...**

For the complete list, see `pkg/parser/datadog_taggable_resources.go`.
