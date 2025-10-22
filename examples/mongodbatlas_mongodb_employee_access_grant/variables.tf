variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}

variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}

variable "project_id" {
  description = "Unique 24-hexadecimal digit string that identifies your project"
  type        = string
}

variable "cluster_name" {
  description = "Name of an existing cluster in your project that you want to grant access to"
  type        = string
}
