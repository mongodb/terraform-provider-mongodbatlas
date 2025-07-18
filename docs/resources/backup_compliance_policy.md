# Resource: mongodbatlas_backup_compliance_policy

`mongodbatlas_backup_compliance_policy` provides a resource that enables you to set up a Backup Compliance Policy resource. [Backup Compliance Policy ](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy) prevents any user, regardless of role, from modifying or deleting specific cluster settings, backups, and backup configurations. When enabled, the Backup Compliance Policy will be applied as the minimum policy for all clusters and backups in the project. It can only be disabled by contacting MongoDB support. This feature is only supported for cluster tiers M10+.

When enabled, the Backup Compliance Policy will be applied as the minimum backup policy to all clusters in a project and will protect all existing snapshots. This will prevent any user, regardless of role, from modifying or deleting existing snapshots prior to expiration. Changes made to existing backup policies will only apply to future snapshots.

-> **NOTE:** Groups and projects are synonymous terms. You might find `groupId` in the official documentation.

-> **IMPORTANT NOTE:** Once a Backup Compliance Policy is enabled, no user, regardless of role, can disable the Backup Compliance Policy via Terraform, or any other method, without contacting MongoDB Support. However, if the Backup Compliance Policy is approved for disablement by MongoDB Support, the resource can be removed and the delete operation will successfully disable the Backup Compliance Policy. With Backup Compliance Policy enabled, some resources defined in Terraform will no longer be modifiable. See the full list of [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy). 

-> **NOTE:** With Backup Compliance Policy enabled, cluster backups are retained after a cluster is deleted and backups can be used normally until retention expiration. When the Backup Compliance Policy is not enabled, Atlas deletes the cluster's associated backup snapshots when a cluster is terminated. By default, a Backup Compliance Policy is not enabled. For more details see [Back Up, Restore, and Archive Data](https://www.mongodb.com/docs/atlas/backup-restore-cluster/). 

-> **NOTE:** To delete an Atlas cluster that has an associated `mongodbatlas_cloud_backup_schedule` resource and an enabled Backup Compliance Policy, first instruct Terraform to remove the `mongodbatlas_cloud_backup_schedule` resource from the state and then use Terraform to delete the cluster. To learn more, see [Delete a Cluster with a Backup Compliance Policy](../guides/delete-cluster-with-backup-compliance-policy.md).


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
      region_name   = var.region
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
    retention_value    = 7
  }
  policy_item_daily {
    frequency_interval = 1        #accepted values = 1 -> every 1 day
    retention_unit     = "days"
    retention_value    = 7
  }
  policy_item_weekly {
    frequency_interval = 1        # accepted values = 1 to 7 -> every 1=Monday,2=Tuesday,3=Wednesday,4=Thursday,5=Friday,6=Saturday,7=Sunday day of the week
    retention_unit     = "weeks"
    retention_value    = 4
  }
  policy_item_monthly {
    frequency_interval = 1        # accepted values = 1 to 28 -> 1 to 28 every nth day of the month  
                                  # accepted values = 40 -> every last day of the month
    retention_unit     = "months"
    retention_value    = 12
  }
  policy_item_yearly {
    frequency_interval = 1        # accepted values = 1 to 12 -> 1st day of nth month  
    retention_unit     = "years"
    retention_value    = 1
  }

}

data "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cloud_backup_schedule.test.project_id
  cluster_name = mongodbatlas_cloud_backup_schedule.test.cluster_name
}

data "mongodbatlas_backup_compliance_policy" "backup_policy" {
  project_id = mongodbatlas_cloud_backup_schedule.test.project_id
}

