---
subcategory: "Organizations"
---

# Resource: mongodbatlas_organization

`mongodbatlas_organization` provides programmatic management (including creation) of a MongoDB Atlas Organization resource.

~> **IMPORTANT NOTE:**  When you establish an Atlas organization using this resource, it automatically generates a set of initial public and private Programmatic API Keys. These key values are vital to store because you'll need to use them to grant access to the newly created Atlas organization. To use this resource, `role_names` for new API Key must have the ORG_OWNER role specified.

~> **IMPORTANT NOTE:** To use this resource, the requesting API Key must have the Organization Owner role. The requesting API Key's organization must be a paying organization. To learn more, see Configure a Paying Organization in the MongoDB Atlas documentation.

## Example Usage

```terraform
resource "mongodbatlas_organization" "test" {
  org_owner_id = "<ORG_OWNER_ID>"
  name = "testCreateORG"
  description = "test API key from Org Creation Test"
  role_names = ["ORG_OWNER"]
}
```

### Further Examples
- [Organization setup - step 1](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.41.1/examples/mongodbatlas_organization/organization-step-1)
- [Organization setup - step 2](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.41.1/examples/mongodbatlas_organization/organization-step-2)
- [Organization import](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.41.1/examples/mongodbatlas_organization/organization-import)

## Argument Reference

* `name` - (Required) The name of the organization.
* `org_owner_id` - (Optional) Unique 24-hexadecimal digit string that identifies the Atlas user that you want to assign the Organization Owner role. This user must be a member of the same organization as the calling API key.  This is only required when authenticating with Programmatic API Keys. [MongoDB Atlas Admin API - Get User By Username](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUserByUsername). This attribute is required in creation and can't be updated later.
* `description` - (Optional) Programmatic API Key description. This attribute is required in creation and can't be updated later.

~> **NOTE:** Creating an organization will return a new API Key pair that can be used to authenticate and manage the new organization  with MongoDB Atlas Terraform modules/blueprints.  These credentials will be used by the `mongodbatlas_organization` resource. In case of importing the resource, these credentials will be empty so the provider credentials will be used instead.

* `role_names` - (Optional) List of Organization roles that the Programmatic API key needs to have. Ensure that you provide at least one role and ensure all roles are valid for the Organization.  You must specify an array even if you are only associating a single role with the Programmatic API key. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#organization-roles) describes the roles that you can assign to a Programmatic API key. This attribute is required in creation and can't be updated later.
* `federation_settings_id` - (Optional) Unique 24-hexadecimal digit string that identifies the federation to link the newly created organization to. If specified, the proposed Organization Owner of the new organization must have the Organization Owner role in an organization associated with the federation. This attribute can't be updated after creation.
* `api_access_list_required` - (Optional) Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.
* `multi_factor_auth_required` - (Optional) Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization. To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.
* `restrict_employee_access` - (Optional) Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure for any deployment in the specified organization without explicit permission. Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues. To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.
* `gen_ai_features_enabled` - (Optional) Flag that indicates whether this organization has access to generative AI features. This setting only applies to Atlas Commercial and defaults to `true`. With this setting on, Project Owners may be able to enable or disable individual AI features at the project level. To learn more, see https://www.mongodb.com/docs/generative-ai-faq/.
* `security_contact` - (Optional) String that specifies a single email address for the specified organization to receive security-related notifications. Specifying a security contact does not grant them authorization or access to Atlas for security decisions or approvals.
* `skip_default_alerts_settings` - (Optional) Flag that indicates whether to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform. Defaults to `true`. 

~> **NOTE:** - If you create an organization with our Terraform provider version >=1.30.0, this field is set to `true` by default.<br> - If you have an existing organization created with our Terraform provider version <1.30.0, this field might be `false`, which is the [API default value](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createorganization). To prevent the creation of future default alerts, set this explicitly to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `org_id` - The organization id.
* `public_key` - Public API key value set for the specified organization API key.
* `private_key` - Redacted private key returned for this organization API key. This key displays unredacted when first created and is saved within the Terraform state file.

## Import

You can import an existing organization using the organization ID, e.g.:

```
$ terraform import mongodbatlas_organization.example 5d09d6a59ccf6445652a444a
```

~> **IMPORTANT:** When importing an existing organization, you should **NOT** specify the creation-only attributes (`org_owner_id`, `description`, `role_names`, `federation_settings_id`) in your Terraform configuration.

See the [Guide: Importing MongoDB Atlas Organizations](../guides/importing-organization) for more information.

For more information about the `mongodbatlas_organization` resource see: [MongoDB Atlas Admin API Organization](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-organizations).
