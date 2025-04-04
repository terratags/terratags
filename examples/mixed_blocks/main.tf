provider "aws" {
  region = "us-west-2"
}

locals {
  common_tags = {
    Environment = "dev"
    Owner       = "team-a"
    Project     = "demo"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
  
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  
  owners = ["099720109477"] # Canonical
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  tags = {
    Name        = "example-instance"
    Environment = "dev"
    Owner       = "team-a"
    Project     = "demo"
  }
}

resource "aws_s3_bucket" "missing_tags" {
  bucket = "my-example-bucket"
  
  tags = {
    Name        = "example-bucket"
    Environment = "dev"
    Owner       = "team-a"
    Project     = "demo"
  }
}

output "instance_ip" {
  value = "10.0.0.1" # Placeholder
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}
