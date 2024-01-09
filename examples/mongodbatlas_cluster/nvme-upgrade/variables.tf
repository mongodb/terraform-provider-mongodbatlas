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

variable "provider_instance_size_name" {
  description = "Atlas cluster provider instance name"
  default     = "M40"
  type        = string
}

variable "provider_volume_type" {
  description = "Atlas cluster provider storage volume name"
  default     = "STANDARD"
  type        = string
}

variable "provider_disk_iops" {
  description = "Atlas cluster provider disk iops"
  default     = 100
  type        = number
}
