# Provider Default Tags Regression Test

This directory contains test cases to ensure terratags correctly handles provider default tags/labels in Terraform plan validation.

## Test Cases

### 1. AWS with Default Tags (`test_default_tags/`)
- Tests AWS provider with `default_tags` configuration
- Validates that tags from `tags_all` field are properly detected
- Should pass validation when default_tags provide required tags

### 2. AWS without Default Tags (`test_no_default_tags/`)
- Tests AWS provider without `default_tags` configuration  
- Validates that regular `tags` field continues to work
- Should pass validation when resource-level tags provide required tags

### 3. Google with Default Labels (`test_google_default_labels/`)
- Tests Google provider with `default_labels` configuration

### 4. Google Beta with Default Labels (`test_google_beta_default_labels/`)
- Tests Google Beta provider with `default_labels` configuration
- Validates that labels from `effective_labels` field are properly detected
- Should pass validation when default_labels provide required labels

## Running Tests

```bash
# Test AWS with default_tags
cd test_default_tags
terraform init
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config="spec-tags.yaml" -plan="plan.json"

# Test AWS without default_tags
cd ../test_no_default_tags
terraform init
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config="spec-tags.yaml" -plan="plan.json"

# Test Google with default_labels (using mock plan.json)
cd ../test_google_default_labels
terratags -config="spec-tags.yaml" -plan="plan.json"

# Test Google Beta with default_labels (using mock plan.json)
cd ../test_google_beta_default_labels
terratags -config="spec-tags.yaml" -plan="plan.json"
```

All tests should output: `All resources have the required tags!`

## Issue Reference

This addresses GitHub issue #80: AWS default_tags are not detected when running terratags on the json plan.

## Provider Support

| Provider | Default Config | Resource Field | Merged Field | Status |
|----------|----------------|----------------|--------------|---------|
| AWS | `default_tags` | `tags` | `tags_all` | ✅ Supported |
| Google | `default_labels` | `labels` | `effective_labels` | ✅ Supported |
| Google Beta | `default_labels` | `labels` | `effective_labels` | ✅ Supported |
| Azure azapi | `default_tags` | `tags` | ? | ❓ Unknown |
| Azure azurerm | N/A | `tags` | N/A | ✅ Working |
| AWSCC | N/A | `tags` (list) | N/A | ✅ Working |
