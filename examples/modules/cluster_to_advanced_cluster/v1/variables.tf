variable "project_id" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "cluster_type" {
  type = string

  validation {
    condition     = contains(["REPLICASET", "SHARDED", "GEOSHARDED"], var.cluster_type)
    error_message = "Valid supported cluster types are \"REPLICASET\", \"SHARDED\" or \"GEOSHARDED\"."
  }
}

variable "instance_size" {
  type = string
}

variable "mongo_db_major_version" {
  type = string
}

variable "provider_name" {
  type = string
}

# OPTIONAL VARIABLES

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
  type = list(object({
    num_shards = number
    zone_name  = string
    regions_config = set(object({
      region_name     = string
      electable_nodes = number
      priority        = number
      read_only_nodes = number
    }))
  }))
  default = [{
    num_shards = 1
    zone_name  = "Zone 1"
    regions_config = [{
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }]
  }]
}
