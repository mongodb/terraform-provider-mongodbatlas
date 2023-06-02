---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: organization"
sidebar_current: "docs-mongodbatlas-resource-organization"
description: |-
    Provides a Organization resource.
---

# Resource: mongodbatlas_organization

`mongodbatlas_organization` provides programmatic management (including creation) of a MongoDB Atlas Organization resource.

~> **IMPORTANT NOTE:**  Creating an organization with this resource creates initial public and private key values that will need to be used in a separate terraform module/blueprint to permit access to new organization.


## Example Usage

```terraform
resource "mongodbatlas_organization" "test" {
  org_owner_id = "6205e5fffff79cde6f"
  name = "testCreateORG"
  description = "test API key from Org Creation Test"
  role_names = ["ORG_OWNER"]
}
```

## Argument Reference

* `name` - (Required) The name of the organization you want to create. (Cannot be changed via this Provider after creation.)
* `org_owner_id` - (Required) Unique 24-hexadecimal digit string that identifies the Atlas user that you want to assign the Organization Owner role. This user must be a member of the same organization as the calling API key.  This is only required when authenticating with Programmatic API Keys.  [https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUserByUsername](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUserByUsername)
* `description` - Programmatic API Key description

~> **NOTE:** Creating an organization will return a new API Key pair that can be used to authenticate and manage the new organization  with MongoDB Atlas Terraform modules/blueprints.  You cannot use the newly created API key pair to manage the newly created organization in the same Terraform module/blueprint that the organization is created in.


* `role_names` - (Required) List of Organization roles that the Programmatic API key needs to have. Ensure you provide: at least one role and ensure all roles are valid for the Organization.  You must specify an array even if you are only associating a single role with the Programmatic API key.
 The following are valid roles:
  * `ORG_OWNER`
  * `ORG_GROUP_CREATOR`
  * `ORG_BILLING_ADMIN`
  * `ORG_READ_ONLY`
  * `ORG_MEMBER` 
 
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `org_id` - The organization id.
* `public_key` - Public API key value set for the specified organization API key.
* `private_key` - Redacted private key returned for this organization API key. This key displays unredacted when first created and is saved within the Terraform state file.


## Import

Organization must be imported using organization ID, e.g.

```
$ terraform import mongodbatlas_organization.my_org 5d09d6a59ccf6445652a444a
```
For more information see: [MongoDB Atlas Admin API Organization](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/createOrganization)  Documentation for more information.
