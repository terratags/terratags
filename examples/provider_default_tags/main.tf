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

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  
  # Only need to specify Name tag, as other required tags come from default_tags
  tags = {
    Name = "example-instance"
  }
}

resource "aws_s3_bucket" "example" {
  bucket = "my-example-bucket"
  
  # No tags needed here, as all required tags come from default_tags
}
