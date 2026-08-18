variable "project_id" {
  description = "Atlas project ID of the source cluster"
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label of the source cluster"
  type        = string
}

variable "snapshot_id" {
  description = "Snapshot ID to restore from. Defaults to the latest completed snapshot when unset."
  type        = string
  default     = null
}

variable "discovery_database_names" {
  description = "Database names to list collections from during snapshot discovery. Leave empty to list only databases."
  type        = list(string)
  default     = []
}

variable "restore_databases" {
  description = "Database names to restore. Restores all collections under each database."
  type        = list(string)
  default     = []
}

variable "restore_collections" {
  description = "Collection namespaces to restore, in database.collection form."
  type        = list(string)
  default     = []
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
