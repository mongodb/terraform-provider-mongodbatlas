---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: third_party_integration"
sidebar_current: "docs-mongodbatlas-datasource-third-party-integration"
description: |-
    Describes all Third-Party Integration Settings in the project.
---

# Data Source: mongodbatlas_third_party_integration

`mongodbatlas_third_party_integration` describe a Third-Party Integration Settings for the given type.

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
  * `region` - Two-letter code that indicates which API URL to use. See the `region` response field of [MongoDB API Third-Party Service Integration documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Third-Party-Integrations/operation/getThirdPartyIntegration) for more details. Datadog will use "US" by default.
* `OPS_GENIE`
  * `api_key` - Your API Key.
  * `region` - Two-letter code that indicates which API URL to use. See the `region` response field of [MongoDB API Third-Party Service Integration documentation](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Third-Party-Integrations/operation/getThirdPartyIntegration) for more details. Opsgenie will use US by default.
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
  * `scheme` - Your Prometheus protocol scheme configured for requests. **Note** This attribute is deprecated as it is not being used.
  * `enabled` - Whether your cluster has Prometheus enabled.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Third-Party-Integrations/operation/createThirdPartyIntegration) Documentation for more information.
