variable "public_key" {
  description = "MongoDB Atlas API public key."
  type        = string
}

variable "private_key" {
  description = "MongoDB Atlas API private key."
  type        = string
}

variable "org_id" {
  description = "The ID of the MongoDB Atlas organization."
  type        = string
}

variable "user_email" {
  description = "The email address of the user to assign."
  type        = string
} 