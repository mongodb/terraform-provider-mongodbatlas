# Point-in-time collection restore

Restore selected databases and collections from continuous Cloud Backup at a point in time into an existing cluster, optionally in a different project or cluster.

This example does not use snapshot IDs. [`snapshot_discovery/`](../snapshot_discovery/README.md) does not apply. For snapshot-based restores, use [`snapshot_restore/`](../snapshot_restore/README.md).

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_backup_collection_restore_job#limitations).

## Restore time

Set **one** of the following. Atlas [collection restore](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/) and [continuous backup restore](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-continuous) use the same split (Date & Time vs Oplog Timestamp in the UI). Do not combine the paths, and do not set `snapshot_id`. The time must fall inside the cluster's continuous backup window.

- **`point_in_time_utc_seconds`**: Wall-clock UNIX epoch seconds. Use this when you know the clock time to restore to.
- **`oplog_ts` and `oplog_inc`**: [Oplog](https://www.mongodb.com/docs/manual/core/replica-set-oplog/) [Timestamp](https://www.mongodb.com/docs/manual/reference/bson-types/#timestamps): seconds since epoch plus the ordinal of operations in that second. Use this when you have a Timestamp from the oplog, or need a specific operation inside a second.

## Prerequisites

- Backup Manager or Project Owner on the source project, and on the target project if they differ.
- A dedicated cluster with continuous Cloud Backup enabled.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Defaults

Omit optional variables to use these values from `variables.tf`:

- `write_strategy`: `CREATE_NEW`. Append restored data; Atlas renames matching namespaces on conflict instead of replacing live data.
- `index_strategy`: `ALL`. Restore all indexes; `NONE` skips indexes; `ALL_EXCEPT_TTL` restores non-TTL indexes only.
- `target_project_id` / `target_cluster_name`: Unset. Restore to the same project and cluster as the backup source.
- `point_in_time_utc_seconds` / `oplog_ts` / `oplog_inc`: Unset. Set one restore-time path above.
- `restore_databases` / `restore_collections`: Unset. Set at least one list with namespaces to restore.
- `create_timeout`: `"3h"`. Maximum time to wait for the restore job to reach SUCCESSFUL during create.

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

Required inputs plus one restore-time path. Lines below match `variables.tf` defaults; omit them to keep defaults.

```hcl
project_id          = "your-prod-project-id"
cluster_name        = "your-prod-cluster-name"
restore_databases   = ["inventory"]
restore_collections = ["orders.refunds"]

# write_strategy      = "CREATE_NEW"            # default
# index_strategy      = "ALL"                   # default
# target_project_id   = "your-staging-project-id"  # default: same as project_id
# target_cluster_name = "staging-app"              # default: same as cluster_name
# write_strategy      = "OVERWRITE_EXISTING"    # optional: replace matching namespaces
```

Wall clock:

```hcl
point_in_time_utc_seconds = 1710000000
```

Oplog Timestamp:

```hcl
oplog_ts  = 1710000000
oplog_inc = 1
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

**3. Plan and apply.**

```bash
terraform plan
terraform apply
```

Create waits until the job `state` is `SUCCESSFUL`. Default create timeout is 3 hours. Check `job_state` and `collection_states` after apply.

For job read-back data sources after apply, see [`restore_job_inspection/`](../restore_job_inspection/README.md).

**4. Destroy.**

```bash
terraform destroy
```

Destroy removes the job from state only. Restored data stays on the cluster. With the default `CREATE_NEW`, drop restored copies yourself when you no longer need them.
