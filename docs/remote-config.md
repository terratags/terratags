# Remote Config Files

Terratags supports loading configuration files from remote locations, enabling centralized config management across multiple repositories.

## Supported Sources

### HTTP/HTTPS URLs

Fetch config directly from any web server:

```bash
terratags -config https://example.com/configs/terratags.yaml -dir ./infra
```

**Example with GitHub raw:**
```bash
terratags -config https://raw.githubusercontent.com/org/configs/main/terratags.yaml -dir ./infra
```

### Git Repositories (HTTPS)

Clone a Git repository and extract the config file:

```bash
# With branch
terratags -config https://github.com/org/configs.git//terratags.yaml?ref=main -dir ./infra

# With tag
terratags -config https://github.com/org/configs.git//terratags.yaml?ref=v1.0.0 -dir ./infra

# With subdirectory
terratags -config https://github.com/org/configs.git//terraform/prod/terratags.yaml?ref=main -dir ./infra
```

### Git Repositories (SSH)

Use SSH authentication for private repositories:

```bash
terratags -config git@github.com:org/configs.git//terratags.yaml?ref=main -dir ./infra
```

## URL Format

### Git URLs

Git URLs follow the Terraform module source convention:

```
<git-url>//<file-path>?ref=<branch-or-tag>
```

- `<git-url>`: Repository URL (HTTPS or SSH)
- `//`: Separator between repo and file path
- `<file-path>`: Path to config file within the repository
- `?ref=`: Optional branch, tag, or commit reference

### Supported File Types

Only these extensions are allowed:
- `.yaml`
- `.yml`
- `.json`

## Authentication

### HTTP/HTTPS

Authentication is handled by your system's HTTP client. For private endpoints, configure appropriate credentials.

### Git HTTPS

Uses your git credential helper:

```bash
# Configure credential storage
git config --global credential.helper store

# Or use credential manager
git config --global credential.helper manager
```

### Git SSH

Uses SSH keys from `~/.ssh/`:

```bash
# Add your key to ssh-agent
ssh-add ~/.ssh/id_rsa

# Test SSH connection
ssh -T git@github.com
```

## Use Cases

### Centralized Configuration

Maintain a single source of truth for tag requirements:

```bash
# All teams use the same config
terratags -config https://github.com/company/standards.git//terraform/tags.yaml?ref=main -dir ./infra
```

### Environment-Specific Configs

Use different configs per environment:

```bash
# Production
terratags -config https://github.com/org/configs.git//prod.yaml?ref=main -dir ./infra

# Development
terratags -config https://github.com/org/configs.git//dev.yaml?ref=main -dir ./infra
```

### Version Pinning

Pin to specific config versions:

```bash
# Use tagged version
terratags -config https://github.com/org/configs.git//terratags.yaml?ref=v2.1.0 -dir ./infra
```

## Examples

See the [remote_config examples](https://github.com/terratags/terratags/tree/main/examples/remote_config) directory for working examples and test scripts.

## Troubleshooting

### "unsupported file type" error

Ensure your URL ends with `.yaml`, `.yml`, or `.json`:

```bash
# ✗ Wrong
terratags -config https://example.com/config

# ✓ Correct
terratags -config https://example.com/config.yaml
```

### "repository not found" error

For Git URLs:
- Verify the repository URL is correct
- Check authentication (SSH keys or git credentials)
- Ensure you have access to the repository

### "failed to clone" error

- Check network connectivity
- Verify SSH keys are loaded (`ssh-add -l`)
- Test git access: `git clone <repo-url>`
