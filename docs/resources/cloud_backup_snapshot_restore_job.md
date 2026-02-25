---
subcategory: "Cloud Backups"
---

# Resource: mongodbatlas_cloud_backup_snapshot_restore_job

`mongodbatlas_cloud_backup_snapshot_restore_job` provides a resource to create a new restore job from a cloud backup snapshot of a specified cluster. The restore job must define one of three delivery types:
* **automated:** Atlas automatically restores the snapshot with snapshotId to the Atlas cluster with name targetClusterName in the Atlas project with targetGroupId.

* **download:** Atlas provides a URL to download a .tar.gz of the snapshot with snapshotId. The contents of the archive contain the data files for your Atlas cluster.

* **pointInTime:**  Atlas performs a Continuous Cloud Backup restore.

-> **Important:** If you specify `deliveryType` : `automated` or `deliveryType` : `pointInTime` in your request body to create an automated restore job, Atlas removes all existing data on the target cluster prior to the restore.

-> **Important:** If you specify `deliveryType` : `automated` or `deliveryType` : `pointInTime` in your
`mongodbatlas_cloud_backup_snapshot_restore_job` resource, you won't be able to delete the snapshot resource in MongoDB Atlas as the Atlas Admin API doesn't support this. The provider will remove the Terraform resource from the state file but won't destroy the MongoDB Atlas resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

### Example automated delivery type

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "MyCluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud backup snapshots

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_advanced_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_advanced_cluster.my_cluster.name
  description       = "myDescription"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id     = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
  delivery_type_config   {
    automated           = true
    target_cluster_name = "MyCluster"
    target_project_id   = "5cf5a45a9ccf6400e60981b6"
  }
}
```

### Example download delivery type

```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "MyCluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud backup snapshots

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_advanced_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_advanced_cluster.my_cluster.name
  description       = "myDescription"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id     = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
  delivery_type_config {
    download = true
  }
}
```

### Example of a point in time restore
```terraform
resource "mongodbatlas_advanced_cluster" "my_cluster" {
  project_id     = "<PROJECT-ID>"
  name           = "MyCluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # enable cloud backup snapshots

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_advanced_cluster.cluster_test.project_id
  cluster_name      = mongodbatlas_advanced_cluster.cluster_test.name
  description       = "My description"
  retention_in_days = "1"
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  count        = (var.point_in_time_utc_seconds == 0 ? 0 : 1)
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config {
    point_in_time             = true
    target_cluster_name       = mongodbatlas_advanced_cluster.cluster_test.name
    target_project_id         = mongodbatlas_advanced_cluster.cluster_test.project_id
    point_in_time_utc_seconds = var.point_in_time_utc_seconds
  }
}
```

### Further Examples
- [Restore from backup snapshot at point in time](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_cloud_backup_snapshot_restore_job/point-in-time)
- [Restore from backup snapshot using an advanced cluster resource](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_cloud_backup_snapshot_restore_job/point-in-time-advanced-cluster)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_name` (String)
- `project_id` (String)

### Optional

- `delivery_type_config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--delivery_type_config))
- `snapshot_id` (String)

### Read-Only

- `cancelled` (Boolean)
- `delivery_url` (List of String)
- `expired` (Boolean)
- `expires_at` (String)
- `failed` (Boolean)
- `finished_at` (String)
- `id` (String) The ID of this resource.
- `snapshot_restore_job_id` (String)
- `timestamp` (String)

<a id="nestedblock--delivery_type_config"></a>
### Nested Schema for `delivery_type_config`

Optional:

- `automated` (Boolean)
- `download` (Boolean)
- `oplog_inc` (Number)
- `oplog_ts` (Number)
- `point_in_time` (Boolean)
- `point_in_time_utc_seconds` (Number)
- `target_cluster_name` (String)
- `target_project_id` (String)

## Import

Cloud Backup Snapshot Restore Job entries can be imported using project project_id, cluster_name and snapshot_id (Unique identifier of the snapshot), in the format `PROJECTID-CLUSTERNAME-JOBID`, e.g.

```
$ terraform import mongodbatlas_cloud_backup_snapshot_restore_job.test 5cf5a45a9ccf6400e60981b6-MyCluster-5d1b654ecf09a24b888f4c79
```

For more information see: [MongoDB Atlas API Reference.](https://docs.atlas.mongodb.com/reference/api/cloud-backup/restore/restores/)
