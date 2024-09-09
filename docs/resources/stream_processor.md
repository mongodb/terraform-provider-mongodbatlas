# Resource: mongodbatlas_stream_processor

`mongodbatlas_stream_processor` provides a Stream Processor resource. The resource lets you create, delete, import, start and stop a stream processor in a stream instance.

**NOTE**: Updating an Atlas Stream Processor is currently not supported. As a result, the following steps are needed to be able to change an Atlas Stream Processor with an Atlas Change Stream Source:
1. Retrieve the value of Change Stream Source Token `changeStreamState` from the computed `stats` attribute in `mongodbatlas_stream_processor` resource or datasource or from the Terraform state file. This takes the form of a [resume token](https://www.mongodb.com/docs/manual/changeStreams/#resume-tokens-from-change-events). The Stream Processor has to be running in the state `STARTED` for the `stats` attribute to be available. However, before you retrieve the value, you should set the `state` to `STOPPED` to get the latest `changeStreamState`.
	- Example:
		```
		{\"changeStreamState\":{\"_data\":\"8266C71670000000012B0429296E1404\"}
		```
2. Update the `pipeline` argument setting `config.StartAfter` with the value retrieved in the previous step. More details in the [MongoDB Collection Change Stream](https://www.mongodb.com/docs/atlas/atlas-stream-processing/sp-agg-source/#mongodb-collection-change-stream) documentation.
	- Example: 
		```
		pipeline = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-cluster.connection_name, "config" = { "startAfter" = { "_data" : "8266C71562000000012B0429296E1404" } } } }, { "$emit" = { "connectionName" : "__testLog" } }])
		```
3. Delete the existing Atlas Stream Processor and then create a new Atlas Stream Processor with updated pipeline parameter and the updated values.  

## Example Usages

```terraform
resource "mongodbatlas_stream_instance" "example" {
  project_id    = var.project_id
  instance_name = "InstanceName"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

resource "mongodbatlas_stream_connection" "example-sample" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "sample_stream_solar"
  type            = "Sample"
}

resource "mongodbatlas_stream_connection" "example-cluster" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "ClusterConnection"
  type            = "Cluster"
  cluster_name    = var.cluster_name
  db_role_to_execute = {
    role = "atlasAdmin"
    type = "BUILT_IN"
  }
}

resource "mongodbatlas_stream_connection" "example-kafka" {
  project_id      = var.project_id
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "KafkaPlaintextConnection"
  type            = "Kafka"
  authentication = {
    mechanism = "PLAIN"
    username  = var.kafka_username
    password  = var.kafka_password
  }
  bootstrap_servers = "localhost:9092,localhost:9092"
  config = {
    "auto.offset.reset" : "earliest"
  }
  security = {
    protocol = "PLAINTEXT"
  }
}

resource "mongodbatlas_stream_processor" "stream-processor-sample-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "sampleProcessorName"
  pipeline       = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-sample.connection_name } }, { "$emit" = { "connectionName" : "__testLog" } }])
  state          = "CREATED"
}

resource "mongodbatlas_stream_processor" "stream-processor-cluster-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "clusterProcessorName"
  pipeline       = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-cluster.connection_name } }, { "$emit" = { "connectionName" : "__testLog" } }])
  state          = "STARTED"
}

resource "mongodbatlas_stream_processor" "stream-processor-kafka-example" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "kafkaProcessorName"
  pipeline       = jsonencode([{ "$source" = { "connectionName" = resource.mongodbatlas_stream_connection.example-cluster.connection_name } }, { "$emit" = { "connectionName" : resource.mongodbatlas_stream_connection.example-kafka.connection_name, "topic" : "example_topic" } }])
  state          = "CREATED"
  options = {
    dlq = {
      coll            = "exampleColumn"
      connection_name = resource.mongodbatlas_stream_connection.example-cluster.connection_name
      db              = "exampleDb"
    }
  }
}

data "mongodbatlas_stream_processors" "example-stream-processors" {
  project_id    = var.project_id
  instance_name = mongodbatlas_stream_instance.example.instance_name
}

data "mongodbatlas_stream_processor" "example-stream-processor" {
  project_id     = var.project_id
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = mongodbatlas_stream_processor.stream-processor-sample-example.processor_name
}

# example making use of data sources
output "stream_processors_state" {
  value = data.mongodbatlas_stream_processor.example-stream-processor.state
}

output "stream_processors_results" {
  value = data.mongodbatlas_stream_processors.example-stream-processors.results
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance_name` (String) Human-readable label that identifies the stream instance.
- `pipeline` (String) Stream aggregation pipeline you want to apply to your streaming data. [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/atlas-stream-processing/stream-aggregation/#std-label-stream-aggregation) contain more information. Using [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) is recommended when settig this attribute. For more details see [Aggregation Pipelines Documentation](https://www.mongodb.com/docs/atlas/atlas-stream-processing/stream-aggregation/)
- `processor_name` (String) Human-readable label that identifies the stream processor.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.

**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.

### Optional

- `options` (Attributes) Optional configuration for the stream processor. (see [below for nested schema](#nestedatt--options))
- `state` (String) The state of the stream processor. Commonly occurring states are 'CREATED', 'STARTED', 'STOPPED' and 'FAILED'. Used to start or stop the Stream Processor. Valid values are `CREATED`, `STARTED` or `STOPPED`. When a Stream Processor is created without specifying the state, it will default to `CREATED` state.

**NOTE** When a stream processor is created, the only valid states are CREATED or STARTED. A stream processor can be automatically started when creating it if the state is set to STARTED.

### Read-Only

- `id` (String) Unique 24-hexadecimal character string that identifies the stream processor.
- `stats` (String) The stats associated with the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/atlas-stream-processing/manage-stream-processor/#view-statistics-of-a-stream-processor) for more information.

<a id="nestedatt--options"></a>
### Nested Schema for `options`

Required:

- `dlq` (Attributes) Dead letter queue for the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/reference/glossary/#std-term-dead-letter-queue) for more information. (see [below for nested schema](#nestedatt--options--dlq))

<a id="nestedatt--options--dlq"></a>
### Nested Schema for `options.dlq`

Required:

- `coll` (String) Name of the collection to use for the DLQ.
- `connection_name` (String) Name of the connection to write DLQ messages to. Must be an Atlas connection.
- `db` (String) Name of the database to use for the DLQ.

# Import 
Stream Processor resource can be imported using the Project ID, Stream Instance name and Stream Processor name, in the format `INSTANCE_NAME-PROJECT_ID-PROCESSOR_NAME`, e.g.
```
$ terraform import mongodbatlas_stream_processor.test yourInstanceName-6117ac2fe2a3d04ed27a987v-yourProcessorName
```

For more information see: [MongoDB Atlas API - Stream Processor](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Streams/operation/createStreamProcessor) Documentation.
