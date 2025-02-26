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
  description = "Unique 24-hexadecimal digit string that identifies your project. Use the `/groups` at https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Projects/operation/listProjects endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups"
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

variable "instance_size" {
  description = "Hardware specification for the instance sizes in this region in this shard. Each instance size has a default storage and memory capacity. Electable nodes and read-only nodes (known as \"base nodes\") within a single shard must use the same instance size. Analytics nodes can scale independently from base nodes within a shard. Both base nodes and analytics nodes can scale independently from their equivalents in other shards."
  type        = string
}

variable "mongo_db_major_version" {
  description = "MongoDB major version of the cluster.  On creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Projects/operation/getProjectLtsVersions.   On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version."
  type        = string
}

variable "provider_name" {
  description = "Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`."
  type        = string
}

# OPTIONAL VARIABLES

variable "disk_size" {
  description = "Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.   This value must be equal for all shards and node types.   This value is not configurable on M0/M2/M5 clusters.   MongoDB Cloud requires this parameter if you set **replicationSpecs**.   If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value.    Storage charge calculations depend on whether you choose the default value or a custom value.   The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier."
  type        = number
  default     = 0
}

variable "auto_scaling_disk_gb_enabled" {
  description = "Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling."
  type        = bool
  default     = false
}

variable "tags" {
  description = "Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster."
  type        = map(string)
  default     = {}
}

variable "replication_specs" {
  description = "List of replication specifications in mongodbatlas_cluster format"
  type = list(object({
    num_shards = number
    zone_name  = string
    regions_config = list(object({
      region_name     = string
      electable_nodes = number
      priority        = number
      read_only_nodes = optional(number, 0)
    }))
  }))
  default = [{
    num_shards = 1
    zone_name  = "Zone 1"
    regions_config = [{
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
    }]
  }]
}

variable "search_nodes_specs" {
  description = "List of settings that configure the search nodes for your cluster. This list is currently limited to defining a single element. Valid values for the \"instance_size\" can be found in the MongoDB Atlas API docs - https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Atlas-Search/operation/createAtlasSearchDeployment."
  type = set(object({
    instance_size = string
    node_count    = number
  }))
  validation {
    condition     = length(var.search_nodes_specs) < 2
    error_message = "You can only define a single element of the specs currently."
  }
}

variable "encryption_at_rest_provider" {
  type = string
}
