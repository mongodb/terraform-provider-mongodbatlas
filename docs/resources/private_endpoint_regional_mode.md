# Resource: private_endpoint_regional_mode

`mongodbatlas_private_endpoint_regional_mode` provides a Private Endpoint Regional Mode resource. This represents a regionalized private endpoint setting for a Project. Enable it to allow region specific private endpoints.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource:
  * Organization Owner
  * Project Owner

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **WARNING:**Your [connection strings](https://www.mongodb.com/docs/atlas/reference/faq/connection-changes/#std-label-connstring-privatelink) to existing multi-region and global sharded clusters change when you enable this setting.  You must update your applications to use the new connection strings. This might cause downtime.

## Example AWS Global Cluster with multiple Private Endpoint

```terraform
resource "mongodbatlas_private_endpoint_regional_mode" "test" {
  project_id = var.atlasprojectid
  enabled    = true
}

resource "mongodbatlas_advanced_cluster" "cluster_atlas" {
  project_id     = var.atlasprojectid
  name           = var.cluster_name
  cluster_type   = "GEOSHARDED"
  backup_enabled = true

  replication_specs { # Shard 1
    zone_name = "Zone 1"

    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = var.atlas_region_east
    }

    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 2
      }
      provider_name = "AWS"
      priority      = 6
      region_name   = var.atlas_region_west
    }
  }

  replication_specs { # Shard 2
    zone_name = "Zone 1"

    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = var.atlas_region_east
    }

    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 2
      }
      provider_name = "AWS"
      priority      = 6
      region_name   = var.atlas_region_west
    }
  }

  depends_on = [
    mongodbatlas_privatelink_endpoint_service.test_west,
    mongodbatlas_privatelink_endpoint_service.test_east,
    mongodbatlas_private_endpoint_regional_mode.test
  ]
}

resource "mongodbatlas_privatelink_endpoint" "test_west" {
  project_id    = var.atlasprojectid
  provider_name = "AWS"
  region        = "US_WEST_1"
}

resource "mongodbatlas_privatelink_endpoint_service" "test_west" {
  project_id          = mongodbatlas_privatelink_endpoint.test_west.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.test_west.private_link_id
  endpoint_service_id = aws_vpc_endpoint.test_west.id
  provider_name       = "AWS"
}

resource "aws_vpc_endpoint" "test_west" {
  provider           = aws.west
  vpc_id             = "vpc-7fc0a543"
  service_name       = mongodbatlas_privatelink_endpoint.test_west.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-de0406d2"]
  security_group_ids = ["sg-3f238186"]
}

resource "mongodbatlas_privatelink_endpoint" "test_east" {
  project_id    = "var.atlasprojectid
  provider_name = "AWS"
  region        = "US_EAST_1"
}

resource "mongodbatlas_privatelink_endpoint_service" "test_east" {
  project_id          = mongodbatlas_privatelink_endpoint.test_east.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.test_east.private_link_id
  endpoint_service_id = aws_vpc_endpoint.test_east.id
  provider_name       = "AWS"
}

resource "aws_vpc_endpoint" "test_east" {
  provider           = aws.east
  vpc_id             = "vpc-345a0cf7"
  service_name       = mongodbatlas_privatelink_endpoint.test_east.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-2d6040ed"]
  security_group_ids = ["sg-681832f3"]
}

```

## Argument Reference
* `project_id` - (Required) Unique identifier for the project.
* `enabled` - (Optional) Flag that indicates whether the regionalized private endpoint setting is enabled for the project.   Set this value to true to create more than one private endpoint in a cloud provider region to connect to multi-region and global Atlas sharded clusters. You can enable this setting only if your Atlas project contains no replica sets. You can't disable this setting if you have:
   * More than one private endpoint in more than one region, or
   * More than one private endpoint in one region and one private endpoint in one or more regions.
You can create only sharded clusters when you enable the regionalized private endpoint setting. You can't create replica sets.

* `timeouts`- (Optional) The duration of time to wait for Cluster to be created, updated, or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Private Endpoint Regional Mode operations is `3h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).

## Additional Reference

In addition to the example shown above, keep in mind:
* `mongodbatlas_advanced_cluster.cluster_atlas.depends_on` - Make your cluster dependent on the project's `mongodbatlas_private_endpoint_regional_mode` as well as any relevant `mongodbatlas_privatelink_endpoint_service` resources.  See an [example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/aws-privatelink-endpoint/cluster-geosharded). 
* `mongodbatlas_advanced_cluster.cluster_atlas.connection_strings` will differ based on the value of `mongodbatlas_private_endpoint_regional_mode.test.enabled`.
* For more information on usage with GCP, see [our Privatelink Endpoint Service documentation: Example with GCP](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint_service#example-with-gcp)
* For more information on usage with Azure, see [our Privatelink Endpoint Service documentation: Examples with Azure](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/privatelink_endpoint_service#example-with-azure)

## Import
Private Endpoint Regional Mode can be imported using project id in format `{project_id}`, e.g.

```
$ terraform import mongodbatlas_private_endpoint_regional_mode.test 1112222b3bf99403840e8934
```

See detailed information for arguments and attributes: **Private Endpoints** [Get Regional Mode](https://www.mongodb.com/docs/atlas/reference/api/private-endpoints-get-regional-mode/) | [Update Regional Mode](https://www.mongodb.com/docs/atlas/reference/api/private-endpoints-update-regional-mode/)
