# Provider Default Tags Support

Terratags integrates with AWS provider's `default_tags` feature. When you define default tags in your AWS provider configuration, Terratags will recognize these tags and consider them when validating resources.

## How Default Tags Work

1. Tags defined in the AWS provider's `default_tags` block are automatically applied to all taggable resources created by that provider
2. Terratags tracks tag inheritance from provider default_tags to individual resources
3. Resources only need to specify tags not covered by default_tags
4. Default tags can be overridden at the resource level if needed

## Example with Default Tags

```terraform
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
```

In this example, the AWS instance will have all four required tags: `Name` from the resource-level tags, and `Environment`, `Owner`, and `Project` from the provider's default_tags.

## Benefits of Using Default Tags

1. **Consistency**: Ensures consistent tagging across all resources
2. **Reduced Duplication**: Eliminates the need to repeat the same tags on every resource
3. **Centralized Management**: Makes it easier to update tags across all resources
4. **Reduced Errors**: Minimizes the chance of missing required tags

## Default Tags Limitations

1. **Provider Specific**: Only works with providers that support default_tags (like AWS)
2. **Override Behavior**: Resource-level tags override default tags with the same key
3. **Module Awareness**: When using modules, be aware of how default tags propagate

## Best Practices

1. **Use for Common Tags**: Use default_tags for tags that should be consistent across all resources
2. **Resource-Specific Tags**: Use resource-level tags for tags that are specific to individual resources
3. **Documentation**: Document which tags are provided by default_tags to avoid confusion
4. **Validation**: Still use Terratags to validate that all required tags are present