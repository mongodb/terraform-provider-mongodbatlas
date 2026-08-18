---
subcategory: "Cloud Backups"
---

# Resource: mongodbatlas_cloud_backup_collection_restore_job

`mongodbatlas_cloud_backup_collection_restore_job` provides a resource to restore selected databases or collections from a Cloud Backup snapshot or point in time into an existing Atlas cluster.

To use this resource, the requesting Service Account or API Key must have the Backup Manager or Project Owner role on the source project and on the target project.

~> **Note:** Create waits until the job `state` is `SUCCESSFUL`. Pending states are `INITIALIZING`, `IN_PROGRESS`, and `FINALIZING`. The default create timeout is 3 hours and is overridable with `timeouts.create`. The job cannot be updated. Destroy only removes the resource from state; Atlas has no public cancel API, so there is no `delete_on_create_timeout`. Provide either `snapshot_id` or point-in-time fields (`point_in_time_utc_seconds`, `oplog_ts`, `oplog_inc`), not both. Per-collection progress is empty while the job is `INITIALIZING`. Cross-organization restore is not available with programmatic credentials; restore across projects in the same organization by setting `target_project_id` and `target_cluster_name`.

## Limitations

Atlas product limits apply. See [Restore from Selected Databases and Collections](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/#limitations) for the full list.

- One collection-restore job per cluster at a time.
- Does not restore Atlas Search indexes. Skips views, time-series, and queryable encryption collections.
- Data transfer uses the public internet even when private endpoints exist.
- NVMe clusters and fallback snapshots are unsupported.
- Up to 100 databases and 100 explicitly listed collections per job.

## Example Usage

The following example restores selected databases and collections from a snapshot into the same cluster, using snapshot namespace discovery data sources and reading back job and per-collection state.

```terraform
locals {
  target_project_id   = coalesce(var.target_project_id, var.project_id)
  target_cluster_name = coalesce(var.target_cluster_name, var.cluster_name)

  completed_snapshots = [
    for snapshot in data.mongodbatlas_cloud_backup_snapshots.source.results :
    snapshot if snapshot.status == "completed"
  ]
  latest_snapshot_created_at = length(local.completed_snapshots) > 0 ? max([
    for snapshot in local.completed_snapshots : snapshot.created_at
  ]) : null
  latest_snapshot_id = length(local.completed_snapshots) > 0 ? one([
    for snapshot in local.completed_snapshots :
    snapshot.id if snapshot.created_at == local.latest_snapshot_created_at
  ]) : null
  snapshot_id = coalesce(var.snapshot_id, local.latest_snapshot_id)
}

data "mongodbatlas_cloud_backup_snapshots" "source" {
  project_id     = var.project_id
  cluster_name   = var.cluster_name
  items_per_page = 100
}

data "mongodbatlas_cloud_backup_snapshot_databases" "source" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  snapshot_id  = local.snapshot_id
}

data "mongodbatlas_cloud_backup_snapshot_database_collections" "source" {
  for_each = toset(var.discovery_database_names)

  project_id    = var.project_id
  cluster_name  = var.cluster_name
  snapshot_id   = local.snapshot_id
  database_name = each.value
}

resource "mongodbatlas_cloud_backup_collection_restore_job" "example" {
  project_id          = var.project_id
  cluster_name        = var.cluster_name
  snapshot_id         = local.snapshot_id
  target_project_id   = local.target_project_id
  target_cluster_name = local.target_cluster_name
  write_strategy      = var.write_strategy
  index_strategy      = var.index_strategy

  databases = [for db in var.restore_databases : {
    source_namespace = db
  }]

  collections = [for ns in var.restore_collections : {
    source_namespace = ns
  }]

  timeouts = {
    create = "3h"
  }
}

# Point-in-time restore instead of snapshot_id:
# resource "mongodbatlas_cloud_backup_collection_restore_job" "pit" {
#   project_id                = var.project_id
#   cluster_name              = var.cluster_name
#   point_in_time_utc_seconds = 1710000000
#   target_project_id         = local.target_project_id
#   target_cluster_name       = local.target_cluster_name
#   write_strategy            = var.write_strategy
#   index_strategy            = var.index_strategy
#   databases = [for db in var.restore_databases : {
#     source_namespace = db
#   }]
#   collections = [for ns in var.restore_collections : {
#     source_namespace = ns
#   }]
# }

data "mongodbatlas_cloud_backup_collection_restore_job" "example" {
  project_id   = mongodbatlas_cloud_backup_collection_restore_job.example.project_id
  cluster_name = mongodbatlas_cloud_backup_collection_restore_job.example.cluster_name
  job_id       = mongodbatlas_cloud_backup_collection_restore_job.example.job_id
}

data "mongodbatlas_cloud_backup_collection_restore_jobs" "example" {
  project_id   = mongodbatlas_cloud_backup_collection_restore_job.example.project_id
  cluster_name = mongodbatlas_cloud_backup_collection_restore_job.example.cluster_name
}

data "mongodbatlas_cloud_backup_collection_restore_job_collections" "example" {
  project_id   = mongodbatlas_cloud_backup_collection_restore_job.example.project_id
  cluster_name = mongodbatlas_cloud_backup_collection_restore_job.example.cluster_name
  job_id       = mongodbatlas_cloud_backup_collection_restore_job.example.job_id
}

output "available_snapshots" {
  value = [for snapshot in data.mongodbatlas_cloud_backup_snapshots.source.results : {
    id           = snapshot.id
    created_at   = snapshot.created_at
    status       = snapshot.status
    snapshot_type = snapshot.snapshot_type
  }]
}

output "snapshot_id" {
  value = local.snapshot_id
}

output "snapshot_databases" {
  value = [for db in data.mongodbatlas_cloud_backup_snapshot_databases.source.results : db.name]
}

output "snapshot_collections" {
  value = {
    for db, ds in data.mongodbatlas_cloud_backup_snapshot_database_collections.source :
    db => [for coll in ds.results : coll.name]
  }
}

output "job_state" {
  value = data.mongodbatlas_cloud_backup_collection_restore_job.example.state
}

output "collection_states" {
  value = [for c in data.mongodbatlas_cloud_backup_collection_restore_job_collections.example.results : {
    source_namespace = c.source_namespace
    state            = c.state
  }]
}
```

### Further Examples
- [Collection Restore Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_backup_collection_restore_job)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_name` (String) Human-readable label that identifies the source cluster for the restore.
- `index_strategy` (String) Strategy for restoring indexes (all, none, or all except TTL).
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.
- `target_cluster_name` (String) Target cluster name.
- `target_project_id` (String) Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.
- `write_strategy` (String) Strategy for writing data on the target (create as new or overwrite existing). With `OVERWRITE_EXISTING`, any writes to the affected databases or collections during the restore will be lost when the existing namespaces are dropped and replaced. To avoid data loss, stop writes to the affected namespaces before starting the restore.

### Optional

- `collection_suffix` (String) Optional suffix applied to restored collection names.
- `collections` (Attributes List) List of collections to restore (up to 100 items). (see [below for nested schema](#nestedatt--collections))
- `database_suffix` (String) Optional suffix applied to restored database names.
- `databases` (Attributes List) List of databases to restore (up to 100 items). (see [below for nested schema](#nestedatt--databases))
- `oplog_inc` (Number) Oplog increment for point-in-time restore.
- `oplog_ts` (Number) Oplog timestamp (seconds part) for point-in-time restore.
- `point_in_time_utc_seconds` (Number) Point-in-time restore time in seconds since UNIX epoch.
- `snapshot_id` (String) ID of the snapshot to restore.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_at` (String) Date and time when the restore job was created (ISO 8601 format in UTC).
- `error_message` (String) Error message when the job has failed or been canceled.
- `finished_at` (String) Date and time when the restore job finished (ISO 8601 format in UTC).
- `index_status` (Attributes) Overall index build status for a collection restore job. (see [below for nested schema](#nestedatt--index_status))
- `job_id` (String) Unique 24-hexadecimal digit string that identifies the collection restore job.
- `restored_documents` (Number) Number of documents restored so far across all supported collections.
- `state` (String) Current state of the collection restore job.
- `total_documents` (Number) Total number of documents across all supported collections in the restore job. This value may initially reflect an estimate based on collection metadata and can change as accurate document counts become available during the restore.

<a id="nestedatt--collections"></a>
### Nested Schema for `collections`

Required:

- `source_namespace` (String) Namespace requested to restore (e.g. database name or `database.collection`).

Optional:

- `target_namespace` (String) Requested target namespace for the restored data; if empty, source namespace is used.


<a id="nestedatt--databases"></a>
### Nested Schema for `databases`

Required:

- `source_namespace` (String) Namespace requested to restore (e.g. database name or `database.collection`).

Optional:

- `target_namespace` (String) Requested target namespace for the restored data; if empty, source namespace is used.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--index_status"></a>
### Nested Schema for `index_status`

Read-Only:

- `failed_collection_count` (Number) Number of collections that failed to build indexes.
- `state` (String) Index build state indicating the status of index creation during or after a restore operation.

## Import

Collection restore job resource can be imported using the source project ID, source cluster name, and job ID, separated by slashes, for example:

```
$ terraform import mongodbatlas_cloud_backup_collection_restore_job.example 650972848269185c55f40ca1/MyCluster/6789abcd1234ef5678901234
```

Import hydrates the current GET state. It does not wait for a running job to finish.

For more information see: [MongoDB Atlas - Restore from Selected Databases and Collections](https://www.mongodb.com/docs/atlas/backup/cloud-backup/restore-from-db-coll/) Documentation.
