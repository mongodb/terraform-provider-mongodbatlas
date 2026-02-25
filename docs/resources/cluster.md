---
subcategory: "Clusters"
---

# Resource: mongodbatlas_cluster

`mongodbatlas_cluster` provides a Cluster resource. The resource lets you create, edit and delete clusters. The resource requires your Project ID.

~> **DEPRECATION:** This resource is deprecated and will be removed in the next major release. Please use `mongodbatlas_advanced_cluster`. For more details, see [our migration guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** A network container is created for a cluster to reside in. To use this container with another resource, such as peering, reference the computed`container_id` attribute on the cluster.

-> **NOTE:** To enable Cluster Extended Storage Sizes use the `is_extended_storage_sizes_enabled` parameter in the [mongodbatlas_project resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/project).

-> **NOTE:** If Backup Compliance Policy is enabled for the project for which this backup schedule is defined, you cannot modify the backup schedule for an individual cluster below the minimum requirements set in the Backup Compliance Policy.  See [Backup Compliance Policy Prohibited Actions and Considerations](https://www.mongodb.com/docs/atlas/backup/cloud-backup/backup-compliance-policy/#configure-a-backup-compliance-policy).

-> **NOTE:** The Low-CPU instance clusters are prefixed with `R`, i.e. `R40`. For complete list of Low-CPU instance clusters see Cluster Configuration Options under each [Cloud Provider](https://www.mongodb.com/docs/atlas/reference/cloud-providers).

~> **IMPORTANT:**
<br> &#8226; Multi Region Cluster: The `mongodbatlas_cluster` resource doesn't return the `container_id` for each region utilized by the cluster. For retrieving the `container_id`, we recommend to use the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource instead.
<br> &#8226; Free tier cluster creation (M0) is supported.
<br> &#8226; Free tier clusters (M0) can be upgraded to dedicated tiers (M10+) via this provider. WARNING WHEN UPGRADING FREE CLUSTERS!!! Any change from free tier to a different instance size will be considered a tenant upgrade. When upgrading from free tier to dedicated simply change the `provider_name` from "TENANT"  to your preferred provider (AWS, GCP, AZURE) and remove the variable `backing_provider_name`, for example if you have an existing free cluster and want to upgrade your Terraform config should be changed from:
```
provider_instance_size_name = "M0"
provider_name               = "TENANT"
backing_provider_name       = "AWS"
```
To:
```
provider_instance_size_name = "M10"
provider_name               = "AWS"
```
<br> &#8226; Changes to cluster configurations can affect costs. Before making changes, please see [Billing](https://docs.atlas.mongodb.com/billing/).   
<br> &#8226; If your Atlas project contains a custom role that uses actions introduced in a specific MongoDB version, you cannot create a cluster with a MongoDB version less than that version unless you delete the custom role.

## Example Usage

### Example AWS cluster

```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "cluster-test"
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "7.0"

  # Provider Settings "block"
  provider_name               = "AWS"
  provider_instance_size_name = "M40"
}
```

### Example Azure cluster.

```terraform
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "test"
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup     = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "7.0"

  # Provider Settings "block"
  provider_name               = "AZURE"
  provider_disk_type_name     = "P6"
  provider_instance_size_name = "M30"
}
```

### Example GCP cluster

```terraform
resource "mongodbatlas_cluster" "test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "test"
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "EASTERN_US"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup                 = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "7.0"

  # Provider Settings "block"
  provider_name               = "GCP"
  provider_instance_size_name = "M30"
}
```

### Example Multi Region cluster

```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id               = "<YOUR-PROJECT-ID>"
  name                     = "cluster-test-multi-region"
  num_shards               = 1
  cloud_backup             = true
  cluster_type             = "REPLICASET"

  # Provider Settings "block"
  provider_name               = "AWS"
  provider_instance_size_name = "M10"

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
    regions_config {
      region_name     = "US_EAST_2"
      electable_nodes = 2
      priority        = 6
      read_only_nodes = 0
    }
    regions_config {
      region_name     = "US_WEST_1"
      electable_nodes = 2
      priority        = 5
      read_only_nodes = 2
    }
  }
}
```

### Example Global cluster

```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"
  num_shards              = 1
  cloud_backup            = true
  cluster_type            = "GEOSHARDED"

  # Provider Settings "block"
  provider_name               = "AWS"
  provider_instance_size_name = "M30"

  replication_specs {
    zone_name  = "Zone 1"
    num_shards = 2
    regions_config {
      region_name     = "US_EAST_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }

  replication_specs {
    zone_name  = "Zone 2"
    num_shards = 2
    regions_config {
      region_name     = "EU_CENTRAL_1"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
}
```
### Example AWS Free Tier cluster
```terraform
resource "mongodbatlas_cluster" "cluster-test" {
  project_id              = "<YOUR-PROJECT-ID>"
  name                    = "cluster-test-global"

  # Provider Settings "block"
  provider_name = "TENANT"
  backing_provider_name = "AWS"
  provider_region_name = "US_EAST_1"
  provider_instance_size_name = "M0"
}
```
### Example - Return a Connection String
Standard
```terraform
output "standard" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].standard
}
# Example return string: standard = "mongodb://cluster-atlas-shard-00-00.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-01.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-02.ygo1m.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-12diht-shard-0"
```
Standard srv
```terraform
output "standard_srv" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].standard_srv
}
# Example return string: standard_srv = "mongodb+srv://cluster-atlas.ygo1m.mongodb.net"
```
Private with Network peering and Custom DNS AWS enabled
```terraform
output "private" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].private
}
# Example return string: private = "mongodb://cluster-atlas-shard-00-00-pri.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-01-pri.ygo1m.mongodb.net:27017,cluster-atlas-shard-00-02-pri.ygo1m.mongodb.net:27017/?ssl=true&authSource=admin&replicaSet=atlas-12diht-shard-0"
```
Private srv with Network peering and Custom DNS AWS enabled
```terraform
output "private_srv" {
    value = mongodbatlas_cluster.cluster-test.connection_strings[0].private_srv
}
# Example return string: private_srv = "mongodb+srv://cluster-atlas-pri.ygo1m.mongodb.net"
```

By endpoint_service_id
```terraform
locals {
  endpoint_service_id = google_compute_network.default.name
  private_endpoints   = try(flatten([for cs in data.mongodbatlas_advanced_cluster.cluster[0].connection_strings : cs.private_endpoint]), [])
  connection_strings = [
    for pe in local.private_endpoints : pe.srv_connection_string
    if contains([for e in pe.endpoints : e.endpoint_id], local.endpoint_service_id)
  ]
}
output "endpoint_service_connection_string" {
  value = length(local.connection_strings) > 0 ? local.connection_strings[0] : ""
}
# Example return string: connection_string = "mongodb+srv://cluster-atlas-pl-0.ygo1m.mongodb.net"
```

Refer to the following for full privatelink endpoint connection string examples:
* [GCP Private Endpoint (Port-Mapped Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)
* [Azure Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/azure)
* [AWS, Private Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster)
* [AWS, Regionalized Private Endpoints](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster-geosharded)


### Further Examples 
- [NVMe Upgrade (Dedicated Cluster)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_cluster/nvme-upgrade)
- [Tenant to Dedicated Upgrade (Cluster)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_cluster/tenant-upgrade)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `project_id` (String)
- `provider_instance_size_name` (String)
- `provider_name` (String)

### Optional

- `accept_data_risks_and_force_replica_set_reconfig` (String) Submit this field alongside your topology reconfiguration to request a new regional outage resistant topology
- `advanced_configuration` (Block List, Max: 1) (see [below for nested schema](#nestedblock--advanced_configuration))
- `auto_scaling_compute_enabled` (Boolean)
- `auto_scaling_compute_scale_down_enabled` (Boolean)
- `auto_scaling_disk_gb_enabled` (Boolean)
- `backing_provider_name` (String)
- `backup_enabled` (Boolean) Clusters running MongoDB FCV 4.2 or later and any new Atlas clusters of any type do not support this parameter
- `bi_connector_config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--bi_connector_config))
- `cloud_backup` (Boolean)
- `cluster_type` (String)
- `disk_size_gb` (Number)
- `encryption_at_rest_provider` (String)
- `labels` (Block Set) (see [below for nested schema](#nestedblock--labels))
- `mongo_db_major_version` (String)
- `num_shards` (Number)
- `paused` (Boolean)
- `pinned_fcv` (Block List, Max: 1) (see [below for nested schema](#nestedblock--pinned_fcv))
- `pit_enabled` (Boolean)
- `provider_auto_scaling_compute_max_instance_size` (String)
- `provider_auto_scaling_compute_min_instance_size` (String)
- `provider_disk_iops` (Number)
- `provider_disk_type_name` (String)
- `provider_encrypt_ebs_volume` (Boolean, Deprecated)
- `provider_region_name` (String)
- `provider_volume_type` (String)
- `redact_client_log_data` (Boolean)
- `replication_factor` (Number)
- `replication_specs` (Block Set) (see [below for nested schema](#nestedblock--replication_specs))
- `retain_backups_enabled` (Boolean) Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster
- `tags` (Block Set) (see [below for nested schema](#nestedblock--tags))
- `termination_protection_enabled` (Boolean)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `version_release_system` (String)

### Read-Only

- `cluster_id` (String)
- `connection_strings` (List of Object) (see [below for nested schema](#nestedatt--connection_strings))
- `container_id` (String)
- `id` (String) The ID of this resource.
- `mongo_db_version` (String)
- `mongo_uri` (String)
- `mongo_uri_updated` (String)
- `mongo_uri_with_options` (String)
- `provider_encrypt_ebs_volume_flag` (Boolean)
- `snapshot_backup_policy` (List of Object) (see [below for nested schema](#nestedatt--snapshot_backup_policy))
- `srv_address` (String)
- `state_name` (String)

<a id="nestedblock--advanced_configuration"></a>
### Nested Schema for `advanced_configuration`

Optional:

- `change_stream_options_pre_and_post_images_expire_after_seconds` (Number)
- `custom_openssl_cipher_config_tls12` (Set of String)
- `default_max_time_ms` (Number)
- `default_read_concern` (String, Deprecated)
- `default_write_concern` (String)
- `fail_index_key_too_long` (Boolean, Deprecated)
- `javascript_enabled` (Boolean)
- `minimum_enabled_tls_protocol` (String)
- `no_table_scan` (Boolean)
- `oplog_min_retention_hours` (Number)
- `oplog_size_mb` (Number)
- `sample_refresh_interval_bi_connector` (Number)
- `sample_size_bi_connector` (Number)
- `tls_cipher_config_mode` (String)
- `transaction_lifetime_limit_seconds` (Number)


<a id="nestedblock--bi_connector_config"></a>
### Nested Schema for `bi_connector_config`

Optional:

- `enabled` (Boolean)
- `read_preference` (String)


<a id="nestedblock--labels"></a>
### Nested Schema for `labels`

Optional:

- `key` (String)
- `value` (String)


<a id="nestedblock--pinned_fcv"></a>
### Nested Schema for `pinned_fcv`

Required:

- `expiration_date` (String)

Read-Only:

- `version` (String)


<a id="nestedblock--replication_specs"></a>
### Nested Schema for `replication_specs`

Required:

- `num_shards` (Number)

Optional:

- `id` (String)
- `regions_config` (Block Set) (see [below for nested schema](#nestedblock--replication_specs--regions_config))
- `zone_name` (String)

<a id="nestedblock--replication_specs--regions_config"></a>
### Nested Schema for `replication_specs.regions_config`

Required:

- `region_name` (String)

Optional:

- `analytics_nodes` (Number)
- `electable_nodes` (Number)
- `priority` (Number)
- `read_only_nodes` (Number)



<a id="nestedblock--tags"></a>
### Nested Schema for `tags`

Required:

- `key` (String)
- `value` (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)


<a id="nestedatt--connection_strings"></a>
### Nested Schema for `connection_strings`

Read-Only:

- `private` (String)
- `private_endpoint` (List of Object) (see [below for nested schema](#nestedobjatt--connection_strings--private_endpoint))
- `private_srv` (String)
- `standard` (String)
- `standard_srv` (String)

<a id="nestedobjatt--connection_strings--private_endpoint"></a>
### Nested Schema for `connection_strings.private_endpoint`

Read-Only:

- `connection_string` (String)
- `endpoints` (List of Object) (see [below for nested schema](#nestedobjatt--connection_strings--private_endpoint--endpoints))
- `srv_connection_string` (String)
- `srv_shard_optimized_connection_string` (String)
- `type` (String)

<a id="nestedobjatt--connection_strings--private_endpoint--endpoints"></a>
### Nested Schema for `connection_strings.private_endpoint.endpoints`

Read-Only:

- `endpoint_id` (String)
- `provider_name` (String)
- `region` (String)




<a id="nestedatt--snapshot_backup_policy"></a>
### Nested Schema for `snapshot_backup_policy`

Read-Only:

- `cluster_id` (String)
- `cluster_name` (String)
- `next_snapshot` (String)
- `policies` (List of Object) (see [below for nested schema](#nestedobjatt--snapshot_backup_policy--policies))
- `reference_hour_of_day` (Number)
- `reference_minute_of_hour` (Number)
- `restore_window_days` (Number)
- `update_snapshots` (Boolean)

<a id="nestedobjatt--snapshot_backup_policy--policies"></a>
### Nested Schema for `snapshot_backup_policy.policies`

Read-Only:

- `id` (String)
- `policy_item` (List of Object) (see [below for nested schema](#nestedobjatt--snapshot_backup_policy--policies--policy_item))

<a id="nestedobjatt--snapshot_backup_policy--policies--policy_item"></a>
### Nested Schema for `snapshot_backup_policy.policies.policy_item`

Read-Only:

- `frequency_interval` (Number)
- `frequency_type` (String)
- `id` (String)
- `retention_unit` (String)
- `retention_value` (Number)

## Import

Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTERNAME`, e.g.

```
$ terraform import mongodbatlas_cluster.my_cluster 1112222b3bf99403840e8934-Cluster0
```

See detailed information for arguments and attributes: [MongoDB API Clusters](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/)
