data "mongodbatlas_cloud_backup_collection_restore_jobs" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
}

output "restore_jobs" {
  value = [for j in data.mongodbatlas_cloud_backup_collection_restore_jobs.this.results : {
    job_id     = j.job_id
    state      = j.state
    created_at = j.created_at
  }]
}

locals {
  job_id = coalesce(var.job_id, data.mongodbatlas_cloud_backup_collection_restore_jobs.this.results[length(data.mongodbatlas_cloud_backup_collection_restore_jobs.this.results) - 1].job_id)
}

output "job_id" {
  value = local.job_id
}

data "mongodbatlas_cloud_backup_collection_restore_job" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  job_id       = local.job_id
}

output "job_state" {
  value = data.mongodbatlas_cloud_backup_collection_restore_job.this.state
}
