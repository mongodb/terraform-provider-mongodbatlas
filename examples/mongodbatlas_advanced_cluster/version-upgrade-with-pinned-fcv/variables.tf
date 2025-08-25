variable "project_id" {
  description = "Atlas project id"
  type        = string
}
variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
}

variable "fcv_expiration_date" {
  description = "Expiration date of the pinned FCV"
  type        = string
}
