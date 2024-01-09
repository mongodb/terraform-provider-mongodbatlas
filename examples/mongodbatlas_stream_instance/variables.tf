variable "public_key" {
  description = "Public API key to authenticate to Atlas"
  type        = string
}
variable "private_key" {
  description = "Private API key to authenticate to Atlas"
  type        = string
}
variable "org_id" {
  description = "Unique 24-hexadecimal digit string that identifies your Atlas Organization"
  type        = string
} 