# Default destination is the source cluster when target_project_id and target_cluster_name are unset.
locals {
  target_project_id   = coalesce(var.target_project_id, var.project_id)
  target_cluster_name = coalesce(var.target_cluster_name, var.cluster_name)
}

resource "mongodbatlas_cloud_backup_collection_restore_job" "this" {
  project_id          = var.project_id
  cluster_name        = var.cluster_name
  snapshot_id         = var.snapshot_id
  target_project_id   = local.target_project_id
  target_cluster_name = local.target_cluster_name
  write_strategy      = var.write_strategy
  index_strategy      = var.index_strategy
  database_suffix     = var.database_suffix
  collection_suffix   = var.collection_suffix

  databases = [for db in var.restore_databases : {
    source_namespace = db
    target_namespace = try(var.database_renames[db], null)
  }]

  collections = [for ns in var.restore_collections : {
    source_namespace = ns
    target_namespace = try(var.collection_renames[ns], null)
  }]

  timeouts = {
    create = var.create_timeout
  }
}

# Illustration only; not referenced by outputs.
data "mongodbatlas_cloud_backup_collection_restore_job" "this" {
  project_id   = mongodbatlas_cloud_backup_collection_restore_job.this.project_id
  cluster_name = mongodbatlas_cloud_backup_collection_restore_job.this.cluster_name
  job_id       = mongodbatlas_cloud_backup_collection_restore_job.this.job_id
}

# Illustration only; not referenced by outputs.
data "mongodbatlas_cloud_backup_collection_restore_jobs" "this" {
  project_id   = mongodbatlas_cloud_backup_collection_restore_job.this.project_id
  cluster_name = mongodbatlas_cloud_backup_collection_restore_job.this.cluster_name
}

data "mongodbatlas_cloud_backup_collection_restore_job_collections" "this" {
  project_id   = mongodbatlas_cloud_backup_collection_restore_job.this.project_id
  cluster_name = mongodbatlas_cloud_backup_collection_restore_job.this.cluster_name
  job_id       = mongodbatlas_cloud_backup_collection_restore_job.this.job_id
}

output "job_state" {
  value = mongodbatlas_cloud_backup_collection_restore_job.this.state
}

output "collection_states" {
  value = [for c in data.mongodbatlas_cloud_backup_collection_restore_job_collections.this.results : {
    source_namespace           = c.source_namespace
    target_namespace           = c.target_namespace
    effective_target_namespace = c.effective_target_namespace
    state                      = c.state
  }]
}
