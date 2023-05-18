variable "public_key" {
  type        = string
  description = "Public Programmatic API key to authenticate to Atlas"
}
variable "private_key" {
  type        = string
  description = "Private Programmatic API key to authenticate to Atlas"
}
variable "atlas_org_id" {
  type        = string
  description = "MongoDB Organization ID"
}

variable "atlas_project_name" {
  type        = string
  description = "MongoDB Project Name"
}

variable "name" {
  type        = string
  description = "MongoDB DataLake Pipeline Name"
  default     = "datalakePipelineName"
}
