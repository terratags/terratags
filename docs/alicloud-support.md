# Alibaba Cloud Support

Terratags supports Alibaba Cloud (AliCloud) resources through the `alicloud` provider. AliCloud uses the same tag format as AWS, making it straightforward to validate tags across your AliCloud infrastructure.

## Supported Resources

Terratags supports a comprehensive list of AliCloud resources that have tagging capabilities, including:

### Compute Services
- `alicloud_instance` - ECS instances
- `alicloud_reserved_instance` - Reserved instances
- `alicloud_ecs_instance_set` - Instance sets
- `alicloud_simple_application_server_instance` - Simple application servers

### Storage Services
- `alicloud_oss_bucket` - Object Storage Service buckets
- `alicloud_oss_bucket_object` - OSS objects

### Database Services
- `alicloud_db_instance` - RDS instances
- `alicloud_mongodb_instance` - MongoDB instances
- `alicloud_redis_tair_instance` - Redis instances
- `alicloud_kvstore_instance` - KVStore instances
- And many more database services...

### Networking
- `alicloud_vpc` - Virtual Private Clouds
- `alicloud_vswitch` - Virtual switches
- `alicloud_security_group` - Security groups
- `alicloud_nat_gateway` - NAT gateways
- `alicloud_eip` - Elastic IP addresses
- `alicloud_slb` - Server Load Balancers

### Container Services
- `alicloud_cs_kubernetes_cluster` - Kubernetes clusters
- `alicloud_cs_managed_kubernetes` - Managed Kubernetes
- `alicloud_cs_serverless_kubernetes` - Serverless Kubernetes

### Security Services
- `alicloud_kms_key` - KMS keys
- `alicloud_bastionhost_instance` - Bastion hosts
- `alicloud_waf_instance` - Web Application Firewall
- `alicloud_cloud_firewall_instance` - Cloud Firewall

And many more services across analytics, messaging, CDN, and other categories.

## Tag Format

AliCloud uses the same tag format as AWS - a simple key-value map:

```hcl
resource "alicloud_instance" "example" {
  availability_zone = "cn-beijing-a"
  instance_type     = "ecs.n4.large"
  image_id          = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  
  tags = {
    Name        = "web-server-01"
    Environment = "production"
    Project     = "ecommerce"
    Owner       = "devops@company.com"
  }
}
```

## Tag Constraints

AliCloud has specific constraints for tag keys and values:

- **Key**: Up to 128 characters, cannot begin with "aliyun", "acs:", "http://", or "https://", cannot be null
- **Value**: Up to 128 characters, cannot begin with "aliyun", "acs:", "http://", or "https://", can be null

## Configuration Example

Here's a complete example of using Terratags with AliCloud resources:

**config.yaml:**
```yaml
required_tags:
  Name: {}
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  Project: {}
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

**main.tf:**
```hcl
resource "alicloud_instance" "web" {
  availability_zone = "cn-beijing-a"
  instance_type     = "ecs.n4.large"
  image_id          = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  
  tags = {
    Name        = "web-server"
    Environment = "prod"
    Project     = "website"
    Owner       = "devops@company.com"
  }
}

resource "alicloud_oss_bucket" "assets" {
  bucket = "company-assets-bucket"
  
  tags = {
    Name        = "assets-bucket"
    Environment = "prod"
    Project     = "website"
    Owner       = "devops@company.com"
  }
}
```

**Validation:**
```bash
terratags -config config.yaml -dir .
```

## Pattern Validation

AliCloud resources support the same advanced pattern validation as other providers:

```yaml
required_tags:
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
  Project:
    pattern: "^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$"
```

## Default Tags

Unlike AWS, AliCloud provider does not support `default_tags` at the provider level. All required tags must be specified at the resource level.

## Volume Tags

ECS instances also support `volume_tags` for tagging attached storage devices:

```hcl
resource "alicloud_instance" "example" {
  # ... other configuration ...
  
  tags = {
    Name        = "web-server"
    Environment = "production"
  }
  
  volume_tags = {
    VolumeType = "SystemDisk"
    Backup     = "Required"
  }
}
```

## Exemptions

You can exempt specific AliCloud resources from tag requirements:

```yaml
exemptions:
  - resource_type: alicloud_oss_bucket
    resource_name: logs_bucket
    exempt_tags: [Owner, Project]
    reason: "Legacy bucket used for system logs only"
  
  - resource_type: alicloud_instance
    resource_name: "*"
    exempt_tags: [Environment]
    reason: "Environment determined by VPC placement"
```

## Integration with Terraform Plan

Terratags can validate AliCloud resources in Terraform plans, including resources created by modules:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

This provides comprehensive validation across your entire AliCloud infrastructure, including resources created by external modules.
