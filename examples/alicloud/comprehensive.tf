resource "alicloud_instance" "web" {
  availability_zone = "cn-beijing-a"
  security_groups   = ["sg-12345"]
  instance_type     = "ecs.n4.large"
  image_id          = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_name     = "web-server"
  vswitch_id        = "vsw-12345"
  
  tags = {
    Name        = "web-server"
    Environment = "prod"
    Project     = "ecommerce"
    Owner       = "devops@company.com"
  }
  
  volume_tags = {
    VolumeType = "SystemDisk"
    Backup     = "Required"
  }
}

resource "alicloud_oss_bucket" "assets" {
  bucket = "company-assets-bucket"
  
  tags = {
    Name        = "assets-bucket"
    Environment = "prod"
    Project     = "ecommerce"
    Owner       = "devops@company.com"
  }
}

resource "alicloud_vpc" "main" {
  vpc_name   = "main-vpc"
  cidr_block = "172.16.0.0/16"
  
  tags = {
    Name        = "main-vpc"
    Environment = "prod"
    Project     = "ecommerce"
    Owner       = "devops@company.com"
  }
}

resource "alicloud_security_group" "web" {
  security_group_name = "web-sg"
  description         = "Security group for web servers"
  vpc_id              = alicloud_vpc.main.id
  
  tags = {
    Name        = "web-security-group"
    Environment = "prod"
    Project     = "ecommerce"
    Owner       = "devops@company.com"
  }
}
