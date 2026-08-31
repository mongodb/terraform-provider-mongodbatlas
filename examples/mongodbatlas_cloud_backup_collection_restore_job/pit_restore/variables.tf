variable "project_id" {
  description = "Atlas project ID of the backup source cluster."
  type        = string
}

variable "cluster_name" {
  description = "Name of the backup source cluster."
  type        = string
}

variable "target_project_id" {
  description = "Atlas project ID of the restore target cluster. When unset, restores to project_id."
  type        = string
  default     = null
}

variable "target_cluster_name" {
  description = "Name of the restore target cluster. When unset, restores to cluster_name."
  type        = string
  default     = null
}

variable "write_strategy" {
  description = "How Atlas handles existing data on the target. CREATE_NEW appends and renames on conflict; OVERWRITE_EXISTING drops and replaces matching namespaces."
  type        = string
  default     = "CREATE_NEW"
}

variable "index_strategy" {
  description = "Which indexes to restore. ALL restores all indexes; NONE skips indexes; ALL_EXCEPT_TTL restores non-TTL indexes only."
  type        = string
  default     = "ALL"
}

variable "restore_databases" {
  description = "Database names to restore (database name only, not database.collection)."
  type        = list(string)
  default     = []
}

variable "restore_collections" {
  description = "Collection namespaces to restore in database.collection form."
  type        = list(string)
  default     = []
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

variable "point_in_time_utc_seconds" {
  description = "Wall-clock UNIX epoch seconds to restore to from continuous backup. Set this or oplog_ts and oplog_inc."
  type        = number
  default     = null
}

variable "oplog_ts" {
  description = "Oplog Timestamp seconds component for continuous backup restore. Use with oplog_inc instead of point_in_time_utc_seconds."
  type        = number
  default     = null
}

variable "oplog_inc" {
  description = "Oplog Timestamp increment component for continuous backup restore. Use with oplog_ts."
  type        = number
  default     = null
}

variable "create_timeout" {
  description = "Maximum time to wait for the restore job to reach SUCCESSFUL during create."
  type        = string
  default     = "3h"
}
