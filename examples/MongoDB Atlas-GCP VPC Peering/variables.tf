variable "public_key" {
  description = "Public API key to authenticate to Atlas"
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
}
variable "mongodbversion" {
  description = "The Major MongoDB Version"
  default     = "4.2"
}
variable "project_id" {
  description = "The Atlas Project Name"
}
variable "gcpprojectid" {
  default = "terraform-gcp-atlas"
}
variable "gcp_cidr" {
  default = "10.128.0.0/20"
}
variable "gcp_region" {
  description = "The GCP Region to use for deployment"
}
variable "atlas_region" {
  description = "The MongoDB Atlas region"
}