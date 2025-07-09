# Terratags Architecture

This document provides an overview of the Terratags architecture, design decisions, and internal workings.

## Table of Contents

- [Overview](#overview)
- [Core Components](#core-components)
- [Data Flow](#data-flow)
- [Configuration System](#configuration-system)
- [Parsing Strategy](#parsing-strategy)
- [Validation Engine](#validation-engine)
- [Reporting System](#reporting-system)
- [Design Decisions](#design-decisions)

## Overview

Terratags is designed as a modular CLI tool that validates tag compliance across Terraform configurations. The architecture follows a pipeline pattern where Terraform files are parsed, resources are extracted, tags are validated, and reports are generated.

```
Input (Terraform Files) → Parser → Validator → Reporter → Output (Results/Reports)
                            ↑
                      Configuration
```

## Core Components

### 1. Configuration System (`pkg/config`)

The configuration system handles:
- Loading configuration files (JSON/YAML)
- Supporting both legacy array format and new object format with patterns
- Compiling regex patterns for validation
- Managing exemptions
- Runtime configuration options

**Key Types:**
```go
type Config struct {
    RequiredTags  map[string]TagRequirement
    Exemptions    []ResourceExemption
    IgnoreTagCase bool
}

type TagRequirement struct {
    Pattern         string
    compiledPattern *regexp.Regexp
}
```

### 2. Parser System (`pkg/parser`)

The parser system is responsible for:
- Parsing Terraform files using HCL (HashiCorp Configuration Language)
- Extracting resources and their tags
- Handling different provider tag formats (AWS, Azure, AWSCC)
- Processing provider default_tags
- Supporting both directory scanning and Terraform plan analysis

**Key Types:**
```go
type Resource struct {
    Type       string
    Name       string
    Tags       map[string]string
    Path       string
    TagSources map[string]TagSource
}

type ProviderConfig struct {
    Type        string
    DefaultTags map[string]string
    Path        string
}
```

### 3. Validation Engine (`pkg/validator`)

The validation engine performs:
- Tag presence validation
- Pattern matching validation using compiled regex
- Exemption processing
- Statistics collection
- Violation tracking

**Key Types:**
```go
type TagViolation struct {
    ResourceType      string
    ResourceName      string
    MissingTags       []string
    PatternViolations []PatternViolation
    IsExempt          bool
}

type TagComplianceStats struct {
    TotalResources           int
    CompliantResources       int
    FullyExemptResources     int
    PartiallyExemptResources int
}
```

### 4. Logging System (`pkg/logging`)

The logging system provides:
- Structured logging using Zap
- Multiple log levels (DEBUG, INFO, WARN, ERROR)
- Custom console formatting
- Print function for always-visible output

## Data Flow

### 1. Initialization Phase
```
main.go → Parse CLI flags → Load Configuration → Initialize Logging
```

### 2. Parsing Phase
```
Directory/Plan Input → HCL Parser → Resource Extraction → Tag Extraction
```

### 3. Validation Phase
```
Resources + Config → Tag Validation → Pattern Validation → Exemption Processing
```

### 4. Reporting Phase
```
Violations + Stats → Console Output → HTML Report (optional) → Exit Code
```

## Configuration System

### Configuration Loading

The configuration system supports two formats for backward compatibility:

**Legacy Format (Array):**
```yaml
required_tags:
  - Name
  - Environment
  - Owner
```

**New Format (Object with Patterns):**
```yaml
required_tags:
  Name: {}
  Environment:
    pattern: "^(dev|test|staging|prod)$"
  Owner:
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

### Pattern Compilation

Regex patterns are compiled once during configuration loading and stored in the `TagRequirement` struct for efficient validation.

### Exemption System

Exemptions allow fine-grained control over which resources are exempt from specific tag requirements:

```yaml
exemptions:
  - resource_type: aws_s3_bucket
    resource_name: logs_bucket
    exempt_tags: [Owner, Project]
    reason: "Legacy bucket used for system logs only"
```

## Parsing Strategy

### HCL vs Regex Parsing

Terratags uses a hybrid approach:
- **HCL Parser**: Primary method for parsing Terraform files
- **Regex Fallback**: Used for complex tag extraction patterns

### Provider Support

Different providers have different tag formats:

**AWS Provider:**
```hcl
tags = {
  Environment = "prod"
  Owner       = "team-a"
}
```

**AWSCC Provider:**
```hcl
tags = [{
  key   = "Environment"
  value = "prod"
}, {
  key   = "Owner"
  value = "team-a"
}]
```

**Azure Providers:**
```hcl
tags = {
  Environment = "prod"
  Owner       = "team-a"
}
```

### Default Tags Handling

Provider default_tags are processed and merged with resource-level tags:

1. Extract provider configurations
2. Identify default_tags blocks
3. Merge with resource tags during validation
4. Track tag sources for reporting

## Validation Engine

### Validation Pipeline

1. **Resource Filtering**: Only validate taggable resources
2. **Tag Presence Check**: Verify required tags exist
3. **Pattern Validation**: Apply regex patterns if defined
4. **Exemption Processing**: Apply exemptions
5. **Statistics Collection**: Track compliance metrics

### Pattern Validation

Pattern validation uses compiled regex patterns:

```go
func (c *Config) ValidateTagValue(tagName, tagValue string) (bool, string) {
    req, found := c.RequiredTags[tagName]
    if !found || req.compiledPattern == nil {
        return true, ""
    }
    
    if req.compiledPattern.MatchString(tagValue) {
        return true, ""
    }
    
    return false, fmt.Sprintf("value '%s' does not match pattern '%s'", 
        tagValue, req.Pattern)
}
```

### Case Sensitivity

The validation engine supports case-insensitive tag name matching through the `IgnoreTagCase` configuration option.

## Reporting System

### Console Output

The console output provides:
- Violation details with resource information
- Pattern violation explanations
- Auto-remediation suggestions
- Compliance statistics
- Exemption summaries

### HTML Reports

HTML reports include:
- Visual compliance indicators
- Detailed resource breakdowns
- Tag source tracking
- Exemption details with reasons
- Interactive filtering and sorting

### Auto-Remediation

The system can generate remediation suggestions:
- Missing tag additions
- Provider default_tags recommendations
- Pattern fix suggestions

## Design Decisions

### 1. Modular Architecture

**Decision**: Separate packages for config, parser, validator, and logging.

**Rationale**: 
- Separation of concerns
- Easier testing and maintenance
- Clear interfaces between components

### 2. HCL + Regex Hybrid Parsing

**Decision**: Use HCL parser as primary with regex fallback.

**Rationale**:
- HCL provides structured parsing
- Regex handles edge cases and complex patterns
- Maintains compatibility with various Terraform syntax styles

### 3. Pattern Compilation at Load Time

**Decision**: Compile regex patterns during configuration loading.

**Rationale**:
- Better performance during validation
- Early error detection for invalid patterns
- Single compilation per pattern

### 4. Backward Compatible Configuration

**Decision**: Support both array and object configuration formats.

**Rationale**:
- Smooth migration path for existing users
- Maintains compatibility with v0.2.x configurations
- Allows gradual adoption of pattern features

### 5. Tag Source Tracking

**Decision**: Track where each tag comes from (resource vs provider default).

**Rationale**:
- Better debugging and reporting
- Helps users understand tag inheritance
- Enables more detailed compliance reports

### 6. Exemption System

**Decision**: Flexible exemption system with reasons.

**Rationale**:
- Real-world compliance often requires exceptions
- Audit trail for why exemptions exist
- Fine-grained control over exemptions

### 7. Statistics and Reporting

**Decision**: Comprehensive statistics collection and HTML reporting.

**Rationale**:
- Provides visibility into compliance posture
- Supports compliance reporting requirements
- Helps identify trends and patterns

## Performance Considerations

### Memory Usage

- Resources are processed in batches
- Compiled patterns are reused
- Tag maps are created on-demand

### CPU Usage

- Regex patterns are compiled once
- File parsing is done sequentially
- Validation is optimized for common cases

### Scalability

Current limitations and considerations:
- Single-threaded file processing
- Memory usage grows with number of resources
- Regex complexity can impact performance

## Error Handling

### Error Categories

1. **Configuration Errors**: Invalid config files, bad patterns
2. **Parsing Errors**: Invalid Terraform syntax, file access issues
3. **Validation Errors**: Pattern compilation failures
4. **Runtime Errors**: File system issues, memory constraints

### Error Propagation

Errors are propagated up the call stack with context:
```go
if err := loadConfig(path); err != nil {
    return fmt.Errorf("failed to load config from %s: %w", path, err)
}
```

## Extension Points

### Adding New Providers

1. Create resource type definitions
2. Implement tag extraction logic
3. Update parser to recognize provider
4. Add provider-specific default_tags handling

### Adding New Validation Rules

1. Extend `TagRequirement` struct
2. Update configuration parsing
3. Implement validation logic
4. Add to validation pipeline

### Adding New Output Formats

1. Create new reporter interface
2. Implement format-specific logic
3. Integrate with main validation flow
4. Add CLI options for new format

This architecture provides a solid foundation for tag validation while maintaining flexibility for future enhancements and provider support.
