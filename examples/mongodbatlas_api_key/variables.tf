variable "atlas_client_id" {
  type        = string
  description = "MongoDB Atlas Service Account Client ID"
}

variable "atlas_client_secret" {
  type        = string
  description = "MongoDB Atlas Service Account Client Secret"
}

variable "org_id" {
  type        = string
  description = "The ID of the organization to create the API key in."
}
