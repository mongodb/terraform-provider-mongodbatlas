variable "public_key" {
  description = "Public API key to authenticate to Atlas"
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
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
}
variable "org_id" {
  description = "MongoDB Organization ID"
}
variable "region" {
  description = "MongoDB Atlas Cluster Region"
}
