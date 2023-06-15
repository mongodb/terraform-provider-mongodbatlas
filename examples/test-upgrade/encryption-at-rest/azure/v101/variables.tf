# azure
variable "client_id" {
  type = string
}
variable "subscription_id" {
  type = string
}
variable "resource_group_name" {
  type = string
}
variable "client_secret" {
  type = string
}
variable "tenant_id" {
  type = string
}
variable "key_vault_name" {
  type = string
}
variable "key_identifier" {
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
