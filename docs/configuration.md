# Configuration

Terratags requires a configuration file that specifies which tags must be present on your AWS resources. This file can be in either YAML or JSON format.

## Required Tags Configuration

### Simple Format (YAML)

```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

### Simple Format (JSON)

```json
{
  "required_tags": [
    "Name",
    "Environment",
    "Owner",
    "Project"
  ]
}
```

### Pattern Validation Format (YAML)

For advanced tag value validation using regular expressions:

```yaml
required_tags:
  Name:
    pattern: "^\\S+$"  # No whitespace
  
  Environment:
    pattern: "^(dev|test|staging|prod)$"  # Specific values only
  
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"  # Email format
  
  Project:
    pattern: "^[A-Z]{2,4}-[0-9]{3,6}$"  # Project code format
  
  CostCenter:
    pattern: "^CC-[0-9]{4}$"  # Cost center format
```

### Mixed Format (YAML)

Combine simple and pattern validation:

```yaml
required_tags:
  # Pattern validation for critical tags
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  
  # Simple validation for others
  Name: {}
  Project: {}
```

For comprehensive pattern matching documentation, see the [Pattern Matching Guide](pattern-matching.md).

## Command Options

Terratags supports the following command-line options:

- `-config`, `-c`: Path to the config file (JSON/YAML) containing required tag keys (required)
- `-dir`, `-d`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`, `-v`: Enable verbose output
- `-log-level`, `-l`: Set logging level: DEBUG, INFO, WARN, ERROR (default: ERROR)
- `-plan`, `-p`: Path to Terraform plan JSON file to analyze
- `-report`, `-r`: Path to output HTML report file
- `-remediate`, `-re`: Show auto-remediation suggestions for non-compliant resources
- `-exemptions`, `-e`: Path to exemptions file (JSON/YAML)
- `-ignore-case`, `-i`: Ignore case when comparing required tag keys
- `-help`, `-h`: Show help message
- `-version`, `-V`: Show version information

## Configuration Best Practices

1. **Start Simple**: Begin with a small set of required tags and gradually expand
2. **Be Consistent**: Use consistent naming conventions for your tags
3. **Document Purpose**: Include comments in your configuration files explaining the purpose of each tag
4. **Version Control**: Keep your configuration files in version control
5. **Team Alignment**: Ensure your team understands the tagging requirements