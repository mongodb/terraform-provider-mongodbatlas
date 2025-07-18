# Data Source: mongodbatlas_cloud_backup_schedule

`mongodbatlas_cloud_backup_schedule` provides a Cloud Backup Schedule datasource. An Atlas Cloud Backup Schedule provides the current cloud backup schedule for the cluster. 

-> **NOTE:** To delete an Atlas cluster that has an associated `mongodbatlas_cloud_backup_schedule` resource and an enabled Backup Compliance Policy, first instruct Terraform to remove the `mongodbatlas_cloud_backup_schedule` resource from the state and then use Terraform to delete the cluster. To learn more, see [Delete a Cluster with a Backup Compliance Policy](../guides/delete-cluster-with-backup-compliance-policy.md).

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud backup snapshots

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_CENTRAL_1"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_advanced_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_advanced_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4

  policy_item_daily {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 14
  }

  copy_settings {
    cloud_provider = "AWS"
    frequencies = ["HOURLY",
		   "DAILY",
		   "WEEKLY",
		   "MONTHLY",
		   "YEARLY",
		   "ON_DEMAND"]
    region_name = "US_EAST_1"
    zone_id = mongodbatlas_advanced_cluster.my_cluster.replication_specs.*.zone_id[0]
    should_copy_oplogs = false
  }
}

data "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cloud_backup_schedule.test.project_id
  cluster_name = mongodbatlas_cloud_backup_schedule.test.cluster_name
  use_zone_id_for_copy_settings = true
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshots backup policy you want to retrieve.
* `use_zone_id_for_copy_settings` - Set this field to `true` to allow the data source to use the latest schema that populates `copy_settings.#.zone_id` instead of the deprecated `copy_settings.#.replication_spec_id`. These fields also enable you to reference cluster zones using independent shard scaling, which no longer supports `replication_spec.*.id`. To learn more, see the [1.18.0 upgrade guide](../guides/1.18.0-upgrade-guide.md#transition-cloud-backup-schedules-for-clusters-to-use-zones).


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
* `reference_hour_of_day` - UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot.
* `reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the `reference_hour_of_day` that Atlas takes the snapshot.
* `restore_window_days` - Specifies a restore window in days for cloud backup to maintain.
* `id_policy` - Unique identifier of the backup policy.
* `policy_item_hourly` - (Optional) Hourly policy item. See [below](#policy_item_hourly)
* `policy_item_daily` - (Optional) Daily policy item. See [below](#policy_item_daily)
* `policy_item_weekly` - (Optional) Weekly policy item. See [below](#policy_item_weekly)
* `policy_item_monthly` - (Optional) Monthly policy item. See [below](#policy_item_monthly)
* `policy_item_yearly` - (Optional) Yearly policy item. See [below](#policy_item_yearly)
* `auto_export_enabled` - Flag that indicates whether MongoDB Cloud automatically exports Cloud Backup Snapshots to the Export Bucket. Value can be one of the following:
  * true - Enables automatic export of cloud backup snapshots to the Export Bucket.
  * false - Disables automatic export of cloud backup snapshots to the Export Bucket. (default)
* `use_org_and_group_names_in_export_prefix` - Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your bucket after it finishes exporting the snapshots. To learn more about the metadata files that Atlas uploads, see [Export Cloud Backup Snapshot](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/#std-label-cloud-provider-snapshot-export).
* `copy_settings` - List that contains a document for each copy setting item in the desired backup policy. See [below](#copy_settings)
* `export` - Policy for automatically exporting Cloud Backup Snapshots. See [below](#export)

### export
* `export_bucket_id` - Unique identifier of the mongodbatlas_cloud_backup_snapshot_export_bucket export_bucket_id value.
* `frequency_type` - Frequency associated with the export snapshot item.

### policy_item_hourly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For hourly policies, the frequency type is defined as `hourly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (hourly in this case). The supported values for hourly policies are `1`, `2`, `4`, `6`, `8` or `12` hours. Note that `12` hours is the only accepted value for NVMe clusters.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.

### policy_item_daily
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For daily policies, the frequency type is defined as `daily`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (daily in this case). The only supported value for daily policies is `1` day.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.  Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the hourly policy item specifies a retention of two days, the daily retention policy must specify two days or greater.

### policy_item_weekly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For weekly policies, the frequency type is defined as `weekly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (weekly in this case). The supported values for weekly policies are `1` through `7`, where `1` represents Monday and `7` represents Sunday.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Weekly policy must have retention of at least 7 days or 1 week. Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the daily policy item specifies a retention of two weeks, the weekly retention policy must specify two weeks or greater.

### policy_item_monthly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For monthly policies, the frequency type is defined as `monthly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (monthly in this case). The supported values for weekly policies are 
  * `1` through `28` where the number represents the day of the month i.e. `1` is the first of the month and `5` is the fifth day of the month.
  * `40` represents the last day of the month (depending on the month).
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Monthly policy must have retention days of at least 31 days or 5 weeks or 1 month. Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the weekly policy item specifies a retention of two weeks, the montly retention policy must specify two weeks or greater.

### policy_item_yearly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For yearly policies, the frequency type is defined as `yearly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (yearly in this case). The supported values for yearly policies are 
  * `1` through `12` the first day of the month where the number represents the month, i.e. `1` is January and `12` is December.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Yearly policy must have retention of at least 1 year.

### copy_settings
* `cloud_provider` - Human-readable label that identifies the cloud provider that stores the snapshot copy. i.e. "AWS" "AZURE" "GCP"
* `frequencies` - List that describes which types of snapshots to copy. i.e. "HOURLY" "DAILY" "WEEKLY" "MONTHLY" "YEARLY" "ON_DEMAND"
* `region_name` - Target region to copy snapshots belonging to replicationSpecId to. Please supply the 'Atlas Region' which can be found under https://www.mongodb.com/docs/atlas/reference/cloud-providers/ 'regions' link
* `zone_id` - Unique 24-hexadecimal digit string that identifies the zone in a cluster. For global clusters, there can be multiple zones to choose from. For sharded clusters and replica set clusters, there is only one zone in the cluster.
* `replication_spec_id` - Unique 24-hexadecimal digit string that identifies the replication object for a zone in a cluster. For global clusters, there can be multiple zones to choose from. For sharded clusters and replica set clusters, there is only one zone in the cluster. To find the Replication Spec Id, consult the replicationSpecs array returned from [Return One Multi-Cloud Cluster in One Project](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getcluster). **(DEPRECATED)** Use `zone_id` instead. To learn more, see the [1.18.0 upgrade guide](../guides/1.18.0-upgrade-guide.md#transition-cloud-backup-schedules-for-clusters-to-use-zones).
* `should_copy_oplogs` - Flag that indicates whether to copy the oplogs to the target region. You can use the oplogs to perform point-in-time restores.

**Note** The parameter deleteCopiedBackups is not supported in terraform please leverage Atlas Admin API or AtlasCLI instead to manage the lifecycle of backup snaphot copies.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/).
