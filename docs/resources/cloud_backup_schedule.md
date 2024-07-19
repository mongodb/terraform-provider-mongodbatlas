# Resource: mongodbatlas_cloud_backup_schedule

`mongodbatlas_cloud_backup_schedule` provides a cloud backup schedule resource. The resource lets you create, read, update and delete a cloud backup schedule.

-> **NOTE** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot modify the backup schedule for an individual cluster below the minimum requirements set in the Backup Compliance Policy.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

-> **NOTE:** When creating a backup schedule you **must either** use the `depends_on` clause to indicate the cluster to which it refers **or** specify the values of `project_id` and `cluster_name` as reference of the cluster resource (e.g. `cluster_name = mongodbatlas_advanced_cluster.my_cluster.name` - see the example below). Failure in doing so will result in an error when executing the plan.

In the Terraform MongoDB Atlas Provider 1.0.0 we have re-architected the way in which Cloud Backup Policies are manged with Terraform to significantly reduce the complexity. Due to this change we've provided multiple examples below to help express how this new resource functions.


## Example Usage - Create a Cluster with 2 Policies Items

You can create a new cluster with `cloud_backup` enabled and then immediately overwrite the default cloud backup policy that Atlas creates by default at the same time with this example.

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_backup_schedule resource

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


  // This will now add the desired policy items to the existing mongodbatlas_cloud_backup_schedule resource
  policy_item_hourly {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 1
  }
  policy_item_daily {
    frequency_interval = 1
    retention_unit     = "days"
    retention_value    = 2
  }
}
```

## Example Usage - Create a Cluster with Cloud Backup Enabled but No Policy Items

You can enable `cloud_backup` in the Cluster resource and then use the `cloud_backup_schedule` resource with no policy items to remove the default policy that Atlas creates when you enable Cloud Backup. This allows you to then create a policy when you are ready to via Terraform.

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_backup_schedule resource

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

}
```

## Example Usage - Add 4 Policies Items To A Cluster With Cloud Backup Previously Enabled but with No Policy Items

If you followed the example to Create a Cluster with Cloud Backup Enabled but No Policy Items and then want to add policy items later to the `mongodbatlas_cloud_backup_schedule` this example shows how.

The cluster already exists with `cloud_backup` enabled
```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_backup_schedule resource

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
  
  // This will now add the desired policy items to the existing mongodbatlas_cloud_backup_schedule resource
  policy_item_hourly {
    frequency_interval = 1        #accepted values = 1, 2, 4, 6, 8, 12 -> every n hours
    retention_unit     = "days"
    retention_value    = 1
  }
  policy_item_daily {
    frequency_interval = 1        #accepted values = 1 -> every 1 day
    retention_unit     = "days"
    retention_value    = 2
  }
  policy_item_weekly {
    frequency_interval = 4        # accepted values = 1 to 7 -> every 1=Monday,2=Tuesday,3=Wednesday,4=Thursday,5=Friday,6=Saturday,7=Sunday day of the week
    retention_unit     = "weeks"
    retention_value    = 3
  }
  policy_item_monthly {
    frequency_interval = 5        # accepted values = 1 to 28 -> 1 to 28 every nth day of the month  
                                  # accepted values = 40 -> every last day of the month
    retention_unit     = "months"
    retention_value    = 4
  }
  policy_item_yearly {
    frequency_interval = 1        # accepted values = 1 to 12 -> 1st day of nth month  
    retention_unit     = "years"
    retention_value    = 1
  }

}
```

## Example Usage - Create a Cluster with Cloud Backup Enabled with Snapshot Distribution

