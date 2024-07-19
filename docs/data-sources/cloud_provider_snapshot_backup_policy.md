---
subcategory: "Deprecated"    
---

 **WARNING:** This data source is deprecated, use `mongodbatlas_cloud_backup_schedule`
 **Note:** This resource have now been fully deprecated as part of v1.10.0 release

# Data Source: mongodbatlas_cloud_provider_snapshot_backup_policy

`mongodbatlas_cloud_provider_snapshot_backup_policy` provides a Cloud Backup Snapshot Backup Policy datasource. An Atlas Cloud Backup Snapshot Policy provides the current snapshot schedule and retention settings for the cluster. 

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

resource "mongodbatlas_cloud_provider_snapshot_backup_policy" "test" {
  project_id   = mongodbatlas_advanced_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_advanced_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4


  policies {
    id = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.id

    policy_item {
      id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.0.id
      frequency_interval = 1
      frequency_type     = "hourly"
      retention_unit     = "days"
      retention_value    = 1
    }
    policy_item {
      id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.1.id
      frequency_interval = 1
      frequency_type     = "daily"
      retention_unit     = "days"
      retention_value    = 2
    }
    policy_item {
      id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
      frequency_interval = 4
      frequency_type     = "weekly"
      retention_unit     = "weeks"
      retention_value    = 3
    }
    policy_item {
      id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id
      frequency_interval = 5
      frequency_type     = "monthly"
      retention_unit     = "months"
      retention_value    = 4
    }
  }
}

data "mongodbatlas_cloud_provider_snapshot_backup_policy" "test" {
  project_id   = mongodbatlas_cloud_provider_snapshot_backup_policy.test.project_id
  cluster_name = mongodbatlas_cloud_provider_snapshot_backup_policy.test.cluster_name
}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshots backup policy you want to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - UTC ISO 8601 formatted point in time when Atlas will take the next snapshot.
* `reference_hour_of_day` - UTC Hour of day between 0 and 23 representing which hour of the day that Atlas takes a snapshot.
* `reference_minute_of_hour` - UTC Minute of day between 0 and 59 representing which minute of the referenceHourOfDay that Atlas takes the snapshot.
* `restore_window_days` - Specifies a restore window in days for cloud backup to maintain.

### Policies
* `policies` - A list of policy definitions for the cluster.
* `policies.#.id` - Unique identifier of the backup policy.

#### Policy Item
* `policies.#.policy_item` - A list of specifications for a policy.
* `policies.#.policy_item.#.id` - Unique identifier for this policy item.
* `policies.#.policy_item.#.frequency_interval` - The frequency interval for a set of snapshots.
* `policies.#.policy_item.#.frequency_type` - A type of frequency (hourly, daily, weekly, monthly).
* `policies.#.policy_item.#.retention_unit` - The unit of time in which snapshot retention is measured (days, weeks, months).
* `policies.#.policy_item.#.retention_value` - The number of days, weeks, or months the snapshot is retained.

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/get-all-schedules/)