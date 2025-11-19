variable "atlas_org_id" {
  description = "Atlas organization id"
  type        = string
}

variable "project_name" {
  description = "Atlas project name"
  type        = string
}

variable "cluster_name" {
  description = "Atlas cluster name"
  type        = string
}

variable "cluster_type" {
  description = "Atlas cluster type"
  type        = string
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

variable "tags" {
  description = "Map of tags to assign to the cluster"
  type        = map(string)
  default     = {}
}
