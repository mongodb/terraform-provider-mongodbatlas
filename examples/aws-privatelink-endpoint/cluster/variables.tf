variable "public_key" {
  description = "The public API key for MongoDB Atlas"
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
}
variable "access_key" {
  description = "The access key for AWS Account"
}
variable "secret_key" {
  description = "The secret key for AWS Account"
}
variable "project_id" {
  description = "Atlas project ID"
}
variable "cluster_name" {
  description = "Atlas cluster name"
  default     = "aws-private-connection"
}