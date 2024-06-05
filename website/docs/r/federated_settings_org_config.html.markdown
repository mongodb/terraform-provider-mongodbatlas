---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_federated_settings_org_config"
sidebar_current: "docs-mongodbatlas-resource-federated-settings-org-config"
description: |-
    Provides a federated settings Organization Configuration.
---

# Resource: mongodbatlas_federated_settings_org_config

`mongodbatlas_federated_settings_org_config` provides an Federated Settings Identity Providers datasource. Atlas Cloud Federated Settings Identity Providers provides federated settings outputs for the configured Identity Providers.


## Example Usage

~> **IMPORTANT** You **MUST** import this resource before you can manage it with this provider. 

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
  federation_settings_id = data.mongodbatlas_federated_settings_org_config.org_connection.id
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration. 
* `org_id` - (Required) Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `domain_allow_list` - List that contains the approved domains from which organization users can log in.
* `post_auth_role_grants` - (Optional) List that contains the default [roles](https://www.mongodb.com/docs/atlas/reference/user-roles/#std-label-organization-roles) granted to users who authenticate through the IdP in a connected organization.

* `domain_restriction_enabled` - (Required) Flag that indicates whether domain restriction is enabled for the connected organization.
* `identity_provider_id` - (Optional) Legacy 20-hexadecimal digit string that identifies the SAML access identity provider that this connected org config is associated with. Removing the attribute or providing the value `""` will detach/remove the SAML identity provider. This id can be found in two ways:
  1. Within the Federation Management UI in Atlas in the Identity Providers tab by clicking the info icon in the IdP ID row of a configured SAML identity provider
  2. `okta_idp_id` on the `mongodbatlas_federated_settings_identity_provider` resource
* `data_access_identity_provider_ids` - (Optional) The collection of unique ids representing the identity providers that can be used for data access in this organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_conflicts` - List that contains the users who have an email address that doesn't match any domain on the allowed list. See [below](#user-conflicts)

### User Conflicts
* `email_address` - Email address of the the user that conflicts with selected domains.
* `federation_settings_id` - Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `first_name` - First name of the the user that conflicts with selected domains.
* `last_name` - Last name of the the user that conflicts with selected domains.
* `user_id` - Name of the Atlas user that conflicts with selected domains.

## Import

FederatedSettingsOrgConfig must be imported using federation_settings_id-org_id, e.g.

```
$ terraform import mongodbatlas_federated_settings_org_config.org_connection 627a9687f7f7f7f774de306f14-627a9683ea7ff7f74de306f14
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)

