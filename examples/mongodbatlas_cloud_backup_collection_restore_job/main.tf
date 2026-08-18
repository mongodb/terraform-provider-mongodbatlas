locals {
  target_project_id   = coalesce(var.target_project_id, var.project_id)
  target_cluster_name = coalesce(var.target_cluster_name, var.cluster_name)
}

data "mongodbatlas_cloud_backup_snapshot_databases" "source" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  snapshot_id  = var.snapshot_id
}

data "mongodbatlas_cloud_backup_snapshot_database_collections" "source" {
  project_id    = var.project_id
  cluster_name  = var.cluster_name
  snapshot_id   = var.snapshot_id
  database_name = var.database_name
}

resource "mongodbatlas_cloud_backup_collection_restore_job" "example" {
  project_id          = var.project_id
  cluster_name        = var.cluster_name
  snapshot_id         = var.snapshot_id
  target_project_id   = local.target_project_id
  target_cluster_name = local.target_cluster_name
  write_strategy      = var.write_strategy
  index_strategy      = var.index_strategy

  collections = [{
    source_namespace = var.source_namespace
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
#   collections = [{
#     source_namespace = var.source_namespace
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

output "snapshot_databases" {
  value = [for db in data.mongodbatlas_cloud_backup_snapshot_databases.source.results : db.name]
}

output "snapshot_collections" {
  value = [for coll in data.mongodbatlas_cloud_backup_snapshot_database_collections.source.results : coll.name]
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
