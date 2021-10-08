output "project_id" {
  value = mongodbatlas_project.project_test.id
}
output "cluster_name" {
  value = mongodbatlas_cluster.cluster_test.name
}
output "snapshot_id" {
  value = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
}
output "snapshot_restore_job_id" {
  value = mongodbatlas_cloud_provider_snapshot_restore_job.test.snapshot_restore_job_id
}
