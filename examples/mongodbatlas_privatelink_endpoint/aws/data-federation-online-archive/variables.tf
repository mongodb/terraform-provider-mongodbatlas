variable "public_key" {
  description = "The public API key for MongoDB Atlas"
  type        = string
}
variable "private_key" {
  description = "The private API key for MongoDB Atlas"
  type        = string
}
variable "access_key" {
  description = "The access key for AWS Account"
  type        = string
}
variable "secret_key" {
  description = "The secret key for AWS Account"
  type        = string
}
variable "project_id" {
  description = "Atlas project ID"
  type        = string
}

variable "service_name" {
  description = "AWS VPC endpoint service name associated to the specific region. Values can be found in https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createdatafederationprivateendpoint."
  type        = string
}
