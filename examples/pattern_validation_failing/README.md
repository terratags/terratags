# Pattern Validation - Failing Example

This example demonstrates Terraform resources with tags that **fail** pattern validation requirements.

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

## Resources with Violations

All resources in this example have non-compliant tag values:

### AWS Instance
- **Name**: `web server 01` ❌ (contains whitespace)
- **Environment**: `Production` ❌ (not in allowed values)
- **Owner**: `DevOps Team` ❌ (not a valid email)
- **Project**: `website` ❌ (doesn't match project format)
- **CostCenter**: `CC123` ❌ (missing dash)

### S3 Bucket
- **Name**: `data bucket` ❌ (contains whitespace)
- **Environment**: `development` ❌ (not in allowed values)
- **Owner**: `data-team` ❌ (not a valid email)
- **Project**: `data-project-1` ❌ (doesn't match project format)
- **CostCenter**: `CostCenter-567` ❌ (wrong format)

### VPC
- **Name**: `main vpc network` ❌ (contains whitespace)
- **Environment**: `PROD` ❌ (case sensitive, not in allowed values)
- **Owner**: `network.admin` ❌ (not a valid email)
- **Project**: `infrastructure` ❌ (doesn't match project format)
- **CostCenter**: `CC-12345` ❌ (too many digits)

### Security Group
- **Name**: `allow http security group` ❌ (contains whitespace)
- **Environment**: `Testing` ❌ (not in allowed values)
- **Owner**: `security` ❌ (not a valid email)
- **Project**: `SEC` ❌ (missing numbers)
- **CostCenter**: `CC-12` ❌ (too few digits)

## Running the Validation

```bash
# From the terratags root directory
terratags -config examples/config-patterns.yaml -dir examples/pattern_validation_failing
```

**Expected Output:**
```
Tag validation issues found:
Resource aws_instance 'web_server' has tag pattern violations:
  - Tag 'Name': value 'web server 01' does not match required pattern '^\\S+$'
  - Tag 'Environment': value 'Production' does not match required pattern '^(dev|test|staging|prod)$'
  - Tag 'Owner': value 'DevOps Team' does not match required pattern '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'
  - Tag 'Project': value 'website' does not match required pattern '^[A-Z]{2,4}-[0-9]{3,6}$'
  - Tag 'CostCenter': value 'CC123' does not match required pattern '^CC-[0-9]{4}$'
...
Summary: 0/4 resources compliant (0.0%)
```

## Common Pattern Violations

This example demonstrates common mistakes:

### 1. Whitespace in Names
❌ `"web server 01"`  
✅ `"web-server-01"`

### 2. Non-Standard Environment Values
❌ `"Production"`, `"development"`, `"Testing"`  
✅ `"prod"`, `"dev"`, `"test"`

### 3. Invalid Email Formats
❌ `"DevOps Team"`, `"data-team"`, `"security"`  
✅ `"devops@company.com"`, `"data@company.com"`

### 4. Incorrect Project Codes
❌ `"website"`, `"data-project-1"`, `"SEC"`  
✅ `"WEB-123456"`, `"DATA-567890"`, `"SEC-123456"`

### 5. Wrong Cost Center Format
❌ `"CC123"`, `"CostCenter-567"`, `"CC-12345"`  
✅ `"CC-1234"`, `"CC-5678"`, `"CC-9012"`

## Remediation Suggestions

To fix these violations:

1. **Remove whitespace** from tag values
2. **Use standard environment names**: dev, test, staging, prod
3. **Use valid email addresses** for ownership
4. **Follow project code format**: 2-4 uppercase letters, dash, 3-6 digits
5. **Use correct cost center format**: CC-XXXX where X is a digit
