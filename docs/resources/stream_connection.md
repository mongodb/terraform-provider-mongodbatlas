---
subcategory: "Streams"
---

# Resource: mongodbatlas_stream_connection

`mongodbatlas_stream_connection` provides a Stream Connection resource. The resource lets you create, edit, and delete stream instance connections.

~> **IMPORTANT:** All arguments including the Kafka authentication password will be stored in the raw state as plaintext. [Read more about sensitive data in state.](https://www.terraform.io/docs/state/sensitive-data.html)


## Example Usage

### Example Cluster Connection

```terraform
resource "mongodbatlas_stream_connection" "test" {
    project_id = var.project_id
    workspace_name = "WorkspaceName"
    connection_name = "ConnectionName"
    type = "Cluster"
    cluster_name = "Cluster0"
}
```

### Further Examples
- [Atlas Stream Connection](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_stream_connection)

### Example Cross Project Cluster Connection

```terraform
resource "mongodbatlas_stream_connection" "test" {
    project_id         = var.project_id
    workspace_name      = "WorskpaceName"
    connection_name    = "ConnectionName"
    type               = "Cluster"
    cluster_name       = "OtherCluster"
    cluster_project_id = var.other_project_id
}
```

### Example Kafka SASL Plaintext Connection

```terraform
resource "mongodbatlas_stream_connection" "test" {
    project_id = var.project_id
    workspace_name = "NewWorkspace"
    connection_name = "KafkaConnection"
    type = "Kafka"
    authentication = {
        mechanism = "SCRAM-256"
        username = "user"
        password = "somepassword"
    }
    security = {
        protocol = "SASL_PLAINTEXT"
    }
    config = {
        "auto.offset.reset": "latest"
    }
    bootstrap_servers = "localhost:9091,localhost:9092"
}    
```

### Example Kafka SASL OAuthbearer Connection

```terraform
resource "mongodbatlas_stream_connection" "example-kafka-oauthbearer" {
    project_id      = var.project_id
    workspace_name  = mongodbatlas_stream_workspace.example.workspace_name
    connection_name = "KafkaOAuthbearerConnection"
    type            = "Kafka"
    authentication = {
        mechanism = "OAUTHBEARER"
        method = "OIDC"
        token_endpoint_url = "https://example.com/oauth/token"
        client_id  = "auth0Client"
        client_secret  = var.kafka_client_secret
        scope = "read:messages write:messages"
        sasl_oauthbearer_extensions = "logicalCluster=lkc-kmom,identityPoolId=pool-lAr"
    }
    bootstrap_servers = "localhost:9092,localhost:9092"
    config = {
        "auto.offset.reset" : "earliest"
    }
    security = {
        protocol = "SASL_PLAINTEXT"
    }
    networking = {
        access = {
        type = "PUBLIC"
        }
    }
}
```

### Example Kafka SASL SSL Connection

```terraform
resource "mongodbatlas_stream_connection" "test" {
    project_id = var.project_id
    workspace_name = "NewWorkspace"
    connection_name = "KafkaConnection"
    type = "Kafka"
    authentication = {
        mechanism = "PLAIN"
        username = "user"
        password = "somepassword"
    }
    security = {
        protocol = "SASL_SSL"
        broker_public_certificate = "-----BEGIN CERTIFICATE-----<CONTENT>-----END CERTIFICATE-----"
    }
    config = {
        "auto.offset.reset": "latest"
    }
    bootstrap_servers = "localhost:9091,localhost:9092"
}
```

### Example AWSLambda Connection

```terraform
resource "mongodbatlas_stream_connection" "test" {
    project_id      = var.project_id
    workspace_name   = "NewWorkspace"
    connection_name = "AWSLambdaConnection"
    type            = "AWSLambda"
    aws             = {
      role_arn = "arn:aws:iam::<AWS_ACCOUNT_ID>:role/lambdaRole"
    }
}

```

### Example Https Connection

```terraform
resource "mongodbatlas_stream_connection" "example-https" {
  project_id      = var.project_id
  workspace_name  = mongodbatlas_stream_workspace.example.workspace_name
  connection_name = "https_connection_tf_new"
  type            = "Https"
  url             = "https://example.com"
  headers = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

### Example Schema Registry Connection with USER_INFO Authentication

```terraform
resource "mongodbatlas_stream_connection" "example-schema-registry" {
  project_id               = var.project_id
  workspace_name           = mongodbatlas_stream_workspace.example.workspace_name
  connection_name          = "SchemaRegistryConnection"
  type                     = "SchemaRegistry"
  schema_registry_provider = "CONFLUENT"
  schema_registry_urls     = ["https://schema-registry.example.com:8081"]
  schema_registry_authentication = {
    type     = "USER_INFO"
    username = "registry-user"
    password = var.schema_registry_password
  }
}
```

### Example Schema Registry Connection with SASL_INHERIT Authentication

```terraform
resource "mongodbatlas_stream_connection" "example-schema-registry-sasl" {
  project_id               = var.project_id
  workspace_name           = mongodbatlas_stream_workspace.example.workspace_name
  connection_name          = "SchemaRegistryConnectionSASL"
  type                     = "SchemaRegistry"
  schema_registry_provider = "CONFLUENT"
  schema_registry_urls     = ["https://schema-registry.example.com:8081"]
  schema_registry_authentication = {
    type = "SASL_INHERIT"
  }
}
```

### Example Usage with Stream Processor

When using a stream connection with a stream processor, the connection must be fully provisioned before the processor can be created. The provider automatically waits for connections to be ready after creation or updates. The example below shows a typical pattern:

```terraform
resource "mongodbatlas_stream_workspace" "example" {
  project_id     = var.project_id
  workspace_name = "ExampleWorkspace"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

# Source connection (Sample data)
resource "mongodbatlas_stream_connection" "source" {
  project_id      = var.project_id
  workspace_name  = mongodbatlas_stream_workspace.example.workspace_name
  connection_name = "sample_stream_solar"
  type            = "Sample"
}

# Sink connection (Atlas Cluster)
resource "mongodbatlas_stream_connection" "sink" {
  project_id      = var.project_id
  workspace_name  = mongodbatlas_stream_workspace.example.workspace_name
  connection_name = "ClusterConnection"
  type            = "Cluster"
  cluster_name    = mongodbatlas_cluster.example.name
  db_role_to_execute = {
    role = "atlasAdmin"
    type = "BUILT_IN"
  }
}

# Stream processor that depends on both connections
resource "mongodbatlas_stream_processor" "example" {
  project_id     = var.project_id
  workspace_name = mongodbatlas_stream_workspace.example.workspace_name
  processor_name = "ExampleProcessor"
  pipeline = jsonencode([
    { "$source" = { "connectionName" = mongodbatlas_stream_connection.source.connection_name } },
    { "$emit" = { "connectionName" = mongodbatlas_stream_connection.sink.connection_name } }
  ])
  state = "STARTED"
}
```

~> **NOTE:** The stream processor resource automatically depends on the stream connections through the `connection_name` references in the pipeline. This ensures proper creation order. The provider waits for each connection to be fully provisioned before returning from create or update operations.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_name` (String)
- `project_id` (String)
- `type` (String)

### Optional

- `authentication` (Attributes) (see [below for nested schema](#nestedatt--authentication))
- `aws` (Attributes) (see [below for nested schema](#nestedatt--aws))
- `bootstrap_servers` (String)
- `cluster_name` (String)
- `cluster_project_id` (String)
- `config` (Map of String)
- `db_role_to_execute` (Attributes) (see [below for nested schema](#nestedatt--db_role_to_execute))
- `headers` (Map of String)
- `instance_name` (String, Deprecated)
- `networking` (Attributes) (see [below for nested schema](#nestedatt--networking))
- `schema_registry_authentication` (Attributes) (see [below for nested schema](#nestedatt--schema_registry_authentication))
- `schema_registry_provider` (String)
- `schema_registry_urls` (List of String)
- `security` (Attributes) (see [below for nested schema](#nestedatt--security))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `url` (String)
- `workspace_name` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--authentication"></a>
### Nested Schema for `authentication`

Optional:

- `client_id` (String)
- `client_secret` (String, Sensitive)
- `mechanism` (String)
- `method` (String)
- `password` (String, Sensitive)
- `sasl_oauthbearer_extensions` (String)
- `scope` (String)
- `token_endpoint_url` (String)
- `username` (String)


<a id="nestedatt--aws"></a>
### Nested Schema for `aws`

Required:

- `role_arn` (String)


<a id="nestedatt--db_role_to_execute"></a>
### Nested Schema for `db_role_to_execute`

Required:

- `role` (String)
- `type` (String)


<a id="nestedatt--networking"></a>
### Nested Schema for `networking`

Required:

- `access` (Attributes) (see [below for nested schema](#nestedatt--networking--access))

<a id="nestedatt--networking--access"></a>
### Nested Schema for `networking.access`

Required:

- `type` (String)

Optional:

- `connection_id` (String)



<a id="nestedatt--schema_registry_authentication"></a>
### Nested Schema for `schema_registry_authentication`

Optional:

- `password` (String, Sensitive)
- `type` (String)
- `username` (String)


<a id="nestedatt--security"></a>
### Nested Schema for `security`

Optional:

- `broker_public_certificate` (String)
- `protocol` (String)


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), and "h" (hours). Default: `20m`.
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), and "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs. Default: `10m`.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), and "h" (hours). Default: `20m`.

## Import

You can import a stream connection resource using the workspace name, project ID, and connection name. The format must be `WORKSPACE_NAME-PROJECT_ID-CONNECTION_NAME`. For example:

```
$ terraform import mongodbatlas_stream_connection.test "DefaultInstance-12251446ae5f3f6ec7968b13-NewConnection"
```

To learn more, see: [MongoDB Atlas API - Stream Connection](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamConnection) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
