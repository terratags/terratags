# Installation

Terratags can be installed using Go's package manager:

```bash
go install github.com/terratags/terratags@latest
```

This will download and install the latest version of Terratags to your Go bin directory.

## Prerequisites

- Go 1.18 or later
- Terraform (for analyzing Terraform configurations)

## Verifying Installation

After installation, you can verify that Terratags is installed correctly by running:

```bash
terratags -version
```

This should display the current version of Terratags.

## Building from Source

If you prefer to build from source:

```bash
git clone https://github.com/terratags/terratags.git
cd terratags
go build
```

This will create a `terratags` binary in your current directory.

## Next Steps

After installation, you'll need to:

1. Create a [configuration file](configuration.md) that defines your required tags
2. Run Terratags against your Terraform code
3. Review the results and fix any non-compliant resources