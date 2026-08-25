# List Cloud Backup snapshots for a cluster, resolve a snapshot_id, and browse databases and collections in that snapshot.
locals {
  completed_snapshots = [
    for snapshot in data.mongodbatlas_cloud_backup_snapshots.this.results :
    snapshot if snapshot.status == "completed"
  ]
  completed_created_at = [for snapshot in local.completed_snapshots : snapshot.created_at]
  /*
    created_at is ISO-8601. max() only accepts numbers, and max([list]) passes a tuple instead of
    spreading values. ISO-8601 sorts lexicographically, so sort(...)[last] is the latest snapshot.
  */
  latest_snapshot_created_at = length(local.completed_created_at) > 0 ? sort(local.completed_created_at)[length(local.completed_created_at) - 1] : null
  latest_snapshot_id = length(local.completed_snapshots) > 0 ? one([
    for snapshot in local.completed_snapshots :
    snapshot.id if snapshot.created_at == local.latest_snapshot_created_at
  ]) : null
  snapshot_id = coalesce(var.snapshot_id, local.latest_snapshot_id)
  selected_snapshot_created_at = try(one([
    for s in local.available_snapshots_all : s.created_at if s.id == local.snapshot_id
  ]), null)
  available_snapshots_all = [
    for snapshot in data.mongodbatlas_cloud_backup_snapshots.this.results : {
      id            = snapshot.id
      created_at    = snapshot.created_at
      status        = snapshot.status
      snapshot_type = snapshot.snapshot_type
    }
  ]
  available_snapshots_newest = [
    for created_at in slice(
      reverse(sort([for s in local.completed_snapshots : s.created_at])),
      0,
      min(var.available_snapshots_limit, length(local.completed_snapshots)),
      ) : [
      for s in local.available_snapshots_all : s if s.created_at == created_at && s.status == "completed"
    ][0]
  ]
}

data "mongodbatlas_cloud_backup_snapshots" "this" {
  project_id     = var.project_id
  cluster_name   = var.cluster_name
  items_per_page = 100
}

data "mongodbatlas_cloud_backup_snapshot_databases" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  snapshot_id  = local.snapshot_id
}

data "mongodbatlas_cloud_backup_snapshot_database_collections" "this" {
  for_each = toset(data.mongodbatlas_cloud_backup_snapshot_databases.this.results[*].name)

  project_id    = var.project_id
  cluster_name  = var.cluster_name
  snapshot_id   = local.snapshot_id
  database_name = each.value
}

output "available_snapshots" {
  description = "Newest completed snapshots for the cluster (up to available_snapshots_limit)."
  value       = local.available_snapshots_newest
}

output "snapshot_id" {
  description = "Snapshot ID used for database and collection listing."
  value       = local.snapshot_id
}

output "snapshot_details" {
  description = "Created time, databases, and collections in the selected snapshot."
  value = {
    created_at = local.selected_snapshot_created_at
    databases = [for db in data.mongodbatlas_cloud_backup_snapshot_databases.this.results : {
      name = db.name
      collections = [
        for coll in try(data.mongodbatlas_cloud_backup_snapshot_database_collections.this[db.name].results, []) :
        coll.name
      ]
    }]
  }
}
