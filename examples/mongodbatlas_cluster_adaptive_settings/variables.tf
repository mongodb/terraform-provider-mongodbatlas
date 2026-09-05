variable "project_id" {
  description = "Unique identifier of the MongoDB Atlas project."
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label for the MongoDB Atlas cluster."
  type        = string
}

variable "load_shedding_enabled" {
  description = "Whether to enable the Load Shedding override."
  type        = bool
  default     = true
}

variable "search_load_shedding_enabled" {
  description = "Whether to enable the Search Load Shedding override."
  type        = bool
  default     = true
}
