# MongoDB Atlas Collection Restore Example

You already have a live Atlas cluster with Cloud Backup. This example restores selected databases and collections from a snapshot into an existing cluster without wiping everything else.

Default destination is the same cluster as the snapshot source. Leave `target_project_id` and `target_cluster_name` unset to restore there. Set both to seed another cluster in the same organization, for example prod snapshots into staging.

Default `write_strategy` is `CREATE_NEW`. Live data stays; Atlas adds restored copies and renames on conflict. That is the inspect-a-copy path, not a production rollback.

You pick what to restore (`restore_databases` and `restore_collections`), where to restore it (source cluster or a target cluster), and how to write it (`CREATE_NEW` or `OVERWRITE_EXISTING`). Index strategy defaults to `ALL` in `variables.tf`.

The example does not create a cluster. Point it at a cluster that already has Cloud Backup and a snapshot that contains the namespaces you want.

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/cloud_backup_collection_restore_job#limitations) and [Restore from Selected Databases and Collections](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/).

## Prerequisites

- MongoDB Atlas Service Account with Backup Manager or Project Owner on the source project, and on the target project if they differ.
- An existing dedicated cluster with Cloud Backup enabled and at least one snapshot.
- The namespaces you list must exist in that snapshot.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`). Do not put keys in Terraform files.

## Usage

**1. Ensure your MongoDB Atlas credentials are set up.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

Create a **terraform.tfvars** file (gitignored):

```hcl
# Source cluster (snapshot comes from here)
project_id               = "your-prod-project-id"
cluster_name             = "your-prod-cluster-name"
discovery_database_names = ["sample_mflix", "inventory"]
restore_databases        = ["inventory"]
restore_collections      = ["sample_mflix.movies", "sample_mflix.comments"]

# Destination: same cluster (default — omit target_project_id and target_cluster_name)
# target_project_id   = "your-staging-project-id"
# target_cluster_name = "staging-app"

# write_strategy = "OVERWRITE_EXISTING"
```

Restore into another cluster in the same organization by uncommenting `target_project_id` and `target_cluster_name`. Cross-organization restore is not available with programmatic credentials.

`restore_databases` restores every collection under each listed database. `restore_collections` restores only the namespaces you name. You can use either list or both in the same job.

For a snapshot restore, leave `snapshot_id` unset to use the latest completed snapshot. Set it explicitly when you need an older snapshot:

```hcl
snapshot_id = "6789abcd1234ef5678901234"
```

For a point-in-time restore from continuous backup, uncomment the PIT resource block in `main.tf`, set `point_in_time_utc_seconds` (and optional `oplog_ts` / `oplog_inc`), and comment out the snapshot-based resource. Do not set both snapshot and PIT fields on the same resource.

**2. (Optional) Discover snapshots and namespaces.**

Skip this step if you already know what to restore. A full `terraform apply` starts the restore job.

This step applies to snapshot restores only. Point-in-time restores have no snapshot ID, so the snapshot discovery data sources do not apply.

To list snapshots and pick one, apply only the snapshots data source:

```bash
terraform apply -target=data.mongodbatlas_cloud_backup_snapshots.source
```

Read `available_snapshots` from the outputs. The example uses the latest completed snapshot when `snapshot_id` is unset. The `snapshot_id` output shows which snapshot the restore will use.

To browse databases and collections in that snapshot, set `discovery_database_names` and apply the discovery data sources:

```bash
terraform apply \
  -target=data.mongodbatlas_cloud_backup_snapshots.source \
  -target=data.mongodbatlas_cloud_backup_snapshot_databases.source \
  -target=data.mongodbatlas_cloud_backup_snapshot_database_collections.source
```

Read `snapshot_databases` and `snapshot_collections` from the outputs, then update `restore_databases` and `restore_collections` in `terraform.tfvars`.

**3. Review the Terraform plan.**

```bash
terraform plan
```

**4. Apply the restore.**

```bash
terraform apply
```

If you set `write_strategy = "OVERWRITE_EXISTING"`, stop writers to the affected namespaces before apply. Overwrite drops and replaces matching namespaces on the target.

Create waits until the job `state` is `SUCCESSFUL`. Default timeout is 3 hours. Atlas keeps running if Terraform times out.

After apply, check `job_state` and `collection_states`. Job `SUCCESSFUL` is not enough if a collection is skipped or unsupported.

**5. Destroy the Terraform resource.**

```bash
terraform destroy
```

Destroy removes the job from state only. Restored data stays on the cluster. For the default `CREATE_NEW`, drop the restored copies yourself when you no longer need them.
