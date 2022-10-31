variable "atlas_org_id" {
  description = "Atlas organization id"
  default     = ""
}
variable "public_key" {
  description = "Public API key to authenticate to Atlas"
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
}
variable "provider_name" {
  description = "Atlas cluster provider name"
  default     = "AWS"
}
variable "backing_provider_name" {
  description = "Atlas cluster backing provider name"
  default     = null
}
variable "provider_instance_size_name" {
  description = "Atlas cluster provider instance name"
  default     = "M10"
}