# v2: During Migration - Both Resources Active

This configuration shows the migration phase where both `mongodbatlas_push_based_log_export` and `mongodbatlas_log_integration` resources are active simultaneously.

## Key Points

- **Distinct Prefix Paths**: The new `mongodbatlas_log_integration` uses a different prefix path (`atlas-logs-new`) to avoid log conflicts
- **Parallel Operation**: Both resources export logs to the same S3 bucket but different paths
- **Validation Period**: Use this phase to verify logs appear correctly at the new path before removing the old resource

## Next Steps

After validating logs are appearing at the new prefix path:
1. Verify log content and format at `atlas-logs-new/`
2. Check that all expected log types are present
3. Once confident, proceed to v3 to remove the old resource

