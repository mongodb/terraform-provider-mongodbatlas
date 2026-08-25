# Read back collection restore jobs on a cluster without creating a new job.
data "mongodbatlas_cloud_backup_collection_restore_jobs" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
}

output "restore_jobs" {
  description = "Trimmed list of collection restore jobs on the cluster."
  value = [for j in data.mongodbatlas_cloud_backup_collection_restore_jobs.this.results : {
    job_id     = j.job_id
    state      = j.state
    created_at = j.created_at
  }]
}

locals {
  jobs   = data.mongodbatlas_cloud_backup_collection_restore_jobs.this.results
  job_id = coalesce(var.job_id, length(local.jobs) > 0 ? local.jobs[length(local.jobs) - 1].job_id : null)
}

output "job_id" {
  description = "Resolved collection restore job ID."
  value       = local.job_id
}

data "mongodbatlas_cloud_backup_collection_restore_job" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  job_id       = local.job_id
}

output "job_state" {
  description = "State of the resolved job from the singular job data source."
  value       = data.mongodbatlas_cloud_backup_collection_restore_job.this.state
}
