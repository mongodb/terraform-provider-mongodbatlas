---
subcategory: "Organizations"
---

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
* `users`- Returns a list of all pending and active MongoDB Cloud users associated with the specified organization.
* `api_access_list_required` - (Optional) Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.
* `multi_factor_auth_required` - (Optional) Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization. To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.
* `restrict_employee_access` - (Optional) Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure for any deployment in the specified organization without explicit permission. Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues. To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.
* `gen_ai_features_enabled` - (Optional) Flag that indicates whether this organization has access to generative AI features. This setting only applies to Atlas Commercial and defaults to `true`. With this setting on, Project Owners may be able to enable or disable individual AI features at the project level. To learn more, see https://www.mongodb.com/docs/generative-ai-faq/.
* `security_contact` - (Optional) String that specifies a single email address for the specified organization to receive security-related notifications. Specifying a security contact does not grant them authorization or access to Atlas for security decisions or approvals.
* `skip_default_alerts_settings` - (Optional) Flag that indicates whether to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform. Defaults to `true`.


### Users
* `id` - Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.
* `org_membership_status` - String enum that indicates whether the MongoDB Cloud user has a pending invitation to join the organization or they are already active in the organization.
* `roles` - Organization- and project-level roles assigned to one MongoDB Cloud user within one organization.
* `team_ids` - List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.
* `username` - Email address that represents the username of the MongoDB Cloud user.
* `country` - Two-character alphabetical string that identifies the MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.
* `invitation_created_at` - Date and time when MongoDB Cloud sent the invitation. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
* `invitation_expires_at` - Date and time when the invitation from MongoDB Cloud expires. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
* `inviter_username` - Username of the MongoDB Cloud user who sent the invitation to join the organization.
* `created_at` - Date and time when MongoDB Cloud created the current account. This value is in the ISO 8601 timestamp format in UTC.
* `first_name` - First or given name that belongs to the MongoDB Cloud user.
* `last_auth` - Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.
* `last_name` - Last name, family name, or surname that belongs to the MongoDB Cloud user.
* `mobile_number` - Mobile phone number that belongs to the MongoDB Cloud user.


~> **NOTE:** - Users with pending invitations created using [`mongodbatlas_project_invitation`](../resources/project_invitation.md) resource or via the deprecated [Invite One MongoDB Cloud User to Join One Project](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createprojectinvitation) endpoint are excluded (or cannot be managed) with this resource. See [MongoDB Atlas API - MongoDB Cloud Users](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-mongodb-cloud-users) for details.
To manage these users with this resource/data source, refer to our [Org Invitation to Cloud User Org Assignment Migration Guide](../guides/atlas-user-management).



~> **NOTE:** - If you create an organization with our Terraform provider version >=1.30.0, this field is set to `true` by default.<br> - If you have an existing organization created with our Terraform provider version <1.30.0, this field might be `false`, which is the [API default value](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createorganization). To prevent the creation of future default alerts, set this explicitly to `true` using the [`mongodbatlas_organization`](../resources/organization.md) resource.


See [MongoDB Atlas API - Organization](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/getOrganization) Documentation for more information.
