---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_cloud_federated_settings_org_config"
sidebar_current: "docs-mongodbatlas-datasource-federated-settings-org-config"
description: |-
    Provides an Federated Settings Organization Configuration.
---

# mongodbatlas_cloud_federated_settings_org_configs

`mongodbatlas_cloud_federated_settings_org_config` provides an Federated Settings Identity Providers datasource. Atlas Cloud Federated Settings Organizational configuration provides federated settings outputs for the configured Organizational configuration.


## Example Usage

```terraform
resource "mongodbatlas_cloud_federated_settings_org_config" "org_connections" {
  federation_settings_id     = "627a9687f7f7f7f774de306f14"
  org_id                     = "627a9683ea7ff7f74de306f14"
  domain_restriction_enabled = false
  domain_allow_list          = ["mydomain.com"]
}

data "mongodbatlas_cloud_federated_settings_org_config" "org_configs_ds" {
  federation_settings_id = data.mongodbatlas_cloud_federated_settings_org_config.org_connections.id
  org_id                 = "627a9683ea7ff7f74de306f14"
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration. 
* `org_id` - Unique 24-hexadecimal digit string that identifies the connected organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### FederatedSettingsOrgConfig
          
* `domain_allow_list` - List that contains the approved domains from which organization users can log in.
* `domain_restriction_enabled` - Flag that indicates whether domain restriction is enabled for the connected organization.
* `identity_provider_id` - Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `post_auth_role_grants` - List that contains the default roles granted to users who authenticate through the IdP in a connected organization. If you provide a postAuthRoleGrants field in the request, the array that you provide replaces the current postAuthRoleGrants.

  ### Role_mappings
* `external_group_name` - Unique human-readable label that identifies the identity provider group to which this role mapping applies.
* `id` - Unique 24-hexadecimal digit string that identifies this role mapping.
* `role_assignments` - Atlas roles and the unique identifiers of the groups and organizations associated with each role.
* `group_id` - Unique identifier of the project that owns this Role Mapping Configuration.
* `role` - Specifies the Role that is attached to the Role Mapping.
### User Conflicts
* `email_address` - Email address of the the user that conflicts with selected domains.
* `federation_settings_id` - Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `first_name` - First name of the the user that conflicts with selected domains.
* `last_name` - Last name of the the user that conflicts with selected domains.
* `user_id` - Name of the Atlas user that conflicts with selected domains.

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