You can enable `cloud_backup` in the Cluster resource and then use the `cloud_backup_schedule` resource with a basic policy for Cloud Backup.

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "clusterTest"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_backup_schedule resource

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
    replication_spec_id = mongodbatlas_advanced_cluster.my_cluster.replication_specs.*.id[0]
    should_copy_oplogs = false
  }

}
```
## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot backup policy you want to retrieve.
* `reference_hour_of_day` - (Optional) UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
* `reference_minute_of_hour` - (Optional) UTC Minutes after `reference_hour_of_day` that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
* `restore_window_days` - (Optional) Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
* `update_snapshots` - (Optional) Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously. 
  
  **Note** This parameter does not return updates on return from API, this is a feature of the MongoDB Atlas Admin API itself and not Terraform.  For more details about this resource see [Cloud Backup Schedule](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cloud-Backups/operation/getBackupSchedule).

* `policy_item_hourly` - (Optional) Hourly policy item
* `policy_item_daily` - (Optional) Daily policy item
* `policy_item_weekly` - (Optional) Weekly policy item
* `policy_item_monthly` - (Optional) Monthly policy item
* `policy_item_yearly` - (Optional) Yearly policy item
* `auto_export_enabled` - Flag that indicates whether automatic export of cloud backup snapshots to the AWS bucket is enabled. Value can be one of the following:
	* true - enables automatic export of cloud backup snapshots to the AWS bucket
 	* false - disables automatic export of cloud backup snapshots to the AWS bucket (default)
* `use_org_and_group_names_in_export_prefix` - Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots. To learn more about the metadata files that Atlas uploads, see [Export Cloud Backup Snapshot](https://www.mongodb.com/docs/atlas/backup/cloud-backup/export/#std-label-cloud-provider-snapshot-export).
### Export
* `export_bucket_id` - Unique identifier of the mongodbatlas_cloud_backup_snapshot_export_bucket export_bucket_id value.
* `frequency_type` - Frequency associated with the export snapshot item.

### Policy Item Hourly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For hourly policies, the frequency type is defined as `hourly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (hourly in this case). The supported values for hourly policies are `1`, `2`, `4`, `6`, `8` or `12` hours. Note that `12` hours is the only accepted value for NVMe clusters.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.

### Policy Item Daily
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For daily policies, the frequency type is defined as `daily`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (daily in this case). The only supported value for daily policies is `1` day.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.  Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the hourly policy item specifies a retention of two days, the daily retention policy must specify two days or greater.

### Policy Item Weekly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For weekly policies, the frequency type is defined as `weekly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (weekly in this case). The supported values for weekly policies are `1` through `7`, where `1` represents Monday and `7` represents Sunday.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Weekly policy must have retention of at least 7 days or 1 week. Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the daily policy item specifies a retention of two weeks, the weekly retention policy must specify two weeks or greater.

### Policy Item Monthly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For monthly policies, the frequency type is defined as `monthly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (monthly in this case). The supported values for weekly policies are 
  * `1` through `28` where the number represents the day of the month i.e. `1` is the first of the month and `5` is the fifth day of the month.
  * `40` represents the last day of the month (depending on the month).
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Monthly policy must have retention days of at least 31 days or 5 weeks or 1 month. Note that for less frequent policy items, Atlas requires that you specify a retention period greater than or equal to the retention period specified for more frequent policy items. For example: If the weekly policy item specifies a retention of two weeks, the montly retention policy must specify two weeks or greater.

### Policy Item Yearly
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For yearly policies, the frequency type is defined as `yearly`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (yearly in this case). The supported values for yearly policies are 
  * `1` through `12` the first day of the month where the number represents the month, i.e. `1` is January and `12` is December.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`. Yearly policy must have retention of at least 1 year.

### Snapshot Distribution
*
* `cloud_provider` - (Required) Human-readable label that identifies the cloud provider that stores the snapshot copy. i.e. "AWS" "AZURE" "GCP"
* `frequencies` - (Required) List that describes which types of snapshots to copy. i.e. "HOURLY" "DAILY" "WEEKLY" "MONTHLY" "ON_DEMAND"
* `region_name` - (Required) Target region to copy snapshots belonging to replicationSpecId to. Please supply the 'Atlas Region' which can be found under https://www.mongodb.com/docs/atlas/reference/cloud-providers/ 'regions' link
* `replication_spec_id` -(Required) Unique 24-hexadecimal digit string that identifies the replication object for a zone in a cluster. For global clusters, there can be multiple zones to choose from. For sharded clusters and replica set clusters, there is only one zone in the cluster. To find the Replication Spec Id, consult the replicationSpecs array returned from [Return One Multi-Cloud Cluster in One Project](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/getCluster).
* `should_copy_oplogs` - (Required) Flag that indicates whether to copy the oplogs to the target region. You can use the oplogs to perform point-in-time restores.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - Timestamp in the number of seconds that have elapsed since the UNIX epoch when Atlas takes the next snapshot.
* `id_policy` - Unique identifier of the backup policy.

## Import

Cloud Backup Schedule entries can be imported using project_id and cluster_name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_schedule.test 5d0f1f73cf09a29120e173cf-MyClusterTest
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/)
