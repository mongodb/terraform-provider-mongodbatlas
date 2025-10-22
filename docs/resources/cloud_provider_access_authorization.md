---
subcategory: "Cloud Provider Access"
---

# Resource: mongodbatlas_cloud_provider_access_authorization

## Cloud Provider Access Configuration Paths

The Terraform MongoDB Atlas Provider offers a two-resource path to perform an authorization for a cloud provider role.
- The first resource, [`mongodbatlas_cloud_provider_access_setup`](cloud_provider_access_setup), only generates the initial configuration (create, delete operations).
- The second resource, `mongodbatlas_cloud_provider_access_authorization`, helps to perform the authorization using the `role_id` of the first resource.

This path is helpful in a multi-provider Terraform file, and allows for a single and decoupled apply. 
See example of this two-resource path option with AWS Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples/mongodbatlas_cloud_provider_access/aws), AZURE Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples/mongodbatlas_cloud_provider_access/azure) and GCP Cloud [here](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples/mongodbatlas_cloud_provider_access/gcp).


## mongodbatlas_cloud_provider_access_authorization

This is the second resource in the two-resource path as described above.

`mongodbatlas_cloud_provider_access_authorization` allows you to authorize an AWS, AZURE or GCP IAM roles in Atlas.

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

## Example Usage with GCP

```terraform

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id = "64259ee860c43338194b0f8e"
  provider_name = "GCP"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id =  mongodbatlas_cloud_provider_access_setup.setup_only.project_id
  role_id    =  mongodbatlas_cloud_provider_access_setup.setup_only.role_id
}
```

### Further Examples
- [AWS Cloud Provider Access](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples/mongodbatlas_cloud_provider_access/aws)
- [Azure Cloud Provider Access](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples/mongodbatlas_cloud_provider_access/azure)
- [GCP Cloud Provider Access](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.1.0/examples/mongodbatlas_cloud_provider_access/gcp)


## Argument Reference

* `project_id` - (Required) The unique ID for the project. **WARNING**: Changing the `project_id` will result in destruction of the existing authorization resource and the creation of a new authorization resource.
* `role_id`    - (Required) The unique ID of this role returned by the mongodb atlas api. **WARNING**: Changing the `role_id` will result in destruction of the existing authorization resource and the creation of a new authorization resource.

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
