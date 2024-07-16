# Data Source: mongodbatlas_federated_settings_org_role_mapping

`mongodbatlas_federated_settings_org_role_mapping` provides an Federated Settings Org Role Mapping datasource. Atlas Cloud Federated Settings Org Role Mapping provides federated settings outputs for the configured Org Role Mapping.


## Example Usage

```terraform
resource "mongodbatlas_federated_settings_org_role_mapping" "org_group_role_mapping_import" {
  federation_settings_id = data.mongodbatlas_federated_settings.federated_settings.id
  org_id                 = "627a9683e7f7f7ff7fe306f14"

  external_group_name = "myGrouptest"

  role_assignments {
    org_id = "627a9683e7f7f7ff7fe306f14"
    roles     = ["ORG_MEMBER","ORG_GROUP_CREATOR","ORG_BILLING_ADMIN"]
  }

  role_assignments {
    group_id = "628aa20db7f7f7f98b81b8"
    roles     = ["GROUP_OWNER","GROUP_DATA_ACCESS_ADMIN","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY"]
  }

  role_assignments {
    group_id = "62b477f7f7f7f5e741489c"
    roles     = ["GROUP_OWNER","GROUP_DATA_ACCESS_ADMIN","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY","GROUP_DATA_ACCESS_READ_WRITE"]
  }
}

data "mongodbatlas_federated_settings_org_role_mapping" "role_mapping" {
  federation_settings_id = mongodbatlas_federated_settings_org_role_mapping.org_group_role_mapping_import.id
  org_id                 = "627a9683e7f7f7ff7fe306f14"
  role_mapping_id        = "627a9673e7f7f7ff7fe306f14"
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
### FederatedSettingsOrgRoleMappings

* `external_group_name` - Unique human-readable label that identifies the identity provider group to which this role mapping applies.
* `id` - Unique 24-hexadecimal digit string that identifies this role mapping.
* `role_assignments` - Atlas roles and the unique identifiers of the groups and organizations associated with each role.
* `group_id` - Unique identifier of the project to which you want the role mapping to apply.
* `role` - Specifies the Role that is attached to the Role Mapping.


For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
