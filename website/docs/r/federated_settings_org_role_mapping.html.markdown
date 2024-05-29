---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: mongodbatlas_federated_settings_org_role_mapping"
sidebar_current: "docs-mongodbatlas-resource-federated-settings-org-role-mapping"
description: |-
    Provides a federated settings Role Mapping resource.
---

# Resource: mongodbatlas_federated_settings_org_role_mapping

`mongodbatlas_federated_settings_org_role_mapping` provides an Role Mapping resource. This allows organization role mapping to be created.

## Example Usage

```terraform
resource "mongodbatlas_federated_settings_org_role_mapping" "org_group_role_mapping_import" {
  federation_settings_id = "627a9687f7f7f7f774de306f14"
  org_id                 = "627a9683e7f7f7ff7fe306f14"
  external_group_name    = "myGrouptest"

  role_assignments {
    org_id = "627a9683e7f7f7ff7fe306f14"
    roles     = ["ORG_MEMBER","ORG_GROUP_CREATOR","ORG_BILLING_ADMIN"]
  }

  role_assignments {
    group_id = "628aa20d7f7f7f7f7098b81b8"
    roles     = ["GROUP_OWNER","GROUP_DATA_ACCESS_ADMIN","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY"]
  }

  role_assignments {
    group_id = "628aa20d7f7f7f7f7078b81b8"
    roles     = ["GROUP_OWNER","GROUP_DATA_ACCESS_ADMIN","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY","GROUP_DATA_ACCESS_READ_WRITE"]
  }
}
```

## Argument Reference

* `federation_settings_id` - (Required) Unique 24-hexadecimal digit string that identifies the federated authentication configuration.
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
* `external_group_name` - Unique human-readable label that identifies the identity provider group to which this role mapping applies.
* `role_assignments` - Atlas roles and the unique identifiers of the groups and organizations associated with each role.
    * `group_id` - Unique identifier of the project to which you want the role mapping to apply.
    * `roles` - Specifies the Roles that are attached to the Role Mapping. Available role IDs can be found on [the User Roles
  Reference](https://www.mongodb.com/docs/atlas/reference/user-roles/).


## Attributes Reference

In addition to all arguments above, the following attributes are exported:
* `role_mapping_id` - Unique 24-hexadecimal digit string that identifies this role mapping.

## Import

FederatedSettingsOrgRoleMapping can be imported using federation_settings_id-org_id-role_mapping_id, e.g.

```
$ terraform import mongodbatlas_federated_settings_org_role_mapping.org_group_role_mapping_import 6287a663c7f7f7f71c441c6c-627a96837f7f7f7e306f14-628ae97f7f7468ea3727
```

For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
