---
page_title: "Migration Guide: Moving to Advanced Cluster v2.0.0"
---

**Objective**: This guide is for users who are currently using `mongodbatlas_advanced_cluster` and want to upgrade to our Terraform Provider v2.0.0 or later from v1.x.


### How to use this guide?
1. [Identify your current setup](#identify-your-current-setup) using the table below.
2. Review the [Overview of configuration changes when upgrading from 1.x](#overview-of-configuration-changes-when-upgrading-from-1x) to understand what changes have been made to `mongodbatlas_advanced_cluster` in provider v2.0.0. 
3. Review [Important Considerations](#important-considerations) before you start updating your configurations.
4. Follow the instructions based on your current setup from #1 above.



## Identify your current setup

| Case &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; | You are currently using `mongodbatlas_advanced_cluster`.. | What this means |
|------|--------------------------|-----------------|
| &nbsp; | &nbsp; | &nbsp; |
| [1(a)](#case-1a-currently-using-mongodbatlas_advanced_cluster-with-num_shards--older-blocks-syntax) | ..with `num_shards` | You are using the [older blocks syntax](#attribute-vs-block-syntax) and legacy sharding configuration. |
| &nbsp; | &nbsp; | &nbsp; |
| [1(b)](#case-1b-currently-using-mongodbatlas_advanced_cluster-without-num_shards--one-replication_specs-per-shard--older-blocks-syntax) | ..without `num_shards` and one `replication_specs` per shard | You are using the [older blocks syntax](#attribute-vs-block-syntax) and the [**new** independent shard scaling](./advanced-cluster-new-sharding-schema.md) configuration. |
| &nbsp; | &nbsp; | &nbsp; |
| [2(a)](#case-2a-currently-using-mongodbatlas_advanced_cluster-preview-for-v200-using-environment-flag-with-num_shards) | ..[**preview for v2.0.0**](https://registry.terraform.io/providers/mongodb/mongodbatlas/1.41.1/docs/resources/advanced_cluster%2520%2528preview%2520provider%25202.0.0%2529) (using environment flag) with `num_shards` | You are using the [new attributes syntax](#attribute-vs-block-syntax) and legacy sharding configuration. |
| &nbsp; | &nbsp; | &nbsp; |
| [2(b)](#case-2b-currently-using-mongodbatlas_advanced_cluster-preview-for-v200-using-environment-flag-without-num_shards) | ..[**preview for v2.0.0**](https://registry.terraform.io/providers/mongodb/mongodbatlas/1.41.1/docs/resources/advanced_cluster%2520%2528preview%2520provider%25202.0.0%2529) (using environment flag) without `num_shards` and one `replication_specs` per shard | You are using the [new attributes syntax](#attribute-vs-block-syntax) and the [**new** independent shard scaling](./advanced-cluster-new-sharding-schema.md) configuration. |
| &nbsp; | &nbsp; | &nbsp; |



If you are still using the deprecated `mongodbatlas_cluster` resource, use [Migration Guide: Cluster → Advanced Cluster instead](./cluster-to-advanced-cluster-migration-guide.md).



## Overview of configuration changes when upgrading from 1.x

1. Below deprecated attributes have been removed:
  - `id`
  - `disk_size_gb`
  - `replication_specs.#.num_shards`
  - `replication_specs.#.id`
  - `advanced_configuration.default_read_concern`
  - `advanced_configuration.fail_index_key_too_long`
  
2. `mongodbatlas_advanced_cluster` now supports **only** the new sharding configuration that allows scaling shards independently. To learn more about independent shard configuration, please review [Advanced Cluster New Sharding Configurations](./advanced-cluster-new-sharding-schema.md)
  
3. Elements `replication_specs` and `region_configs` are now [list attributes instead of blocks](#attribute-vs-block-syntax) so they are an array of objects.

4. Elements `connection_strings`, `timeouts`, `advanced_configuration`, `bi_connector_config`, `pinned_fcv`, `electable_specs`, `read_only_specs`, `analytics_specs`, `auto_scaling` and `analytics_auto_scaling` are now single [attributes instead of blocks](#attribute-vs-block-syntax) so they are an object. If there are references to them, the index such as`[0]` or `.0` are dropped.


5. Elements `tags` and `labels` are now `maps` instead of `blocks`.  For example,
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
becomes:
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


## Important Considerations
1. Before you migrate, create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state). The state file will update to the new format and the old format will no longer be supported.
2. For nested attributes that have been removed, such as `replication_specs.#.num_shards` etc, Terraform may NOT throw an explicit error even if these attributes are left in the configuration. This is a [known Terraform issue](https://github.com/hashicorp/terraform-plugin-framework/issues/1210). Users should ensure to remove any such attributes from the configuration to avoid any confusion.
3. For any of the migration paths mentioned in this document, do NOT apply until the plan is empty. If it shows other changes, you must update the `mongodbatlas_advanced_cluster` configuration until it matches the original configuration.
4. Please refer to our [Considerations and Best Practices](../resources/advanced_cluster.md#considerations-and-best-practices) section for additional guidance on this resource.




## Case 1(a): Currently using `mongodbatlas_advanced_cluster` with `num_shards` & older blocks syntax

After you upgrade to v2.0.0+ from v1.x.x, when you run `terraform plan`, you may run into syntax errors such as `Required attribute "replication_specs" not specified` or `Unexpected block: Blocks of type "advanced_configuration" are not expected here`. **This is expected** since the resource definition hasn't been updated yet using the latest schema.
At this point, you need to update the configuration by following all of below steps at once and finally running terraform apply:

- **Step #1:** Apply **all** configuration [changes mentioned here](#overview-of-configuration-changes-when-upgrading-from-1x) **until there are no errors and no planned changes**. 
This may include:
    - Adjust syntax of [blocks to attributes](#attribute-vs-block-syntax).
    - Remove `num_shards` and restructure `replication_specs` list to represent one element per shard instead. This ensures you use the new sharding configuration with independent scaling going forward. To learn more, please review [Advanced Cluster New Sharding Configurations](./advanced-cluster-new-sharding-schema.md).
    - Remove references to deprecated attributes mentioned [in this section](#overview-of-configuration-changes-when-upgrading-from-1x).
    
**Recommended:** Use the [Atlas CLI plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform?tab=readme-ov-file#2-advancedclustertov2-adv2v2) to generate the `mongodbatlas_advanced_cluster` resource definition as it will generate a clean configuration while keeping the original Terraform expressions. Please ensure to review [plugin limitations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform/blob/main/docs/command_adv2v2.md#limitations).
  
- **Step #2:**  Ensure all deprecated attributes (and their references) mentioned [above](#overview-of-configuration-changes-when-upgrading-from-1x) are removed.

- **Step #3:** Even though there are no plan changes shown at this point, run `terraform apply`. This will update the `mongodbatlas_advanced_cluster` state to support the new schema.

### Example for migrating a `SHARDED` cluster using `num_shards` (see [complete example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_advanced_cluster/symmetric-sharded-cluster/v1.x.x/README.md)):
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  ...
  disk_size_gb   = 10 # remove this and set per shard for inner specs

  replication_specs {   # update syntax to a list attribute of objects instead of blocks
    num_shards = 2  # remove this & add another replication_spec element for the second shard

    region_configs {    # update syntax to a list attribute of objects instead of blocks
      electable_specs { # update syntax to an attribute instead of a block
        instance_size = "M30"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "zone n1"
  }
  advanced_configuration { # update syntax to an attribute instead of a block
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags { # update syntax for tags and labels to maps instead of blocks in v2.0.0+
    key   = "environment"
    value = "dev"
  }
}
```
- After:
```hcl
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  ...

  replication_specs = [{    # notice the list of objects and the "=" sign
    region_configs = [{
      electable_specs = {   # notice the "=" sign
        instance_size = "M30"
        disk_iops     = 3000
        node_count    = 3
        disk_size_gb  = 10 # this is now set at spec level
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }]
    },
    {
      region_configs = [{
        electable_specs = {
          instance_size = "M30"
          disk_iops     = 3000
          node_count    = 3
          disk_size_gb  = 10
        }
        provider_name = "AWS"
        priority      = 7
        region_name   = "EU_WEST_1"
      }]
  }]

  advanced_configuration = {    # notice the "=" sign
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags = {      # notice the "=" sign and map syntax
    environment = "dev"
  }
}
```


### Example for migrating a `REPLICASET` cluster using `num_shards = 1` (see [complete example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_advanced_cluster/replicaset/v1.x.x/README.md)):
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type = "REPLICASET"
  ...
  
  replication_specs {   # update syntax to a list attribute of objects instead of blocks
    num_shards = 1    # remove this

    region_configs { # update syntax to a list attribute of objects instead of blocks
      electable_specs { # update syntax to an attribute instead of a block
        instance_size = var.provider_instance_size_name
        node_count    = 3
      }
      provider_name = var.provider_name
      region_name   = "US_EAST_1"
      priority      = 7
    }
  }

  tags {    # update syntax for tags and labels to maps instead of blocks in v2.0.0+
    key   = "environment"
    value = "dev"
  }
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type = "REPLICASET"

  replication_specs = [ # notice the list of objects and the "=" sign
    {
      region_configs = [
        {
          electable_specs = {   # notice the "=" sign
            instance_size = var.provider_instance_size_name
            node_count    = 3
          }
          provider_name = var.provider_name
          region_name   = "US_EAST_1"
          priority      = 7
        }
      ]
    }
  ]

  tags = {  # notice the "=" sign and map syntax
    environment = "dev"
  }
}
```

### Example for migrating a `GEOSHARDED` cluster using `num_shards`:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"
  disk_size_gb = 15

  replication_specs {
    num_shards = 2  # remove this & add another replication_spec element for the second shard
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
    zone_name = "zone 1"
  }

  replication_specs {
    num_shards = 2  # remove this & add another replication_spec element for the second shard
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "zone 2"
  }
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
  }]
}
```


## Case 1(b): Currently using `mongodbatlas_advanced_cluster` without `num_shards`,  one `replication_specs` per shard & older blocks syntax

After you upgrade to v2.0.0+ from v1.x.x, when you run `terraform plan`, you may run into syntax errors such as `Required attribute "replication_specs" not specified` or `Unexpected block: Blocks of type "advanced_configuration" are not expected here`. **This is expected** since the resource definition hasn't been updated yet using the latest schema.
At this point, you need to update the configuration by following all of below steps **at once** and finally running `terraform apply`:

- **Step #1:** Apply **all** configuration [changes mentioned here except #2](#overview-of-configuration-changes-when-upgrading-from-1x) since you are already using the new sharding configuration, **until there are no errors and no planned changes**. This may include:
    - Adjust syntax of [blocks to attributes](#attribute-vs-block-syntax).
    - Remove references to deprecated attributes mentioned [in this section](#overview-of-configuration-changes-when-upgrading-from-1x).
    
**Recommended:** Use the [Atlas CLI plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform?tab=readme-ov-file#2-advancedclustertov2-adv2v2) to generate the `mongodbatlas_advanced_cluster` resource definition as it will generate a clean configuration while keeping the original Terraform expressions. Please ensure to review [plugin limitations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform/blob/main/docs/command_adv2v2.md#limitations).
  
- **Step #2:**  Ensure all deprecated attributes (and their references) mentioned [above](#overview-of-configuration-changes-when-upgrading-from-1x) are removed.

- **Step #3:** Even though there are no plan changes shown at this point, run `terraform apply`. This will update the `mongodbatlas_advanced_cluster` state to support the new schema.

### Example for migrating a `SHARDED` cluster using new sharding configuration (see [complete example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_advanced_cluster/asymmetric-sharded-cluster/main.tf)):
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  ...
  disk_size_gb   = 10 # remove this and set per shard for inner specs

  replication_specs {   # shard 1 - M30 instance size
    region_configs {    # update syntax to a list attribute of objects instead of blocks
      electable_specs { # update syntax to an attribute
        instance_size = "M30"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # shard 2 - M40 instance size
    region_configs {
      electable_specs {
        instance_size = "M40"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  advanced_configuration { # update syntax to an attribute instead of a block
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags { # update syntax for tags and labels to maps instead of blocks in v2.0.0+
    key   = "environment"
    value = "dev"
  }
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  ...

  replication_specs = [ # notice the list of objects and the "=" sign
    { # shard 1 - M30 instance size
      region_configs = [
        {
          electable_specs = {   # notice the "=" sign
            instance_size = "M30"
            disk_iops     = 3000
            node_count    = 3
            disk_size_gb  = 10   # this is now set at spec level
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    },
    { # shard 2 - M40 instance size
      region_configs = [
        {
          electable_specs = {
            instance_size = "M40"
            disk_iops     = 3000
            node_count    = 3
            disk_size_gb  = 10
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "EU_WEST_1"
        }
      ]
    }
  ]

  advanced_configuration = {    # notice the "=" sign
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags = {       # notice the "=" sign and map syntax
    environment = "dev"
  }
}
```

### Example for migrating a `GEOSHARDED` cluster using new sharding configuration:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
    zone_name = "zone 1"
  }

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }
    zone_name = "zone 1"
  }

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "zone 2"
  }

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "zone 2"
  }
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
  }]
}
```

### Example for migrating a `REPLICASET` cluster:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = "true"
  disk_size_gb           = 60

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M10"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_WEST_2"
    }
  }
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = "true"

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb  = 60
      }
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb  = 60
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
  }]
}
```







## Case 2(a): Currently using `mongodbatlas_advanced_cluster` **preview** for v2.0.0 (using environment flag) with `num_shards`

In this case, your configuration is already using the [new attributes syntax](#attribute-vs-block-syntax). Next, you need to remove `num_shards` and restructure the `replication_specs` list to represent one element per shard. Additionally, remove references of the other attributes that have been removed in v2.0.0.

- **Step #1:** Remove the `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER` environment variable, as it is no longer required.

- **Step #2:** Apply relevant configuration [changes mentioned here](#overview-of-configuration-changes-when-upgrading-from-1x) **until there are no errors and no planned changes**. 
This may include:
    - Remove `num_shards` and restructure `replication_specs` list to represent one element per shard. This ensures you use the new sharding configuration with independent scaling going forward. To learn more, please review [Advanced Cluster New Sharding Configurations](./advanced-cluster-new-sharding-schema.md).
    - Remove references to deprecated attributes mentioned [in this section](#overview-of-configuration-changes-when-upgrading-from-1x).

**Recommended:** Use the [Atlas CLI plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform?tab=readme-ov-file#2-advancedclustertov2-adv2v2) to generate the `mongodbatlas_advanced_cluster` resource definition as it will generate a clean configuration while keeping the original Terraform expressions. Please ensure to review [plugin limitations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform/blob/main/docs/command_adv2v2.md#limitations).

- **Step #3:**  Ensure all deprecated attributes (and their references) mentioned [above](#overview-of-configuration-changes-when-upgrading-from-1x) are removed.

- **Step #4:** Even though there are no plan changes shown at this point, run `terraform apply`. This will update the `mongodbatlas_advanced_cluster` state to support the new schema.

### Example for migrating a `SHARDED` cluster using `num_shards` & legacy sharding configuration:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "SHARDED"
  disk_size_gb = 15

  replication_specs = [{
    num_shards = 2
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      auto_scaling = {
        compute_enabled           = true
        compute_max_instance_size = "M20"
        disk_gb_enabled           = true
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
  }]
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "SHARDED"

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      auto_scaling = {
        compute_enabled           = true
        compute_max_instance_size = "M20"
        disk_gb_enabled           = true
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      auto_scaling = {
        compute_enabled           = true
        compute_max_instance_size = "M20"
        disk_gb_enabled           = true
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
  }]
}
```

### Example for migrating a `GEOSHARDED` cluster using `num_shards` & legacy sharding configuration:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"
  disk_size_gb = 15

  replication_specs = [{
    num_shards = 2
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }]
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }]
}
```

### Example for migrating a `REPLICASET` cluster using `num_shards`:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = "true"
  disk_size_gb  = 60

  replication_specs = [{
    num_shards = 1
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
  }]
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = "true"
  disk_size_gb  = 60

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
  }]
}
```


## Case 2(b): Currently using `mongodbatlas_advanced_cluster` **preview** for v2.0.0 (using environment flag) without `num_shards`

In this case, after you upgrade to v2.0.0+ from v1.x.x and run `terraform plan`, you should not run into any errors or unexpected planned changes **unless you are using any of the deprecated attributes** mentioned [here](#overview-of-configuration-changes-when-upgrading-from-1x) as you are already using the [new attribute](#attribute-vs-block-syntax) syntax and the new sharding configuration.

- **Step #1:** Remove the `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER` environment variable, as it is no longer required.

- **Step #2:**  Ensure all deprecated attributes (and their references) mentioned [above](#overview-of-configuration-changes-when-upgrading-from-1x) are removed.

- **Step #3:** Even though there are no plan changes shown at this point, run `terraform apply`. This will update the `mongodbatlas_advanced_cluster` state to support the new schema.


### Example for migrating a `SHARDED` cluster using new sharding configuration:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "SHARDED"
  disk_size_gb = 15

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
  }]
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "SHARDED"

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb = 15
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb = 15
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
  }]
}
```

### Example for migrating a `GEOSHARDED` cluster using new sharding configuration:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"
  disk_size_gb = 15

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
  }]
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  backup_enabled = false
  cluster_type   = "GEOSHARDED"

  replication_specs = [{
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb = 15
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb = 15
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb = 15
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb = 15
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
    }]
    zone_name = "zone 1"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb = 15
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb = 15
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
    }, {
    region_configs = [{
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb = 15
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb = 15
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "EU_WEST_1"
    }]
    zone_name = "zone 2"
  }]
}
```

### Example for migrating a `REPLICASET` cluster using new sharding configuration:
- Before:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = "true"
  disk_size_gb  = 60

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
  }]
}
```

