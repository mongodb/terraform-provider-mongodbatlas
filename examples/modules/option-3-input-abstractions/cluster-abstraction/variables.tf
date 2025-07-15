variable "single_region" {
  description = "Configuration for a single-region cluster. If set, the module will use this for cluster creation."
  type = object({
    provider_name  = string
    region_name    = string
    instance_size  = string
    ebs_volume_type = optional(string)
    disk_size_gb = optional(number)
    disk_iops = optional(number)
    node_count     = number
    read_only_node_count = optional(number, 0)
    analytics_specs = optional(object({
      instance_size   = string
      ebs_volume_type = optional(string)
      disk_size_gb    = optional(number)
      disk_iops       = optional(number)
      node_count      = number
    }))
  })
  default = null
}

variable "auto_scaling" {
  description = "Configuration for auto-scaling."
  type = object({
    disk_gb_enabled           = bool
    compute_enabled           = bool
    compute_max_instance_size = string
    compute_min_instance_size = string
  })
  default = {
    disk_gb_enabled           = true
    compute_enabled           = true
    compute_max_instance_size = "M60" // TODO do we want to keep this as the default?
    compute_min_instance_size = "M30"
  }
}

variable "project_id" {
  description = "The MongoDB Atlas project ID."
  type        = string
}

variable "name" {
  description = "The name of the cluster."
  type        = string
}

variable "cluster_type" {
  description = "Type of cluster (REPLICASET, SHARDED, GEOSHARDED)."
  type        = string
  default     = "REPLICASET"
}

variable "mongo_db_major_version" {
  description = "The major MongoDB version."
  type        = string
  default     = "8.0"
}