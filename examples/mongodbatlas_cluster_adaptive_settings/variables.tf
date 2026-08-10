variable "project_id" {
  description = "Unique identifier of the MongoDB Atlas project."
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label for the MongoDB Atlas cluster."
  type        = string
}

variable "overload_protection_enabled" {
  description = "Whether to enable the overload protection override."
  type        = bool
  default     = true
}

variable "search_overload_protection_enabled" {
  description = "Whether to enable the Search overload protection override."
  type        = bool
  default     = true
}
