---
page_title: "Migration Guide: Advanced Cluster (v1.x → v2.0.0)"
---

# Migration Guide: Advanced Cluster (v1.x → v2.0.0)

This guide helps you migrate from the legacy schema of `mongodbatlas_advanced_cluster` resource to the new schema introduced in v2.0.0 of the provider. The new implementation uses:
 
1. The recommended Terraform Plugin Framework, which, in addition to providing a better user experience and other features, adds support for the `moved` block between different resource types.
2. New sharding configurations that supports scaling shards independently (see the [Migration Guide: Advanced Cluster New Sharding Configurations](advanced-cluster-new-sharding-schema#migration-sharded)).

~> **IMPORTANT:** Preview of the new schema was already released in versions 1.29.0 and later which could be enabled by setting the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`. If you are already using the new schema preview with the new sharding configurations **and not using deprecated attributes**, you would not be required to make any additional changes except that the mentioned environment variable is no longer required.

## Configuration changes when upgrading from v1.x

In this section you can find the configuration changes between the legacy `mongodbatlas_advanced_cluster` and the new one released in v2.0.0.

1. Below deprecated attributes have been removed:
  - `id`
  - `disk_size_gb`
  - `replication_specs.#.num_shards`
  - `replication_specs.#.id`
  - `advanced_configuration.default_read_concern`
  - `advanced_configuration.fail_index_key_too_long`

2. Elements `replication_specs` and `region_configs` are now list attributes instead of blocks so they are an array of objects. If there is only one object, it still needs to be in an array. For example,
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

3. `mongodbatlas_advanced_cluster` now supports only the new sharding configuration that allows scaling shards independently. If your configuration defines the num_shards attribute (removed in 2.0.0), please also see the [Migration Guide: Advanced Cluster New Sharding Configurations](advanced-cluster-new-sharding-schema#migration-sharded).

4. Elements `connection_strings`, `timeouts`, `advanced_configuration`, `bi_connector_config`, `pinned_fcv`, `electable_specs`, `read_only_specs`, `analytics_specs`, `auto_scaling` and `analytics_auto_scaling` are now single attributes instead of blocks so they are an object. For example,
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

5. Elements `tags` and `labels` are now `maps` instead of `blocks`. For example,
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

6. `id` attribute which was an internal encoded resource identifier has been removed. Use `cluster_id` instead.

### Configuration changes when upgrading `data.mongodbatlas_advanced_cluster` and `data.mongodbatlas_advanced_clusters` from v1.x

1. Below deprecated attributes have been removed (same as resource):
  - `id`
  - `disk_size_gb`
  - `replication_specs.#.num_shards`
  - `replication_specs.#.id`
  - `advanced_configuration.default_read_concern`
  - `advanced_configuration.fail_index_key_too_long`

2. Deprecated attribute `use_replication_spec_per_shard` has been removed. The data sources will now return only the new sharding configuration of the clusters.

3. `id` attribute which was an internal encoded resource identifier has been removed. Use `cluster_id` instead.


## How to migrate

If you currently use `mongodbatlas_cluster`, see our [Migration Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/cluster-to-advanced-cluster-migration-guide).

If you currently use `mongodbatlas_advanced_cluster` with the preview of the new schema [released in version 1.29.0](https://registry.terraform.io/providers/mongodb/mongodbatlas/1.29.0/docs/resources/advanced_cluster%2520%2528preview%2520provider%25202.0.0%2529) **and you are not using deprecated attributes**, you would not be required to make any additional changes except removing the `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true` environment variable, which is no longer required.

If you currently use `mongodbatlas_advanced_cluster` from v1.x.x of our provider, we recommend that you do the following steps:

~> **IMPORTANT:** Before you migrate, create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state). The state file will update to the new format and the old format will no longer be supported.

After you upgrade to v2.0.0+ from v1.x.x, when you run `terraform plan`, syntax errors will return as expected since the definition file hasn't been updated yet using the latest schema. At this point, you need to update the configuration by following all of below steps at once and finally running `terraform apply`:

- **Step #1:** Apply definition changes [explained on this page](#configuration-changes-when-upgrading-from-v1x) until there are no errors and no planned changes. 
  - **[Recommended]** You can also use the [Atlas CLI plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform?tab=readme-ov-file#2-advancedclustertov2-adv2v2) to generate the `mongodbatlas_advanced_cluster` resource definition. This is the recommended method as it will generate a clean configuration while keeping the original Terraform expressions. Please be aware of the [plugin limitations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform/blob/main/docs/command_adv2v2.md#limitations).

- **Step #2:** Remove any deprecated attributes (and their references) mentioned [above](#configuration-changes-when-upgrading-from-v1x). 

~> NOTE:  For nested attributes that have been removed, such as `replication_specs.#.num_shards` etc, Terraform may NOT throw an explicit error even if these attributes are left in the configuration. This is a [known Terraform issue](https://github.com/hashicorp/terraform-plugin-framework/issues/1210). Users should ensure to remove any such attributes from the configuration to avoid any confusion.

~> **IMPORTANT:** Don't apply until the plan is empty. If it shows other changes, you must update the `mongodbatlas_advanced_cluster` configuration until it matches the original configuration.

- **Step #3:** Even though there are no plan changes shown at this point, run `terraform apply`. This will update the `mongodbatlas_advanced_cluster` state to support the new schema.

## Important notes

Please refer to our [Considerations and Best Practices](../resources/advanced_cluster.md#considerations-and-best-practices) section for additional guidance on this resource.
