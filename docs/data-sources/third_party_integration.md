# Data Source: mongodbatlas_third_party_integration

`mongodbatlas_third_party_integration` describes a Third-Party Integration Settings for the given type.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform

resource "mongodbatlas_third_party_integration" "test_datadog" {
	project_id = "<PROJECT-ID>"
  type = "DATADOG"
	api_key = "<API-KEY>"
	region = "<REGION>"
}

data "mongodbatlas_third_party_integration" "test" {
	project_id = mongodbatlas_third_party_integration.test_datadog.project_id
  type = "DATADOG"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get all Third-Party service integrations
* `type`       - (Required) Third-Party service integration type
     * PAGER_DUTY
     * DATADOG
     * OPS_GENIE
     * VICTOR_OPS
     * WEBHOOK
     * MICROSOFT_TEAMS
     * PROMETHEUS

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier of the integration.

Additional values based on Type

* `PAGER_DUTY`
  * `service_key` - Your Service Key.
* `DATADOG`
  * `api_key` - Your API Key.
  * `region` - Two-letter code that indicates which API URL to use. See the `region` response field of [MongoDB API Third-Party Service Integration documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getthirdpartyintegration) for more details. Datadog will use "US" by default.
  * `send_collection_latency_metrics` - Toggle sending collection latency metrics that includes database names and collection name sand latency metrics on reads, writes, commands, and transactions.
  * `send_database_metrics` - Toggle sending database metrics that includes database names and metrics on the number of collections, storage size, and index size.
  * `send_user_provided_resource_tags` - Toggle sending user provided group and cluster resource tags with the datadog metrics.
* `OPS_GENIE`
  * `api_key` - Your API Key.
  * `region` - Two-letter code that indicates which API URL to use. See the `region` response field of [MongoDB API Third-Party Service Integration documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getthirdpartyintegration) for more details. Opsgenie will use US by default.
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
  * `password` - Your Prometheus password.
  * `service_discovery` - Indicates which service discovery method is used, either file or http.
  * `enabled` - Whether your cluster has Prometheus enabled.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Third-Party-Integrations/operation/createThirdPartyIntegration) Documentation for more information.
