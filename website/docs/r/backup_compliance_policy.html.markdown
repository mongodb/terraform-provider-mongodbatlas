---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: backup_compliance_policy"
sidebar_current: "docs-mongodbatlas-resource-backup-compliance-policy"
description: |-
    Provides a Backup Compliance Policy resource.
---
# Resource: mongodbatlas_backup_compliance_policy

`mongodbatlas_backup_compliance_policy` provides a resource that enables you to setup a Backup Compliance Policy resource. Prevent any user, regardless of role, from modifying or deleting specific cluster configurations and backups. When enabled, the Backup Compliance Policy will be applied as the minimum policy for all clusters and can only be disabled by contacting MongoDB support. Only supported for clusters M10 or higher.

When enabled, Backup Compliance Policy will be applied as the minimum backup policy to all clusters in your Project and will protect all existing snapshots. This will prevent any user, regardless of role, from modifying or deleting existing snapshots prior to expiration. Changes made to existing backup compliance policies will only apply to future snapshots.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** If you enable a Backup Compliance Policy, no user, regardless of role, can disable the Backup Compliance Policy without contacting MongoDB support, delete a backup snapshot, decrease the retention time for a snapshot after it's taken, disable Cloud Backup, modify the backup policy for an individual cluster below the minimum requirements set in the Backup Compliance Policy, or delete the Atlas project if any snapshots exist. For full list of impacts and more details see [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

## Example Usage

```terraform
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"


  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  provider_backup_enabled     = true // enable cloud backup snapshots
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

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

}

data "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cloud_backup_schedule.test.project_id
  cluster_name = mongodbatlas_cloud_backup_schedule.test.cluster_name
}

data "mongodbatlas_backup_compliance_policy" "backup_policy" {
  project_id = mongodbatlas_cloud_backup_schedule.test.id
}

resource "mongodbatlas_backup_compliance_policy" "backup_policy" {
  project_id                 = "<PROJECT-ID>"
  authorized_email           = "user@email.com"
  copy_protection_enabled    = false
  pit_enabled                = false
  encryption_at_rest_enabled = false

  restore_window_days = 7

  on_demand_policy_item {
		  frequency_interval = 0
		  retention_unit     = "days"
		  retention_value    = 3
		}
		
		policy_item_hourly {
			frequency_interval = 6
			retention_unit     = "days"
			retention_value    = 7
		  }
	  
		policy_item_daily {
			frequency_interval = 0
			retention_unit     = "days"
			retention_value    = 7
		  }
	  
		  policy_item_weekly {
			frequency_interval = 0
			retention_unit     = "weeks"
			retention_value    = 4
		  }
	  
		  policy_item_monthly {
			frequency_interval = 0
			retention_unit     = "months"
			retention_value    = 12
		  }

}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `authorized_email` - (Email address of the user who authorized to updated the Backup Compliance Policy settings.
* `copy_protection_enabled` - Flag that indicates whether to enable additional backup copies for the cluster. If unspecified, this value defaults to false.
* `pit_enabled` - Flag that indicates whether the cluster uses Continuous Cloud Backups with a Backup Compliance Policy. If unspecified, this value defaults to false.
* `encryption_at_rest_enabled` - Flag that indicates whether Encryption at Rest using Customer Key Management is required for all clusters with a Backup Compliance Policy. If unspecified, this value defaults to false.
* `reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.
* `restore_window_days` - Number of previous days that you can restore back to with Continuous Cloud Backup with a Backup Compliance Policy. You must specify a positive, non-zero integer, and the maximum retention window can't exceed the hourly retention time. This parameter applies only to Continuous Cloud Backups with a Backup Compliance Policy.
*  `state` - Label that indicates the state of the Backup Compliance Policy settings. MongoDB Cloud ignores this setting when you enable or update the Backup Compliance Policy settings.
* `updated_date` - ISO 8601 timestamp format in UTC that indicates when the user updated the Data Protection Policy settings. MongoDB Cloud ignores this setting when you enable or update the Backup Compliance Policy settings.
* `updated_user` - Email address that identifies the user who updated the Backup Compliance Policy settings. MongoDB Cloud ignores this email setting when you enable or update the Backup Compliance Policy settings.

### On Demand Policy Item
* `on_demand_policy_item.0.policy_item` - A list of specifications for a policy.
* `on_demand_policy_item.#.policy_item.#.id` - Unique identifier for this policy item.
* `on_demand_policy_item.#.policy_item.#.frequency_interval` - The frequency interval for a set of snapshots.
* `on_demand_policy_item.#.policy_item.#.frequency_type` - A type of frequency (hourly, daily, weekly, monthly).
* `on_demand_policy_item.#.policy_item.#.retention_unit` - The unit of time in which snapshot retention is measured (days, weeks, months).
* `policies.#.policy_item.#.retention_value` - The number of days, weeks, or months the snapshot is retained.
* 
#### Scheduled Policy Items
* `scheduled_policy_items.#.policy_item` - A list of specifications for a policy.
* `scheduled_policy_items.#.policy_item.#.id` - Unique identifier for this policy item.
* `scheduled_policy_items.#.policy_item.#.frequency_interval` - The frequency interval for a set of snapshots.
* `scheduled_policy_items.#.policy_item.#.frequency_type` - A type of frequency (hourly, daily, weekly, monthly).
* `scheduled_policy_items.#.policy_item.#.retention_unit` - The unit of time in which snapshot retention is measured (days, weeks, months).
* `scheduled_policy_items.#.policy_item.#.retention_value` - The number of days, weeks, or months the snapshot is retained.

## Import

Backup Compliance Policy entries can be imported using project project_id  in the format `project_id`, e.g.

```
$ terraform import mongodbatlas_backup_compliance_policy.backup_policy 5d0f1f73cf09a29120e173cf
```

For more information see: [MongoDB Atlas API Reference](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Cloud-Backups/operation/updateDataProtectionSettings) and [Backup Compliance Policy Prohibited Actions](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#prohibited-actions)



