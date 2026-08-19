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
  for_each = toset(var.discovery_database_names)

  project_id    = var.project_id
  cluster_name  = var.cluster_name
  snapshot_id   = local.snapshot_id
  database_name = each.value
}

output "available_snapshots" {
  value = [for snapshot in data.mongodbatlas_cloud_backup_snapshots.this.results : {
    id            = snapshot.id
    created_at    = snapshot.created_at
    status        = snapshot.status
    snapshot_type = snapshot.snapshot_type
  }]
}

output "snapshot_id" {
  value = local.snapshot_id
}

output "snapshot_databases" {
  value = [for db in data.mongodbatlas_cloud_backup_snapshot_databases.this.results : db.name]
}

output "snapshot_collections" {
  value = {
    for db, ds in data.mongodbatlas_cloud_backup_snapshot_database_collections.this :
    db => [for coll in ds.results : coll.name]
  }
}
