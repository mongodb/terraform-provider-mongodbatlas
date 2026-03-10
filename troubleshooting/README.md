# Troubleshooting

Common issues and debugging guidance for the MongoDB Atlas Terraform Provider.

## Enabling Debug Logging

Set the `TF_LOG` environment variable to get detailed provider logs:

```bash
export TF_LOG=DEBUG
terraform apply 2>&1 | tee terraform-debug.log
```

This outputs HTTP requests, responses, and internal provider operations, which are essential for diagnosing most issues.

For enhanced visibility into HTTP requests and responses when communicating with the MongoDB Atlas API, see [Enhanced Network Logging](network-logging.md).

## Common Issues

### Random ordering of elements in TypeList attributes during `terraform plan`

This occurs when dynamically adding objects to an attribute list (e.g., using `dynamic`). Terraform's `dynamic` block can bring objects into the schema in any order.

**Solutions:**

1. Define a static list of objects in your resource:

```terraform
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

2. Use a `type = list()` variable when using `dynamic`:

```terraform
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

### Service Account rate limiting

If you encounter a rate limit error when using Service Accounts:

```
Error: error initializing provider: oauth2: cannot fetch token: 429 Too Many Requests
Response: {"detail":"Resource /api/oauth/token is limited to 50 requests every 1 minutes.","error":429,...}
```

Atlas enforces rate limiting for each combination of IP address and SA client. Each Terraform operation generates a new token used for the duration of that operation.

**Solutions:**

- Contact [MongoDB Support](https://support.mongodb.com/) to request a rate limit increase.
- Create separate Service Accounts for different environments or CI/CD pipelines.
- Distribute Terraform executions across different IP addresses.
- Add retry logic to your automation workflows.

For more details, see [Service Account configuration](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/provider-configuration#service-account-recommended).

### Drift detection on `mongodbatlas_alert_configuration` `integration_id`

If you see unexpected drift on `notification.#.integration_id`:

```
~ notification {
          - integration_id  = "xxxxxxxxxxxxxxxxxxxxxxxx" -> null
```

This affects provider versions 1.16.0 to 1.19.0. The Atlas API returns a computed value for `integration_id` when none is set. See the [1.20.0 upgrade guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.20.0-upgrade-guide) for details.
