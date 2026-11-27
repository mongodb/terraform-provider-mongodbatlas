---
subcategory: "Projects"
---

# Resource: mongodbatlas_third_party_integration

`mongodbatlas_third_party_integration` Provides a Third-Party Integration Settings for the given type.

-> **NOTE:** Slack integrations now use the OAuth2 verification method and must be initially configured, or updated from a legacy integration, through the Atlas third-party service integrations page. Legacy tokens will soon no longer be supported.[Read more about slack setup](https://www.mongodb.com/docs/atlas/tutorial/third-party-service-integrations/)

~> **IMPORTANT** Each project can only have one configuration per integration `type`.

~> **IMPORTANT:** All arguments including the secrets will be stored in the raw state as plain-text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html) To keep a secret out of state, use its write-only alternative instead. See [Write-only secrets](#write-only-secrets).


## Example Usage

```terraform

resource "mongodbatlas_third_party_integration" "test_datadog" {
  project_id = "<PROJECT-ID>"
  type = "DATADOG"
  api_key = "<API-KEY>"
  region = "<REGION>"
}

```

### Example with a write-only secret

`password_wo` is never written to state or plan files. Here the password comes from an [ephemeral](https://developer.hashicorp.com/terraform/language/manage-sensitive-data/ephemeral) input variable, so Terraform doesn't keep it in state or plan files at all.

```terraform
variable "prometheus_password" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "mongodbatlas_third_party_integration" "prometheus" {
  project_id          = "<PROJECT-ID>"
  type                = "PROMETHEUS"
  user_name           = "<USERNAME>"
  password_wo         = var.prometheus_password
  password_wo_version = 1
  service_discovery   = "http"
  enabled             = true
}
```

### Further Examples
- [Third-Party Integration Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.22.0/examples/mongodbatlas_third_party_integration)

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get all Third-Party service integrations, also known as `groupId` in the official documentation
* `type`       - (Required) Third-Party Integration Settings type 
     * PAGER_DUTY
     * DATADOG
     * OPS_GENIE
     * VICTOR_OPS
     * WEBHOOK
     * MICROSOFT_TEAMS
     * PROMETHEUS
       

* `PAGER_DUTY`
  * `service_key` - Your Service Key.
  * `region` (Required) - PagerDuty region that indicates the API Uniform Resource Locator (URL) to use, either "US" or "EU". PagerDuty will use "US" by default.    
* `DATADOG`
  * `api_key` - Your API Key.
  * `region` (Required) - Two-letter code that indicates which API URL to use. See the `region` request parameter of [MongoDB API Third-Party Service Integration documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createthirdpartyintegration) for more details.
  * `send_collection_latency_metrics` - Toggle sending collection latency metrics that includes database names and collection names and latency metrics on reads, writes, commands, and transactions. Default: `false`.
  * `send_database_metrics` - Toggle sending database metrics that includes database names and metrics on the number of collections, storage size, and index size. Default: `false`.
  * `send_user_provided_resource_tags` - Toggle sending user provided group and cluster resource tags with the Datadog metrics. Default: `false`.
  * `send_query_stats_metrics` - Toggle sending query shape metrics that includes query hash and metrics on latency, execution frequency, documents returned, and timestamps. Default: `false`.
  * `send_sharding_metrics` - Toggle sending sharding metrics that includes sharding distribution and chunk metrics per cluster, shard, and collection. Default: `false`.
* `OPS_GENIE`
  * `api_key` - Your API Key.
  * `region` (Required) - Two-letter code that indicates which API URL to use. See the `region` request parameter of [MongoDB API Third-Party Service Integration documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-createthirdpartyintegration) for more details.
* `VICTOR_OPS`
  * `api_key` - 	Your API Key.
  * `routing_key` - An optional field for your Routing Key.
* `WEBHOOK`
  * `url` - Your webhook URL.
  * `secret` - An optional field for your webhook secret.
* `MICROSOFT_TEAMS`
  * `microsoft_teams_webhook_url` -  Your Microsoft Teams incoming webhook URL.
* `PROMETHEUS`
  * `user_name` - Your Prometheus username.
  * `password`  - Your Prometheus password.
  * `service_discovery` - Indicates which service discovery method is used, either file or http.
  * `enabled` - Whether your cluster has Prometheus enabled.

-> **NOTE:** For certain attributes with default values, it's recommended to explicitly set them back to their default instead of removing them from the configuration. For example, if `send_collection_latency_metrics` is set to `true` and you want to revert to the default (`false`), set it to `false` rather than removing it.

### Write-only secrets

Each secret argument has a [write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral/write-only) alternative that is never written to state or plan files. Write-only arguments require Terraform 1.11 or later.

| Secret | Write-only argument | Version argument |
|--------|---------------------|------------------|
| `api_key` | `api_key_wo` | `api_key_wo_version` |
| `service_key` | `service_key_wo` | `service_key_wo_version` |
| `routing_key` | `routing_key_wo` | `routing_key_wo_version` |
| `secret` | `secret_wo` | `secret_wo_version` |
| `microsoft_teams_webhook_url` | `microsoft_teams_webhook_url_wo` | `microsoft_teams_webhook_url_wo_version` |
| `password` | `password_wo` | `password_wo_version` |

* `<secret>_wo` - (Optional) Value of `<secret>`, passed as a write-only argument. Conflicts with `<secret>`, and must be set together with `<secret>_wo_version`. Where the integration type requires `<secret>`, set either `<secret>` or `<secret>_wo`. While `<secret>_wo` is set, `<secret>` is empty in state.
* `<secret>_wo_version` - (Optional) Integer, at least `1`, that triggers an update of `<secret>_wo`. To rotate the secret, change `<secret>_wo` and increment this value in the same edit. Changing `<secret>_wo` on its own has no effect, since Terraform can't detect changes to a value that is not in state.

Like `api_key`, `api_key_wo` is sent on every update, not only when `api_key_wo_version` changes. Set it from a value that is the same on every run, such as a secret read from a secrets manager, not one generated on each run.

## Attributes Reference

* `id` - Unique identifier of the integration.

## Import

Third-Party Integration Settings can be imported using project ID and the integration type, in the format `project_id`-`type`, e.g.

```
$ terraform import mongodbatlas_third_party_integration.test_datadog 1112222b3bf99403840e8934-DATADOG
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Third-Party-Integrations/operation/createThirdPartyIntegration) Documentation for more information.
