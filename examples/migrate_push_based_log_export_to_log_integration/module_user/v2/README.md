# v2: Module User - Migration Complete

This configuration shows the final state after migration where the old resource has been removed.

## Key Points

- `skip_push_based_log_export = true` removes the old resource.
- Only `mongodbatlas_log_integration` is now active.
- Can optionally use the same prefix path as before.

## Migration Complete

At this point:
- Logs are being exported every 1 minute (instead of 5 minutes).
- You have control over which log types to export.
- The old `push_based_log_export` has been removed.

