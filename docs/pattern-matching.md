# Pattern Matching Guide

This guide provides comprehensive information about using pattern matching in Terratags to validate tag values with regular expressions.

## Overview

Pattern matching allows you to enforce specific formats, naming conventions, and business rules for your tag values using regular expressions. This goes beyond simple presence validation to ensure tag values meet your organization's standards.

## Basic Concepts

### What is Pattern Matching?

Pattern matching uses regular expressions (regex) to validate that tag values conform to specific formats. For example:

- Ensure environment tags only contain approved values (`dev`, `test`, `prod`)
- Validate email addresses for ownership tags
- Enforce project code formats (`ABC-123`)
- Prevent whitespace in resource names

### When to Use Pattern Matching

Use pattern matching when you need to:

- **Standardize Values**: Ensure consistent naming across resources
- **Enforce Business Rules**: Implement organizational policies
- **Prevent Errors**: Catch common mistakes early
- **Improve Compliance**: Meet regulatory or audit requirements
- **Maintain Quality**: Ensure clean, consistent infrastructure

## Configuration Formats

### Simple Format (No Patterns)

Basic presence validation without value checking:

```yaml
required_tags:
  - Name
  - Environment
  - Owner
```

### Pattern Format

Advanced validation with regex patterns:

```yaml
required_tags:
  Name:
    pattern: "^[a-zA-Z0-9-]+$"
  
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

### Mixed Format

Combine both approaches:

```yaml
required_tags:
  # Pattern validation
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  # Simple validation
  Team: {}
  Description: {}
```

## Common Patterns Library

### Environment Tags

Restrict to approved environment names:

```yaml
Environment:
  pattern: "^(dev|test|staging|prod)$"
```

**Examples:**
- ✅ Matches: `dev`, `test`, `staging`, `prod`
- ❌ Rejects: `development`, `production`, `DEV`, `Test`

### Email Addresses

Validate email format for ownership:

```yaml
Owner:
  pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

**Examples:**
- ✅ Matches: `devops@company.com`, `team.lead@company.com`
- ❌ Rejects: `username`, `user@domain`, `@company.com`

### Project Codes

Enforce structured project identifiers:

```yaml
Project:
  pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"
```

**Examples:**
- ✅ Matches: `WEB-123456`, `DATA-567890`, `INFRA-890123`
- ❌ Rejects: `web-123`, `PROJECT`, `ABC-12`

### Cost Centers

Standardize cost center format:

```yaml
CostCenter:
  pattern: "^CC-[0-9]{4}$"
```

**Examples:**
- ✅ Matches: `CC-1234`, `CC-5678`, `CC-9012`
- ❌ Rejects: `CC123`, `CC-12345`, `cc-1234`

### Resource Names

Prevent whitespace and special characters:

```yaml
Name:
  pattern: "^[a-zA-Z0-9][a-zA-Z0-9-_]*[a-zA-Z0-9]$"
```

**Examples:**
- ✅ Matches: `web-server-01`, `data-bucket`, `main-vpc`
- ❌ Rejects: `web server`, `-web-server`, `api-gateway-`

### Version Numbers

Validate semantic versioning:

```yaml
Version:
  pattern: "^v?[0-9]+\\.[0-9]+\\.[0-9]+$"
```

**Examples:**
- ✅ Matches: `1.0.0`, `v2.1.3`, `10.15.2`
- ❌ Rejects: `1.0`, `v1`, `1.0.0-beta`

### IP Addresses

Validate IPv4 addresses:

```yaml
IPAddress:
  pattern: "^([0-9]{1,3}\\.){3}[0-9]{1,3}$"
```

**Examples:**
- ✅ `192.168.1.1`, `10.0.0.1`, `172.16.0.1`
- ❌ `192.168.1`, `256.1.1.1`, `not-an-ip`

### AWS Resource ARNs

Validate ARN format:

```yaml
SourceARN:
  pattern: "^arn:aws:[a-zA-Z0-9-]+:[a-zA-Z0-9-]*:[0-9]{12}:.+$"
```

