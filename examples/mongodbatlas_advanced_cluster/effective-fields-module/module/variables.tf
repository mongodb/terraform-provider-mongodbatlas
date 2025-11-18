variable "atlas_org_id" {
  type        = string
  description = "The ID of your MongoDB Atlas organization"
}

variable "project_name" {
  type        = string
  description = "The name of the Atlas project to create"
}

variable "cluster_name" {
  type        = string
  description = "The name of the cluster"
}

variable "cluster_type" {
  type        = string
  description = "The type of cluster (REPLICASET, SHARDED, or GEOSHARDED)"
  default     = "REPLICASET"

  validation {
    condition     = contains(["REPLICASET", "SHARDED", "GEOSHARDED"], var.cluster_type)
    error_message = "cluster_type must be one of: REPLICASET, SHARDED, GEOSHARDED"
  }
}

variable "replication_specs" {
  type = list(object({
    zone_name = optional(string)
    region_configs = list(object({
      priority      = number
      provider_name = string
      region_name   = string
      electable_specs = object({
        instance_size = string
        node_count    = number
        disk_iops     = optional(number)
        disk_size_gb  = optional(number)
      })
      analytics_specs = optional(object({
        instance_size = string
        node_count    = number
        disk_iops     = optional(number)
        disk_size_gb  = optional(number)
      }))
      read_only_specs = optional(object({
        instance_size = string
        node_count    = number
        disk_iops     = optional(number)
        disk_size_gb  = optional(number)
      }))
      auto_scaling = optional(object({
        disk_gb_enabled            = optional(bool, false)
        compute_enabled            = optional(bool, false)
        compute_scale_down_enabled = optional(bool, false)
        compute_min_instance_size  = optional(string)
        compute_max_instance_size  = optional(string)
      }))
      analytics_auto_scaling = optional(object({
        disk_gb_enabled            = optional(bool, false)
        compute_enabled            = optional(bool, false)
        compute_scale_down_enabled = optional(bool, false)
        compute_min_instance_size  = optional(string)
        compute_max_instance_size  = optional(string)
      }))
    }))
  }))
  description = "List of replication specifications for the cluster"
}

variable "enable_auto_scaling" {
  type        = bool
  description = "Enable auto-scaling for electable nodes. When true, auto_scaling configuration from replication_specs will be used"
  default     = false
}

variable "enable_analytics_auto_scaling" {
  type        = bool
  description = "Enable auto-scaling for analytics nodes. When true, analytics_auto_scaling configuration from replication_specs will be used"
  default     = false
}

variable "tags" {
  type        = map(string)
  description = "Map of tags to apply to the cluster"
  default     = {}
}
