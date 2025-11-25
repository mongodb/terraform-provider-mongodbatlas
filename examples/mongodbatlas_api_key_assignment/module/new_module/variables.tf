variable "org_id" {
  description = "Atlas Organization ID"
  type        = string
}

variable "project_id" {
  description = "Atlas Project ID"
  type        = string
}

variable "role_names" {
  description = "List of project-level roles to assign to the API key."
  type        = list(string)
}
