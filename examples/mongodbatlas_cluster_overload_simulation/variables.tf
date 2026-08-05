variable "project_id" {
  description = "Unique identifier of the MongoDB Atlas project."
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label for the MongoDB Atlas cluster."
  type        = string
  default     = "cluster-overload-simulation"
}

variable "cloud_provider" {
  description = "Cloud provider for the MongoDB Atlas cluster."
  type        = string
  default     = "AWS"
}

variable "region_name" {
  description = "Cloud provider region for the MongoDB Atlas cluster."
  type        = string
  default     = "US_EAST_1"
}

variable "instance_size" {
  description = "Instance size for the MongoDB Atlas cluster."
  type        = string
  default     = "M10"
}

variable "duration_seconds" {
  description = "Duration of the overload protection simulation in seconds."
  type        = number
  default     = 900

  validation {
    condition     = contains([900, 3600, 28800, 86400], var.duration_seconds)
    error_message = "Duration must be one of 900, 3600, 28800, or 86400 seconds."
  }
}
