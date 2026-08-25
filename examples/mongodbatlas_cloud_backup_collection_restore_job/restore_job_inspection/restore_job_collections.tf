# Per-collection read-back for the resolved collection restore job.
data "mongodbatlas_cloud_backup_collection_restore_job_collections" "this" {
  project_id   = var.project_id
  cluster_name = var.cluster_name
  job_id       = local.job_id
}

output "collection_states" {
  description = "Per-collection state for the resolved job."
  value = [for c in data.mongodbatlas_cloud_backup_collection_restore_job_collections.this.results : {
    source_namespace           = c.source_namespace
    target_namespace           = c.target_namespace
    effective_target_namespace = c.effective_target_namespace
    state                      = c.state
    index_status               = c.index_status
    documents                  = "${c.restored_documents} / ${c.total_documents}"
  }]
}

locals {
  collections      = data.mongodbatlas_cloud_backup_collection_restore_job_collections.this.results
  source_namespace = coalesce(var.source_namespace, length(local.collections) > 0 ? local.collections[length(local.collections) - 1].source_namespace : null)
}

data "mongodbatlas_cloud_backup_collection_restore_job_collection" "this" {
  project_id       = var.project_id
  cluster_name     = var.cluster_name
  job_id           = local.job_id
  source_namespace = local.source_namespace
}

output "collection_state" {
  description = "State of one collection from the singular collection data source."
  value       = data.mongodbatlas_cloud_backup_collection_restore_job_collection.this.state
}
