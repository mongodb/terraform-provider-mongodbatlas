variable "atlas_public_key" {
  description = "MongoDB Atlas public API key for authentication"
  type        = string
}

variable "atlas_private_key" {
  description = "MongoDB Atlas private API key for authentication"
  type        = string
  sensitive   = true
}

variable "atlas_project_id" {
  description = "MongoDB Atlas project ID where the cloud provider access will be configured"
  type        = string
}
