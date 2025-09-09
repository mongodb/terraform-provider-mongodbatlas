---
subcategory: "Cloud Provider Access"
---

# Resource: Cloud Provider Access Configuration Paths

The Terraform MongoDB Atlas Provider offers the following path to perform an authorization for a cloud provider role -

* A Two Resource path: consisting of `mongodbatlas_cloud_provider_access_setup` and `mongodbatlas_cloud_provider_access_authorization`. The first resource, `mongodbatlas_cloud_provider_access_setup`, only generates
the initial configuration (create, delete operations). The second resource, `mongodbatlas_cloud_provider_access_authorization`, helps to perform the authorization using the role_id of the first resource. This path is helpful in a multi-provider Terraform file, and allows for a single and decoupled apply. See example of this Two Resource path option with AWS Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/aws) and AZURE Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/azure). 


-> **IMPORTANT** If you want to move from the single resource path to the two resources path, see the [Migration Guide](../guides/0.9.1-upgrade-guide#migration-to-cloud-provider-access-setup)


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

## Example Usage with GCP

```terraform

resource "mongodbatlas_cloud_provider_access_setup" "test_role" {
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "GCP"
}

```

### Further Examples
- [AWS Cloud Provider Access](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/aws)
- [Azure Cloud Provider Access](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/azure)
- [GCP Cloud Provider Access](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/gcp)


## Argument Reference

* `project_id` - (Required) The unique ID for the project
* `provider_name` - (Required) The cloud provider for which to create a new role. Currently, AWS, AZURE and GCP are supported. **WARNING** Changing the `provider_name` will result in destruction of the existing resource and the creation of a new resource.
* `azure_config` - azure related configurations 
   * `atlas_azure_app_id` - Azure Active Directory Application ID of Atlas. This property is required when `provider_name = "AZURE".`
   * `service_principal_id`- UUID string that identifies the Azure Service Principal. This property is required when `provider_name = "AZURE".`
   * `tenant_id`          - UUID String that identifies the Azure Active Directory Tenant ID. This property is required when `provider_name = "AZURE".`
* `timeouts`- (Optional) The duration of time to wait for the resource to be created. The default timeout is `1h`. The timeout value is defined by a signed sequence of decimal numbers with a time unit suffix such as: `1h45m`, `300s`, `10m`, etc. The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).
* `delete_on_create_timeout`- (Optional) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.

## Attributes Reference

* `id` - Unique identifier used by terraform for internal management.
* `aws_config` - aws related arn roles 
   * `atlas_assumed_role_external_id` - Unique external ID Atlas uses when assuming the IAM role in your AWS account.
   * `atlas_aws_account_arn`          - ARN associated with the Atlas AWS account used to assume IAM roles in your AWS account.
* `gcp_config` - gcp related configuration
  * `status` - The status of the GCP cloud provider access setup. See [MongoDB Atlas API](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getgroupcloudprovideraccess#operation-getgroupcloudprovideraccess-200-body-application-vnd-atlas-2023-01-01-json-gcp-object-status).
  * `service_account_for_atlas` - The GCP service account email that Atlas uses.
* `created_date`                   - Date on which this role was created.
* `last_updated_date`                - Date and time when this Azure Service Principal was last updated. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
* `role_id`                        - Unique ID of this role.

-> **NOTE:** For more details on how attributes are used to enable access to cloud provider accounts see [AWS example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/aws) and [Azure example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_cloud_provider_access/azure). 

## Import: mongodbatlas_cloud_provider_access_setup
For consistency is has the same format as the regular mongodbatlas_cloud_provider_access resource 
can be imported using project ID and the provider name and mongodbatlas role id, in the format 
`project_id`-`provider_name`-`role_id`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_access_setup.my_role 1112222b3bf99403840e8934-AWS-5fc17d476f7a33224f5b224e
```

## mongodbatlas_cloud_provider_access_authorization

This is the second resource in the two-resource path as described above.
`mongodbatlas_cloud_provider_access_authorization`  Allows you to authorize an AWS or AZURE IAM roles in Atlas.

-> **IMPORTANT:** Changes to `project_id` or `role_id` will result in the destruction and recreation of the authorization resource. This action happens as these fields uniquely identify the authorization and cannot be modified in-place.

## Example Usage with AWS
```terraform

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
   project_id = "64259ee860c43338194b0f8e"
   provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
   project_id =  mongodbatlas_cloud_provider_access_setup.setup_only.project_id
   role_id    =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id

   aws {
      iam_assumed_role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/test-user-role"
   }
}

```

## Example Usage with Azure

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

```


## Argument Reference

* `project_id` - (Required) The unique ID for the project. **WARNING**: Changing the `project_id` will result in destruction of the existing authorization resource and the creation of a new authorization resource.
* `role_id`    - (Required) Unique ID of this role returned by mongodb atlas api. **WARNING**: Changing the `role_id` will result in destruction of the existing authorization resource and the creation of a new authorization resource.

Conditional 
* `aws`
   * `iam_assumed_role_arn` - (Required) ARN of the IAM Role that Atlas assumes when accessing resources in your AWS account. This value is required after the creation (register of the role) as part of [Set Up Unified AWS Access](https://docs.atlas.mongodb.com/security/set-up-unified-aws-access/#set-up-unified-aws-access).
* `azure`
   * `atlas_azure_app_id` - (Required) Azure Active Directory Application ID of Atlas.
   * `service_principal_id` - (Required) UUID string that identifies the Azure Service Principal.
   * `tenant_id` - (Required) UUID String that identifies the Azure Active Directory Tenant ID.

## Attributes Reference

* `id`               - Unique identifier used by terraform for internal management.
* `authorized_date`  - Date on which this role was authorized.
* `feature_usages`   - Atlas features this AWS IAM role is linked to.
* `gcp`
   * `service_account_for_atlas` - Email address for the Google Service Account created by Atlas.



## Import mongodbatlas_cloud_provider_access_authorization

You cannot import the Cloud Provider Access Authorization resource into Terraform. Instead, if the associated role is already authorized, you can recreate the resource without any adverse effects.
