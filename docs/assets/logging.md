# Logging

Terratags provides a flexible logging system that allows you to control the verbosity and detail of output during execution. This document explains how to configure and use the logging functionality.

## Log Levels

Terratags supports the following log levels, in order of increasing severity:

- `DEBUG`: Detailed information, typically useful only for diagnosing problems
- `INFO`: Confirmation that things are working as expected
- `WARN`: Indication that something unexpected happened, but the process can continue
- `ERROR`: Due to a more serious problem, the process couldn't perform a specific function

## Setting the Log Level

You can set the log level using the `-log-level` or `-l` command-line flag:

```bash
terratags -config config.yaml -log-level DEBUG
```

Or using the short form:

```bash
terratags -c config.yaml -l DEBUG
```

The default log level is `ERROR` if not specified.

## Examples

### Basic Usage

```bash
# Run with default ERROR level logging
terratags -config config.yaml

# Run with INFO level for more detailed output
terratags -config config.yaml -log-level INFO

# Run with DEBUG level for maximum verbosity
terratags -config config.yaml -log-level DEBUG
```

### Troubleshooting

If you're experiencing issues with Terratags, running with the `DEBUG` log level can provide additional information to help diagnose the problem:

```bash
terratags -config config.yaml -log-level DEBUG
```

This will output detailed information about each step of the process, including file parsing, tag validation, and any errors encountered.

## Log Output Format

Terratags uses a custom logging format where the log level appears before the timestamp. Here are examples of how logs appear at different levels:

### DEBUG Level Output

```
DEBUG   2025-05-18T20:33:48.268-0400    Found AWSCC tags attribute in awscc_s3_bucket name
DEBUG   2025-05-18T20:33:48.268-0400    Found AWSCC tag key: Name with value: test
```

### INFO Level Output

```
INFO    2025-05-18T20:33:48.265-0400    Loaded configuration with 4 required tags
INFO    2025-05-18T20:33:48.266-0400    Validating Terraform directory: ../awscc_examples/
INFO    2025-05-18T20:33:48.267-0400    Found 1 Terraform files to analyze
INFO    2025-05-18T20:33:48.267-0400    Analyzing file: ../awscc_examples/main.tf
```

### WARN Level Output

```
WARN    2025-05-18T20:33:48.269-0400    Resource aws_s3_bucket.example is missing required tag: Environment
```

### ERROR Level Output

```
ERROR   2025-05-18T20:33:48.270-0400    Failed to parse configuration file: invalid syntax at line 15
```

## Logging Implementation

Terratags uses a custom logging implementation based on [Zap](https://github.com/uber-go/zap), which provides:

- High-performance, structured logging
- Configurable output formats
- Multiple log levels