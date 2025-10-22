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
variable "data_lake" {
  description = "The datalake name"
  type        = string
}
variable "org_id" {
  description = "MongoDB Organization ID"
  type        = string
}
variable "region" {
  description = "MongoDB Atlas Cluster Region"
  type        = string
}
