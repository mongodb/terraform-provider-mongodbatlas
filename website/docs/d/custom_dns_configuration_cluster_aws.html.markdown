---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: custom_dns_configuration_cluster_aws"
sidebar_current: "docs-mongodbatlas-datasource-custom_dns_configuration_cluster_aws"
description: |-
    Describes a Custom DNS Configuration for Atlas Clusters on AWS.
---

# mongodbatlas_custom_dns_configuration_cluster_aws

`mongodbatlas_custom_dns_configuration_cluster_aws` describes a Custom DNS Configuration for Atlas Clusters on AWS.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

```hcl
resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
	project_id                  = "<project-id>"
	enabled                     = true
}

data "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
    project_id = mongodbatlas_custom_dns_configuration_cluster_aws.test.id
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enabled` - Indicates whether the project's clusters deployed to AWS use custom DNS.


See detailed information for arguments and attributes: [MongoDB API Custom DNS Configuration for Atlas Clusters on AWS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-get)