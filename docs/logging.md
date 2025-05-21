# Logging

Terratags provides a flexible logging system that allows you to control the verbosity and detail of output during execution.

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

## Log Output Format

Here are examples of how logs appear at different levels:

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
WARN    2025-05-18T20:33:48.269-0400    Warning: Some blocks in main.tf couldn't be parsed, but we'll continue with what we can parse
WARN    2025-05-18T20:33:48.269-0400    Error parsing provider blocks in main.tf: invalid syntax
```

### ERROR Level Output

```
ERROR   2025-05-18T20:33:48.270-0400    Failed to parse configuration file: invalid syntax at line 15
```

## Logging in CI/CD Environments

When using Terratags in CI/CD pipelines, consider the following:

1. Use the default ERROR level for normal operation to keep logs clean
2. Use the INFO level for more detailed output when debugging pipeline issues
3. Consider using the `-report` option to generate an HTML report for better visualization of results
4. Redirect logs to a file for later analysis if needed

Example GitHub Actions workflow:

```yaml
- name: Validate Tags
  run: |
    terratags -config config.yaml -dir ./infra -log-level INFO -report report.html
  continue-on-error: true

- name: Upload Report
  uses: actions/upload-artifact@v3
  with:
    name: tag-validation-report
    path: report.html
```

## Technical Implementation

Terratags uses a logging implementation based on [Zap](https://github.com/uber-go/zap).