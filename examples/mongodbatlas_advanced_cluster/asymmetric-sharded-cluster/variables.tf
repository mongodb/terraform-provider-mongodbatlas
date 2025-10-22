variable "atlas_org_id" {
  description = "Atlas organization id"
  type        = string
}
variable "atlas_client_id" {
  description = "MongoDB Atlas Service Account Client ID"
  type        = string
}
variable "atlas_client_secret" {
  description = "MongoDB Atlas Service Account Client Secret"
  type        = string
}
variable "cluster_name" {
  description = "Atlas cluster name"
  type        = string
  default     = "AsymmetricShardedCluster"
}
