provider "aws" {
  region = "us-west-2"
}

# Example module with proper tags
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "3.14.0"

  name = "my-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  tags = {
    Name        = "MainVPC"
    Environment = "Production"
    Owner       = "Network"
    Project     = "Infrastructure"
    CostCenter  = "CC123"
  }
}

# Example module with missing tags
module "ec2_instance" {
  source  = "terraform-aws-modules/ec2-instance/aws"
  version = "3.5.0"

  name = "web-server"

  ami                    = "ami-0c55b159cbfafe1f0"
  instance_type          = "t2.micro"
  key_name               = "user1"
  monitoring             = true
  vpc_security_group_ids = ["sg-12345678"]
  subnet_id              = module.vpc.public_subnets[0]

  tags = {
    Name        = "WebServer"
    Environment = "Production"
    # Missing Owner tag
    # Missing Project tag
    CostCenter  = "CC456"
  }
}

# Example module with all required tags
module "s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"
  version = "3.3.0"

  bucket = "my-s3-bucket"
  acl    = "private"

  versioning = {
    enabled = true
  }

  tags = {
    Name        = "S3Bucket"
    Environment = "Production"
    Owner       = "DataTeam"
    Project     = "DataLake"
    CostCenter  = "CC789"
  }
}
