---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: third_party_integration"
sidebar_current: "docs-mongodbatlas-datasource-third-party-integration"
description: |-
    Describes all Third-Party Integration Settings in the project.
---

# mongodbatlas_third_party_integration

`mongodbatlas_third_party_integration` describe a Third-Party Integration Settings for the given type.

-> **NOTE:** Groups and projects are synonymous terms. You may find `groupId` in the official documentation.

## Example Usage

```terraform

resource "mongodbatlas_third_party_integration" "test_flowdock" {
	project_id = "<PROJECT-ID>"
	type = "FLOWDOCK"
	flow_name = "<FLOW-NAME>"
	api_token = "<API-TOKEN>"
	org_name =  "<ORG-NAME>"
}

data "mongodbatlas_third_party_integration" "test" {
	project_id = mongodbatlas_third_party_integration.test_flowdock.project_id
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to get all Third-Party service integrations
* `type`       - (Required) Third-Party service integration type
     * PAGER_DUTY
     * DATADOG
     * NEW_RELIC
     * OPS_GENIE
     * VICTOR_OPS
     * FLOWDOCK
     * WEBHOOK
     * MICROSOFT_TEAMS
     * PROMETHEUS


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Unique identifier used for terraform for internal manages and can be used to import.
* `type` -  Property equal to its own integration type

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
    * `password` - Your Prometheus password.
    * `service_discovery` - Indicates which service discovery method is used, either file or http.
    * `scheme` - Your Prometheus protocol scheme configured for requests.
    * `enabled` - Whether your cluster has Prometheus enabled.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/third-party-integration-settings-get-one/) Documentation for more information.