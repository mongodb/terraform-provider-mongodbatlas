variable "atlas_org_id" {
  description = "Atlas organization id"
  type        = string
}
variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
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