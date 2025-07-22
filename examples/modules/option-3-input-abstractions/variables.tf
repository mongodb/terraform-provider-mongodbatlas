variable "region_configs" {
  description = "List of region configurations for a replica set cluster. Each object defines a region."
  type = list(object({
    provider_name        = string
    region_name          = string
    instance_size        = string
    priority             = optional(number, 7) # required if you have more than one region
    ebs_volume_type      = optional(string)
    disk_size_gb         = optional(number)
    disk_iops            = optional(number)
    electable_node_count = number
    read_only_node_count = optional(number, 0)
    analytics_specs = optional(object({
      instance_size   = string
      ebs_volume_type = optional(string)
      disk_size_gb    = optional(number)
      disk_iops       = optional(number)
      node_count      = number
    }))
  }))
  default = []
}

variable "shards" {
  description = "List of shard configurations for a sharded cluster. Each object defines a shard and its regions."
  type = list(object({
    zone_name = optional(string)
    region_configs = list(object({
      provider_name        = string
      region_name          = string
      instance_size        = string
      priority             = optional(number, 7) # required if you have more than one region
      ebs_volume_type      = optional(string)
      disk_size_gb         = optional(number)
      disk_iops            = optional(number)
      electable_node_count = number
      read_only_node_count = optional(number, 0)
      analytics_specs = optional(object({
        instance_size   = string
        ebs_volume_type = optional(string)
        disk_size_gb    = optional(number)
        disk_iops       = optional(number)
        node_count      = number
      }))
    }))
  }))
  default = []
}

variable "auto_scaling" {
  description = "Configuration for auto-scaling."
  type = object({
    disk_gb_enabled            = bool
    compute_enabled            = bool
    compute_scale_down_enabled = optional(bool)
    compute_max_instance_size  = optional(string)
    compute_min_instance_size  = optional(string)
  })
  default = {
    disk_gb_enabled = false # defaults to false, default to true would imply being opinionated on a max_instance_size which can vary significantly
    compute_enabled = false
  }
}

variable "analytics_auto_scaling" {
  description = "Configuration for analytics auto-scaling."
  type = object({
    disk_gb_enabled            = bool
    compute_enabled            = bool
    compute_scale_down_enabled = optional(bool)
    compute_max_instance_size  = optional(string)
    compute_min_instance_size  = optional(string)
  })
  default = {
    disk_gb_enabled = false # defaults to false, default to true would imply being opinionated on a max_instance_size which can vary significantly
    compute_enabled = false
  }
}
