# v1: Module User - Initial Upgrade

This configuration shows how to upgrade to the new module version while keeping both resources active.

## Key Points

- `skip_push_based_log_export = false` keeps the old resource active
- Both log export configurations are running in parallel
- Use distinct prefix paths to avoid log conflicts

## Next Steps

1. Apply this configuration
2. Validate logs appear at the new prefix path (`atlas-logs-new`)
3. Once validated, proceed to v2 to remove the old resource

