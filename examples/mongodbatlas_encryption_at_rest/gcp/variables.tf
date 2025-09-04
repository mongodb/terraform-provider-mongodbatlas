variable "atlas_public_key" {
  description = "MongoDB Atlas public API key for authentication"
  type        = string
}

variable "atlas_private_key" {
  description = "MongoDB Atlas private API key for authentication"
  type        = string
  sensitive   = true
}

variable "atlas_project_id" {
  description = "MongoDB Atlas project ID where the cloud provider access will be configured"
  type        = string
}

variable "gcp_project_id" {
  description = "Google Cloud Platform project ID where KMS resources will be created"
  type        = string
}

variable "key_ring_name" {
  description = "Name of the Google Cloud KMS key ring to create"
  type        = string
  default     = "atlas-key-ring"
}

variable "crypto_key_name" {
  description = "Name of the Google Cloud KMS crypto key to create for MongoDB Atlas encryption"
  type        = string
  default     = "atlas-crypto-key"
}

variable "location" {
  description = "Google Cloud region where the KMS key ring will be created"
  type        = string
  default     = "us-central1"
}
