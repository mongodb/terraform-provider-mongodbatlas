variable "public_key" {
  type        = string
  description = "The public key of the API key."
}

variable "private_key" {
  type        = string
  description = "The private key of the API key."
}

variable "org_id" {
  type        = string
  description = "The ID of the organization to create the API key in."
}

variable "project_name" {
  type        = string
  description = "The name of the project to create the API key in."
}
