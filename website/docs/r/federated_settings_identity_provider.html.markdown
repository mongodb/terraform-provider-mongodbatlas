---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_cloud_federated_settings_identity_provider"
sidebar_current: "docs-mongodbatlas-cloud-federated-settings-identity-provider"
description: |-
    Provides an Federated Settings Identity Provider Resource.
---

# mongodbatlas_cloud_federated_settings_identity_provider

`mongodbatlas_cloud_federated_settings_identity_provider` provides an Atlas Cloud Federated Settings Identity Provider resource provides a subset of settings to be maintained post import of the existing resource.
## Example Usage

```terraform
resource "mongodbatlas_cloud_federated_settings_identity_provider" "identity_provider" {
  federation_settings_id     = "627a9687f7f7f7f774de306f14"
  name = "mongodb_federation_test"
  associated_domains           = ["yourdomain.com"]
  sso_debug_enabled = true
  status = "ACTIVE"
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `name` - Human-readable label that identifies the identity provider.
* `associated_domains` - List that contains the domains associated with the identity provider.
* `sso_debug_enabled` - Flag that indicates whether the identity provider has SSO debug enabled.
* `status`- String enum that indicates whether the identity provider is active or not.
Accepted values are ACTIVE or INACTIVE.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:


### FederatedSettingsIdentityProvider

* `okta_idp_id` - Unique 20-hexadecimal digit string that identifies the IdP.

## Import

Identity Provider must be imported using federation_settings_id-okta_idp_id, e.g.

```
$ terraform import mongodbatlas_cloud_federated_settings_identity_provider.identity_provider 6287a663c660f52b1c441c6c-0oad4fas87jL5Xnk1297
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)