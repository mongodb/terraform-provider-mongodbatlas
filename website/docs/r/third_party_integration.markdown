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

-> **WARNING:** This field type has values (NEW_RELIC, FLOWDOCK) that are deprecated and will be removed in 1.9.0 release release

-> **NOTE:** Slack integrations now use the OAuth2 verification method and must be initially configured, or updated from a legacy integration, through the Atlas third-party service integrations page. Legacy tokens will soon no longer be supported.[Read more about slack setup](https://docs.atlas.mongodb.com/tutorial/third-party-service-integrations/)

~> **IMPORTANT** Each project can only have one configuration per {INTEGRATION-TYPE}.

~> **IMPORTANT:** All arguments including the secrets will be stored in the raw state as plain-text. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)


## Example Usage

```terraform

resource "mongodbatlas_third_party_integration" "test_flowdock" {
	project_id = "<PROJECT-ID>"
	type = "FLOWDOCK"
	flow_name = "<FLOW-NAME>"
	api_token = "<API-TOKEN>"
	org_name =  "<ORG-NAME>"
}

```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get all Third-Party service integrations
* `type`       - (Required) Third-Party Integration Settings type 
     * PAGER_DUTY
     * DATADOG
     * NEW_RELIC
     * OPS_GENIE
     * VICTOR_OPS
     * FLOWDOCK
     * WEBHOOK
     * MICROSOFT_TEAMS
     * PROMETHEUS

Additional values based on Type

* `PAGER_DUTY`
  * `service_key` - Your Service Key.
* `DATADOG`
  * `api_key` - Your API Key.
  * `region` - Indicates which API URL to use, either US or EU. Datadog will use US by default.    
* `NEW_RELIC`
  * `license_key` - Your License Key.
  * `account_id`  - Unique identifier of your New Relic account.
  * `write_token` - Your Insights Insert Key.
  * `read_token`  - Your Insights Query Key.
* `OPS_GENIE`
  * `api_key` - Your API Key.
  * `region` -  Indicates which API URL to use, either US or EU. Opsgenie will use US by default.
* `VICTOR_OPS`
  * `api_key` - 	Your API Key.
  * `routing_key` - An optional field for your Routing Key.
* `FLOWDOCK`
  * `flow_name` - Your Flowdock Flow name.
  * `api_token` - Your API Token.
  * `org_name` - Your Flowdock organization name.
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

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-create/) Documentation for more information.
