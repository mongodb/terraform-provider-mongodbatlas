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
  default = "check-bcp-schedule-deletion"
}

variable "instance_size" {
  type    = string
  default = "M10"
}