**Examples:**
- ✅ `arn:aws:s3:::my-bucket`, `arn:aws:iam::123456789012:role/MyRole`
- ❌ `arn:aws:s3`, `not-an-arn`

## Advanced Pattern Techniques

### Case-Insensitive Matching

Use the `--ignore-case` flag for case-insensitive tag name matching (patterns themselves remain case-sensitive):

```bash
terratags -config config.yaml -dir ./terraform --ignore-case
```

### Optional Components

Use `?` for optional parts:

```yaml
Version:
  pattern: "^v?[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9]+)?$"
```

Matches: `1.0.0`, `v1.0.0`, `1.0.0-beta`, `v2.1.3-alpha`

### Character Classes

Use character classes for flexibility:

```yaml
Name:
  pattern: "^[a-zA-Z][a-zA-Z0-9_-]{2,30}[a-zA-Z0-9]$"
```

- Must start with a letter
- Can contain letters, numbers, underscores, hyphens
- Must end with letter or number
- Length between 4-32 characters

### Alternation

Use `|` for multiple valid formats:

```yaml
Environment:
  pattern: "^(dev|development|test|testing|stage|staging|prod|production)$"
```

### Quantifiers

Control repetition with quantifiers:

- `*` - Zero or more
- `+` - One or more  
- `?` - Zero or one
- `{n}` - Exactly n times
- `{n,}` - n or more times
- `{n,m}` - Between n and m times

```yaml
Project:
  pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"  # 2-4 letters, 3-6 digits
```

## Error Handling and Debugging

### Understanding Error Messages

When validation fails, Terratags provides detailed error messages:

```
Resource aws_instance 'web_server' has tag pattern violations:
  - Tag 'Environment': value 'Production' does not match required pattern '^(dev|test|staging|prod)$'
  - Tag 'Owner': value 'DevOps Team' does not match required pattern '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'
```

### Common Pattern Errors

