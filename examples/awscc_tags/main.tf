terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
    }
    awscc = {
      source  = "hashicorp/awscc"
    }
  }
}

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
  # Note: AWSCC provider doesn't support default_tags
}

# AWS resource with tags in map format
resource "aws_s3_bucket" "example" {
  bucket = "example-bucket"
  
  # AWS provider uses map format for tags
  tags = {
    Name        = "Example Bucket"
    Environment = "Test"
    # Owner and Project come from default_tags
  }
}

# AWSCC resource with tags in list-of-maps format
resource "awscc_apigateway_rest_api" "example" {
  name = "example-api"
  
  # AWSCC provider uses list of maps with key/value pairs for tags
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