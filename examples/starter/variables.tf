variable "public_key" {
  description = "Public API key to authenticate to Atlas"
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
}
variable "user" {
  description = "MongoDB Atlas User"
}
variable "password" {
  description = "MongoDB Atlas User Password"
}
variable "database_name" {
  description = "The Database in the cluster"
}
variable "org_id" {
  description = "MongoDB Organization ID"
}
variable "region" {
  description = "MongoDB Atlas Cluster Region"
}
variable "mongodbversion" {
  description = "The Major MongoDB Version"
}
variable "project_name" {
  description = "The Atlas Project Name"
}
