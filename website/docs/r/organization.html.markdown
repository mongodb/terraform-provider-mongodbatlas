---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: organization"
sidebar_current: "docs-mongodbatlas-resource-organization"
description: |-
    Provides a Organization resource.
---

# Resource: mongodbatlas_organization

`mongodbatlas_organization` provides programmatic management (including creation) of a MongoDB Atlas Organization resource.

~> **IMPORTANT NOTE:**  When you establish an Atlas organization using this resource, it automatically generates a set of initial public and private Programmatic API Keys. These key values are vital to store because you'll need to use them to grant access to the newly created Atlas organization. To use this resource, `role_names` for new API Key must have the ORG_OWNER role specified.

~> **IMPORTANT NOTE:** To use this resource, the requesting API Key must have the Organization Owner role. The requesting API Key's organization must be a paying organization. To learn more, see Configure a Paying Organization in the MongoDB Atlas documentation.

~> **NOTE** Import command is currently not supported for this resource.

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
* `org_owner_id` - (Required) Unique 24-hexadecimal digit string that identifies the Atlas user that you want to assign the Organization Owner role. This user must be a member of the same organization as the calling API key.  This is only required when authenticating with Programmatic API Keys. [MongoDB Atlas Admin API - Get User By Username](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUserByUsername)
* `description` - (Required) Programmatic API Key description

~> **NOTE:** Creating an organization will return a new API Key pair that can be used to authenticate and manage the new organization  with MongoDB Atlas Terraform modules/blueprints.  You cannot use the newly created API key pair to manage the newly created organization in the same Terraform module/blueprint that the organization is created in.

* `role_names` - (Required) List of Organization roles that the Programmatic API key needs to have. Ensure that you provide at least one role and ensure all roles are valid for the Organization.  You must specify an array even if you are only associating a single role with the Programmatic API key. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#organization-roles) describes the roles that you can assign to a Programmatic API key.
* `federation_settings_id` - (Optional) Unique 24-hexadecimal digit string that identifies the federation to link the newly created organization to. If specified, the proposed Organization Owner of the new organization must have the Organization Owner role in an organization associated with the federation.
* `api_access_list_required` - (Optional) Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.
* `multi_factor_auth_required` - (Optional) Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization. To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.
* `restrict_employee_access` - (Optional) Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure for any deployment in the specified organization without explicit permission. Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues. To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `org_id` - The organization id.
* `public_key` - Public API key value set for the specified organization API key.
* `private_key` - Redacted private key returned for this organization API key. This key displays unredacted when first created and is saved within the Terraform state file.


For more information see: [MongoDB Atlas Admin API Organization](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/createOrganization)  Documentation for more information.
