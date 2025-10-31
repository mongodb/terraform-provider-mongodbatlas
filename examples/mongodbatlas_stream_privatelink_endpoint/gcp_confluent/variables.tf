variable "project_id" {
  description = "The MongoDB Atlas project ID"
  type        = string
}

variable "gcp_region" {
  description = "The GCP region where resources will be created"
  type        = string
  default     = "us-west1"
}

variable "confluent_dns_domain" {
  description = "The DNS domain for the Confluent cluster"
  type        = string
  default     = "example.confluent.cloud"
}

variable "confluent_dns_subdomains" {
  description = "List of DNS subdomains for the Confluent cluster"
  type        = list(string)
  default     = ["subdomain1", "subdomain2"]
}
