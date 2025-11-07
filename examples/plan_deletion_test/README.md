# Plan Deletion Test

This example tests that terratags correctly handles Terraform plans containing resource deletions.

## Test Plan Structure

The `test_plan.json` contains:
- 1 resource being created (should be validated)
- 2 resources being deleted (should be skipped)

## Expected Behavior

When running terratags on a plan with deletions:
- Resources with `"actions": ["delete"]` are skipped
- Resources with `"after": null` are skipped
- Only non-deleted resources are validated

## Running the Test

```bash
# From repository root
./terratags -config examples/config.yaml -plan examples/plan_deletion_test/test_plan.json -verbose
```

Expected output:
- Found 1 direct resources and 0 module resources
- All resources have the required tags!
