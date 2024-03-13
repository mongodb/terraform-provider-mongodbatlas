---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: stream connection"
sidebar_current: "docs-mongodbatlas-resource-stream-connection"
description: |-
    Provides a Stream Connection resource.
---

# Resource: mongodbatlas_stream_connection

`mongodbatlas_stream_connection` provides a Stream Connection resource. The resource lets you create, edit, and delete stream instance connections.

-> **NOTE:** Atlas Stream Processing is currently in preview, you'll need to set the environment variable `MONGODB_ATLAS_ENABLE_PREVIEW=true` to use this resource. To learn more about stream processing and existing limitations, see the [Atlas Stream Processing Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/overview/#atlas-stream-processing-overview).

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

### Example Kafka Plaintext Connection

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
        protocol = "PLAINTEXT"
    }
    config = {
        "auto.offset.reset": "latest"
    }
    bootstrap_servers = "localhost:9091,localhost:9092"
}    
```

### Example Kafka SSL Connection

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
        protocol = "SSL"
        broker_public_certificate = "-----BEGIN CERTIFICATE-----<CONTENT>-----END CERTIFICATE-----"
    }
    config = {
        "auto.offset.reset": "latest"
    }
    bootstrap_servers = "localhost:9091,localhost:9092"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `instance_name` - (Required) Human-readable label that identifies the stream instance.
* `connection_name` - (Required) Human-readable label that identifies the stream connection.
* `type` - (Required) Type of connection. Can be either `Cluster` or `Kafka`.

If `type` is of value `Cluster` the following additional arguments are defined:
* `cluster_name` - Name of the cluster configured for this connection.
* `db_role_to_execute` - The name of a Built in or Custom DB Role to connect to an Atlas Cluster. See [DBRoleToExecute](#DBRoleToExecute).

If `type` is of value `Kafka` the following additional arguments are defined:
* `authentication` - User credentials required to connect to a Kafka cluster. Includes the authentication type, as well as the parameters for that authentication mode. See [authentication](#authentication).
* `bootstrap_servers` - Comma separated list of server addresses.
* `config` - A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters.
* `security` - Properties for the secure transport connection to Kafka. For SSL, this can include the trusted certificate to use. See [security](#security).

### Authentication

* `mechanism` - Style of authentication. Can be one of `PLAIN`, `SCRAM-256`, or `SCRAM-512`.
* `username` - Username of the account to connect to the Kafka cluster.
* `password` - Password of the account to connect to the Kafka cluster.

### Security

* `broker_public_certificate` - A trusted, public x509 certificate for connecting to Kafka over SSL. String value of the certificate must be defined in the attribute.
* `protocol` - Describes the transport type. Can be either `PLAINTEXT` or `SSL`.

### DBRoleToExecute

* `role` - The name of the role to use. Can be a built in role or a custom role.
* `type` - Type of the DB role. Can be either BUILT_IN or CUSTOM.

## Import

You can import a stream connection resource using the instance name, project ID, and connection name. The format must be `INSTANCE_NAME-PROJECT_ID-CONNECTION_NAME`. For example:

```
$ terraform import mongodbatlas_stream_connection.test "DefaultInstance-12251446ae5f3f6ec7968b13-NewConnection"
```

To learn more, see: [MongoDB Atlas API - Stream Connection](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/createStreamConnection) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
