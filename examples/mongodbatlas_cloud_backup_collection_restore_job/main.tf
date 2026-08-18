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
