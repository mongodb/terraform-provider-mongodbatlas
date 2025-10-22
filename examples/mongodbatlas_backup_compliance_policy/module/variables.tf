variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}

variable "project_id" {
  type        = string
  description = "Atlas Project ID"
}

variable "user_id" {
  description = "Atlas User ID"
  type        = string
}

variable "cluster_name" {
  type    = string
  default = "check-bcp-module-deletion"
}

variable "instance_size" {
  type    = string
  default = "M10"
}
