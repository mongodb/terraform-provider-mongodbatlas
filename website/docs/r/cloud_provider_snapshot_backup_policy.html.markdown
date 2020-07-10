---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: cloud_provider_snapshot_backup_policy"
sidebar_current: "docs-mongodbatlas-resource-cloud-provider-snapshot-backup-policy"
description: |-
    Provides a Cloud Backup Snapshot Policy resource.
---

# mongodbatlas_cloud_provider_snapshot_backup_policy

`mongodbatlas_cloud_provider_snapshot_backup_policy` provides a resource that enables you to view and modify the snapshot schedule and retention settings for an Atlas cluster with Cloud Backup enabled.  A default policy is created automatically when Cloud Backup is enabled for the cluster.

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
# Examples Modifying Polices
When Cloud Backup is enabled for a cluster MongoDB Atlas automatically creates a default Cloud Backup schedule for the cluster with four policy items; hourly, daily, weekly, and monthly. Because of this default creation this provider automatically saves the Cloud Backup Snapshot Policy into the Terraform state. If the default works well for you then you do not need to do anything other than create a cluster with Cloud Backup enabled and your Terraform state will have this information if you need it. However,  if you want the policy to be different than the default simply follow the next examples.

## Example Usage - Create a Cluster and Modify the 4 Default Policies Simultaneously
This cluster has already been created and is here as an example

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
      id                 = 5f0747cad187d8609a72f546
      frequency_interval = 4
      frequency_type     = "weekly"
      retention_unit     = "days"
      retention_value    = 3
    }
  }
}
```

-> **NOTE:** This is the id MongoDB Atlas returns for the policy item we want to keep. Here it is hard coded because you need to either use the actual value from the Terraform state or look to map the policy item you want to keep to it's placement in the state file array that was imported in when the cluster was originally created.

Summarized: to use the state file value for a policy item you need to determine the array placement # of the same frequency_type you want to keep. With this Terraform configuration the anothers policy items will be removed.

~> **IMPORTANT:** For example in the state file we are using the weekly policy is the third policy in the array so it could be referred to with `mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id` instead of the hard coded value, but it's not recommended when the cluster presents changes or make `terraform refresh` due to once this will applied, the cluster state will remove the rest of the items, so the posicion of the array will change to position 0.

## Example Usage - Create a Cluster and Modify 3 Default Policies and Remove 1 Default Policy Simultaneously

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

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
    #   frequency_interval = 4
    #   frequency_type     = "weekly"
    #   retention_unit     = "weeks"
    #   retention_value    = 3
    # }

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

-> **NOTE:** (See text in first example for more details on the default.) If you want the Cloud Backup Snapshot Policy to vary in the number of policies from the default when creating the cluster, perhaps you want to remove one policy item and modify the remaining three, simply follow this example here to remove a policy and modify three.

~> **IMPORTANT:** If we decide to remove item 2 as our above example marked with `#` we need to consider that once the cluster being modified or makes a `terraform refresh` the item 2 will be replaced with the 3, so it could cause inconsistency. We recommend using hardcoded id value to handle these situations. (See text in the first example for more details on it)



## Example Usage - Remove 3 Default Policies Items After the Cluster Has Already Been Created

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
      id                 = 5f0747cad187d8609a72f546
      frequency_interval = 5
      frequency_type     = "monthly"
      retention_unit     = "months"
      retention_value    = 4
    }
  }
}
```

-> **NOTE:** (See text in first example for more details on the default.) If you want the Cloud Backup Snapshot Policy to vary in number of policies for a cluster that was already created/imported, perhaps you want to remove three policy items and modify the remaining policy, simply follow the above example here.

~> **IMPORTANT:** Note in this example we decided to remove the first 3 items so we can't use `mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id` this sentence to retrieve the monthly id value of the cluster state due to once the cluster being modified or makes a `terraform refresh` will cause that the three items will remove from the state, so we will get an error due to the index 3 doesn't exists any more and our monthly policy item is moved to the first place of the array.(See text in the first example for more details on it)

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot backup policy you want to retrieve.
* `reference_hour_of_day` - (Optional) UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
* `reference_minute_of_hour` - (Optional) UTC Minutes after referenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
* `restore_window_days` - (Optional) Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
* `update_snapshots` - (Optional) Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
* `policies` - (Required) Contains a document for each backup policy item in the desired updated backup policy.
* `policies.#.id` - (Required) Unique identifier of the backup policy that you want to update. policies.#.id is a value obtained via the mongodbatlas_cluster resource. provider_backup_enabled of the mongodbatlas_cluster resource must be set to true. See the example above for how to refer to the mongodbatlas_cluster resource for policies.#.id

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

Cloud Backup Snapshot Policy entries can be imported using project project_id and cluster_name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cloud_provider_snapshot_backup_policy.test 5d0f1f73cf09a29120e173cf-MyClusterTest
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/schedule/modify-one-schedule/)