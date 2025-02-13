variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
  default     = ""
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
  default     = ""
}

# v1 & v2 variables
variable "project_id" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "cluster_type" {
  type = string
}

variable "mongo_db_major_version" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}

# NEW variable
variable "replication_specs_new" {
  type = list(object({
    num_shards = optional(number, 1)
    zone_name  = optional(string, "Zone 1")

    region_configs = list(object({
      region_name   = string
      provider_name = string
      priority      = optional(number, 7)

      auto_scaling = optional(object({
        disk_gb_enabled = optional(bool, false)
      }), null)
      read_only_specs = optional(object({
        node_count      = number
        instance_size   = string
        disk_size_gb    = optional(number, null)
        ebs_volume_type = optional(string, null)
        disk_iops       = optional(number, null)
      }), null)
      analytics_specs = optional(object({
        node_count      = number
        instance_size   = string
        disk_size_gb    = optional(number, null)
        ebs_volume_type = optional(string, null)
        disk_iops       = optional(number, null)
      }), null)
      electable_specs = object({
        node_count      = number
        instance_size   = string
        disk_size_gb    = optional(number, null)
        ebs_volume_type = optional(string, null)
        disk_iops       = optional(number, null)
      })
    }))
  }))
  description = "List of replication specifications for different regions"
  default     = []
}
