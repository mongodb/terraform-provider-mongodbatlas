---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_cloud_provider_access"
sidebar_current: "docs-mongodbatlas-resource-cloud-provider-access"
description: |-
    Provides a Cloud Provider Access settings resource for registration, authorization, and deauthorization
---

# Resource: Cloud Provider Access Configuration Paths

The Terraform MongoDB Atlas Provider offers two either-or/mutually exclusive paths to perform an authorization for a cloud provider role -

* A Two Resource path: consisting of `mongodbatlas_cloud_provider_access_setup` and `mongodbatlas_cloud_provider_access_authorization`. The first resource, `mongodbatlas_cloud_provider_access_setup`, only generates
the initial configuration (create, delete operations). The second resource, `mongodbatlas_cloud_provider_access_authorization`, helps to perform the authorization using the role_id of the first resource. This path is helpful in a multi-provider Terraform file, and allows for a single and decoupled apply. See example of this Two Resource path option with AWS Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-cloud-provider-access/aws) and AZURE Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-cloud-provider-access/azure). 

* A Single Resource path: using the `mongodbatlas_cloud_provider_access` that at provision time sets up all the required configuration for a given provider, then with a subsequent update it can perform the authorize of the role. Note this path requires two `terraform apply` commands, once for setup and once for auth. This resource supports only `AWS`.
* A Two Resource path: consisting of `mongodbatlas_cloud_provider_access_setup` and `mongodbatlas_cloud_provider_access_authorization`. The first resource, `mongodbatlas_cloud_provider_access_setup`, only generates
the initial configuration (create, delete operations). The second resource, `mongodbatlas_cloud_provider_access_authorization`, helps to perform the authorization using the role_id of the first resource. This path is helpful in a multi-provider Terraform file, and allows for a single and decoupled apply. See example of this Two Resource path option with AWS Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-cloud-provider-access/aws) and AZURE Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-cloud-provider-access/azure). 

* A Single Resource path: using the `mongodbatlas_cloud_provider_access` that at provision time sets up all the required configuration for a given provider, then with a subsequent update it can perform the authorize of the role. Note this path requires two `terraform apply` commands, once for setup and once for auth. This resource supports only `AWS`.
**WARNING:** The resource `mongodbatlas_cloud_provider_access` is deprecated and will be removed in version v1.14.0, use the Two Resource path instead.

