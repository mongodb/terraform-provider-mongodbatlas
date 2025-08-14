variable "project_id" {
  description = "The MongoDB Atlas project ID"
  type        = string
}

variable "user_email" {
  description = "The email address of the user"
  type        = string
}

variable "public_key" {
  description = "Atlas API public key"
  type        = string
  default     = ""
}

variable "private_key" {
  description = "Atlas API private key"
  type        = string
  default     = ""
}
