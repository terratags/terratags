# CI/CD Integration

Integrating Terratags into your CI/CD pipeline helps enforce tag compliance across your infrastructure. This page provides examples of how to integrate Terratags with popular CI/CD platforms.

## GitHub Actions

Add Terratags to your GitHub Actions workflow:

```yaml
name: Validate Tags

on:
  pull_request:
    paths:
      - '**.tf'

jobs:
  validate-tags:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
          
      - name: Install Terratags
        run: go install github.com/terratags/terratags@latest
        
      - name: Validate Tags
        run: terratags -config config.yaml -dir ./infra
```

## GitLab CI

Add Terratags to your GitLab CI pipeline:

```yaml
stages:
  - validate

validate-tags:
  stage: validate
  image: golang:1.24
  script:
    - go install github.com/terratags/terratags@latest
    - terratags -config config.yaml -dir ./infra
  only:
    changes:
      - "**/*.tf"
```

## Azure DevOps

Add Terratags to your Azure DevOps pipeline:

```yaml
trigger:
  paths:
    include:
    - '**/*.tf'

pool:
  vmImage: 'ubuntu-latest'

steps:
- task: GoTool@0
  inputs:
    version: '1.24'

- script: |
    go install github.com/terratags/terratags@latest
    terratags -config config.yaml -dir ./infra
  displayName: 'Validate Tags'
```

## Jenkins

Add Terratags to your Jenkinsfile:

```groovy
pipeline {
    agent {
        docker {
            image 'golang:1.24'
        }
    }
    
    stages {
        stage('Validate Tags') {
            when {
                changeset "**/*.tf"
            }
            steps {
                sh 'go install github.com/terratags/terratags@latest'
                sh 'terratags -config config.yaml -dir ./infra'
            }
        }
    }
}
```

## CircleCI

Add Terratags to your CircleCI configuration:

```yaml
version: 2.1
jobs:
  validate-tags:
    docker:
      - image: cimg/go:1.24
    steps:
      - checkout
      - run:
          name: Install Terratags
          command: go install github.com/terratags/terratags@latest
      - run:
          name: Validate Tags
          command: terratags -config config.yaml -dir ./infra

workflows:
  version: 2
  terraform-workflow:
    jobs:
      - validate-tags:
          filters:
            paths:
              - "**/*.tf"
```

## Best Practices for CI/CD Integration

1. **Fail Fast**: Configure your pipeline to fail early if tag validation fails
2. **Generate Reports**: Use the `-report` flag to generate HTML reports for each build
3. **Artifact Storage**: Store the generated reports as build artifacts for easy access
4. **Selective Validation**: Use path filters to only run validation when Terraform files change
5. **Pre-commit Hooks**: Consider adding Terratags as a pre-commit hook for local validation before pushing

## Example: Complete GitHub Actions Workflow

Here's a more complete example for GitHub Actions that includes report generation and artifact storage:

```yaml
name: Terraform Tag Validation

on:
  pull_request:
    paths:
      - '**.tf'
  push:
    branches:
      - main
    paths:
      - '**.tf'

jobs:
  validate-tags:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
          
      - name: Install Terratags
        run: go install github.com/terratags/terratags@latest
        
      - name: Validate Tags
        run: terratags -config config.yaml -dir ./infra -report tag-report.html
        
      - name: Upload Report
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: tag-validation-report
          path: tag-report.html
```