resource "mongodbatlas_backup_compliance_policy" "backup_policy" {
  project_id                 = "<PROJECT-ID>"
  authorized_email           = "user@email.com"
  authorized_user_first_name = "First"
  authorized_user_last_name  = "Last"
  copy_protection_enabled    = false
  pit_enabled                = false
  encryption_at_rest_enabled = false

  restore_window_days = 7

  on_demand_policy_item {
		  frequency_interval = 1
		  retention_unit     = "days"
		  retention_value    = 3
		}
		
		policy_item_hourly {
			frequency_interval = 1
			retention_unit     = "days"
			retention_value    = 7
		  }
	  
		policy_item_daily {
			frequency_interval = 1
			retention_unit     = "days"
			retention_value    = 7
		  }
	  
		  policy_item_weekly {
			frequency_interval = 1
			retention_unit     = "weeks"
			retention_value    = 4
		  }
	  
		  policy_item_monthly {
			frequency_interval = 1
			retention_unit     = "months"
			retention_value    = 12
		  }

	          policy_item_yearly {
	            frequency_interval = 1
	            retention_unit     = "years"
	            retention_value    = 1
	          }
      
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `authorized_email` - (Required) Email address of a security or legal representative for the Backup Compliance Policy who is authorized to update the Backup Compliance Policy settings.
* `authorized_user_first_name` - (Required) First name of the user who authorized to update the Backup Compliance Policy settings.
* `authorized_user_last_name` - (Required) Last name of the user who authorized to update the Backup Compliance Policy settings.
* `copy_protection_enabled` - Flag that indicates whether to enable additional backup copies for the cluster. If unspecified, this value defaults to false.
* `pit_enabled` - Flag that indicates whether the cluster uses Continuous Cloud Backups with a Backup Compliance Policy. If unspecified, this value defaults to false.
* `encryption_at_rest_enabled` - Flag that indicates whether Encryption at Rest using Customer Key Management is required for all clusters with a Backup Compliance Policy. If unspecified, this value defaults to false.
* `restore_window_days` - Number of previous days that you can restore back to with Continuous Cloud Backup with a Backup Compliance Policy. You must specify a positive, non-zero integer, and the maximum retention window can't exceed the hourly retention time. This parameter applies only to Continuous Cloud Backups with a Backup Compliance Policy.
* `on_demand_policy_item` - Specifications for [on-demand policy](#on-demand-policy-item).
* `policy_item_hourly` - Scheduled policy using an hourly frequency type, see [block fields](#policy-item-hourly).
* `policy_item_daily` - Scheduled policy using a daily frequency type, see [block fields](#policy-item-daily).
* `policy_item_weekly` - Scheduled policy using a weekly frequency type, see [block fields](#policy-item-weekly).
* `policy_item_monthly` - Scheduled policy using a monthly frequency type, see [block fields](#policy-item-monthly).
* `policy_item_yearly` - Scheduled policy using a yearly frequency type, see [block fields](#policy-item-yearly).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*  `state` - Label that indicates the state of the Backup Compliance Policy settings. MongoDB Cloud ignores this setting when you enable or update the Backup Compliance Policy settings.
* `updated_date` - ISO 8601 timestamp format in UTC that indicates when the user updated the Data Protection Policy settings. MongoDB Cloud ignores this setting when you enable or update the Backup Compliance Policy settings.
* `updated_user` - Email address that identifies the user who updated the Backup Compliance Policy settings. MongoDB Cloud ignores this email setting when you enable or update the Backup Compliance Policy settings.

### On Demand Policy Item
* `id` - Unique identifier of the backup policy item.
* `frequency_type` - Frequency associated with the backup policy item. For hourly policies, the frequency type is defined as `ondemand`. Note that this is a read-only value and not required in plan files - its value is implied from the policy resource type.
* `frequency_interval` - Desired frequency of the new backup policy item specified by `frequency_type` (hourly in this case). The supported values for hourly policies are `1`, `2`, `4`, `6`, `8` or `12` hours. Note that `12` hours is the only accepted value for NVMe clusters.
* `retention_unit` - Scope of the backup policy item: `days`, `weeks`, `months`, or `years`.
* `retention_value` - Value to associate with `retention_unit`.

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
  * `1` through `28` where the number represents the day of the month (i.e. `1` is the first of the month and `5` is the fifth day of the month).
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

## Import

Backup Compliance Policy entries can be imported using project project_id  in the format `project_id`, e.g.

```
$ terraform import mongodbatlas_backup_compliance_policy.backup_policy 5d0f1f73cf09a29120e173cf
```

For more information see: [MongoDB Atlas API Reference](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cloud-Backups/operation/updateDataProtectionSettings) and [Backup Compliance Policy Prohibited Actions](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#prohibited-actions).


