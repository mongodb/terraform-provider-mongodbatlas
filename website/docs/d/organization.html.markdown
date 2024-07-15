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

  
See [MongoDB Atlas API - Organization](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/getOrganization) Documentation for more information.
