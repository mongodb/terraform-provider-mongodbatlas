variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}
variable "org_id" {
  description = "MongoDB Atlas Organization ID"
  type        = string
}
variable "dev_project_name" {
  description = "Name of the development MongoDB Atlas project"
  type        = string
}
variable "prod_project_name" {
  description = "Name of the production MongoDB Atlas project"
  type        = string
}
