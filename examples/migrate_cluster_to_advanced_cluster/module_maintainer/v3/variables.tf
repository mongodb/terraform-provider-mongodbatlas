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

# OPTIONAL VARIABLES
variable "instance_size" {
  type    = string
  default = "" # optional in v3
}
variable "provider_name" {
  type    = string
  default = "" # optional in v3
}

variable "disk_size" {
  type    = number
  default = 0
}

variable "auto_scaling_disk_gb_enabled" {
  type    = bool
  default = false
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "replication_specs" {
  description = "List of replication specifications in legacy mongodbatlas_cluster format"
  default     = []
  type = list(object({
    num_shards = number
    zone_name  = string
    regions_config = set(object({
      region_name     = string
      electable_nodes = number
      priority        = number
      read_only_nodes = optional(number, 0)
    }))
  }))
}

variable "replication_specs_new" {
  description = "List of replication specifications using new mongodbatlas_advanced_cluster format"
  default     = []
  type = list(object({
    zone_name = optional(string, "Zone 1")

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
}
