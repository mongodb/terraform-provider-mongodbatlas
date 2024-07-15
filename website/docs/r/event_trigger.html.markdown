# Resource: mongodbatlas_event_trigger

`mongodbatlas_event_trigger` provides a Event Trigger resource. 

Note: If the `app_id` changes in the mongodbatlas_event_trigger resource, it will force a replacement and delete itself from the old Atlas App Services app if it still exists then create itself in the new  Atlas App Services app. See [Atlas Triggers](https://www.mongodb.com/docs/atlas/app-services/triggers/) to learn more.   

## Example Usages

### Example Usage: Database Trigger with Function
```terraform
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

### Example Usage: Database Trigger with EventBridge
```terraform
resource "mongodbatlas_event_trigger" "test" {
  project_id = "PROJECT ID"
  app_id = "APPLICATION ID"
  name = "NAME OF THE TRIGGER"
  type = "DATABASE"
  disabled = false
  unordered = false
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

### Example Usage: Authentication Trigger
```terraform
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

### Example Usage: Scheduled Trigger
```terraform
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
    * For more details on `project_id` and `app_id` see: https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/#section/Project-and-Application-IDs
* `name` - (Required) The name of the trigger.
* `type` - (Required) The type of the trigger. Possible Values: `DATABASE`, `AUTHENTICATION`,`SCHEDULED`
* `function_id` - (Optional) The ID of the function associated with the trigger.
* `disabled` - (Optional) Default: `false` If `true`, the trigger is disabled.

* `config_operation_type` - Required for `AUTHENTICATION` type. The [authentication operation type](https://docs.mongodb.com/realm/triggers/authentication-triggers/#std-label-authentication-event-operation-types) to listen for. Possible Values: `LOGIN`, `CREATE`, `DELETE`
* `config_providers` - Required for `AUTHENTICATION` type. A list of one or more [authentication provider](https://docs.mongodb.com/realm/authentication/providers/) id values. The trigger will only listen for authentication events produced by these providers.

* `config_operation_types` - Required for `DATABASE` type. The [database event operation types](https://docs.mongodb.com/realm/triggers/database-triggers/#std-label-database-events) to listen for. This must contain at least one value. Possible Values: `INSERT`, `UPDATE`, `REPLACE`, `DELETE`
* `config_database` - Required for `DATABASE` type. The name of the MongoDB database to watch.
* `config_collection` - Optional for `DATABASE` type. The name of the MongoDB collection that the trigger watches for change events. The collection must be part of the specified database.
* `config_service_id` - Required for `DATABASE` type. The ID of the MongoDB Service associated with the trigger.
* `config_match` - Optional for `DATABASE` type. A [$match](https://docs.mongodb.com/manual/reference/operator/aggregation/match/) expression document that MongoDB Realm includes in the underlying change stream pipeline for the trigger. This is useful when you want to filter change events beyond their operation type. The trigger will only fire if the expression evaluates to true for a given change event.
* `config_project` - Optional for `DATABASE` type. A [$project](https://docs.mongodb.com/manual/reference/operator/aggregation/project/) expression document that Realm uses to filter the fields that appear in change event objects.
* `config_full_document` - Optional for `DATABASE` type. If true, indicates that `UPDATE` change events should include the most current [majority-committed](https://docs.mongodb.com/manual/reference/read-concern-majority/) version of the modified document in the fullDocument field.
* `unordered` - Only Available for Database Triggers. If true, event ordering is disabled and this trigger can process events in parallel. If false, event ordering is enabled and the trigger executes serially.

* `config_schedule` - Required for `SCHEDULED` type. A [cron expression](https://docs.mongodb.com/realm/triggers/cron-expressions/) that defines the trigger schedule.

* `event_processors` - (Optional) An object where each field name is an event processor ID and each value is an object that configures its corresponding event processor. The following event processors are supported: `AWS_EVENTBRIDGE` For an example configuration object, see [Send Trigger Events to AWS EventBridge](https://docs.mongodb.com/realm/triggers/eventbridge/#std-label-event_processor_example).
* `event_processors.0.aws_eventbridge.config_account_id` - (Optional) AWS Account ID.
* `event_processors.0.aws_eventbridge.config_region` - (Optional) Region of AWS Account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Terraform's unique identifier used internally for state management.
* `trigger_id` - The unique ID of the trigger.
* `function_name` - The name of the function associated with the trigger.

## Import

Event trigger can be imported using project ID, App ID and Trigger ID, in the format `project_id`--`app_id`-`trigger_id`, e.g.

```
$ terraform import mongodbatlas_event_trigger.test 1112222b3bf99403840e8934--testing-example--1112222b3bf99403840e8934
```
For more details on this resource see [Triggers resource](https://www.mongodb.com/docs/atlas/app-services/admin/api/v3/#tag/triggers) in Atlas App Services Documentation. 
