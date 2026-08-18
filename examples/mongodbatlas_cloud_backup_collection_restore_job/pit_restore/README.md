# Point-in-time collection restore

Restore selected databases and collections from continuous Cloud Backup into an existing cluster at a specific point in time.

This example does not use snapshot IDs. [`snapshot_discovery/`](../snapshot_discovery/README.md) does not apply. For snapshot-based restores, use [`snapshot_restore/`](../snapshot_restore/README.md).

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_backup_collection_restore_job#limitations).

## Prerequisites

- Backup Manager or Project Owner on the source project, and on the target project if they differ.
- A dedicated cluster with continuous Cloud Backup enabled.
- A restore time within the continuous backup window.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Defaults

Omit optional variables to use these values from `variables.tf`:

- `write_strategy`: `CREATE_NEW` — append restored data; Atlas renames matching namespaces on conflict instead of replacing live data
- `index_strategy`: `ALL` — restore all indexes; `NONE` skips indexes; `ALL_EXCEPT_TTL` restores non-TTL indexes only
- `target_project_id` / `target_cluster_name`: unset — restore to the same project and cluster as the backup source
- `oplog_ts` / `oplog_inc`: unset — use `point_in_time_utc_seconds` only
- `restore_databases` / `restore_collections`: `[]` — set at least one list with namespaces to restore

To drop and replace matching namespaces on the target, set `write_strategy = "OVERWRITE_EXISTING"` and stop writers before apply.

## Namespace targets

- Default: restore under source names. Use per-key maps (`database_renames`, `collection_renames`) or global suffixes (`database_suffix`, `collection_suffix`) to rename; see commented tfvars below.
- Do not combine maps and suffixes without checking [product docs](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/) for precedence.
- Check `effective_target_namespace` in outputs when Atlas adjusts names on conflict.

## Usage

**1. Set credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2. Create `terraform.tfvars`.**

Required inputs plus optional overrides. Lines below match `variables.tf` defaults; omit them to keep defaults.

```hcl
project_id                = "your-prod-project-id"
cluster_name              = "your-prod-cluster-name"
point_in_time_utc_seconds = 1710000000
restore_databases         = ["inventory"]
restore_collections       = ["orders.refunds"]

# write_strategy      = "CREATE_NEW"            # default
# index_strategy      = "ALL"                   # default
# target_project_id   = "your-staging-project-id"  # default: same as project_id
# target_cluster_name = "staging-app"              # default: same as cluster_name
# oplog_ts            = 1710000000                # default: unset; alternative to point_in_time_utc_seconds
# oplog_inc           = 1                         # default: unset; use with oplog_ts
# write_strategy      = "OVERWRITE_EXISTING"    # optional: replace matching namespaces
```

Per-key rename example (commented):

```hcl
# database_renames = {
#   inventory = "inventory_staging"
# }
# collection_renames = {
#   "orders.refunds" = "orders.refunds_restored"
# }
```

Global suffix example (commented):

```hcl
# database_suffix   = "_staging"
# collection_suffix = "_restored"
```

Do not set `snapshot_id` or other snapshot fields on the resource.

**3. Plan and apply.**

```bash
terraform plan
terraform apply
```

Create waits until the job `state` is `SUCCESSFUL`. Default create timeout is 3 hours. Check `job_state` and `collection_states` after apply.

**4. Destroy.**

```bash
terraform destroy
```

Destroy removes the job from state only. Restored data stays on the cluster. With the default `CREATE_NEW`, drop restored copies yourself when you no longer need them.
