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
      node_count      = number
      ebs_volume_type = optional(string)
      disk_size_gb    = optional(number)
      disk_iops       = optional(number)
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
      electable_node_count = number
      priority             = optional(number, 7) # required if you have more than one region
      ebs_volume_type      = optional(string)
      disk_size_gb         = optional(number)
      disk_iops            = optional(number)
      read_only_node_count = optional(number, 0)
      analytics_specs = optional(object({
        instance_size   = string
        node_count      = number
        ebs_volume_type = optional(string)
        disk_size_gb    = optional(number)
        disk_iops       = optional(number)
      }))
    }))
  }))
  default = []
}

# no inputs associated to auto-scaling are defined as module is specific to no auto-scaling