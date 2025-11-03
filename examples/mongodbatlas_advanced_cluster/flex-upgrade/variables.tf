variable "atlas_org_id" {
  description = "Atlas organization id"
  type        = string
}
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
variable "provider_name" {
  description = "Atlas cluster provider name"
  default     = "AWS"
  type        = string
}
variable "backing_provider_name" {
  description = "Atlas cluster backing provider name"
  type        = string
}
variable "provider_instance_size_name" {
  description = "Atlas cluster provider instance name"
  default     = "M10"
  type        = string
}

variable "node_count" {
  description = "Number of nodes in the cluster"
  default     = 3
  type        = number
}