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
  description = "Unique 24-hexadecimal digit string that identifies your project. Use the `/groups` at https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-listprojects endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups"
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label that identifies this cluster."
  type        = string
}

variable "cluster_type" {
  description = "Configuration of nodes that comprise the cluster."
  type        = string
}

variable "mongo_db_major_version" {
  description = "MongoDB major version of the cluster.  On creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Projects/operation/getProjectLtsVersions.   On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version."
  type        = string
}

variable "tags" {
  description = "Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster."
  type        = map(string)
  default     = {}
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
  description = "List of replication specifications using new mongodbatlas_advanced_cluster format"
  default     = []
}
