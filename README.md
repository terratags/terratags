# Terratags

Terratags is a Go-based utility for validating tags in Terraform configurations. It helps review that all taggable resources have the required tags.

## Features

- Validates tags on Terraform resources and modules
- Supports both direct resource blocks and module blocks
- Configurable required tags via JSON or YAML configuration
- Works with AWS and AWSCC resources
- Supports AWS provider default_tags
- Can be integrated into CI/CD pipelines
- Comprehensive list of taggable resources

## Installation

### From Source

```bash
git clone https://github.com/quixoticmonk/terratags.git
cd terratags
go build -o terratags ./cmd/terratags
```

### Using Go Install

```bash
go install github.com/quixoticmonk/terratags/cmd/terratags@latest
```

## Usage

```bash
terratags -config <config_file.json|yaml> [-dir <terraform_directory>] [-verbose]
```

### Options

- `-config`: Path to the configuration file (JSON or YAML) containing required tag keys (required)
- `-dir`: Path to the Terraform directory to analyze (default: current directory)
- `-verbose`: Enable verbose output

### Example

```bash
terratags -config config.yaml -dir ./terraform -verbose
```

## Configuration

Terratags uses a configuration file to define which tags are required. The configuration file can be in JSON or YAML format.

### JSON Example

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

### YAML Example

```yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

## GitHub Actions Integration

You can integrate Terratags into your CI/CD pipeline using GitHub Actions. Here's an example workflow:

```yaml
name: Terraform Tag Validation

on:
  pull_request:
    paths:
      - '**.tf'

jobs:
  validate-tags:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install Terratags
        run: go install github.com/quixoticmonk/terratags/cmd/terratags@latest

      - name: Validate Tags
        run: terratags -config ./config.yaml -dir ./terraform
```

## Examples

The repository includes example Terraform configurations in the `examples` directory:

- `examples/resource_blocks`: Examples of direct resource blocks with various tag configurations
- `examples/module_blocks`: Examples of module blocks with various tag configurations
- `examples/provider_default_tags`: Example using AWS provider default_tags
- `examples/config.json` and `examples/config.yaml`: Example configuration files

## Supported Resources

Terratags includes a comprehensive list of AWS and AWSCC resources that support tagging. The current version supports:

- **730+ AWS provider resources** that accept tags
- **710+ AWSCC provider resources** that accept tags

This list is automatically generated from the provider schemas and is regularly updated to ensure compatibility with the latest AWS resources.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Updating the Resource List

To update the list of taggable resources:

```bash
# Run the update script
go run scripts/update_resources.go
```

This script will:
1. Create a temporary Terraform configuration with AWS and AWSCC providers
2. Initialize the providers
3. Extract the provider schemas
4. Parse the schemas to identify resources with `tags` attributes
5. Generate an updated `aws_taggable_resources.go` file

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
