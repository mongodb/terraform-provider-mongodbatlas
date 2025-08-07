variable "org_id" {
  description = "The MongoDB Atlas organization ID"
  type        = string
}

variable "user_email" {
  description = "The email address of the user"
  type        = string
}

variable "user_id" {
  description = "The user ID"
  type        = string
}

variable "public_key" {
  description = "Atlas API public key"
  type        = string
}

variable "private_key" {
  description = "Atlas API private key"
  type        = string
}
