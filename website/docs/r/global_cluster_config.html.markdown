---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: global_cluster_config"
sidebar_current: "docs-mongodbatlas-resource-global-cluster-config"
description: |-
    Provides a Global Cluster Configuration resource.
---

# mongodbatlas_global_cluster_config

`mongodbatlas_global_cluster_config` provides a Global Cluster Configuration resource.


-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.


## Examples Usage

### Example Global cluster

```hcl
	resource "mongodbatlas_cluster" "test" {
		project_id              = "<YOUR-PROJECT-ID>"
		name                    = "<CLUSTER-NAME>"
		disk_size_gb            = 80
		backup_enabled          = false
		provider_backup_enabled = true
		cluster_type            = "GEOSHARDED"
		
		//Provider Settings "block"
		provider_name               = "AWS"
		provider_disk_iops          = 240
		provider_instance_size_name = "M30"
		
		replication_specs {
			zone_name  = "Zone 1"
			num_shards = 1
			regions_config {
			region_name     = "EU_CENTRAL_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}
		
		replication_specs { 
			zone_name  = "Zone 2"
			num_shards = 1
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
```

### Example AWS cluster

```hcl
resource "mongodbatlas_cluster" "cluster-test" {
  project_id   = "<YOUR-PROJECT-ID>"
  name         = "cluster-test"
  num_shards   = 1

  replication_factor           = 3
  backup_enabled               = true
  auto_scaling_disk_gb_enabled = true
  mongo_db_major_version       = "4.0"

  //Provider Settings "block"
  provider_name               = "AWS"
  disk_size_gb                = 100
  provider_disk_iops          = 300
  provider_encrypt_ebs_volume = false
  provider_instance_size_name = "M40"
  provider_region_name        = "US_EAST_1"
}

resource "mongodbatlas_global_cluster_config" "config" {
	project_id   = mongodbatlas_cluster.test.project_id
	cluster_name = mongodbatlas_cluster.test.name

	managed_namespaces {
		db               = "mydata"
		collection       = "publishers"
		custom_shard_key = "city"
	}

	custom_zone_mappings {
		location = "CA"
		zone     = "Zone 1"
	}
}
```


## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `cluster_name - (Required) The name of the Global Cluster.
*  `managed_namespaces` - (Optional) Add a managed namespaces to a Global Cluster. For more information about managed namespaces, see [Global Clusters](https://docs.atlas.mongodb.com/reference/api/global-clusters/). See [Managed Namespace](#managed-namespace) below for more details.
*  `custom_zone_mappings` - (Optional) Each element in the list maps one ISO location code to a zone in your Global Cluster. See [Custom Zone Mapping](#custom-zone-mapping) below for more details.

### Managed Namespace

* `collection` -	(Required) The name of the collection associated with the managed namespace.
* `custom_shard_key` - (Required)	The custom shard key for the collection. Global Clusters require a compound shard key consisting of a location field and a user-selected second key, the custom shard key.
* `db` - (Required) The name of the database containing the collection.


### Custom Zone Mapping

* `location` - (Required) The ISO location code to which you want to map a zone in your Global Cluster. You can find a list of all supported location codes [here](https://cloud.mongodb.com/static/atlas/country_iso_codes.txt).
* `zone` - (Required) The name of the zone in your Global Cluster that you want to map to location.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `custom_zone_mapping` - A map of all custom zone mappings defined for the Global Cluster. Atlas automatically maps each location code to the closest geographical zone. Custom zone mappings allow administrators to override these automatic mappings. If your Global Cluster does not have any custom zone mappings, this document is empty.

## Import

Database users can be imported using project ID and cluster name, in the format `PROJECTID-CLUSTER_NAME`, e.g.

```
$ terraform import mongodbatlas_global_cluster_config.config 1112222b3bf99403840e8934-my-cluster
```

See detailed information for arguments and attributes: [MongoDB API Global Clusters](https://docs.atlas.mongodb.com/reference/api/global-clusters/)