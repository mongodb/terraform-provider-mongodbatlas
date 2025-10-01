---
subcategory: "Projects"
---

# Resource: mongodbatlas_third_party_integration

`mongodbatlas_third_party_integration` Provides a Third-Party Integration Settings for the given type.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

-> **NOTE:** Slack integrations now use the OAuth2 verification method and must be initially configured, or updated from a legacy integration, through the Atlas third-party service integrations page. Legacy tokens will soon no longer be supported.[Read more about slack setup](https://docs.atlas.mongodb.com/tutorial/third-party-service-integrations/)

~> **IMPORTANT** Each project can only have one configuration per {INTEGRATION-TYPE}.

~> **IMPORTANT:** All arguments including the secrets will be stored in the raw state as plain-text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)


## Example Usage

```terraform

resource "mongodbatlas_third_party_integration" "test_datadog" {
  project_id = "<PROJECT-ID>"
  type = "DATADOG"
  api_key = "<API-KEY>"
  region = "<REGION>"
}

```

### Further Examples
- [Third-Party Integration Examples](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v1.41.1/examples/mongodbatlas_third_party_integration)

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get all Third-Party service integrations
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
  * `send_user_provided_resource_tags` - Toggle sending user provided group and cluster resource tags with the datadog metrics. Default: `false`.
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

## Attributes Reference

* `id` - Unique identifier of the integration.

## Import

Third-Party Integration Settings can be imported using project ID and the integration type, in the format `project_id`-`type`, e.g.

```
$ terraform import mongodbatlas_third_party_integration.test_datadog 1112222b3bf99403840e8934-DATADOG
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Third-Party-Integrations/operation/createThirdPartyIntegration) Documentation for more information.
