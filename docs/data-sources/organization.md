# Data Source: mongodbatlas_organization

`mongodbatlas_organization` describes all MongoDB Atlas Organizations. This represents organizations that have been created.

## Example Usage

### Using project_id attribute to query
```terraform

data "mongodbatlas_organization" "test" {
  org_id = "<org_id>"
}
```

## Argument Reference

* `org_id` - Unique 24-hexadecimal digit string that identifies the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - Human-readable label that identifies the organization.
* `id` - Unique 24-hexadecimal digit string that identifies the organization.
* `is_deleted` - Flag that indicates whether this organization has been deleted.
* `api_access_list_required` - (Optional) Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.
* `multi_factor_auth_required` - (Optional) Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization. To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.
* `restrict_employee_access` - (Optional) Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure for any deployment in the specified organization without explicit permission. Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues. To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.
* `gen_ai_features_enabled` - (Optional) Flag that indicates whether this organization has access to generative AI features. This setting only applies to Atlas Commercial and defaults to `true`. With this setting on, Project Owners may be able to enable or disable individual AI features at the project level. To learn more, see https://www.mongodb.com/docs/generative-ai-faq/.
* `skip_default_alerts_settings` - (Optional) Flag that indicates whether to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform. Defaults to `true`.

    ~> **NOTE:** 
    - *If you created an organization with our Terraform provider version >=1.30.0, this field will be set to `true` by default.*
    - *If you have an organization created with our Terraform provider version <1.30.0, this field might be `false`, which is the [API default value](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Organizations/operation/createOrganization). To prevent creation of future default alerts, you can explicitly set this to `true` using the [`mongodbatlas_organization`](../resources/organization.md) resource.*
  
See [MongoDB Atlas API - Organization](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/getOrganization) Documentation for more information.
