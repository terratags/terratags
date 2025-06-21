# Standard Examples for Documentation

This document defines the canonical examples to be used consistently across all Terratags documentation.

## Standard Examples by Pattern Type

### Environment Values
**Pattern**: `^(dev|test|staging|prod)$`
**Standard Examples**:
- `dev` (development)
- `test` (testing)
- `staging` (staging)
- `prod` (production)

### Project Codes
**Pattern**: `^[A-Z]{2,4}-[0-9]{3,6}$`
**Standard Examples**:
- `WEB-123456` (web applications)
- `DATA-567890` (data projects)  
- `SEC-123456` (security projects)
- `INFRA-890123` (infrastructure)
- `API-456789` (API projects)

### Cost Centers
**Pattern**: `^CC-[0-9]{4}$`
**Standard Examples**:
- `CC-1234` (engineering)
- `CC-5678` (operations)
- `CC-9012` (security)
- `CC-3456` (data)

### Email Addresses
**Pattern**: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`
**Standard Examples**:
- `devops@company.com` (DevOps team)
- `team.lead@company.com` (Team lead)
- `security@company.com` (Security team)
- `data.team@company.com` (Data team)

### Resource Names (No Whitespace)
**Pattern**: `^\\S+$`
**Standard Examples**:
- `web-server-01` (web server)
- `data-bucket` (S3 bucket)
- `main-vpc` (VPC)
- `allow-http-sg` (security group)

### Version Numbers
**Pattern**: `^v?[0-9]+\\.[0-9]+\\.[0-9]+$`
**Standard Examples**:
- `1.0.0` (without prefix)
- `v2.1.3` (with prefix)
- `10.15.2` (multi-digit)

## Usage Guidelines

1. **Always use these exact examples** in documentation
2. **Don't create new examples** without updating this standard
3. **Test all examples** with actual patterns before documenting
4. **Update this file first** when adding new pattern types

## Validation

All examples in this document have been tested with their corresponding patterns to ensure they work correctly.
