---
page_title: "Migration Guide: Advanced Cluster (v1.x → v2.0.0)"
---

This guide helps you migrate from the legacy schema of `mongodbatlas_advanced_cluster` resource to the new schema introduced in v2.0.0 of the provider. The new implementation uses the recommended Terraform Plugin Framework, which, in addition to providing a better user experience and other features, adds support for the `moved` block between different resource types.

~> **IMPORTANT:** Preview of the new schema was already released in versions 1.29.0 and later which could be enabled by setting the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`. If you are already using the new schema preview and not using deprecated attributes, you would not be required to make any additional changes except that the mentioned environment variable is no longer required.

## Configuration changes when upgrading from v1.x

In this section you can find the configuration changes between the legacy `mongodbatlas_advanced_cluster` and the new one released in v2.0.0.

1. Elements `replication_specs` and `region_configs` are now list attributes instead of blocks so they are an array of objects. If there is only one object, it still needs to be in an array. For example,
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
```terraform
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

## How to migrate

If you currently use `mongodbatlas_cluster`, see our [Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

If you currently use `mongodbatlas_advanced_cluster` from v1.x.x of our provider, we recommend that you do the following steps:

~> **IMPORTANT:** Before you migrate, create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state). The state file will update to the new format and the old format will no longer be supported.

1. After you upgrade to v2.0.0+ from v1.x.x, when you run `terraform plan`, syntax errors will return as expected since the definition file hasn't been updated yet using the latest schema.
2. At this point, you can apply definition changes [explained on this page](#configuration-changes) until there are no errors and no planned changes. **Important**: Don't apply until the plan is empty. If it shows other changes, you must update the `mongodbatlas_advanced_cluster` configuration until it matches the original configuration.
3. Run `terraform apply` to apply the changes. Although there are no plan changes shown to the user, the `mongodbatlas_advanced_cluster` state will be updated to support the new schema.

## Important notes

Please refer to our [Considerations and Best Practices](#considerations-and-best-practices) section for additional guidance on this resource.
