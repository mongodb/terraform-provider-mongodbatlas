variable "project_id" {
  description = "Atlas project ID of the source cluster"
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label of the source cluster"
  type        = string
}

variable "snapshot_id" {
  description = "Snapshot ID to restore from"
  type        = string
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

variable "database_renames" {
  description = "Optional map from source database name to target database name. Keys must match entries in restore_databases."
  type        = map(string)
  default     = {}

  validation {
    condition = alltrue([
      for key in keys(var.database_renames) : contains(var.restore_databases, key)
    ])
    error_message = "database_renames keys must match entries in restore_databases."
  }
}

variable "collection_renames" {
  description = "Optional map from source collection namespace (database.collection) to target namespace (database.collection). Keys must match entries in restore_collections."
  type        = map(string)
  default     = {}

  validation {
    condition = alltrue([
      for key in keys(var.collection_renames) : contains(var.restore_collections, key)
    ])
    error_message = "collection_renames keys must match entries in restore_collections."
  }
}

variable "database_suffix" {
  description = "Optional suffix applied to all restored database names."
  type        = string
  default     = null
}

variable "collection_suffix" {
  description = "Optional suffix applied to all restored collection names."
  type        = string
  default     = null
}
