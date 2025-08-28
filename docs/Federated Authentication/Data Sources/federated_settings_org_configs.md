# Data Source: mongodbatlas_federated_settings_org_configs

`mongodbatlas_federated_settings_org_configs` provides an Federated Settings Identity Providers datasource. Atlas Cloud Federated Settings Identity Providers provides federated settings outputs for the configured Identity Providers.


## Example Usage

```terraform
resource "mongodbatlas_federated_settings_org_config" "org_connection" {
  federation_settings_id            = "627a9687f7f7f7f774de306f14"
  org_id                            = "627a9683ea7ff7f74de306f14"
  data_access_identity_provider_ids = ["64d613677e1ad50839cce4db"]
  domain_restriction_enabled        = false
  domain_allow_list                 = ["mydomain.com"]
  post_auth_role_grants             = ["ORG_MEMBER"]
  identity_provider_id              = "0oaqyt9fc2ySTWnA0357"
}

data "mongodbatlas_federated_settings_org_configs" "org_configs_ds" {
  federation_settings_id = data.mongodbatlas_federated_settings_org_config.org_connection.federation_settings_id
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration. 
* `page_num` - (Optional)  	The page to return. Defaults to `1`.
* `items_per_page` - (Optional) Number of items to return per page, up to a maximum of 500. Defaults to `100`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - Includes cloudProviderSnapshot object for each item detailed in the results array section.
* `totalCount` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.

### FederatedSettingsOrgConfigs

* `domain_allow_list` - List that contains the approved domains from which organization users can log in.
* `domain_restriction_enabled` - Flag that indicates whether domain restriction is enabled for the connected organization.
* `identity_provider_id` - Legacy 20-hexadecimal digit string that identifies the SAML access identity provider that this connected org config is associated with. This id can be found in two ways:
  1. Within the Federation Management UI in Atlas in the Identity Providers tab by clicking the info icon in the IdP ID row of a configured SAML identity provider
  2. `okta_idp_id` on the `mongodbatlas_federated_settings_identity_provider` resource
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `post_auth_role_grants` - List that contains the default roles granted to users who authenticate through the IdP in a connected organization.
* `data_access_identity_provider_ids` - The collection of unique ids representing the identity providers that can be used for data access in this organization.
* `role_mappings` - Role mappings that are configured in this organization. See [below](#role_mappings)
* `user_conflicts` - List that contains the users who have an email address that doesn't match any domain on the allowed list. See [below](#user-conflicts)

  ### Role_mappings
* `external_group_name` - Unique human-readable label that identifies the identity provider group to which this role mapping applies.
* `id` - Unique 24-hexadecimal digit string that identifies this role mapping.
* `role_assignments` - Atlas roles and the unique identifiers of the groups and organizations associated with each role.
* `group_id` - Unique identifier of the project to which you want the role mapping to apply.
* `role` - Specifies the Role that is attached to the Role Mapping.
### User Conflicts
* `email_address` - Email address of the the user that conflicts with selected domains.
* `federation_settings_id` - Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `first_name` - First name of the the user that conflicts with selected domains.
* `last_name` - Last name of the the user that conflicts with selected domains.
* `user_id` - Name of the Atlas user that conflicts with selected domains.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