1. **Escaping Issues**: Remember to escape backslashes in YAML (`\\` instead of `\`)
2. **Anchoring**: Use `^` and `$` to match the entire string
3. **Case Sensitivity**: Patterns are case-sensitive by default
4. **Special Characters**: Escape regex special characters when matching literally

### Testing Patterns

Test your patterns before deployment:

```bash
# Test with passing examples
terratags -config config.yaml -dir examples/pattern_validation_passing

# Test with failing examples
terratags -config config.yaml -dir examples/pattern_validation_failing

# Generate detailed report
terratags -config config.yaml -dir examples/failing -report debug-report.html
```

## Best Practices

### 1. Start Simple

Begin with basic patterns and add complexity gradually:

```yaml
# Start with this
Environment:
  pattern: "^(dev|prod)$"

# Expand as needed
Environment:
  pattern: "^(dev|development|test|testing|staging|prod|production)$"
```

### 2. Use Meaningful Patterns

Make patterns reflect real business requirements:

```yaml
# Good: Reflects actual project structure
Project:
  pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"

# Bad: Too restrictive without business justification
Project:
  pattern: "^PROJECT-[0-9]{3}$"
```

### 3. Document Your Patterns

Add comments explaining pattern requirements:

```yaml
required_tags:
  # Project code format: 2-4 uppercase letters, dash, 3-6 digits
  # Examples: WEB-123456, DATA-567890, INFRA-890123
  Project:
    pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"
```

### 4. Test Thoroughly

Create test cases for both valid and invalid values:

```yaml
# Test cases in comments
# Valid: user@company.com, team.lead@example.org
# Invalid: username, user@domain, @company.com
Owner:
  pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

### 5. Consider Migration

Plan migration from simple to pattern validation:

```yaml
# Phase 1: Simple validation
required_tags:
  - Environment
  - Owner

# Phase 2: Add patterns gradually
required_tags:
  Environment:
    pattern: "^(dev|test|prod)$"
  Owner: {}  # Still simple validation

# Phase 3: Full pattern validation
required_tags:
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

## Performance Considerations

### Pattern Complexity Impact

Pattern validation performance depends on regex complexity:

- **Simple patterns** (literal strings, basic character classes): Fastest
- **Medium patterns** (alternation, quantifiers): Good performance  
- **Complex patterns** (nested groups, lookaheads): May impact performance with large files

### Best Practices for Performance

1. **Use Specific Patterns**: Prefer `^(dev|test|prod)$` over `.*dev.*`
2. **Avoid Excessive Backtracking**: Be careful with patterns like `(a+)+b`
3. **Test at Scale**: Validate performance with large Terraform configurations
4. **Profile When Needed**: Use verbose mode to identify slow validations

### Optimization Tips

```yaml
# Good: Specific and efficient
Environment:
  pattern: "^(dev|test|staging|prod)$"

# Avoid: Overly broad and potentially slow
Environment:
  pattern: ".*"
```

## Integration Examples

### CI/CD Pipeline

```yaml
# GitHub Actions example
- name: Validate Tags
  run: |
    terratags -config .terratags.yaml -dir ./terraform
    if [ $? -ne 0 ]; then
      echo "Tag validation failed. Please fix tag patterns."
      exit 1
    fi
```

### Pre-commit Hook

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/terratags/terratags
    rev: v0.3.0
    hooks:
      - id: terratags
        args: [--config=.terratags.yaml, --remediate]
```

### Terraform Plan Integration

```bash
# Validate against Terraform plan
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
terratags -config config.yaml -plan plan.json
```

### GitLab CI Integration

```yaml
# .gitlab-ci.yml
validate-tags:
  stage: validate
  script:
    - terratags -config .terratags.yaml -dir ./terraform
  rules:
    - changes:
        - "**/*.tf"
```

### Jenkins Pipeline Integration

```groovy
pipeline {
    agent any
    stages {
        stage('Validate Tags') {
            steps {
                sh 'terratags -config config.yaml -dir ./terraform'
            }
        }
    }
}
```

## Troubleshooting

### Pattern Not Matching

1. **Check Escaping**: Ensure backslashes are properly escaped in YAML
2. **Verify Anchors**: Use `^` and `$` to match entire string
3. **Test Online**: Use regex testing tools to validate patterns
4. **Check Case**: Patterns are case-sensitive by default

### Performance Considerations

1. **Simple Patterns**: Use simple patterns when possible
2. **Avoid Backtracking**: Be careful with complex nested patterns
3. **Test at Scale**: Validate performance with large Terraform files

### Common Regex Gotchas

1. **Greedy Matching**: `.*` matches as much as possible
2. **Escaping Dots**: Use `\\.` to match literal dots
3. **Word Boundaries**: Use `\\b` for word boundaries in YAML
4. **Unicode**: Go regex supports Unicode by default

## Reference

### Regex Quick Reference

| Pattern | Description | Example |
|---------|-------------|---------|
| `^` | Start of string | `^dev` |
| `$` | End of string | `prod$` |
| `.` | Any character | `a.c` matches `abc` |
| `*` | Zero or more | `ab*` matches `a`, `ab`, `abb` |
| `+` | One or more | `ab+` matches `ab`, `abb` |
| `?` | Zero or one | `ab?` matches `a`, `ab` |
| `{n}` | Exactly n | `a{3}` matches `aaa` |
| `{n,m}` | Between n and m | `a{2,4}` matches `aa`, `aaa`, `aaaa` |
| `[abc]` | Character class | `[abc]` matches `a`, `b`, or `c` |
| `[a-z]` | Character range | `[a-z]` matches any lowercase letter |
| `[^abc]` | Negated class | `[^abc]` matches anything except `a`, `b`, `c` |
| `\|` | Alternation | `cat\|dog` matches `cat` or `dog` |
| `()` | Grouping | `(ab)+` matches `ab`, `abab` |
| `\\` | Escape character | `\\.` matches literal dot |

### Go Regex Documentation

For complete regex syntax, see the [Go regexp documentation](https://golang.org/pkg/regexp/syntax/).

### Testing Tools

- [Regex101](https://regex101.com/) - Online regex tester
- [RegExr](https://regexr.com/) - Interactive regex learning tool
- [RegexPal](https://www.regexpal.com/) - Simple regex tester
