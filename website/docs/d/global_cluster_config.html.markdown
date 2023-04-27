---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: global_cluster_config"
sidebar_current: "docs-mongodbatlas-datasource-global-cluster-config"
description: |-
    Describes the Global Cluster Configuration.
---

# Data Source: mongodbatlas_global_cluster_config

`mongodbatlas_global_cluster_config` describes all managed namespaces and custom zone mappings associated with the specified Global Cluster.


-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.


## Example Usage

```terraform
	resource "mongodbatlas_cluster" "test" {
		project_id              = "<YOUR-PROJECT-ID>"
		name                    = "<CLUSTER-NAME>"
		cloud_backup            = true
		cluster_type            = "GEOSHARDED"

		//Provider Settings "block"
		provider_name               = "AWS"
		provider_instance_size_name = "M30"

		replication_specs {
			zone_name  = "Zone 1"
			num_shards = 2
			regions_config {
			region_name     = "EU_CENTRAL_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}

		replication_specs {
			zone_name  = "Zone 2"
			num_shards = 2
			regions_config {
			region_name     = "US_EAST_2"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}
	}

	resource "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_cluster.test.project_id
		cluster_name = mongodbatlas_cluster.test.name

		managed_namespaces {
			db 				 = "mydata"
			collection 		 = "publishers"
			custom_shard_key = "city"
		}

		custom_zone_mappings {
			location ="CA"
			zone =  "Zone 1"
		}
	}

    data "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_global_cluster_config.config.project_id
		cluster_name = mongodbatlas_global_cluster_config.config.cluster_name
	}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `cluster_name - (Required) The name of the Global Cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `custom_zone_mapping` - A map of all custom zone mappings defined for the Global Cluster. Atlas automatically maps each location code to the closest geographical zone. Custom zone mappings allow administrators to override these automatic mappings. If your Global Cluster does not have any custom zone mappings, this document is empty.
*  `managed_namespaces` - Add a managed namespaces to a Global Cluster. For more information about managed namespaces, see [Global Clusters](https://docs.atlas.mongodb.com/reference/api/global-clusters/). See [Managed Namespace](#managed-namespace) below for more details.

### Managed Namespace

* `collection` -	(Required) The name of the collection associated with the managed namespace.
* `custom_shard_key` - (Required)	The custom shard key for the collection. Global Clusters require a compound shard key consisting of a location field and a user-selected second key, the custom shard key.
* `db` - (Required) The name of the database containing the collection.
* `is_custom_shard_key_hashed` - Specifies whether the custom shard key for the collection is [hashed](https://docs.mongodb.com/manual/reference/method/sh.shardCollection/#hashed-shard-keys). If omitted, defaults to `false`. If `false`, Atlas uses [ranged sharding](https://docs.mongodb.com/manual/core/ranged-sharding/). This is only available for Atlas clusters with MongoDB v4.4 and later.
* `is_shard_key_unique` - Specifies whether the underlying index enforces a unique constraint. If omitted, defaults to false. You cannot specify true when using [hashed shard keys](https://docs.mongodb.com/manual/core/hashed-sharding/#std-label-sharding-hashed).


See detailed information for arguments and attributes: [MongoDB API Global Clusters](https://docs.atlas.mongodb.com/reference/api/global-clusters/)
