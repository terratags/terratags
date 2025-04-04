provider "aws" {
  region = "us-west-2"
}

# Example with all required tags
resource "aws_instance" "web_server" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"

  tags = {
    Name        = "WebServer"
    Environment = "Production"
    Owner       = "DevOps"
    Project     = "Website"
    CostCenter  = "CC123"
  }
}

# Example with missing tags
resource "aws_s3_bucket" "data_bucket" {
  bucket = "example-data-bucket"

  tags = {
    Name        = "DataBucket"
    # Missing Environment tag
    # Missing Owner tag
    # Missing Project tag
    CostCenter  = "CC456"
  }
}

# Example with all required tags
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name        = "MainVPC"
    Environment = "Production"
    Owner       = "Network"
    Project     = "Infrastructure"
    CostCenter  = "CC789"
  }
}

# Example with no tags
resource "aws_security_group" "allow_http" {
  name        = "allow_http"
  description = "Allow HTTP inbound traffic"
  vpc_id      = "vpc-12345"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # No tags defined
}
