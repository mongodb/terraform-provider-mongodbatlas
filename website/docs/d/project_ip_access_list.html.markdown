---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: project_ip_access_list"
sidebar_current: "docs-mongodbatlas-datasource-project-ip-access-list"
description: |-
    Provides an IP Access List resource.
---

# Data Source: mongodbatlas_project_ip_access_list

`mongodbatlas_project_ip_access_list` describes an IP Access List entry resource. The access list grants access from IPs, CIDRs or AWS Security Groups (if VPC Peering is enabled) to clusters within the Project.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

~> **IMPORTANT:**
When you remove an entry from the access list, existing connections from the removed address(es) may remain open for a variable amount of time. How much time passes before Atlas closes the connection depends on several factors, including how the connection was established, the particular behavior of the application or driver using the address, and the connection protocol (e.g., TCP or UDP). This is particularly important to consider when changing an existing IP address or CIDR block as they cannot be updated via the Provider (comments can however), hence a change will force the destruction and recreation of entries.   


## Example Usage

### Using CIDR Block
```terraform
resource "mongodbatlas_project_ip_access_list" "test" {
  project_id = "<PROJECT-ID>"
  cidr_block = "1.2.3.4/32"
  comment    = "cidr block for tf acc testing"
}

data "mongodbatlas_project_ip_access_list" "test" {
	project_id = mongodbatlas_project_ip_access_list.test.project_id
	cidr_block = mongodbatlas_project_ip_access_list.test.cidr_block
}
```

### Using IP Address
```terraform
resource "mongodbatlas_project_ip_access_list" "test" {
  project_id = "<PROJECT-ID>"
  ip_address = "2.3.4.5"
  comment    = "ip address for tf acc testing"
}

data "mongodbatlas_project_ip_access_list" "test" {
	project_id = mongodbatlas_project_ip_access_list.test.project_id
	ip_address = mongodbatlas_project_ip_access_list.test.ip_address
}
```

### Using an AWS Security Group
```terraform
resource "mongodbatlas_network_container" "test" {
  project_id       = "<PROJECT-ID>"
  atlas_cidr_block = "192.168.208.0/21"
  provider_name    = "AWS"
  region_name      = "US_EAST_1"
}

resource "mongodbatlas_network_peering" "test" {
  project_id             = "<PROJECT-ID>"
  container_id           = mongodbatlas_network_container.test.container_id
  accepter_region_name   = "us-east-1"
  provider_name          = "AWS"
  route_table_cidr_block = "172.31.0.0/16"
  vpc_id                 = "vpc-0d93d6f69f1578bd8"
  aws_account_id         = "232589400519"
}

resource "mongodbatlas_project_ip_access_list" "test" {
  project_id         = "<PROJECT-ID>"
  aws_security_group = "sg-0026348ec11780bd1"
  comment            = "TestAcc for awsSecurityGroup"

  depends_on = [mongodbatlas_network_peering.test]
}

data "mongodbatlas_project_ip_access_list" "test" {
	project_id = mongodbatlas_project_ip_access_list.test.project_id
	aws_security_group = mongodbatlas_project_ip_access_list.test.aws_security_group
}
```

~> **IMPORTANT:** In order to use AWS Security Group(s) VPC Peering must be enabled like in the above example.

## Argument Reference

* `project_id` - (Required) Unique identifier for the project to which you want to add one or more access list entries.
* `aws_security_group` - (Optional) Unique identifier of the AWS security group to add to the access list.
* `cidr_block` - (Optional) Range of IP addresses in CIDR notation to be added to the access list.
* `ip_address` - (Optional) Single IP address to be added to the access list.

-> **NOTE:** One of the following attributes must set:  `aws_security_group`, `cidr_block`  or `ip_address`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used by Terraform for internal management and can be used to import.
* `comment` - Comment to add to the access list entry.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/access-lists/)
