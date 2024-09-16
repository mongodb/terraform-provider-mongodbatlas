variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}

variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
}

variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}

variable "cluster_name" {
  description = "Name of an existing cluster in your project that you want to grant access to"
  type        = string
}
