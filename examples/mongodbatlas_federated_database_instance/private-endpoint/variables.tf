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

variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "federated_instance_name" {
  description = "Name of the MongoDB Atlas Federated Database Instance"
  type        = string
  default     = "federated-instance-privatelink"
}

variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "atlas_region" {
  description = "Atlas region for the private endpoint in uppercase underscore format (e.g. US_EAST_1). Must match the AWS region. See https://www.mongodb.com/docs/atlas/data-federation/adf-overview/regions/ for the mapping between AWS and Atlas regions."
  type        = string
}

variable "vpce_service_name" {
  description = "AWS PrivateLink service name for Atlas Data Federation in your region. See https://www.mongodb.com/docs/atlas/data-federation/tutorial/config-private-endpoint/ for the value for your region."
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_cidr" {
  description = "CIDR block for the subnet"
  type        = string
  default     = "10.0.1.0/24"
}

variable "availability_zone" {
  description = "Availability zone for the subnet"
  type        = string
  default     = "us-east-1a"
}
