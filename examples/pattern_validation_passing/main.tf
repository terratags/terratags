# Example Terraform configuration with PASSING pattern validation
# All tag values match the required patterns

resource "aws_instance" "web_server" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"

  tags = {
    Name        = "web-server-01"
    Environment = "prod"
    Owner       = "devops@company.com"
    Project     = "WEB-123456"
    CostCenter  = "CC-5678"
  }
}

resource "aws_s3_bucket" "data_bucket" {
  bucket = "company-data-bucket-prod"

  tags = {
    Name        = "data-bucket"
    Environment = "staging"
    Owner       = "data.team@company.com"
    Project     = "DATA-567890"
    CostCenter  = "CC-9012"
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name        = "main-vpc"
    Environment = "dev"
    Owner       = "network@company.com"
    Project     = "NET-890123"
    CostCenter  = "CC-3456"
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
    Name        = "allow-http-sg"
    Environment = "test"
    Owner       = "security@company.com"
    Project     = "SEC-123456"
    CostCenter  = "CC-7890"
  }
}
