# Pattern Validation - Passing Example

This example demonstrates Terraform resources with tags that **pass** all pattern validation requirements.

## Configuration Used

This example uses the pattern validation configuration from `../config-patterns.yaml`:

```yaml
required_tags:
  Name:
    pattern: "^\\S+$"                                                    # No whitespace
  
  Environment:
    pattern: "^(dev|test|staging|prod)$"                                 # Specific values only
  
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"        # Valid email format
  
  Project:
    pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"                                  # Project code format
  
  CostCenter:
    pattern: "^CC-[0-9]{4}$"                                            # Cost center format
```

## Resources

All resources in this example have compliant tag values:

### AWS Instance
- **Name**: `web-server-01` ✅ (no whitespace)
- **Environment**: `prod` ✅ (matches allowed values)
- **Owner**: `devops@company.com` ✅ (valid email)
- **Project**: `WEB-123456` ✅ (matches project format)
- **CostCenter**: `CC-5678` ✅ (matches cost center format)

### S3 Bucket
- **Name**: `data-bucket` ✅ (no whitespace)
- **Environment**: `staging` ✅ (matches allowed values)
- **Owner**: `data.team@company.com` ✅ (valid email)
- **Project**: `DATA-567890` ✅ (matches project format)
- **CostCenter**: `CC-9012` ✅ (matches cost center format)

### VPC
- **Name**: `main-vpc` ✅ (no whitespace)
- **Environment**: `dev` ✅ (matches allowed values)
- **Owner**: `network@company.com` ✅ (valid email)
- **Project**: `NET-890123` ✅ (matches project format)
- **CostCenter**: `CC-3456` ✅ (matches cost center format)

### Security Group
- **Name**: `allow-http-sg` ✅ (no whitespace)
- **Environment**: `test` ✅ (matches allowed values)
- **Owner**: `security@company.com` ✅ (valid email)
- **Project**: `SEC-123456` ✅ (matches project format)
- **CostCenter**: `CC-7890` ✅ (matches cost center format)

## Running the Validation

```bash
# From the terratags root directory
terratags -config examples/config-patterns.yaml -dir examples/pattern_validation_passing
```

**Expected Output:**
```
All resources have the required tags!
```

## Key Takeaways

This example shows how to structure tag values to comply with pattern validation:

1. **No Whitespace**: Use hyphens or underscores instead of spaces
2. **Controlled Values**: Use predefined environment values (dev, test, staging, prod)
3. **Email Format**: Use proper email addresses for ownership
4. **Structured Codes**: Follow consistent project and cost center formats
5. **Case Sensitivity**: Patterns are case-sensitive by default
