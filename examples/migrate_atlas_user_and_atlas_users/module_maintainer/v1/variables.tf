variable "username" {
  type        = string
  description = "MongoDB Atlas Username (email)"
}

variable "org_id" {
  type        = string
  description = "MongoDB Atlas Organization ID"
}

variable "project_id" {
  type        = string
  description = "MongoDB Atlas Project ID (optional)"
  default     = null
}
