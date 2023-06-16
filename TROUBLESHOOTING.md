# Troubleshooting

The following are some of the common issues/errors encountered when using Terraform Provider for MongoDB Atlas :


## Issue: Order of element objects in a TypeList attribute randomly changes everytime user runs `terraform plan`:

### Cause:
This problem occurs if the user is trying to dynamically add object to an attribute List, for example, by using `dynamic`. This is a known Terraform behavior as `dynamic` can attempt bring objects into schema in any order. 

This can be resolved by:

1. Defining a static list of objects in your resource as shown in the example below:

```
resource "mongodbatlas_advanced_cluster" "main" {
  name         = "advanced-cluster-1"
  project_id   = "64258fba5c9...e5e94617e"
  cluster_type = "REPLICASET"

  replication_specs {
    region_configs {
      electable_specs {
        instance_size = "M20"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "US_EAST_1"
    }

    region_configs {
      electable_specs {
        instance_size = "M20"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 6
      region_name   = "EU_WEST_1"
    }
  }
}
```

2. User can try to maintain order when using `dynamic` by using a `type = list()` variable as shown in the example below:

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

