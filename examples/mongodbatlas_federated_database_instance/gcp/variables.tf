variable "project_id" {
  description = "MongoDB Atlas Project ID"
  type        = string
}

variable "federated_instance_name" {
  description = "Name for the Federated Database Instance"
  type        = string
  default     = "gcp-federated-instance"
}
