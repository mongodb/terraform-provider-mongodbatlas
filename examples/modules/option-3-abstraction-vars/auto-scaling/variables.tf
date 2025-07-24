variable "region_configs" {
  description = "List of region configurations for a replica set cluster. Each object defines a region."
  type = list(object({
    provider_name        = string
    region_name          = string
    priority             = optional(number, 7) # required if you have more than one region
    ebs_volume_type      = optional(string)
    disk_iops            = optional(number)
    electable_node_count = number
    read_only_node_count = optional(number, 0)
    analytics_specs = optional(object({
      ebs_volume_type = optional(string)
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
      priority             = optional(number, 7) # required if you have more than one region
      ebs_volume_type      = optional(string)
      disk_iops            = optional(number)
      electable_node_count = number
      read_only_node_count = optional(number, 0)
      analytics_specs = optional(object({
        ebs_volume_type = optional(string)
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
    compute_scale_down_enabled = optional(bool)
    compute_max_instance_size  = string
    compute_min_instance_size  = string
  })
}

variable "analytics_auto_scaling" {
  description = "Configuration for analytics auto-scaling."
  type = object({
    compute_scale_down_enabled = optional(bool)
    compute_max_instance_size  = string
    compute_min_instance_size  = string
  })
}

variable "tags_recommended" {
  description = "Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster."
  type        = map(string)
  validation {
    condition = (
      contains(keys(var.tags_recommended), "department") &&
      contains(keys(var.tags_recommended), "team_name") &&
      contains(keys(var.tags_recommended), "application_name") &&
      contains(keys(var.tags_recommended), "environment") &&
      contains(keys(var.tags_recommended), "version") &&
      contains(keys(var.tags_recommended), "email_contact") &&
      contains(keys(var.tags_recommended), "criticality")
    )
    error_message = "You must provide the following tags with non-empty values: department, team_name, application_name, environment, version, email_contact, criticality."
  }
}
