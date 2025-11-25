variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
  default     = ""
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
  sensitive   = true
  default     = ""
}
variable "org_id" {
  type        = string
  description = "MongoDB Organization ID"
}
variable "project_name" {
  type        = string
  description = "The MongoDB Atlas Project Name"
}
variable "cluster_name" {
  type        = string
  description = "The MongoDB Atlas Cluster Name"
}
variable "cloud_provider" {
  type        = string
  description = "The cloud provider to use, must be AWS, GCP or AZURE"
}
variable "region" {
  type        = string
  description = "MongoDB Atlas Cluster Region, must be a region for the provider given"
}
variable "dbuser" {
  type        = string
  description = "MongoDB Atlas Database User Name"
}
variable "dbuser_password" {
  type        = string
  description = "MongoDB Atlas Database User Password"
}
variable "database_name" {
  type        = string
  description = "The database in the cluster to limit the database user to, the database does not have to exist yet"
}
variable "ip_address" {
  type        = string
  description = "The IP address that the cluster will be accessed from, can also be a CIDR range or AWS security group"
}



