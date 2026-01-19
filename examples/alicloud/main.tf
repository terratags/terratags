resource "alicloud_instance" "example" {
  availability_zone = "cn-beijing-a"
  security_groups   = ["sg-12345"]
  instance_type     = "ecs.n4.large"
  image_id          = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_name     = "example-instance"
  vswitch_id        = "vsw-12345"
  
  tags = {
    Name        = "example-instance"
    Environment = "prod"
    Project     = "demo"
    Owner       = "devops@company.com"
  }
}

resource "alicloud_oss_bucket" "example" {
  bucket = "example-bucket-12345"
  
  tags = {
    Name        = "example-bucket"
    Environment = "test"
    Project     = "demo"
    Owner       = "devops@company.com"
  }
}
