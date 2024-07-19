---
subcategory: "Deprecated"    
---

**WARNING:** This resource is deprecated, use `mongodbatlas_cloud_backup_schedule`
**Note:** This resource have now been fully deprecated as part of v1.10.0 release

# Resource: mongodbatlas_cloud_provider_snapshot_backup_policy

`mongodbatlas_cloud_provider_snapshot_backup_policy` provides a resource that enables you to view and modify the snapshot schedule and retention settings for an Atlas cluster with Cloud Backup enabled.  A default policy is created automatically when Cloud Backup is enabled for the cluster.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

# Examples - Modifying Polices
When Cloud Backup is enabled for a cluster MongoDB Atlas automatically creates a default Cloud Backup schedule for the cluster with four policy items; hourly, daily, weekly, and monthly. Because of this default creation this provider automatically saves the Cloud Backup Snapshot Policy into the Terraform state when a cluster is created/modified to use Cloud Backup. If the default works well for you then you do not need to do anything other than create a cluster with Cloud Backup enabled and your Terraform state will have this information if you need it. However, if you want the policy to be different than the default we've provided some examples to help below.

## Example Usage - Create a Cluster and Modify the 4 Default Policies Simultaneously

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "MyCluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_provider_snapshot_backup_policy resource

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

  //Keep all 4 default policies but modify the units and values
  //Could also just reflect the policy defaults here for later management
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
```

~> **IMPORTANT:**   `policies.#.policy_item.#.id` is obtained when the cluster is created. The example here shows the default order of the default policy when Cloud Backup is enabled (`cloud_backup` is set to true).  The default policy is viewable in the Terraform State file.

## Example Usage - Create a Cluster and Modify 3 Default Policies and Remove 1 Default Policy Simultaneously

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "MyCluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_provider_snapshot_backup_policy resource

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

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
    #   frequency_interval = 4
    #   frequency_type     = "weekly"
    #   retention_unit     = "weeks"
    #   retention_value    = 3
    # }

    policy_item {
      id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id
      frequency_interval = 5
      frequency_type     = "monthly"
      retention_unit     = "months"
      retention_value    = 4
    }
  }
}
```

-> **NOTE:** If you want the Cloud Backup Snapshot Policy to vary in the number of policies from the default when creating the cluster, perhaps you want to remove one policy item and modify the remaining three, simply follow this example here to remove a policy and modify three.

~> **IMPORTANT:** If we decide to remove the 3rd item, as in our above example marked with `#`, we need to consider that once the cluster is modified or  `terraform refresh` is run the item `2` in the array will be replaced with content of the 4th item, so it could cause an inconsistency. This may be avoided by using hardcoded id values which will better handle this situation. (See below for an example of a hardcoded value)

## Example Usage - Remove 3 Default Policies Items After the Cluster Has Already Been Created and Modify One Policy

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "MyCluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # must be enabled in order to use cloud_provider_snapshot_backup_policy resource

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

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.0.id
    #   frequency_interval = 1
    #   frequency_type     = "hourly"
    #   retention_unit     = "days"
    #   retention_value    = 1
    # }

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.1.id
    #   frequency_interval = 1
    #   frequency_type     = "daily"
    #   retention_unit     = "days"
    #   retention_value    = 2
    # }

    # Item removed
    # policy_item {
    #   id                 = mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
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

-> **NOTE:** In this example we decided to remove the first 3 items so we can't use `mongodbatlas_advanced_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id` to retrieve the monthly id value of the cluster state due to once the cluster being modified or makes a `terraform refresh` will cause that the three items will remove from the state, so we will get an error due to the index 3 doesn't exists any more and our monthly policy item is moved to the first place of the array.  So we use `5f0747cad187d8609a72f546`, which is an example of an id MongoDB Atlas returns for the policy item we want to keep. Here it is hard coded because you need to either use the actual value from the Terraform state or look to map the policy item you want to keep to it's current placement in the state file array.

## Argument Reference

* `project_id` - (Required) The unique identifier of the project for the Atlas cluster.
* `cluster_name` - (Required) The name of the Atlas cluster that contains the snapshot backup policy you want to retrieve.
* `reference_hour_of_day` - (Optional) UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items.
* `reference_minute_of_hour` - (Optional) UTC Minutes after referenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
* `restore_window_days` - (Optional) Number of days back in time you can restore to with point-in-time accuracy. Must be a positive, non-zero integer.
* `update_snapshots` - (Optional) Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.

### Policies
* `policies` - (Required) Contains a document for each backup policy item in the desired updated backup policy.
* `policies.#.id` - (Required) Unique identifier of the backup policy that you want to update. policies.#.id is a value obtained via the mongodbatlas_advanced_cluster resource. `cloud_backup` of the mongodbatlas_advanced_cluster resource must be set to true. See the example above for how to refer to the mongodbatlas_advanced_cluster resource for policies.#.id

#### Policy Item
* `policies.#.policy_item` - (Required) Array of backup policy items.
* `policies.#.policy_item.#.id` - (Required) Unique identifier of the backup policy item. `policies.#.policy_item.#.id` is a value obtained via the mongodbatlas_advanced_cluster resource. `cloud_backup` of the mongodbatlas_advanced_cluster resource must be set to true. See the example above for how to refer to the mongodbatlas_advanced_cluster resource for policies.#.policy_item.#.id
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