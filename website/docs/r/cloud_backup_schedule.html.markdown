---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_backup_schedule"
sidebar_current: "docs-mongodbatlas-resource-cloud-backup-schedule"
description: |-
    Provides a Cloud Backup Schedule resource.
---

# mongodbatlas_cloud_backup_schedule

`mongodbatlas_cloud_backup_schedule` provides a cloud backup schedule resource. The resource lets you create, read, update and delete policies items.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage - Create a Cluster with no policies (no policies will get the default Policies Items)

```hcl
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
  provider_disk_iops          = 100
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4

}
```

~> **IMPORTANT:**   `policies.#.policy_item.#.id` is obtained when the cluster is created. 

## Example Usage - Create a Cluster with 2 policies (It will overwrite the default Policies Items which is usually 4 Policies Items)

```hcl
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
  provider_disk_iops          = 100
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
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
  }
}
```

## Example Usage - Create a cluster with policies and update the existent policies

### First step - Create a cluster with policies(will overwrite the default policies items)
```hcl
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
  provider_disk_iops          = 100
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
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

### Second step(Update) - Update policies by removing 3 policies and keep one policy

-> **NOTE:** In this example we decided to remove the first 3 items so we can't use `mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id` to retrieve the monthly id value of the cluster state due to once the cluster being modified or makes a `terraform refresh` will cause that the three items will remove from the state, so we will get an error due to the index 3 doesn't exists any more and our monthly policy item is moved to the first place of the array.  So we use `5f0747cad187d8609a72f546`, which is an example of an id MongoDB Atlas returns for the policy item we want to keep. Here it is hard coded because you need to either use the actual value from the Terraform state or look to map the policy item you want to keep to it's current placement in the state file array.

```hcl
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "<PROJECT-ID>"
  name         = "clusterTest"
  disk_size_gb = 5

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup     = true // must be enabled in order to use cloud_backup_schedule resource
  provider_disk_iops          = 100
}

resource "mongodbatlas_cloud_backup_schedule" "test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  reference_hour_of_day    = 3
  reference_minute_of_hour = 45
  restore_window_days      = 4
  
  policies {
    id = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.id

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.0.id
    #   frequency_interval = 1
    #   frequency_type     = "hourly"
    #   retention_unit     = "days"
    #   retention_value    = 1
    # }

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.1.id
    #   frequency_interval = 1
    #   frequency_type     = "daily"
    #   retention_unit     = "days"
    #   retention_value    = 2
    # }

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
    #   frequency_interval = 4
    #   frequency_type     = "weekly"
    #   retention_unit     = "weeks"
    #   retention_value    = 3
    # }

    policy_item {
      id                 = "5f0747cad187d8609a72f546"
      frequency_interval = 5
      frequency_type     = "monthly"
      retention_unit     = "months"
      retention_value    = 4
    }
  }

}
```

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot backup policy you want to retrieve.
* `reference_hour_of_day` - (Optional) UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
* `reference_minute_of_hour` - (Optional) UTC Minutes after referenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
* `restore_window_days` - (Optional) Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
* `update_snapshots` - (Optional) Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
* `policies` - (Optional) Contains a document for each backup policy item in the desired updated backup policy.
* `policies.#.id` - (Optional) Unique identifier of the backup policy that you want to update. policies.#.id is a value obtained via the mongodbatlas_cluster resource. cloud_backup of the mongodbatlas_cluster resource must be set to true. See the example above for how to refer to the mongodbatlas_cluster resource for policies.#.id

#### Policy Item
* `policies.#.policy_item` - (Optional) Array of backup policy items.
* `policies.#.policy_item.#.id` - (Optional) Unique identifier of the backup policy item. `policies.#.policy_item.#.id` is a value obtained via the mongodbatlas_cluster resource. `cloud_backup` of the mongodbatlas_cluster resource must be set to true. See the example above for how to refer to the mongodbatlas_cluster resource for policies.#.policy_item.#.id . **NOTE** If not specified, it might create a policy item if the policies items are empty.
* `policies.#.policy_item.#.frequency_interval` - (Required) Desired frequency of the new backup policy item specified by frequencyType.
* `policies.#.policy_item.#.frequency_type` - (Required) Frequency associated with the backup policy item. One of the following values: hourly, daily, weekly or monthly.
* `policies.#.policy_item.#.retention_unit` - (Required) Scope of the backup policy item: days, weeks, or months.
* `policies.#.policy_item.#.retention_value` - (Required) Value to associate with retentionUnit.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cluster_id` - Unique identifier of the Atlas cluster.
* `next_snapshot` - Timestamp in the number of seconds that have elapsed since the UNIX epoch when Atlas takes the next snapshot.

## Import

Cloud Backup Schedule entries can be imported using project_id and cluster_name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_schedule.test 5d0f1f73cf09a29120e173cf-MyClusterTest
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/)