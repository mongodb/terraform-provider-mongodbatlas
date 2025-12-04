# Module Usage Example

This example demonstrates how to use the effective fields modules from a user's perspective. It shows that switching between `module_existing` and `module_effective_fields` requires **only** changing the module source path - all other parameters remain identical.

## Key Takeaway

**Module migration is transparent to users:**
- Change module source from `../module_existing` to `../module_effective_fields`, or change module version
- No other configuration changes needed
- All input variables stay the same
- All outputs remain available (with enhanced functionality in module_effective_fields)

## Usage

### Configure Credentials

Set environment variables for MongoDB Atlas authentication and organization ID.

### Deploy with module_effective_fields (recommended)

The module source in `main.tf` is already set to `module_effective_fields`. Run standard Terraform commands to deploy.

### Switch to module_existing

To see the old approach with `lifecycle.ignore_changes`, edit the module source in `main.tf` to reference `../module_existing` and run `terraform init -reconfigure`. No other configuration changes needed.

## Output Differences

### Both modules provide:
- `cluster_id`, `cluster_state`, `connection_strings`
- `replication_specs` - Hardware specifications

### module_effective_fields additionally provides:
- `auto_scaling_enabled`, `analytics_auto_scaling_enabled` - Auto-scaling status flags

### Behavioral difference:

**module_existing:**
- `*_specs` returns actual provisioned values (may be auto-scaled)
- Cannot distinguish configured from auto-scaled values

**module_effective_fields (Phase 1 - backward compatible):**
- `*_specs` returns actual provisioned values (same as module_existing)
- `effective_*_specs` attributes also available with actual values
- Seamless migration with no output changes

**module_effective_fields (Phase 2 - breaking change, prepares for v3.x):**
- If data source uses `use_effective_fields = true`:
  - `*_specs` returns configured values (client-provided intent)
  - `effective_*_specs` returns actual provisioned values (Atlas-managed reality)
  - **BREAKING CHANGE:** Must switch from `*_specs` to `effective_*_specs` for actual values
  - Prepares for provider v3.x where this becomes default behavior
  - Clear separation between what you configured and what Atlas provisioned

## Complete Migration Example

To migrate from `module_existing` to `module_effective_fields`, only update the module source path. All input variables and configuration remain identical. See `main.tf` for the complete implementation.

## Testing

Use standard Terraform commands to validate, plan, and apply the configuration. View outputs using `terraform output` after deployment.
