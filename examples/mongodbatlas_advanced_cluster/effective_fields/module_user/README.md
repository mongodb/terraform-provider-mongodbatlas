# Module Usage Example

This example demonstrates how to use the effective fields modules from a user's perspective. It shows that switching between `module_existing` and `module_effective_fields` requires **only** changing the module source path - all other parameters remain identical.

## Key Takeaway

**Module migration is transparent to users:**
- Change module source from `../module_existing` to `../module_effective_fields`
- No other configuration changes needed
- All input variables stay the same
- All outputs remain available (with enhanced functionality in module_effective_fields)

## Usage

### Configure Credentials

```bash
export MONGODB_ATLAS_CLIENT_ID="your-client-id"
export MONGODB_ATLAS_CLIENT_SECRET="your-client-secret"
export TF_VAR_atlas_org_id="your-org-id"
```

### Deploy with module_effective_fields (recommended)

The module source in `main.tf` is already set to `module_effective_fields`:

```bash
terraform init
terraform plan
terraform apply
```

### Switch to module_existing

To see the old approach with `lifecycle.ignore_changes`:

1. Edit `main.tf` line 11:
   ```terraform
   source = "../module_existing"  # Changed from ../module_effective_fields
   ```

2. Run:
   ```bash
   terraform init -reconfigure
   terraform plan
   terraform apply
   ```

**That's it!** No other changes needed.

## Output Differences

### Both modules provide:
- `cluster_id` - Atlas cluster ID
- `cluster_state` - Current cluster state
- `connection_strings` - Connection strings (sensitive)
- `configured_specs` - Hardware specifications

### module_effective_fields additionally provides:
- `effective_specs` - Actual provisioned specs (may differ from configured due to auto-scaling)
- `auto_scaling_enabled` - Auto-scaling status
- `analytics_auto_scaling_enabled` - Analytics auto-scaling status

### Behavioral difference:

**module_existing:**
- `configured_specs` returns state values (may be auto-scaled)
- Cannot distinguish configured from auto-scaled values

**module_effective_fields:**
- `configured_specs` returns configuration values (constant)
- `effective_specs` returns actual provisioned values (may be auto-scaled)
- Clear separation between what you configured and what Atlas provisioned

## Complete Migration Example

### Before (using module_existing)

```terraform
module "atlas_cluster" {
  source = "../module_existing"  # v1.0

  atlas_org_id = var.atlas_org_id
  project_name = "MyProject"
  cluster_name = "my-cluster"
  cluster_type = "REPLICASET"
  # ... rest of configuration
}
```

### After (using module_effective_fields)

```terraform
module "atlas_cluster" {
  source = "../module_effective_fields"  # v2.0 - ONLY change needed

  atlas_org_id = var.atlas_org_id
  project_name = "MyProject"
  cluster_name = "my-cluster"
  cluster_type = "REPLICASET"
  # ... same configuration, no changes needed
}
```

## Testing

Validate the configuration:

```bash
terraform init
terraform validate
terraform fmt -check
```

View the plan:

```bash
terraform plan
```

View outputs after apply:

```bash
terraform output
terraform output -json configured_specs | jq
terraform output -json effective_specs | jq  # Only with module_effective_fields
```
