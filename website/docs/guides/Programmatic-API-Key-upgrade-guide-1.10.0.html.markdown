---
layout: "mongodbatlas"
page_title: "Upgrade Guide for Terraform MongoDB Atlas Provider Programmatic API Key Resource in v1.10.0"
sidebar_current: "docs-mongodbatlas-guides-Programmatic-API-Key-upgrade-guide"
description: |-
MongoDB Atlas Provider : Upgrade and Information Guide
---

# MongoDB Atlas Provider: Programmatic API Key Upgrade Guide in v1.10.0
In Terraform MongoDB Atlas Provider v1.10.0, some improvements were introduced which mainly focus on the MongoDB Atlas Programmatic API Keys (PAK) handling. This guide will help you to transition smoothly from the previous version which this resource was first released (v1.8.0) to the new version (v1.10.0).

For comprehensive Upgrade Guide on all v1.10.0 modifications see [here](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.10.0-upgrade-guide). 

## Changes Overview
* `api_keys` parameter is deprecated from the `mongodbatlas_project` resource.
* The `mongodbatlas_project_api_key` resource is extended to include a `project_assignment` parameter.

## Upgrade Steps

**1. Backup Current Terraform State File**

Before you begin any modification process, it is always important to backup your current Terraform State file. This is necessary in case anything goes wrong, and you need to revert to the previous state.

**2. Saftely remove resource `mongodbatlas_project` from Terraform State**

Locate the `mongodbatlas_project` resource block containtaining the `api_keys` parameter in your state file that you wish to migrate to new workflow introduced in v1.10.0 of Terraform Provider for MongoDB Atlas. From there you can remove it from Terraform State file. This means that assignments of the `api_keys` parameter as well as project resource itself will be preserved in the actual infrastructure, but Terraform will no longer manage them.

For example if this was your current `mongodbatlas_project` resource block:
```
resource "mongodbatlas_project" "test" {
  name   = "projectName"
  org_id = var.org_id

  api_keys {
    api_key_id = mongodbatlas_api_key.orgKey1.api_key_id
    role_names = ["GROUP_OWNER"]
  }
  
  depends_on = [
    mongodbatlas_api_key.orgKey1
  ]

}
```

Then you would remove from state file with: 

```
$ terraform state rm mongodbatlas_project.test
```

**3. Update Terraform Scripts**

Now, open your Terraform scripts (i.e. main.tf file). Locate and remove the `api_keys` parameter from the `mongodbatlas_project` resource block from them as well. This is to make sure the parameter is no longer present in your scripts after you've already removed it from the state file. At this point you also want to include the new `mongodbatlas_project_api_key` resource block as well to assign key at the project level.

For example, the revised Terraform script should look like:

```
resource "mongodbatlas_project" "test" {
  name   = "projectName"
  org_id = var.org_id
}

resource "mongodbatlas_project_api_key" "test2" {
  description = "test create and assign"
  project_id  = mongodbatlas_project.test.project_id
  role_names  = ["GROUP_OWNER"]
}

```

**4. Import `mongodbatatlas_project` and "mongodbatlas_project_api_key" resouces back into Terraform State**

Again there should be no impact to real world resources, we are simply updating Terraform State to reflect current infrastructure environment. Hence these changes are made _without_ destroying and recreating them.

`mongodbatatlas_project` must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_project.test 5d09d6a59ccf6445652a444a
```

API Keys must be imported using org ID, API Key ID e.g.

```
$ terraform import mongodbatlas_project_api_key.test2 5d09d6a59ccf6445652a444a-6576974933969669
```

**5. Review & Apply the Changes**
Finally, run `$ terraform plan` to verify the import was successful. Ideally this should show: "No changes. Infrastructure is up-to-date." In such cases, you may choose not to run `$ terraform apply`. But in general, after making changes to your Terraform configurations, it's a good practice to run `$ terraform apply` to ensure your infrastructure matches your configuration after you have fully reviewed and are confidence in any proposed changed to your infrastructure .  

the changes with terraform apply. This will make Terraform update its state to reflect the current infrastructure.

**6. Review the Changes**
After applying the changes, review them to ensure everything has worked as expected. If you encounter any discrepancies or issues, use your backup to restore the previous state and investigate the cause of the problem before trying again.

By following these steps, you will be able to upgrade smoothly and efficiently to the new Programmatic API Key workflow for the Terraform MongoDB Atlas Provider introduced in v1.10.0! As always if you run into any issues or need to report any bugs feel free to open a new issue on our [GitHub repo](https://github.com/mongodb/terraform-provider-mongodbatlas/issues/new/choose). 

## Examples
We provide three examples in the [atlas-api-key folder](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/atlas-api-key) (under examples) to help you understand the new changes:

* "Create and Assign PAK Together": This example demonstrates how to create a Programmatic API Key and assign it to a project within the same resource block. This is a good place to start if you're used to creating and assigning keys at the same time.

* "Create and Assign PAK to Multiple Projects": This example shows how to create a PAK and assign it to several projects at once using the new `project_assignment` parameter. This is useful if you have multiple projects that require the same key.

* "Create and Assign PAK Separately" (Deprecated): This is the older method of creating and assigning PAKs, which is now deprecated. This example remains for reference and to help you understand the changes that have been made.

Before making any changes, please ensure you have thoroughly read and understood these examples, and that your current Terraform scripts align with the new PAK workflow.

Remember, your scripts will still work with deprecated features for now, but it's best to upgrade as soon as possible to benefit from the latest enhancements. Code removal is planned for v1.12.0 at which point prior PAK workflow will no longer function.

Lastly, in MongoDB Atlas, all PAKs are Organization API keys. Once created, a PAK is linked at the organization level with an 'Organization Member' role. However, these Organization API keys can also be assigned to one or more projects within the organization. When a PAK is assigned to a specific project, it essentially takes on the 'Project Owner' role for that particular project. This enables the key to perform operations at the project level, in addition to the organization level. The flexibility of PAKs provides a powerful mechanism for fine-grained access and control, once their functioning is clearly understood.

If you have any questions or face any issues during the migration, feel free to reach out to us by creating a GitHub Issue or PR in our repo. Thank you.  
