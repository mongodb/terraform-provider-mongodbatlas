---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_federated_settings_role_mappings"
sidebar_current: "docs-mongodbatlas-datasource-federated-settings-role-mappings"
description: |-
    Provides a federated settings Role Mapping datasource.
---

# Data Source: mongodbatlas_federated_settings_org_role_mappings

`mongodbatlas_federated_settings_org_role_mappings` provides an Federated Settings Org Role Mapping datasource. Atlas Cloud Federated Settings Org Role Mapping provides federated settings outputs for the configured Org Role Mapping.


## Example Usage

```terraform
resource "mongodbatlas_federated_settings_org_role_mapping" "org_group_role_mapping_import" {
  federation_settings_id = ""
  org_id                 = "627a9683e7f7f7ff7fe306f14"
  group_id               = "628aa20d7f7f7f7f7098b81b8"
  external_group_name    = "myGrouptest"
  organization_roles     = ["ORG_OWNER", "ORG_MEMBER", "ORG_BILLING_ADMIN", "ORG_GROUP_CREATOR", "ORG_READ_ONLY"]
  group_roles            = ["GROUP_OWNER","GROUP_CLUSTER_MANAGER","GROUP_DATA_ACCESS_ADMIN","GROUP_DATA_ACCESS_READ_WRITE","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY","GROUP_READ_ONLY"]
}

data "mongodbatlas_federated_settings_org_role_mappings" "role_mappings" {
  federation_settings_id = mongodbatlas_federated_settings_org_role_mapping.org_group_role_mapping_import.id
  org_id                 = "627a9683e7f7f7ff7fe306f14"
  page_num = 1
  items_per_page = 5
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `page_num` - (Optional)  	The page to return. Defaults to `1`.
* `items_per_page` - (Optional) Number of items to return per page, up to a maximum of 500. Defaults to `100`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `results` - Includes cloudProviderSnapshot object for each item detailed in the results array section.
* `totalCount` - Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.

### FederatedSettingsOrgRoleMappings

* `external_group_name` - Unique human-readable label that identifies the identity provider group to which this role mapping applies.
* `id` - Unique 24-hexadecimal digit string that identifies this role mapping.
* `role_assignments` - Atlas roles and the unique identifiers of the groups and organizations associated with each role.
* `group_id` - Unique identifier of the project to which you want the role mapping to apply.
* `role` - Specifies the Role that is attached to the Role Mapping.


For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
