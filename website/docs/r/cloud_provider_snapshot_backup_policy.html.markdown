---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_provider_snapshot_backup_policy"
sidebar_current: "docs-mongodbatlas-resource-cloud-provider-snapshot-backup-policy"
description: |-
    Provides a Cloud Provider Snapshot Backup Policy resource.
---

# mongodbatlas_cloud_provider_snapshot_backup_policy

`mongodbatlas_cloud_provider_snapshot_backup_policy` provides a resource that enables you to view and modify the snapshot schedule and retention settings for an Atlas cluster with Cloud Provider Snapshots enabled.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  provider_backup_enabled     = true // must be enabled in order to use cloud_provider_snapshot_backup_policy resource
  provider_disk_iops          = 100
  provider_encrypt_ebs_volume = false
}

resource "mongodbatlas_cloud_provider_snapshot_backup_policy" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4


  policies {
    id = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.id

    policy_item {
      id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.0.id
      frequency_interval = 1
      frequency_type     = "hourly"
      retention_unit     = "days"
      retention_value    = 1
    }
    policy_item {
      id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.1.id
      frequency_interval = 1
      frequency_type     = "daily"
      retention_unit     = "days"
      retention_value    = 2
    }
    policy_item {
      id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
      frequency_interval = 4
      frequency_type     = "weekly"
      retention_unit     = "weeks"
      retention_value    = 3
    }
    policy_item {
      id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id
      frequency_interval = 5
      frequency_type     = "monthly"
      retention_unit     = "months"
      retention_value    = 4
    }
  }
}
```

~> **IMPORTANT:** If setting to true for an existing cluster or imported cluster be sure to run terraform refresh after applying to enable modification of the Cloud Provider Snapshot Backup Policy going forward.

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot backup policy you want to retrieve.
* `reference_hour_of_day` - (Optional) UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
* `reference_minute_of_hour` - (Optional) UTC Minutes after referenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
* `restore_window_days` - (Optional) Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
* `update_snapshots` - (Optional) Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
* `policies` - (Required) Contains a document for each backup policy item in the desired updated backup policy.
* `policies.#.id` - (Required) Unique identifier of the backup policy that you want to update.

#### Policy Item
* `policies.#.policy_item` - (Required) Array of backup policy items.
* `policies.#.policy_item.#.id` - (Required) Unique identifier of the backup policy item. `policies.#.policy_item.#.id` is a value obtained via the mongodbatlas_cluster resource. provider_backup_enabled of the mongodbatlas_cluster resource must be set to true. See the example above for how to refer to the mongodbatlas_cluster resource forpolicies.#.policy_item.#.id
* `policies.#.policy_item.#.frequency_interval` - (Required) Desired frequency of the new backup policy item specified by frequencyType.
* `policies.#.policy_item.#.frequency_type` - (Required) Frequency associated with the backup policy item. One of the following values: hourly, daily, weekly or monthly.
* `policies.#.policy_item.#.retention_unit` - (Required) Scope of the backup policy item: days, weeks, or months.
* `policies.#.policy_item.#.retention_value` - (Required) Value to associate with retentionUnit.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - Timestamp in the number of seconds that have elapsed since the UNIX epoch when Atlas takes the next snapshot.

## Import

Cloud Provider Snapshot Backup Policy entries can be imported using project project_id, cluster_name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_snapshot_backup_policy.test 5d0f1f73cf09a29120e173cf-MyClusterTest
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-provider-snapshot-schedule-modify-one/)