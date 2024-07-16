# Resource: mongodbatlas_project

`mongodbatlas_project` provides a Project resource. This allows project to be created.

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot delete the Atlas project if any snapshots exist.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

## Example Usage

```terraform
data "mongodbatlas_roles_org_id" "test" {
}

resource "mongodbatlas_project" "test" {
  name   = "project-name"
  org_id = data.mongodbatlas_roles_org_id.test.org_id
  project_owner_id = "<OWNER_ACCOUNT_ID>"

  teams {
    team_id    = "5e0fa8c99ccf641c722fe645"
    role_names = ["GROUP_OWNER"]

  }
  teams {
    team_id    = "5e1dd7b4f2a30ba80a70cd4rw"
    role_names = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
  }

  limits {
    name = "atlas.project.deployment.clusters"
    value = 26
  }

  limits {
    name = "atlas.project.deployment.nodesPerPrivateLinkRegion"
    value = 51
  }

  is_collect_database_specifics_statistics_enabled = true
  is_data_explorer_enabled                         = true
  is_extended_storage_sizes_enabled                = true
  is_performance_advisor_enabled                   = true
  is_realtime_performance_panel_enabled            = true
  is_schema_advisor_enabled                        = true
}
```

## Argument Reference

* `name` - (Required) The name of the project you want to create.
* `org_id` - (Required) The ID of the organization you want to create the project within.
* `project_owner_id` - (Optional) Unique 24-hexadecimal digit string that identifies the Atlas user account to be granted the [Project Owner](https://docs.atlas.mongodb.com/reference/user-roles/#mongodb-authrole-Project-Owner) role on the specified project. If you set this parameter, it overrides the default value of the oldest [Organization Owner](https://docs.atlas.mongodb.com/reference/user-roles/#mongodb-authrole-Organization-Owner).
* `tags` - (Optional) Map that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the project. See [below](#tags).
* `with_default_alerts_settings` - (Optional) It allows users to disable the creation of the default alert settings. By default, this flag is set to true.
* `is_collect_database_specifics_statistics_enabled` - (Optional) Flag that indicates whether to enable statistics in [cluster metrics](https://www.mongodb.com/docs/atlas/monitor-cluster-metrics/) collection for the project. By default, this flag is set to true.
* `is_data_explorer_enabled` - (Optional) Flag that indicates whether to enable Data Explorer for the project. If enabled, you can query your database with an easy to use interface.  When Data Explorer is disabled, you cannot terminate slow operations from the [Real-Time Performance Panel](https://www.mongodb.com/docs/atlas/real-time-performance-panel/#std-label-real-time-metrics-status-tab) or create indexes from the [Performance Advisor](https://www.mongodb.com/docs/atlas/performance-advisor/#std-label-performance-advisor). You can still view Performance Advisor recommendations, but you must create those indexes from [mongosh](https://www.mongodb.com/docs/mongodb-shell/#mongodb-binary-bin.mongosh). By default, this flag is set to true.
* `is_extended_storage_sizes_enabled` - (Optional) Flag that indicates whether to enable extended storage sizes for the specified project. Clusters with extended storage sizes must be on AWS or GCP, and cannot span multiple regions. When extending storage size, initial syncs and cross-project snapshot restores will be slow. This setting should only be used as a measure of temporary relief; consider sharding if more storage is required.
* `is_performance_advisor_enabled` - (Optional) Flag that indicates whether to enable Performance Advisor and Profiler for the project. If enabled, you can analyze database logs to recommend performance improvements. By default, this flag is set to true.
* `is_realtime_performance_panel_enabled` - (Optional) Flag that indicates whether to enable Real Time Performance Panel for the project. If enabled, you can see real time metrics from your MongoDB database. By default, this flag is set to true.
* `is_schema_advisor_enabled` - (Optional) Flag that indicates whether to enable Schema Advisor for the project. If enabled, you receive customized recommendations to optimize your data model and enhance performance. Disable this setting to disable schema suggestions in the [Performance Advisor](https://www.mongodb.com/docs/atlas/performance-advisor/#std-label-performance-advisor) and the [Data Explorer](https://www.mongodb.com/docs/atlas/atlas-ui/#std-label-atlas-ui). By default, this flag is set to true.
* `region_usage_restrictions` - (Optional - set value to GOV_REGIONS_ONLY) Designates that this project can be used for government regions only.  If not set the project will default to standard regions.   You cannot deploy clusters across government and standard regions in the same project. AWS is the only cloud provider for AtlasGov.  For more information see [MongoDB Atlas for Government](https://www.mongodb.com/docs/atlas/government/api/#creating-a-project).

### Tags

 ```terraform
 tags = {
    Owner       = "Terraform"
    Environment = "Example"
    Team        = "tf-experts"
  }
  
  lifecycle {
    ignore_changes = [
      tags["CostCenter"] # useful if `CostCenter` key is managed outside terraform
    ]
  }
```

Key-value pairs between 1 to 255 characters in length for tagging and categorizing the project.

To learn more, see [Resource Tags](https://www.mongodb.com/docs/atlas/tags/).

### Teams
Teams attribute is optional

~> **NOTE:** Atlas limits the number of users to a maximum of 100 teams per project and a maximum of 250 teams per organization.

* `team_id` - (Required) The unique identifier of the team you want to associate with the project. The team and project must share the same parent organization.

* `role_names` - (Required) Each string in the array represents a project role you want to assign to the team. Every user associated with the team inherits these roles. You must specify an array even if you are only associating a single role with the team. The [MongoDB Documentation](https://www.mongodb.com/docs/atlas/reference/user-roles/#organization-roles) describes the roles a user can have.

~> **NOTE:** Project created by API Keys must belong to an existing organization.

### Limits
`limits` allows one to configure a variety of limits to a Project. The limits attribute is optional.

* `name` - (Required) Human-readable label that identifies this project limit. See [Project Limit Documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Projects/operation/setProjectLimit) under `limitName` parameter to find all the limits that can be defined.

* `value` - (Required) Amount to set the limit to. Use the [Project Limit Documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Projects/operation/setProjectLimit) under `limitName` parameter to verify the override limits. 


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project id.
* `created` - The ISO-8601-formatted timestamp of when Atlas created the project.
* `cluster_count` - The number of Atlas clusters deployed in the project.
* `ip_addresses` - IP addresses in a project categorized by services. See [IP Addresses](#ip-addresses).


### IP Addresses

* `services.clusters.#.cluster_name` - Human-readable label that identifies the cluster.
* `services.clusters.#.inbound` - List of inbound IP addresses associated with the cluster. If your network allows outbound HTTP requests only to specific IP addresses, you must allow access to the following IP addresses so that your application can connect to your Atlas cluster.
* `services.clusters.#.outbound` - List of outbound IP addresses associated with the cluster. If your network allows inbound HTTP requests only from specific IP addresses, you must allow access from the following IP addresses so that your Atlas cluster can communicate with your webhooks and KMS.

## Import

Project must be imported using project ID, e.g.

```
$ terraform import mongodbatlas_project.my_project 5d09d6a59ccf6445652a444a
```
For more information see: [MongoDB Atlas Admin API Projects](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Projects) and [MongoDB Atlas Admin API Teams](https://docs.atlas.mongodb.com/reference/api/teams/) Documentation for more information.
