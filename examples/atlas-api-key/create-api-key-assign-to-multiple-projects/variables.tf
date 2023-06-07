variable "public_key" {
  description = "Public API key to authenticate to Atlas"
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
}
variable "project_id" {
  description = "Atlas project id"
  default     = ""
}
variable "additional_project_id" {
  description = "Atlas project id"
  default     = ""
}
