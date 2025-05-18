# AWSCC Provider Support

Terratags now supports the AWS Cloud Control (AWSCC) provider's tag format, which differs from the standard AWS provider tag format.

## Tag Format Differences

### AWS Provider Tag Format

The AWS provider uses a map of key-value pairs for tags:

```hcl
resource "aws_s3_bucket" "example" {
  bucket = "example-bucket"
  
  tags = {
    Name        = "Example Bucket"
    Environment = "Test"
    Owner       = "DevOps"
    Project     = "Terratags"
  }
}
```

### AWSCC Provider Tag Format

The AWSCC provider uses a list of maps with `key` and `value` fields:

```hcl
resource "awscc_apigateway_rest_api" "example" {
  name = "example-api"
  
  tags = [
    {
      key   = "Name"
      value = "Example API"
    },
    {
      key   = "Environment"
      value = "Test"
    },
    {
      key   = "Owner"
      value = "API Team"
    },
    {
      key   = "Project"
      value = "Terratags"
    }
  ]
}
```

## Default Tags Support

**Important**: The AWSCC provider does not support `default_tags`. Each AWSCC resource must specify all required tags directly in its `tags` attribute.

```hcl
provider "aws" {
  region = "us-west-2"
  
  # AWS provider supports default_tags
  default_tags {
    tags = {
      Owner       = "DevOps"
      Project     = "Terratags"
    }
  }
}

provider "awscc" {
  region = "us-west-2"
  # AWSCC provider doesn't support default_tags
}
```

## Validation

Terratags validates AWSCC resources by:

1. Detecting resources with the `awscc_` prefix
2. Parsing the list-of-maps tag format
3. Validating that all required tags are present
4. Reporting any missing tags

## Example

Here's a complete example showing both AWS and AWSCC resources with their respective tag formats:

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    awscc = {
      source  = "hashicorp/awscc"
      version = "~> 0.67"
    }
  }
}

provider "aws" {
  region = "us-west-2"
  
  default_tags {
    tags = {
      Owner       = "DevOps"
      Project     = "Terratags"
    }
  }
}

provider "awscc" {
  region = "us-west-2"
}

# AWS resource with tags in map format
resource "aws_s3_bucket" "example" {
  bucket = "example-bucket"
  
  tags = {
    Name        = "Example Bucket"
    Environment = "Test"
    # Owner and Project come from default_tags
  }
}

# AWSCC resource with tags in list-of-maps format
resource "awscc_apigateway_rest_api" "example" {
  name = "example-api"
  
  tags = [
    {
      key   = "Name"
      value = "Example API"
    },
    {
      key   = "Environment"
      value = "Test"
    },
    {
      key   = "Owner"
      value = "API Team"
    },
    {
      key   = "Project"
      value = "Terratags"
    }
  ]
}
```

## Best Practices for AWSCC Resources

1. **Specify All Tags**: Since AWSCC doesn't support default_tags, make sure to specify all required tags directly on each resource
2. **Consistent Keys**: Use consistent tag keys across both AWS and AWSCC resources
3. **Case Sensitivity**: Be aware that tag keys are case-sensitive
4. **Validation**: Use Terratags to validate that all required tags are present on AWSCC resources