-> **IMPORTANT** If you want to move from the single resource path to the two resources path see the [migration guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/0.9.1-upgrade-guide#migration-to-cloud-provider-access-setup)


## mongodbatlas_cloud_provider_access_setup

This is the first resource in the two-resource path as described above.

`mongodbatlas_cloud_provider_access_setup` Allows you to only register AWS or AZURE IAM roles in Atlas.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage with AWS

```terraform

resource "mongodbatlas_cloud_provider_access_setup" "test_role" {
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "AWS"
}

```

## Example Usage with Azure

```terraform

resource "mongodbatlas_cloud_provider_access_setup" "test_role" {
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "AZURE"
   azure_config {
      atlas_azure_app_id = "9f2deb0d-be22-4524-a403-df531868bac0"
      service_principal_id = "22f1d2a6-d0e9-482a-83a4-b8dd7dddc2c1"
      tenant_id = "91402384-d71e-22f5-22dd-759e272cdc1c"
   }
}

```

## Argument Reference

* `project_id` - (Required) The unique ID for the project
* `provider_name` - (Required) The cloud provider for which to create a new role. Currently only AWS and AZURE are supported. **WARNING** Changing the `provider_name`` will result in destruction of the existing resource and the creation of a new resource.
* `azure_config` - azure related configurations 
   * `atlas_azure_app_id` - Azure Active Directory Application ID of Atlas. This property is required when `provider_name = "AZURE".`
   * `service_principal_id`- UUID string that identifies the Azure Service Principal. This property is required when `provider_name = "AZURE".`
   * `tenant_id`          - UUID String that identifies the Azure Active Directory Tenant ID. This property is required when `provider_name = "AZURE".`

## Attributes Reference

* `id` - Unique identifier used by terraform for internal management.
* `aws_config` - aws related arn roles 
   * `atlas_assumed_role_external_id` - Unique external ID Atlas uses when assuming the IAM role in your AWS account.
   * `atlas_aws_account_arn`          - ARN associated with the Atlas AWS account used to assume IAM roles in your AWS account.
* `created_date`                   - Date on which this role was created.
* `last_updated_date`                - Date and time when this Azure Service Principal was last updated. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `role_id`                        - Unique ID of this role.

## Import: mongodbatlas_cloud_provider_access_setup
For consistency is has the same format as the regular mongodbatlas_cloud_provider_access resource 
can be imported using project ID and the provider name and mongodbatlas role id, in the format 
`project_id`-`provider_name`-`role_id`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_access_setup.my_role 1112222b3bf99403840e8934-AWS-5fc17d476f7a33224f5b224e
```

## mongodbatlas_cloud_provider_authorization

This is the second resource in the two-resource path as described above.
`mongodbatlas_cloud_provider_access_authorization`  Allows you to authorize an AWS or AZURE IAM roles in Atlas.

## Example Usage with AWS
```terraform

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "AWS"
}


resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id =  mongodbatlas_cloud_provider_access_setup.setup_only.project_id
   role_id    =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id

   aws_config {
      atlas_aws_account_arn = "arn:aws:iam::772401394250:role/test-user-role"
   }
}

```


```terraform

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "AZURE"
   azure_config {
      atlas_azure_app_id = "9f2deb0d-be22-4524-a403-df531868bac0"
      service_principal_id = "22f1d2a6-d0e9-482a-83a4-b8dd7dddc2c1"
      tenant_id = "91402384-d71e-22f5-22dd-759e272cdc1c"
	}
}


resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id =  mongodbatlas_cloud_provider_access_setup.setup_only.project_id
   role_id    =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id

   azure {
      atlas_azure_app_id = "9f2deb0d-be22-4524-a403-df531868bac0"
      service_principal_id = "22f1d2a6-d0e9-482a-83a4-b8dd7dddc2c1"
      tenant_id = "91402384-d71e-22f5-22dd-759e272cdc1c"
   }
}


## Argument Reference

* `project_id` - (Required) The unique ID for the project
* `role_id`    - (Required) Unique ID of this role returned by mongodb atlas api

Conditional 
* `aws`
   * `iam_assumed_role_arn` - (Required) ARN of the IAM Role that Atlas assumes when accessing resources in your AWS account. This value is required after the creation (register of the role) as part of [Set Up Unified AWS Access](https://docs.atlas.mongodb.com/security/set-up-unified-aws-access/#set-up-unified-aws-access).
   

## Attributes Reference

* `id`               - Unique identifier used by terraform for internal management.
* `authorized_date`  - Date on which this role was authorized.
* `feature_usages`   - Atlas features this AWS IAM role is linked to.


## mongodbatlas_cloud_provider_access

**WARNING:** The resource `mongodbatlas_cloud_provider_access` is deprecated and will be removed in version v1.14.0, use the Two Resource path instead.

`mongodbatlas_cloud_provider_access` Allows you to register and authorize AWS IAM roles in Atlas. This is the resource to use for the single resource path described above.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** The update of the argument iam_assumed_role_arn is one step in a procedure to create unified AWS access for Atlas services. For the complete procedure, see [Set Up Unified AWS Access](https://docs.atlas.mongodb.com/security/set-up-unified-aws-access/#set-up-unified-aws-access).

## Example Usage

```terraform

resource "mongodbatlas_cloud_provider_access" "test_role" {
   project_id = "64259ee860c43338194b0f8e"
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "AWS"
}

```

## Argument Reference

* `project_id` - (Required) The unique ID for the project
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

```terraform

resource "mongodbatlas_cloud_provider_access" "test_role" {
   project_id = "<PROJECT-ID>"
   provider_name = "AWS"
   iam_assumed_role_arn = "arn:aws:iam::772401394250:role/test-user-role"
}

```

## Import

The Cloud Provider Access resource can be imported using project ID and the provider name and mongodbatlas role id, in the format `project_id`-`provider_name`-`role_id`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_access.my_role 1112222b3bf99403840e8934-AWS-5fc17d476f7a33224f5b224e
```



See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/cloud-provider-access-create-one-role/) Documentation for more information.
