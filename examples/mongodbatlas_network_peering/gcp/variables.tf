variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}
variable "project_id" {
  description = "The Atlas Project Name"
  type        = string
}
variable "gcpprojectid" {
  default = "terraform-gcp-atlas"
  type    = string
}
variable "gcp_cidr" {
  default = "10.128.0.0/20"
  type    = string
}
variable "gcp_region" {
  description = "The GCP Region to use for deployment"
  type        = string
}
variable "atlas_region" {
  description = "The MongoDB Atlas region"
  type        = string
}
