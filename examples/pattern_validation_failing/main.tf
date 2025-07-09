# Example Terraform configuration with FAILING pattern validation
# Tag values that violate the required patterns

resource "aws_instance" "web_server" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"

  tags = {
    Name        = "web server 01"
    Environment = "Production"
    Owner       = "DevOps Team"
    Project     = "website"
    CostCenter  = "CC123"
  }
}

resource "aws_s3_bucket" "data_bucket" {
  bucket = "company-data-bucket-dev"

  tags = {
    Name        = "data bucket"
    Environment = "development"
    Owner       = "data-team"
    Project     = "data-project-1"
    CostCenter  = "CostCenter-567"
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name        = "main vpc network"
    Environment = "PROD"
    Owner       = "network.admin"
    Project     = "infrastructure"
    CostCenter  = "CC-12345"
  }
}

resource "aws_security_group" "allow_http" {
  name_prefix = "allow-http-"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "allow http security group"
    Environment = "Testing"
    Owner       = "security"
    Project     = "SEC"
    CostCenter  = "CC-12"
  }
}
