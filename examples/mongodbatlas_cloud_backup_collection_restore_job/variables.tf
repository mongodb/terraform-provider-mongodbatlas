variable "project_id" {
  description = "Atlas project ID of the source cluster"
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label of the source cluster"
  type        = string
}

variable "snapshot_id" {
  description = "Snapshot ID to restore from. Omit this and set point_in_time_utc_seconds for a PIT restore instead."
  type        = string
}

variable "database_name" {
  description = "Database name used to list collections in the snapshot"
  type        = string
}

variable "source_namespace" {
  description = "Source namespace to restore, in database.collection form"
  type        = string
}

variable "target_project_id" {
  description = "Atlas project ID of the target cluster. Defaults to the source project."
  type        = string
  default     = null
}

variable "target_cluster_name" {
  description = "Human-readable label of the target cluster. Defaults to the source cluster."
  type        = string
  default     = null
}

variable "write_strategy" {
  description = "How to write restored data on the target. CREATE_NEW or OVERWRITE_EXISTING."
  type        = string
  default     = "CREATE_NEW"
}

variable "index_strategy" {
  description = "How to restore indexes. ALL, NONE, or ALL_EXCEPT_TTL."
  type        = string
  default     = "ALL"
}
