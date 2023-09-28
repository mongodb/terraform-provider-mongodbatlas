variable "atlas_org_id" {
  description = "Atlas organization id"
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

# project
variable "project_name" {
  description = "Atlas project name"
  default     = "TenantUpgradeTest"
  type        = string
}

#cluster
variable "cluster_name" {
  description = "Atlas cluster name"
  default     = "cluster"
  type        = string
}

variable "cluster_type" {
  description = "Atlas cluster type"
  default     = "REPLICASET"
  type        = string
}

variable "num_shards" {
    description = "Atlas cluster number of shards"
    default     = 1
    type        = number
}

variable "priority" {
    description = "Atlas cluster priority"
    default     = 7
    type        = number
}

variable "read_only_nodes" {
    description = "Atlas cluster number of read only nodes"
    default     = 0
    type        = number
}
variable "electable_nodes" {
    description = "Atlas cluster number of electable nodes"
    default     = 3
    type        = number
}

variable "auto_scaling_disk_gb_enabled" {
    description = "Atlas cluster auto scaling disk enabled"
    default     = false
    type        = bool
}

variable "disk_size_gb" {
    description = "Atlas cluster disk size in GB"
    default     = 10
    type        = number
}
variable "provider_name" {
  description = "Atlas cluster provider name"
  default     = "AWS"
  type        = string
}
variable "backing_provider_name" {
  description = "Atlas cluster backing provider name"
  default = "AWS"
  type        = string
}
variable "provider_instance_size_name" {
  description = "Atlas cluster provider instance name"
  default     = "M10"
  type        = string
}

variable "region" {
    description = "Atlas cluster region"
    default     = "US_EAST_1"
    type        = string
}
variable "aws_region"{
    description = "AWS region"
    default     = "us-east-1"
    type        = string
}

variable "mongo_version" {
  description = "Atlas cluster version"
  default     = "4.4"
  type        = string
}


variable "user" {
  description = "MongoDB Atlas User"
  type        = list(string)
  default     = ["dbuser1", "dbuser2"]
}
variable "db_passwords" {
  description = "MongoDB Atlas User Password"
  type        = list(string)
}
variable "database_names" {
  description = "The Database in the cluster"
  type        = list(string)
}

# database user
variable "role_name" {
    description = "Atlas database user role name"
    default     = "readWrite"
    type        = string
}

# IP Access List
variable "cidr_block" {
    description = "IP Access List CIDRs"
    type        = list(string)
}

variable "ip_address" {
    description = "IP Access List IP Addresses"
    type        = list(string)
}
# aws

variable "aws_vpc_cidr_block" {
    description = "AWS VPC CIDR block"
    default = "10.0.0.0/16"
    type        = string
}

# aws vpc
variable "aws_vpc_ingress" {
  description = "AWS VPC ingress CIDR block"
  type        = string
}

variable "aws_vpc_egress" {
  description = "AWS VPC egress CIDR block"
  type        = string
}

variable "aws_route_table_cidr_block" {
    description = "AWS route table CIDR block"
    default     = "0.0.0.0/0"
    type        = string
}

variable "aws_subnet_cidr_block1" {
    description = "AWS subnet CIDR block"
    type        = string
}
variable "aws_subnet_cidr_block2" {
  description = "AWS subnet CIDR block"
  type        = string
}

variable "aws_subnet_availability_zone1" {
    description = "AWS subnet availability zone"
    default     = "us-east-1a"
    type        = string
}
variable "aws_subnet_availability_zone2" {
  description = "AWS subnet availability zone"
  default     = "us-east-1b"
  type        = string
}

variable "aws_sg_ingress_from_port" {
    description = "AWS security group ingress from port"
    default     = 27017
    type        = number
}

variable "aws_sg_ingress_to_port" {
    description = "AWS security group ingress to port"
    default     = 27017
    type        = number
}

variable "aws_sg_ingress_protocol" {
    description = "AWS security group ingress protocol"
    default     = "tcp"
    type        = string
}

variable "aws_sg_egress_from_port" {
    description = "AWS security group egress from port"
    default     = 0
    type        = number
}

variable "aws_sg_egress_to_port" {
    description = "AWS security group egress to port"
    default     = 0
    type        = number
}

variable "aws_sg_egress_protocol" {
    description = "AWS security group egress protocol"
    default     = "-1"
    type        = string
}

variable "db_users" {
    description = "Atlas database users"
    type        = list(string)
}