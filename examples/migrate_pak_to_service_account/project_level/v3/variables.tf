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
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "project_service_account_name" {
  description = "Name for the Project Service Account"
  type        = string
  default     = "example-project-service-account"
}

variable "project_roles" {
  description = "Project roles for the Service Account"
  type        = list(string)
  default     = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

variable "cidr_block" {
  description = "CIDR block for IP access list entry"
  type        = string
  default     = "192.168.1.100/32"
}

variable "secret_expires_after_hours" {
  description = "Number of hours after which the Service Account secret expires"
  type        = number
  default     = 2160 # 90 days
}
