---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_federated_settings_org_config"
sidebar_current: "docs-mongodbatlas-resource-federated-settings-org-config"
description: |-
    Provides an Federated Settings Organization Configuration.
---

# mongodbatlas_federated_settings_org_config

`mongodbatlas_federated_settings_org_config` provides an Federated Settings Identity Providers datasource. Atlas Cloud Federated Settings Identity Providers provides federated settings outputs for the configured Identity Providers.


## Example Usage

```terraform
resource "mongodbatlas_federated_settings_org_config" "org_connection" {
  federation_settings_id     = "627a9687f7f7f7f774de306f14"
  org_id                     = "627a9683ea7ff7f74de306f14"
  domain_restriction_enabled = false
  domain_allow_list          = ["mydomain.com"]
  dentity_provider_id = "0oad4fas87jL7f75Xnk1297"
}

data "mongodbatlas_federated_settings_org_configs" "org_configs_ds" {
  federation_settings_id = data.mongodbatlas_federated_settings_org_config.org_connection.id
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration. 
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

### FederatedSettingsOrgConfig
          
* `domain_allow_list` - List that contains the approved domains from which organization users can log in.
* `domain_restriction_enabled` - Flag that indicates whether domain restriction is enabled for the connected organization.

## Import

FederatedSettingsOrgConfig must be imported using federation_settings_id-org_id, e.g.

```
$ terraform import mongodbatlas_federated_settings_org_config.org_connection 6287a663c7f7f7f71c441c6c-627a96837f7f7f7e306f14-628ae97f7f7468ea3727

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)

