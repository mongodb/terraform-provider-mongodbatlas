variable "public_key" {
  description = "Public Programmatic API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private Programmatic API key to authenticate to Atlas"
  type        = string
}
variable "org_id" {
  type        = string
  description = "MongoDB Organization ID"
}
variable "project_name" {
  type        = string
  description = "The MongoDB Atlas Project Name"
}





