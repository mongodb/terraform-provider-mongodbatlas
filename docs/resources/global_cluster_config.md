# Resource: mongodbatlas_global_cluster_config

`mongodbatlas_global_cluster_config` provides a Global Cluster Configuration resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** This resource can only be used with Atlas-managed clusters. See doc for `global_cluster_self_managed_sharding` attribute in [`mongodbatlas_advanced_cluster` resource](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) for more information.

~> **IMPORTANT:** You can update a Global Cluster Configuration to add new custom zone mappings and managed namespaces. However, once configured, you can't modify or partially delete custom zone mappings (you must remove them all at once). You can add or remove, but can't modify, managed namespaces. Any update that changes an existing managed namespace results in an error. [Read more about Global Cluster Configuration](https://www.mongodb.com/docs/atlas/global-clusters/). For more details, see [Global Clusters API](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-global-clusters)

## Examples Usage

### Example Global cluster

```terraform
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = "<YOUR-PROJECT-ID>"
  name           = "<CLUSTER-NAME>"
  cluster_type   = "GEOSHARDED"
  backup_enabled = true

  replication_specs {
    zone_name = "Zone 1"

    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_CENTRAL_1"
    }
  }

  replication_specs {
    zone_name = "Zone 2"

    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_2"
    }
  }
}

resource "mongodbatlas_global_cluster_config" "config" {
	project_id = mongodbatlas_advanced_cluster.test.project_id
	cluster_name = mongodbatlas_advanced_cluster.test.name

	managed_namespaces {
		db 				 = "mydata"
		collection 		 = "publishers"
		custom_shard_key = "city"
					is_custom_shard_key_hashed = false
					is_shard_key_unique = false
	}

	custom_zone_mappings {
		location ="CA"
		zone =  "Zone 1"
	}
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project.
* `cluster_name` - (Required) The name of the Global Cluster.
*  `managed_namespaces` - (Optional) Add a managed namespaces to a Global Cluster. For more information about managed namespaces, see [Global Clusters](https://docs.atlas.mongodb.com/reference/api/global-clusters/). See [Managed Namespace](#managed-namespace) below for more details.
*  `custom_zone_mappings` - (Optional) Each element in the list maps one ISO location code to a zone in your Global Cluster. See [Custom Zone Mapping](#custom-zone-mapping) below for more details.

### Managed Namespace

* `collection` -	(Required) The name of the collection associated with the managed namespace.
* `custom_shard_key` - (Required)	The custom shard key for the collection. Global Clusters require a compound shard key consisting of a location field and a user-selected second key, the custom shard key.
* `db` - (Required) The name of the database containing the collection.
* `is_custom_shard_key_hashed` - (Optional) Specifies whether the custom shard key for the collection is [hashed](https://docs.mongodb.com/manual/reference/method/sh.shardCollection/#hashed-shard-keys). If omitted, defaults to `false`. If `false`, Atlas uses [ranged sharding](https://docs.mongodb.com/manual/core/ranged-sharding/). This is only available for Atlas clusters with MongoDB v4.4 and later.
* `is_shard_key_unique` - (Optional) Specifies whether the underlying index enforces a unique constraint. If omitted, defaults to false. You cannot specify true when using [hashed shard keys](https://docs.mongodb.com/manual/core/hashed-sharding/#std-label-sharding-hashed).

### Custom Zone Mapping

* `location` - (Required) The ISO location code to which you want to map a zone in your Global Cluster. You can find a list of all supported location codes [here](https://cloud.mongodb.com/static/atlas/country_iso_codes.txt).
* `zone` - (Required) The name of the zone in your Global Cluster that you want to map to location.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `custom_zone_mapping_zone_id` - A map of all custom zone mappings defined for the Global Cluster to `replication_specs.*.zone_id`. Atlas automatically maps each location code to the closest geographical zone. Custom zone mappings allow administrators to override these automatic mappings. If your Global Cluster does not have any custom zone mappings, this document is empty.
* `custom_zone_mapping` - (Deprecated) A map of all custom zone mappings defined for the Global Cluster to `replication_specs.*.id`. This attribute is deprecated, use `custom_zone_mapping_zone_id` instead. This attribute is not set when a cluster uses independent shard scaling. To learn more, see the [Sharding Configuration guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

## Import

Global Clusters can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTER_NAME`, e.g.

```
$ terraform import mongodbatlas_global_cluster_config.config 1112222b3bf99403840e8934-Cluster0
```

See detailed information for arguments and attributes: [MongoDB API Global Clusters](https://docs.atlas.mongodb.com/reference/api/global-clusters/)
