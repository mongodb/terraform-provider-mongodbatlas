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
- [Atlas Stream Connection](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.5.0/examples/mongodbatlas_stream_connection)

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
    instance_name   = mongodbatlas_stream_instance.example.instance_name
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
  workspace_name   = mongodbatlas_stream_instance.example.instance_name
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
  workspace_name           = mongodbatlas_stream_instance.example.instance_name
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
  workspace_name           = mongodbatlas_stream_instance.example.instance_name
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
resource "mongodbatlas_stream_instance" "example" {
  project_id    = var.project_id
  instance_name = "ExampleInstance"
  data_process_region = {
    region         = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
}

# Source connection (Sample data)
resource "mongodbatlas_stream_connection" "source" {
  project_id      = var.project_id
  workspace_name  = mongodbatlas_stream_instance.example.instance_name
  connection_name = "sample_stream_solar"
  type            = "Sample"
}

# Sink connection (Atlas Cluster)
resource "mongodbatlas_stream_connection" "sink" {
  project_id      = var.project_id
  workspace_name  = mongodbatlas_stream_instance.example.instance_name
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
  instance_name  = mongodbatlas_stream_instance.example.instance_name
  processor_name = "ExampleProcessor"
  pipeline = jsonencode([
    { "$source" = { "connectionName" = mongodbatlas_stream_connection.source.connection_name } },
    { "$emit" = { "connectionName" = mongodbatlas_stream_connection.sink.connection_name } }
  ])
  state = "STARTED"
}
```

~> **NOTE:** The stream processor resource automatically depends on the stream connections through the `connection_name` references in the pipeline. This ensures proper creation order. The provider waits for each connection to be fully provisioned before returning from create or update operations.

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `instance_name` - (Deprecated) Label that identifies the stream processing workspace. Attribute is deprecated and will be removed in following major versions in favor of `workspace_name`.
* `workspace_name` - (Optional) Label that identifies the stream processing workspace. Conflicts with `instance_name`.
* `connection_name` - (Required) Label that identifies the stream connection. In the case of the Sample type, this is the name of the sample source.
* `type` - (Required) Type of connection. Can be `AWSLambda`, `Cluster`, `Https`, `Kafka`, `Sample`, or `SchemaRegistry`.

~> **NOTE:** Either `workspace_name` or `instance_name` must be provided, but not both. These fields are functionally identical and `workspace_name` is an alias for `instance_name`. `workspace_name` should be used instead of `instance_name`.

If `type` is of value `Cluster` the following additional arguments are defined:
* `cluster_name` - Name of the cluster configured for this connection.
* `db_role_to_execute` - The name of a Built in or Custom DB Role to connect to an Atlas Cluster. See [DBRoleToExecute](#DBRoleToExecute).
* `cluster_project_id` - Unique 24-hexadecimal digit string that identifies the project that contains the configured cluster. Required if the ID does not match the project containing the streams instance. You must first enable the organization setting.

If `type` is of value `Kafka` the following additional arguments are defined:
* `authentication` - User credentials required to connect to a Kafka cluster. Includes the authentication type, as well as the parameters for that authentication mode. See [authentication](#authentication).
* `bootstrap_servers` - Comma separated list of server addresses.
* `config` - A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters.
* `security` - Properties for the secure transport connection to Kafka. For SASL_SSL, this can include the trusted certificate to use. See [security](#security).
* `networking` - Networking Access Type can either be `PUBLIC` (default) or `VPC`. See [networking](#networking).

If `type` is of value `AWSLambda` the following additional arguments are defined:
* `aws` - The configuration for AWS Lambda connection. See [AWS](#AWS)

If `type` is of value `Https` the following additional attributes are defined:
* `url` - URL of the HTTPs endpoint that will be used for creating a connection.
* `headers` - A map of key-value pairs for optional headers.

If `type` is of value `SchemaRegistry` the following additional arguments are defined:
* `schema_registry_provider` - The Schema Registry provider. Must be set to `CONFLUENT`.
* `schema_registry_urls` - List of Schema Registry endpoint URLs used by this connection. Each URL must use the http or https scheme and specify a valid host and optional port.
* `schema_registry_authentication` - Authentication configuration for Schema Registry. See [Schema Registry Authentication](#schema-registry-authentication).

### Authentication

* `mechanism` - Method of authentication. Value can be `PLAIN`, `SCRAM-256`, or `SCRAM-512`.
* `method` - SASL OAUTHBEARER authentication method. Value must be OIDC.
* `username` - Username of the account to connect to the Kafka cluster.
* `password` - Password of the account to connect to the Kafka cluster.
* `token_endpoint_url` -  OAUTH issuer (IdP provider) token endpoint HTTP(S) URI used to retrieve the token.
* `client_id` - Public identifier for the Kafka client.
* `client_secret` - Secret known only to the Kafka client and the authorization server.
* `scope` - Scope of the access request to the broker specified by the Kafka clients.
* `sasl_oauthbearer_extensions` - Additional information to provide to the Kafka broker.

### Security

* `broker_public_certificate` - A trusted, public x509 certificate for connecting to Kafka over SSL. String value of the certificate must be defined in the attribute.
* `protocol` - Describes the transport type. Can be either `SASL_PLAINTEXT` or `SASL_SSL`.

### DBRoleToExecute

* `role` - The name of the role to use. Value can be  `atlasAdmin`, `readWriteAnyDatabase`, or `readAnyDatabase` if `type` is set to `BUILT_IN`, or the name of a user-defined role if `type` is set to `CUSTOM`.
* `type` - Type of the DB role. Can be either BUILT_IN or CUSTOM.

### Networking
* `access` - Information about the networking access. See [access](#access).

### Access
* `type` - Selected networking type. Either `PUBLIC`, `VPC` or `PRIVATE_LINK`. Defaults to `PUBLIC`.
* `connection_id` - Id of the Private Link connection when type is `PRIVATE_LINK`.

### AWS
* `role_arn` - Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account.

### Schema Registry Authentication
* `type` - Authentication type discriminator. Specifies the authentication mechanism for Confluent Schema Registry. Valid values are `USER_INFO` or `SASL_INHERIT`.
  * `USER_INFO` - Uses username and password authentication for Confluent Schema Registry.
  * `SASL_INHERIT` - Inherits the authentication configuration from Kafka for the Confluent Schema Registry.
* `username` - Username for the Schema Registry. Required when `type` is `USER_INFO`.
* `password` - Password for the Schema Registry. Required when `type` is `USER_INFO`.

## Import

You can import a stream connection resource using the workspace name, project ID, and connection name. The format must be `WORKSPACE_NAME-PROJECT_ID-CONNECTION_NAME`. For example:

```
$ terraform import mongodbatlas_stream_connection.test "DefaultInstance-12251446ae5f3f6ec7968b13-NewConnection"
```

To learn more, see: [MongoDB Atlas API - Stream Connection](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamConnection) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
