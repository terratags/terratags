# Exemptions Configuration

Exemptions allow you to exclude specific resources or resource types from certain tag requirements. Create a YAML or JSON file with your exemptions.

## Exemption Fields

- `resource_type`: The AWS resource type (e.g., aws_s3_bucket, aws_instance)
- `resource_name`: The name of the specific resource to exempt. Use "*" to exempt all resources of the specified type
- `exempt_tags`: List of tags that are not required for this resource
- `reason`: A description explaining why this exemption exists

## YAML Example

```yaml
exemptions:
  - resource_type: aws_s3_bucket
    resource_name: logs_bucket
    exempt_tags: [Owner, Project]
    reason: "Legacy bucket used for system logs only"
  
  - resource_type: aws_dynamodb_table
    resource_name: "*"
    exempt_tags: [Environment]
    reason: "DynamoDB tables use environment from provider default_tags"
```

## JSON Example

```json
{
  "exemptions": [
    {
      "resource_type": "aws_s3_bucket",
      "resource_name": "logs_bucket",
      "exempt_tags": ["Owner", "Project"],
      "reason": "Legacy bucket used for system logs only"
    },
    {
      "resource_type": "aws_dynamodb_table",
      "resource_name": "*",
      "exempt_tags": ["Environment"],
      "reason": "DynamoDB tables use environment from provider default_tags"
    }
  ]
}
```

## When to Use Exemptions

Exemptions are useful in several scenarios:

1. **Legacy Resources**: Older resources that cannot be easily updated
2. **Special Purpose Resources**: Resources with a specific purpose that don't fit the standard tagging model
3. **Resource Limitations**: Some AWS resources have limitations on tagging
4. **Default Tag Inheritance**: Resources that inherit tags from other sources

## Best Practices for Exemptions

1. **Document Reasons**: Always include a clear reason for each exemption
2. **Regular Review**: Periodically review exemptions to see if they're still necessary
3. **Minimize Use**: Use exemptions sparingly to maintain consistent tagging
4. **Specific Scope**: Make exemptions as specific as possible (prefer specific resource names over wildcards)
5. **Version Control**: Keep your exemptions file in version control