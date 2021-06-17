---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: event_trigger"
sidebar_current: "docs-mongodbatlas-resource-event-trigger"
description: |-
    Provides a Event Trigger resource.
---

# mongodbatlas_event_trigger

`mongodbatlas_event_trigger` provides a Event Trigger resource. 

## Example Usages

### Example Usage: DATABASE with function
```hcl
resource "mongodbatlas_event_trigger" "test" {
  project_id = "PROJECT ID"
  app_id = "APPLICATION ID"
  name = "NAME OF THE TRIGGER"
  type = "DATABASE"
  function_id = "FUNCTION ID"
  disabled = false
  config_operation_types = ["INSERT", "UPDATE"]
  config_database = "DATABASE NAME"
  config_collection = "COLLECTION NAME"
  config_service_id = "SERVICE ID"
  config_match = <<-EOF
{
  "updateDescription.updatedFields": {
    "status": "blocked"
  }
}
EOF
  config_project = "{\"updateDescription.updatedFields\":{\"status\":\"blocked\"}}"
  config_full_document = false
  config_full_document_before = false
  event_processors {
    aws_eventbridge {
      config_account_id = "AWS ACCOUNT ID"
      config_region = "AWS REGIOn"
    }
  }
}
```

### Example Usage: DATABASE with eventbridge
```hcl
resource "mongodbatlas_event_trigger" "test" {
  project_id = "PROJECT ID"
  app_id = "APPLICATION ID"
  name = "NAME OF THE TRIGGER"
  type = "DATABASE"
  disabled = false
  config_operation_types = ["INSERT", "UPDATE"]
  config_operation_type = "LOGIN"
  config_providers = ["anon-user"]
  config_database = "DATABASE NAME"
  config_collection = "COLLECTION NAME"
  config_service_id = "1"
  config_match = "{\"updateDescription.updatedFields\":{\"status\":\"blocked\"}}"
  config_project = "{\"updateDescription.updatedFields\":{\"status\":\"blocked\"}}"
  config_full_document = false
  config_full_document_before = false
  config_schedule = "*"
  event_processors {
    aws_eventbridge {
      config_account_id = "AWS ACCOUNT ID"
      config_region = "AWS REGIOn"
    }
  }
}
```

### Example Usage: AUTHENTICATION
```hcl
resource "mongodbatlas_event_trigger" "test" {
  project_id = "PROJECT ID"
  app_id = "APPLICATION ID"
  name = "NAME OF THE TRIGGER"
  type = "AUTHENTICATION"
  function_id = "1"
  disabled = false
  config_operation_type = "LOGIN"
  config_providers = ["anon-user"]
}
```

### Example Usage: SCHEDULED with eventbridge
```hcl
resource "mongodbatlas_event_trigger" "test" {
  project_id = "PROJECT ID"
  app_id = "APPLICATION ID"
  name = "NAME OF THE TRIGGER"
  type = "SCHEDULED"
  function_id = "1"
  disabled = false
  config_schedule = "*"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the trigger.
* `app_id` - (Required) The ObjectID of your application.
* `name` - (Required) The name of the trigger.
* `type` - (Required) The type of the trigger. Possible Values: `DATABASE`, `AUTHENTICATION`
* `function_id` - (Optional) The ID of the function associated with the trigger.
* `disabled` - (Optional) Default: `false` If `true`, the trigger is disabled.
* `config_operation_types` - (Optional) The [database event operation types](https://docs.mongodb.com/realm/triggers/database-triggers/#std-label-database-events) to listen for. This must contain at least one value. Required for `DATABASE` type. Possible Values: `INSERT`, `UPDATE`, `REPLACE`, `DELETE`
* `config_operation_type` - (Optional) The [authentication operation type](https://docs.mongodb.com/realm/triggers/authentication-triggers/#std-label-authentication-event-operation-types) to listen for. Required for `AUTHENTICATION` type. Possible Values: `LOGIN`, `CREATE`, `DELETE`
* `config_providers` - (Optional) A list of one or more [authentication provider](https://docs.mongodb.com/realm/authentication/providers/) id values. The trigger will only listen for authentication events produced by these providers. Required for `AUTHENTICATION` type.
* `config_database` - (Optional) The name of the MongoDB database that contains the watched collection. Required for `DATABASE` type.
* `config_collection` - (Optional) The name of the MongoDB collection that the trigger watches for change events. The collection must be part of the specified database. Required for `DATABASE` type.
* `config_service_id` - (Optional) The ID of the MongoDB Service associated with the trigger. Required for `DATABASE` type.
* `config_match` - (Optional) A [$match](https://docs.mongodb.com/manual/reference/operator/aggregation/match/) expression document that MongoDB Realm includes in the underlying change stream pipeline for the trigger. This is useful when you want to filter change events beyond their operation type. The trigger will only fire if the expression evaluates to true for a given change event. Required for `DATABASE` type.
* `config_project` - (Optional) A [$project](https://docs.mongodb.com/manual/reference/operator/aggregation/project/) expression document that Realm uses to filter the fields that appear in change event objects. Required for `DATABASE` type.
* `config_full_document` - (Optional) If true, indicates that `UPDATE` change events should include the most current [majority-committed](https://docs.mongodb.com/manual/reference/read-concern-majority/) version of the modified document in the fullDocument field.
* `config_schedule` - (Optional) A [cron expression](https://docs.mongodb.com/realm/triggers/cron-expressions/) that defines the trigger schedule.
* `event_processors` - (Optional) An object where each field name is an event processor ID and each value is an object that configures its corresponding event processor. The following event processors are supported: `AWS_EVENTBRIDGE` For an example configuration object, see [Send Trigger Events to AWS EventBridge](https://docs.mongodb.com/realm/triggers/eventbridge/#std-label-event_processor_example).
* `event_processors.0.aws_eventbridge.config_account_id` - (Optional) AWS Account ID.
* `event_processors.0.aws_eventbridge.config_region` - (Optional) Region of AWS Account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `trigger_id` - The unique ID of the trigger.
* `function_name` - The name of the function associated with the trigger.

## Import

Event trigger can be imported using project ID, App ID and Trigger ID, in the format `project_id`--`app_id`-`trigger_id`, e.g.

```
$ terraform import mongodbatlas_event_trigger.test 1112222b3bf99403840e8934--testing-example--1112222b3bf99403840e8934
```
