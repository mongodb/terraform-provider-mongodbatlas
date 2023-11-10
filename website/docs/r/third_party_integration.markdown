---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: third_party_integration"
sidebar_current: "docs-mongodbatlas-datasource-third-party-integration"
description: |-
     Provides a Third-Party Integration Settings resource.
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
  * `region` (Required) - Indicates which API URL to use, either "US", "EU", "US3", or "US5". Datadog will use "US" by default.    
* `OPS_GENIE`
  * `api_key` - Your API Key.
  * `region` (Required) -  Indicates which API URL to use, either "US" or "EU". OpsGenie will use "US" by default.
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
  * `scheme` - Your Prometheus protocol scheme configured for requests.
  * `enabled` - Whether your cluster has Prometheus enabled.

## Attributes Reference

* `id` - Unique identifier used by terraform for internal management, which can also be used to import.

## Import

Third-Party Integration Settings can be imported using project ID and the integration type, in the format `project_id`-`type`, e.g.

```
$ terraform import mongodbatlas_database_user.my_user 1112222b3bf99403840e8934-OPS_GENIE
```

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Third-Party-Integrations/operation/createThirdPartyIntegration) Documentation for more information.
