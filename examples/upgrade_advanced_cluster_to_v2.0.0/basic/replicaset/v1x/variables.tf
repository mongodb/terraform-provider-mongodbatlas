variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project."
  type        = string
}

variable "cluster_name" {
  description = "Human-readable label that identifies this cluster."
  type        = string
}

variable "instance_size" {
  description = "Instance size for nodes."
  type        = string
  default     = "M10"
}

variable "disk_size_gb" {
  description = "Disk size (GB) for nodes."
  type        = number
  default     = 60
}

variable "node_count_electable" {
  description = "Electable node count."
  type        = number
  default     = 3
}
