# Google Cloud Provider Support

Terratags supports the Google Cloud provider for validating labels on GCP resources.

## Overview

Google Cloud Platform uses **labels** instead of **tags** for resource metadata. Terratags treats labels the same way as tags for validation purposes, ensuring consistent tag/label compliance across AWS, Azure, and Google Cloud.

## Supported Features

- ✅ Label validation on 244+ Google Cloud resources
- ✅ Provider-level `default_labels` support
- ✅ Pattern matching for label values
- ✅ HTML report generation
- ✅ Terraform plan validation
- ✅ Module resource validation

## Label Format

Google Cloud resources use a map of key/value pairs for labels:

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
    environment = "production"
    project     = "terratags"
    name        = "example-instance"
  }
}
```

## Provider Default Labels

The Google provider supports `default_labels` at the provider level, similar to AWS `default_tags`:

```terraform
provider "google" {
  project = "my-project-id"
  region  = "us-central1"

  default_labels = {
    environment = "production"
    owner       = "team-a"
  }
}

resource "google_storage_bucket" "example" {
  name     = "example-bucket"
  location = "US"

  labels = {
    name    = "example-bucket"
    project = "demo"
  }
}
```

In this example, the bucket will have all four labels:
- `name` and `project` from resource-level labels
- `environment` and `owner` from provider's `default_labels`

## Label Inheritance

Terratags tracks label sources and inheritance:

1. **Provider default_labels**: Applied to all resources created by the provider
2. **Resource labels**: Specified directly on the resource
3. **Module labels**: Inherited from module blocks

Resources only need to specify labels not covered by `default_labels`.

## Validation Example

### Configuration File

```yaml
required_tags:
  name: {}
  environment:
    pattern: "^(dev|test|staging|prod)$"
  owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  project: {}
```

### Terraform File

```terraform
provider "google" {
  project = "my-project-id"
  region  = "us-central1"

  default_labels = {
    environment = "prod"
    owner       = "devops@company.com"
  }
}

# Compliant - has all required labels
resource "google_compute_disk" "compliant" {
  name = "test-disk"
  type = "pd-ssd"
  zone = "us-central1-a"
  size = 10

  labels = {
    name    = "test-disk"
    project = "demo"
  }
}

# Non-compliant - missing project label
resource "google_storage_bucket" "non_compliant" {
  name     = "my-bucket"
  location = "US"

  labels = {
    name = "my-bucket"
  }
}
```

### Validation Output

```bash
$ terratags -config config.yaml -dir ./gcp-infra

Tag validation issues found:
Resource google_storage_bucket 'non_compliant' is missing required tags: project

Summary: 1/2 resources compliant (50.0%)
```

## Supported Resources

Terratags supports 244+ Google Cloud resources that have labels support, including:

- Compute Engine (instances, disks, images)
- Cloud Storage (buckets)
- BigQuery (datasets, tables)
- Cloud SQL (instances)
- GKE (clusters, node pools)
- Cloud Functions
- Cloud Run
- And many more...

For the complete list, see the [Supported Providers](providers.md) documentation.

## Key Differences from AWS/Azure

| Feature | AWS | Azure | Google Cloud |
|---------|-----|-------|--------------|
| Terminology | tags | tags | labels |
| Provider defaults | default_tags | default_tags (azapi only) | default_labels |
| Format | Map | Map | Map |
| Validation | ✅ | ✅ | ✅ |

## Usage Examples

### Basic Validation

```bash
terratags -config config.yaml -dir ./gcp-infra
```

### Generate Report

```bash
terratags -config config.yaml -dir ./gcp-infra -report gcp-report.html
```

### Validate Terraform Plan

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

### With Exemptions

```bash
terratags -config config.yaml -dir ./gcp-infra -exemptions exemptions.yaml
```

## Best Practices

1. **Use default_labels**: Define common labels at the provider level
2. **Pattern validation**: Use regex patterns to enforce label value formats
3. **Consistent naming**: Use the same label keys across AWS, Azure, and GCP
4. **Documentation**: Document your labeling strategy
5. **Automation**: Integrate terratags into CI/CD pipelines

## Limitations

- Provider aliases are not tested and behavior cannot be guaranteed
- Labels must follow [GCP label requirements](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements)

## See Also

- [Configuration](configuration.md) - Configure required labels
- [Pattern Matching](pattern-matching.md) - Validate label values with regex
- [Default Tags](default-tags.md) - Learn about default_labels inheritance
- [Examples](examples.md) - More usage examples
