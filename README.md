# Terratags

Terratags is a Go-based utility for validating tags in Terraform configurations. It helps ensure that all taggable resources have the required tags.

## Features

- Validates tags on Terraform resources and modules
- Supports both direct resource blocks and module blocks
- Configurable required tags via JSON or YAML configuration
- Works with AWS, Azure, and GCP resources
- Can be integrated into CI/CD pipelines

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
- `examples/config.json` and `examples/config.yaml`: Example configuration files

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
