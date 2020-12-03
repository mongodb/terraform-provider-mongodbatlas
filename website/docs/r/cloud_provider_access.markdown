---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_cloud_provider_access"
sidebar_current: "docs-mongodbatlas-datasource-cloud-provider-access"
description: |-
    Provides a Cloud Provider Access settings resource for registration, authorization, and deauthorization
---

# mongodbatlas_cloud_provider_access

`mongodbatlas_cloud_provider_access` Allows you to register and authorize AWS IAM roles in Atlas.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** The update of the argument iam_assumed_role_arn is one step in a procedure to create unified AWS access for Atlas services. For the complete procedure, see [Set Up Unified AWS Access](https://docs.atlas.mongodb.com/security/set-up-unified-aws-access/#set-up-unified-aws-access).

## Example Usage

```hcl

resource "mongodbatlas_cloud_provider_access" "test_role" {
   project_id = "<PROJECT-ID>"
   provider_name = "AWS"
}

```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get all Cloud Provider Access 
* `provider_name` - (Required) The cloud provider for which to create a new role. Currently only AWS is supported.
* `iam_assumed_role_arn` - (Optional) - ARN of the IAM Role that Atlas assumes when accessing resources in your AWS account. This value is required after the creation (register of the role) as part of [Set Up Unified AWS Access](https://docs.atlas.mongodb.com/security/set-up-unified-aws-access/#set-up-unified-aws-access).


## Attributes Reference

* `id` - Unique identifier used by terraform for internal management.
* `atlas_assumed_role_external_id` - Unique external ID Atlas uses when assuming the IAM role in your AWS account.
* `atlas_aws_account_arn`          - ARN associated with the Atlas AWS account used to assume IAM roles in your AWS account.
* `authorized_date`                - Date on which this role was authorized.
* `created_date`                   - Date on which this role was created.
* `feature_usages`                 - Atlas features this AWS IAM role is linked to.
* `provider_name`                  - Name of the cloud provider. Currently limited to AWS.
* `role_id`                        - Unique ID of this role.

## Authorize role

Once the resource is created add the field `iam_assumed_role_arn` see [Set Up Unified AWS Access](https://docs.atlas.mongodb.com/security/set-up-unified-aws-access/#set-up-unified-aws-access) , and execute a new `terraform apply` this will create a PATCH request.

```hcl

resource "mongodbatlas_cloud_provider_access" "test_role" {
   project_id = "<PROJECT-ID>"
   provider_name = "AWS"
   iam_assumed_role_arn = "arn:aws:iam::520983883852:role/mongodb_ec2_s3"
}

```

## Import

The Cloud Provider Access resource can be imported using project ID and the provider name and mongodbatlas role id, in the format `project_id`-`provider_name`-`role_id`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_access.my_role 1112222b3bf99403840e8934-AWS-5fc17d476f7a33224f5b224e
```

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/cloud-provider-access-create-one-role/) Documentation for more information.