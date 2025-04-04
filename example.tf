resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"

  tags = {
    Name        = "example-instance"
    Environment = "dev"
    Owner       = "team-a"
    Project     = "demo"
    CostCenter  = "123456"
  }
}

resource "aws_s3_bucket" "missing_tags" {
  bucket = "my-example-bucket"

  tags = {
    Name = "example-bucket"
    # Missing required tags: Environment, Owner, Project
  }
}
