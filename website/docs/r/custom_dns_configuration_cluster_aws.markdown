---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: custom_dns_configuration_cluster_aws"
sidebar_current: "docs-mongodbatlas-resource-custom_dns_configuration_cluster_aws"
description: |-
    Provides a Custom DNS Configuration for Atlas Clusters on AWS resource.
---

# mongodbatlas_custom_dns_configuration_cluster_aws

`mongodbatlas_custom_dns_configuration_cluster_aws` provides a Custom DNS Configuration for Atlas Clusters on AWS resource. This represents a Custom DNS Configuration for Atlas Clusters on AWS that can be updated in an Atlas project.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource:
  * Organization Owner
  * Project Owner

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.


## Example Usage

```hcl
resource "mongodbatlas_custom_dns_configuration_cluster_aws" "test" {
  project_id    = "<PROJECT-ID>"
  enabled = true
}
```

## Argument Reference

* `project_id` - Required 	Unique identifier for the project.
* `enabled` - (Required) Indicates whether the project's clusters deployed to AWS use custom DNS. If `true`, the `Get All Clusters` and `Get One Cluster` endpoints return the `connectionStrings.private` and `connectionStrings.privateSrv` fields for clusters deployed to AWS .


## Import
Custom DNS Configuration for Atlas Clusters on AWS must be imported using auditing ID, e.g.

```
$ terraform import mongodbatlas_custom_dns_configuration_cluster_aws.test 1112222b3bf99403840e8934
```

See detailed information for arguments and attributes: [MongoDB API Custom DNS Configuration for Atlas Clusters on AWS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns)