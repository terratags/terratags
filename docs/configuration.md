# Configuration

Terratags requires a configuration file that specifies which tags must be present on your AWS resources. This file can be in either YAML or JSON format.

## Required Tags Configuration

### YAML Format

```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

### JSON Format

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

## Command Options

Terratags supports the following command-line options:

- `-config`, `-c`: Path to the config file (JSON/YAML) containing required tag keys (required)
- `-dir`, `-d`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`, `-v`: Enable verbose output
- `-plan`, `-p`: Path to Terraform plan JSON file to analyze
- `-report`, `-r`: Path to output HTML report file
- `-remediate`, `-m`: Show auto-remediation suggestions for non-compliant resources
- `-exemptions`, `-e`: Path to exemptions file (JSON/YAML)
- `-version`, `-V`: Show version information

## Configuration Best Practices

1. **Start Simple**: Begin with a small set of required tags and gradually expand
2. **Be Consistent**: Use consistent naming conventions for your tags
3. **Document Purpose**: Include comments in your configuration files explaining the purpose of each tag
4. **Version Control**: Keep your configuration files in version control
5. **Team Alignment**: Ensure your team understands the tagging requirements