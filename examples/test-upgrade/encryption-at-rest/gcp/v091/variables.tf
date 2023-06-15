# gcp
variable "service_account_key" {
  type = string
}
variable "gcp_key_version_resource_id" {
  type = string
}

# encryption at rest
variable "project_name" {
  description = "Atlas project name"
  type        = string
}
variable "org_id" {
  description = "The organization ID"
  type        = string
}
