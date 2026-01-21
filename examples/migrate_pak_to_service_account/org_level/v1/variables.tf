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

variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "org_roles" {
  description = "Organization roles for the API Key"
  type        = list(string)
  default     = ["ORG_MEMBER"]
}

variable "project_roles" {
  description = "Project roles for the API Key assignment"
  type        = list(string)
  default     = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

variable "cidr_block" {
  description = "CIDR block for IP access list entry"
  type        = string
  default     = "192.168.1.100/32"
}
