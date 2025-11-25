# Troubleshooting

The following are some of the common issues/errors encountered when using Terraform Provider for MongoDB Atlas:


## Issue: The order of element objects in a TypeList attribute randomly changes every time a user runs `terraform plan`:

### Cause:
This issue occurs if the user tries to dynamically add objects to an attribute list (for example, by using `dynamic`). This is a known Terraform behavior, as `dynamic` can attempt to bring objects into the schema in any order. 

This can be resolved by:

1. Defining a static list of objects in your resource as shown in the example below:

```
resource "mongodbatlas_advanced_cluster" "main" {
  name         = "advanced-cluster-1"
  project_id   = "64258fba5c9...e5e94617e"
  cluster_type = "REPLICASET"

  replication_specs = [
    {
      region_configs = [
        {
          electable_specs = {
            instance_size = "M20"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 7
          region_name   = "US_EAST_1"
        },
        {
          electable_specs = {
            instance_size = "M20"
            node_count    = 1
          }
          provider_name = "AWS"
          priority      = 6
          region_name   = "EU_WEST_1"
        }
      ]
    }
  ]
}
```

2. Using a `type = list()` variable when using `dynamic` as shown in the example below:

```
variable "region_configs_list" {
  description = "List of region_configs"
  type = list(object({
    provider_name = string
    priority      = number
    region_name   = string
    electable_specs = list(object({
      instance_size = string
      node_count    = number
    }))
  }))
  default = [{
    provider_name = "AWS",
    priority      = 7,
    region_name   = "US_EAST_1",
    electable_specs = [{
      instance_size = "M20"
      node_count    = 1
    }]
    }
  ]
}

```

## Issue: `mongodbatlas_alert_configuration` resource attribute `notification.#.integration_id` causes drift detection (non-empty plan) as shown below if no value is set:
```
~ notification {
          - integration_id  = "xxxxxxxxxxxxxxxxxxxxxxxx" -> null
```

### Cause:
Due to recent updates in the Atlas API [Alert Configuration](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getalertconfiguration), a computed value for `integration_id` might be returned. This affects MongoDB Atlas Provider versions **1.16.0 to 1.19.0**. Please review the Breaking Changes section in our [1.20.0 upgrade guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.20.0-upgrade-guide) to learn more.
