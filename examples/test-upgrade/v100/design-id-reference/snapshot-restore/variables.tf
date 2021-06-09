variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  default     = ""
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  default     = ""
}
variable "project_name" {
  description = "Atlas project name"
  default     = ""
}
variable "cluster_name" {
  description = "Cluster name"
  default     = ""
}
variable "description" {
  description = "Description"
  default     = ""
}
variable "retention_in_days" {
  description = "Retention in days"
  default     = ""
}