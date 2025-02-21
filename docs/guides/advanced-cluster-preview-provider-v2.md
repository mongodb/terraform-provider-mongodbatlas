---
page_title: "Migration Guide: Advanced Cluster (Preview for MongoDB Atlas Provider v2)"
---

# Migration Guide: Advanced Cluster (Preview for MongoDB Atlas Provider v2)

**Objective**: This guide explains the changes introduced for the `mongodbatlas_advanced_cluster` resource in the Preview for MongoDB Atlas Provider v2 of the and how to migrate to it.

 `mongodbatlas_advanced_cluster` in the Preview for MongoDB Atlas Provider v2 is implemented using the recommended [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework). This improves the overall user experience and provides a more consistent and predictable behavior. It supports the latest Terraform features and best practices, including support for `moved` block between different resource types, for more info see the [Migration Guide: Cluster to Advanced Cluster](cluster-to-advanced-cluster-migration-guide#moved-block).

The [resource doc](../resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529) contains all the details about the `mongodbatlas_advanced_cluster` resource in the Preview for MongoDB Atlas Provider v2.

## Enable the Preview for MongoDB Atlas Provider v2 for `mongodbatlas_advanced_cluster`

In order to enable the Preview for MongoDB Atlas Provider v2 for `mongodbatlas_advanced_cluster`, set the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`. This will allow you to use the new `mongodbatlas_advanced_cluster` resource. You can also define the environment variable in your local development environment so your tools can use the new format and help you with linting and auto-completion.

This environment variable only affects the `mongodbatlas_advanced_cluster` resource and corresponding data sources, it doesn't affect other resources. `mongodbatlas_advanced_cluster` definition will use the new format and new features like `moved block` from `mongodbatlas_cluster`to `mongodbatlas_advanced_cluster` will be available.

## Main Changes

1. Elements `replication_specs` and `region_configs` are now list attributes instead of blocks so they they are an array of objects. If there is only one object, it still needs to be in an array. For example,
```terraform
replication_specs {
  region_configs {
    electable_specs {
      instance_size = "M10"
      node_count    = 1
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "US_WEST_1"
  }
  region_configs {
    electable_specs {
      instance_size = "M10"
      node_count    = 2
    }
    provider_name = "AWS"
    priority      = 6
    region_name   = "US_EAST_1"
  }
}
```
goes to:
```
replication_specs = [
  {
    region_configs = [
      {
        electable_specs = {
          instance_size = "M10"
          node_count    = 1
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = "US_WEST_1"
      },
      {
        electable_specs = {
          instance_size = "M10"
          node_count    = 2
        }
        provider_name = "AWS"
        priority      = 6
        region_name   = "US_EAST_1"
      }
    ]
  }
]
```

2. Elements `connection_strings`, `timeouts`, `advanced_configuration`, `bi_connector_config`, `pinned_fcv`, `electable_specs`, `read_only_specs`, `analytics_specs`, `auto_scaling` and `analytics_auto_scaling` are now single attributes instead of blocks so they are an object. For example,
```terraform 
advanced_configuration {
  default_write_concern = "majority"
  javascript_enabled    = true
}  
```
goes to:
```terraform
advanced_configuration = {
  default_write_concern = "majority"
  javascript_enabled    = true
}  
```
If there are references to them, `[0]` or `.0` are dropped. For example,
```terraform
output "standard" {
    value = mongodbatlas_advanced_cluster.cluster.connection_strings[0].standard
}
output "javascript_enabled" {
  value = mongodbatlas_advanced_cluster.cluster.advanced_configuration.0.javascript_enabled
}
```
goes to:
```terraform
output "standard" {
    value = mongodbatlas_advanced_cluster.cluster.connection_strings.standard
}
output "javascript_enabled" {
  value = mongodbatlas_advanced_cluster.cluster.advanced_configuration.javascript_enabled
}
```

3. Elements `tags` and `labels` are now `maps` instead of `blocks`. For example,
```terraform
tags {
  key   = "env"
  value = "dev"
}
tags {
  key   = "tag 2"
  value = "val"
}
tags {
  key   = var.tag_key
  value = "another_val"
}

```
goes to:
```terraform
tags = {
  env           = "dev"         # key strings without blanks can be enclosed in quotes but not required
  "tag 2"       = "val"         # enclose key strings with blanks in quotes
  (var.tag_key) = "another_val" # enclose key expressions in brackets so they can be evaluated
}
```

## Best Practices Before Migrating
Before doing any migration create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).

## How to migrate `mongodbatlas_advanced_cluster` to Preview for MongoDB Atlas Provider v2 

The process to migrate from current `mongodbatlas_advanced_cluster` to the one in Preview for MongoDB Atlas Provider v2 is as follows:
- Before starting, run `terraform plan` to make sure that there are no planned changes.
- Set environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true` in order to use the Preview for MongoDB Atlas Provider v2.
- Run `terraform plan` and you'll see errors as definition file hasn't been updated yet.
- Apply definition changes explained in previous section until there are no errors and no planned changes. **Important**: Don't apply until plan is empty. If it shows other changes, you will need to keep updating the `mongodbatlas_advanced_cluster` configuration until it matches the original configuration.
- Run `terraform apply` to apply the changes. Although there are no plan changes shown to the user, the `mongodbatlas_advanced_cluster` state will be updated to support the Preview for MongoDB Atlas Provider v2.
