variable "project_id" { type = string }
variable "name" { type = string }
variable "provider_name" { type = string }
variable "region_name" { type = string }
variable "instance_size" { type = string }
variable "tags" { type = map(string) default = {} }
