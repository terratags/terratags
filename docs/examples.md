# Examples

This page provides practical examples of how to use Terratags in various scenarios.

## Configuration Examples

### Basic Required Tags Configuration (YAML)

```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

### Required Tags with Descriptions (YAML)

```yaml
required_tags:
  - key: Name
    description: "Identifies the resource"
  - key: Environment
    description: "Deployment environment (dev, test, prod)"
  - key: Owner
    description: "Team or individual responsible for the resource"
  - key: Project
    description: "Project or application name"
```

### Exemptions Configuration

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

## Terraform Examples

### AWS Provider with Default Tags

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
```

### Resource with Tags

```terraform
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  tags = {
    Name = "example-instance"
    Environment = "production"
    Owner = "team-b"
    Project = "website"
  }
}
```

### Resource with Default Tags

```terraform
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  # Only need to specify Name tag, as other required tags come from default_tags
  tags = {
    Name = "example-instance"
  }
}
```

### Module with Tags

```terraform
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "3.14.0"
  
  name = "my-vpc"
  cidr = "10.0.0.0/16"
  
  tags = {
    Name = "my-vpc"
    Environment = "production"
    Owner = "team-b"
    Project = "website"
  }
}
```

## Command Examples

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

### Verbose Output

```bash
terratags -config config.yaml -dir ./infra -verbose
```

## Sample HTML Reports

### Module Blocks Report

<iframe src="../assets/reports/module_blocks.html" width="100%" height="600px" style="border: 1px solid #ddd; border-radius: 5px;"></iframe>

### Resource Blocks Report

<iframe src="../assets/reports/resource_blocks.html" width="100%" height="600px" style="border: 1px solid #ddd; border-radius: 5px;"></iframe>

### Provider Default Tags Report

<iframe src="../assets/reports/provider_default_tags.html" width="100%" height="600px" style="border: 1px solid #ddd; border-radius: 5px;"></iframe>

### AWSCC Resources Report

<iframe src="../assets/reports/awscc_report.html" width="100%" height="600px" style="border: 1px solid #ddd; border-radius: 5px;"></iframe>

This report shows how Terratags handles AWSCC resources, including the new "Excluded" category for resources with non-compliant tag schemas.

## Real-World Scenarios

### Scenario 1: Multi-Environment Deployment

For a project with multiple environments, you might have different tag requirements for each environment:

```yaml
# dev-config.yaml
required_tags:
  - Name
  - Environment
  - Owner
```

```yaml
# prod-config.yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
  - CostCenter
  - DataClassification
```

You can then validate each environment with the appropriate configuration:

```bash
terratags -config dev-config.yaml -dir ./infra/environments/dev
terratags -config prod-config.yaml -dir ./infra/environments/prod
```

### Scenario 2: Gradual Tag Implementation

When implementing tagging policies gradually, you might start with a subset of required tags and add more over time:

```yaml
# phase1-config.yaml
required_tags:
  - Name
  - Environment
```

```yaml
# phase2-config.yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

You can use exemptions to gradually roll out the new requirements:

```yaml
# phase2-exemptions.yaml
exemptions:
  - resource_type: "*"
    resource_name: "*"
    exempt_tags: [Project]
    reason: "Project tag requirement being phased in"
```

```bash
terratags -config phase2-config.yaml -dir ./infra -exemptions phase2-exemptions.yaml
```

As teams update their resources, you can remove exemptions until all resources comply with the full tagging policy.