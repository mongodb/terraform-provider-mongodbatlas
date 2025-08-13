variable "project_id" {
  description = "The MongoDB Atlas project ID"
  type        = string
}

variable "team_id" {
  description = "The MongoDB Atlas team ID"
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