- After:
```hcl
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type           = "REPLICASET"
  retain_backups_enabled = "true"

  replication_specs = [{
    region_configs = [{
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
        disk_size_gb  = 60
      }
      analytics_specs = {
        instance_size = "M10"
        node_count    = 1
        disk_size_gb  = 60
      }
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_2"
    }]
  }]
}
```


### Attribute vs Block Syntax
In older versions, some attributes (like `replication_specs` and `region_configs`) were written as nested **blocks**, for example:

```hcl
replication_specs {
  region_configs {
    provider_name = "AWS"
    region_name   = "US_EAST_1"
  }
}
replication_specs {
  region_configs {
    provider_name = "AWS"
    region_name   = "US_EAST_2"
  }
}

advanced_configuration {
  default_write_concern = "majority"
  javascript_enabled    = true
}  
```

Newer Terraform plugin implementations, requires these to be expressed as **attributes** (with an `=` sign) instead using a list or map syntax:

```hcl
replication_specs = [{
  region_configs = [{
    provider_name = "AWS"
    region_name   = "US_EAST_1"
  }]
},
{
  region_configs = [{
    provider_name = "AWS"
    region_name   = "US_EAST_1"
  }]
}]

advanced_configuration = {
  default_write_concern = "majority"
  javascript_enabled    = true
} 
```
For single attributes such as `advanced_configuration` above, if you have references to them, the index such as`[0]` or `.0` will need to be dropped.
