variable "org_id" {
  description = "The ID of the MongoDB Atlas organization"
  type        = string
}

variable "team_id_1" {
  description = "The ID of the first team"
  type        = string
}

variable "team_1_roles" {
  description = "Roles to assign to the first team in the project"
  type        = list(string)
}

variable "team_id_2" {
  description = "The ID of the second team"
  type        = string
}

variable "team_2_roles" {
  description = "Roles to assign to the second team in the project"
  type        = list(string)
}

variable "public_key" {
  description = "Public key for MongoDB Atlas API"
  type    = string
  default = ""
}
variable "private_key" {
  description = "Private key for MongoDB Atlas API"
  type    = string
  default = ""
}
