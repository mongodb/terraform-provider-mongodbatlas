# Atlas Organization ID 
variable "atlas_org_id" {
  type        = string
  description = "Atlas organization id"
}
# Atlas Project Name
variable "atlas_project_name" {
  type        = string
  description = "Atlas project name"
}

# Atlas Project Environment
variable "environment" {
  type        = string
  description = "The environment to be built"
}

# Cluster Instance Size Name 
variable "cluster_instance_size_name" {
  type        = string
  description = "Cluster instance size name"
}

# Cloud Provider to Host Atlas Cluster
variable "cloud_provider" {
  type        = string
  description = "AWS or GCP or Azure"
}

# Atlas Region
variable "atlas_region" {
  type        = string
  description = "AWS region where resources will be created"
}

# MongoDB Version 
variable "mongodb_version" {
  type        = string
  description = "MongoDB Version"
}

# IP Adress Access
variable "ip_address" {
  type = string
  description = "IP address that will be used to access Atlas cluster"
}

# AWS Region
variable "aws_region" {
  type = string
  description = "AWS Region"
}
