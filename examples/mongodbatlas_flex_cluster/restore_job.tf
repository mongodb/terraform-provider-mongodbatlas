data "mongodbatlas_flex_restore_job" "restore_job" {
  project_id     = var.project_id
  name           = mongodbatlas_flex_cluster.example-cluster.name
  restore_job_id = var.restore_job_id
}

data "mongodbatlas_flex_restore_jobs" "restore_jobs" {
  project_id = var.project_id
  name       = mongodbatlas_flex_cluster.example-cluster.name
}

output "mongodbatlas_flex_restore_job" {
  value = data.mongodbatlas_flex_restore_job.restore_job.name
}

output "mongodbatlas_flex_restore_jobs" {
  value = [for restore_job in data.mongodbatlas_flex_restore_jobs.restore_jobs.results : restore_job.restore_job_id]
}
