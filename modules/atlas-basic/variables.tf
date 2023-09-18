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
  default = "AWS"
  type        = string
}
variable "provider_instance_size_name" {
  description = "Atlas cluster provider instance name"
  default     = "M10"
  type        = string
}

variable "region" {
    description = "Atlas cluster region"
    default     = "US_EAST_1"
    type        = string
}

variable "version" {
  description = "Atlas cluster version"
  default     = "4.4"
  type        = string
}


variable "user" {
  description = "MongoDB Atlas User"
  type        = list(string)
  default     = ["dbuser1", "dbuser2"]
}
variable "password" {
  description = "MongoDB Atlas User Password"
  type        = list(string)
}
variable "database_name" {
  description = "The Database in the cluster"
  type        = list(string)
}