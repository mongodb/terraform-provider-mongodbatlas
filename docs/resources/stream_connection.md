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
    instance_name = "InstanceName"
    connection_name = "ConnectionName"
    type = "Cluster"
    cluster_name = "Cluster0"
}
```

### Further Examples
- [Atlas Stream Connection](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.0.1/examples/mongodbatlas_stream_connection)

### Example Cross Project Cluster Connection

```terraform
resource "mongodbatlas_stream_connection" "test" {
    project_id         = var.project_id
    instance_name      = "InstanceName"
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
    instance_name = "NewInstance"
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
    instance_name = "NewInstance"
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
    instance_name   = "NewInstance"
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
  instance_name   = mongodbatlas_stream_instance.example.instance_name
  connection_name = "https_connection_tf_new"
  type            = "Https"
  url             = "https://example.com"
  headers = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `instance_name` - (Required) Human-readable label that identifies the stream instance.
* `connection_name` - (Required) Human-readable label that identifies the stream connection. In the case of the Sample type, this is the name of the sample source.
* `type` - (Required) Type of connection. Can be `AWSLambda`, `Cluster`, `Https`, `Kafka` or `Sample`.

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

## Import

You can import a stream connection resource using the instance name, project ID, and connection name. The format must be `INSTANCE_NAME-PROJECT_ID-CONNECTION_NAME`. For example:

```
$ terraform import mongodbatlas_stream_connection.test "DefaultInstance-12251446ae5f3f6ec7968b13-NewConnection"
```

To learn more, see: [MongoDB Atlas API - Stream Connection](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamConnection) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
