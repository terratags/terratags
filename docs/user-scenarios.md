# User Scenarios

### Scenario 1: Multi-Environment Deployment

For a project with multiple environments, you might have different tag requirements for each environment:

```yaml
# dev-config.yaml
required_tags:
  - Name
  - Environment
  - Owner
```

```yaml
# prod-config.yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
  - CostCenter
  - DataClassification
```

You can then validate each environment with the appropriate configuration:

```bash
terratags -config dev-config.yaml -dir ./infra/environments/dev
terratags -config prod-config.yaml -dir ./infra/environments/prod
```

### Scenario 2: Gradual Tag Implementation

When implementing tagging policies gradually, you might start with a subset of required tags and add more over time:

```yaml
# phase1-config.yaml
required_tags:
  - Name
  - Environment
```

```yaml
# phase2-config.yaml
required_tags:
  - Name
  - Environment
  - Owner
  - Project
```

You can use exemptions to gradually roll out the new requirements:

```yaml
# phase2-exemptions.yaml
exemptions:
  - resource_type: "*"
    resource_name: "*"
    exempt_tags: [Project]
    reason: "Project tag requirement being phased in"
```

```bash
terratags -config phase2-config.yaml -dir ./infra -exemptions phase2-exemptions.yaml
```

As teams update their resources, you can remove exemptions until all resources comply with the full tagging policy.