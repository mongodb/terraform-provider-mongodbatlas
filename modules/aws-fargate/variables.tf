variable availability_zones {
  description = "List of Availability Zones for VPC subnets. You can specify two Availability Zones."
  type = list(string)
}

variable private_subnet1_cidr {
  description = "CIDR block for private subnet, located in Availability Zone 1."
  type = string
  default = "10.0.0.0/19"
}

variable private_subnet2_cidr {
  description = "CIDR block for private subnet, located in Availability Zone 2."
  type = string
  default = "10.0.32.0/19"
}

variable public_subnet1_cidr {
  description = "CIDR block for public (DMZ) subnet 1, located in Availability Zone 1."
  type = string
  default = "10.0.128.0/20"
}

variable public_subnet2_cidr {
  description = "CIDR block for public (DMZ) subnet 2, located in Availability Zone 2."
  type = string
  default = "10.0.144.0/20"
}

variable web_access_cidr {
  description = "CIDR block to allow access to web application."
  type = string
}

variable vpccidr {
  description = "CIDR block for VPC."
  type = string
  default = "10.0.0.0/16"
}

variable profile {
  description = "Secret with name cfn/atlas/profile/{Profile}."
  type = string
  default = "default"
}

variable org_id {
  description = "MongoDB Cloud Organization ID."
  type = string
}

variable project_name {
  description = "Name of project."
  type = string
  default = "aws-quickstart"
}

variable cluster_name {
  description = "Name of cluster as it appears in Atlas. Once created, the name can't be changed."
  type = string
  default = "Cluster-1"
}

variable database_user_name {
  description = "MongoDB Atlas database user name."
  type = string
  default = "testUser"
}

variable database_password {
  description = "MongoDB Atlas database user password."
  type = string
}

variable cluster_instance_size {
  description = "Atlas provides different cluster tiers, each with a default storage capacity and RAM size. The cluster you select is used for all the data-bearing hosts in your cluster tier. See https://docs.atlas.mongodb.com/reference/amazon-aws/#amazon-aws."
  type = string
  default = "M10"
}

variable cluster_region {
  description = "AWS Region where Atlas DB Cluster will run."
  type = string
  default = "US_EAST_1"
}

variable cluster_mongo_db_major_version {
  description = "Version of MongoDB."
  type = string
  default = "5.0"
}

variable database_user_role_database_name {
  description = "Database user role database name."
  type = string
  default = "admin"
}

variable activate_mongo_db_resources {
  description = "Enter Yes to activate MongoDB Atlas CloudFormation resource types. If you already activated resources in your AWS Region, enter No."
  type = string
  default = "Yes"
}

variable qss3_bucket_name {
  description = "Name of the S3 bucket for your copy of the deployment assets. Keep the default name unless you are customizing the template. Changing the name updates code  references to point to a new location."
  type = string
  default = "aws-quickstart"
}

variable qss3_key_prefix {
  description = "S3 key prefix that is used to simulate a folder for your copy of the  deployment assets. Keep the default prefix unless you are customizing  the template. Changing the prefix updates code references to point to  a new location."
  type = string
  default = "quickstart-mongodb-atlas-mean-stack-aws-fargate-integration/"
}

variable qss3_bucket_region {
  description = "AWS Region where the S3 bucket (QSS3BucketName) is hosted. Keep  the default Region unless you are customizing the template. Changing the Region  updates code references to point to a new location. When using your own bucket,  specify the Region."
  type = string
  default = "us-east-1"
}

variable client_service_ecr_image_uri {
  description = "ECR image URI of client application."
  type = string
}

variable server_service_ecr_image_uri {
  description = "ECR image URI of server application."
  type = string
}


// atlas-basic parameters

variable "region" {
    description = "Atlas cluster region"
    default     = "US_EAST_1"
    type        = string
}
variable "environmentId" {
  description = "Environment name to display it in the resources tags"
  type        = string
}

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
variable "password" {
  description = "MongoDB Atlas User Password"
  type        = list(string)
}
variable "database_name" {
  description = "The Database in the cluster"
  type        = list(string)
}
