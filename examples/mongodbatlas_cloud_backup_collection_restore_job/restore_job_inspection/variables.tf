variable "project_id" {
  description = "Atlas project ID of the cluster that owns the collection restore jobs."
  type        = string
}

variable "cluster_name" {
  description = "Name of the cluster that owns the collection restore jobs."
  type        = string
}

variable "job_id" {
  description = "Collection restore job ID. When unset, defaults to the last job_id from mongodbatlas_cloud_backup_collection_restore_jobs."
  type        = string
  default     = null
}

variable "source_namespace" {
  description = "Namespace in database.collection form. When unset, defaults to the last source_namespace from mongodbatlas_cloud_backup_collection_restore_job_collections."
  type        = string
  default     = null
}
