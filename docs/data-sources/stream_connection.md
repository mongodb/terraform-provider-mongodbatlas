# Data Source: mongodbatlas_stream_connection

`mongodbatlas_stream_connection` describes a stream connection.

## Example Usage

```terraform
data "mongodbatlas_stream_connection" "example" {
    project_id = "<PROJECT_ID>"
    instance_name = "<INSTANCE_NAME>"
    connection_name = "<CONNECTION_NAME>"
}
```

## Argument Reference

* `project_id` - (Required) Unique 24-hexadecimal digit string that identifies your project.
* `instance_name` - (Required) Human-readable label that identifies the stream instance.
* `connection_name` - (Required) Human-readable label that identifies the stream connection. In the case of the Sample type, this is the name of the sample source.

## Attributes Reference

* `type` - Type of connection. Can be `AWSLambda`, `Cluster`, `Https`, `Kafka` or `Sample`.

If `type` is of value `Cluster` the following additional attributes are defined:
* `cluster_name` - Name of the cluster configured for this connection.
* `db_role_to_execute` - The name of a Built in or Custom DB Role to connect to an Atlas Cluster. See [DBRoleToExecute](#DBRoleToExecute).
* `cluster_project_id` - Unique 24-hexadecimal digit string that identifies the project that contains the configured cluster. Required if the ID does not match the project containing the streams instance. You must first enable the organization setting.

If `type` is of value `Kafka` the following additional attributes are defined:
* `authentication` - User credentials required to connect to a Kafka cluster. Includes the authentication type, as well as the parameters for that authentication mode. See [authentication](#authentication).
* `bootstrap_servers` - Comma separated list of server addresses.
* `config` - A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters.
* `security` - Properties for the secure transport connection to Kafka. For SASL_SSL, this can include the trusted certificate to use. See [security](#security).
* `networking` - Networking Access Type can either be `PUBLIC` (default) or `VPC`. See [networking](#networking).

If `type` is of value `AWSLambda` the following additional attributes are defined:
* `aws` - The configuration for AWS Lambda connection. See [AWS](#AWS)

If `type` is of value `Https` the following additional attributes are defined:
* `url` - URL of the HTTPs endpoint that will be used for creating a connection.
* `headers` - A map of key-value pairs for optional headers.

### Authentication

* `mechanism` - Style of authentication. Can be one of `PLAIN`, `SCRAM-256`, or `SCRAM-512`.
* `username` - Username of the account to connect to the Kafka cluster.
* `password` - Password of the account to connect to the Kafka cluster.

### Security

* `broker_public_certificate` - A trusted, public x509 certificate for connecting to Kafka over SSL. String value of the certificate must be defined in the attribute.
* `protocol` - Describes the transport type. Can be either `SASL_PLAINTEXT` or `SASL_SSL`.

### DBRoleToExecute

* `role` - The name of the role to use. Can be a built in role or a custom role.
* `type` - Type of the DB role. Can be either BUILT_IN or CUSTOM.

### Networking
* `access` - Information about the networking access. See [access](#access).

### Access
* `type` - Selected networking type. Either `PUBLIC`, `VPC` or `PRIVATE_LINK`. Defaults to `PUBLIC`.
* `connection_id` - Id of the Private Link connection when type is `PRIVATE_LINK`.

### AWS
* `role_arn` - Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account. 

To learn more, see: [MongoDB Atlas API - Stream Connection](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Streams/operation/getStreamConnection) Documentation.
The [Terraform Provider Examples Section](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/mongodbatlas_stream_instance/atlas-streams-user-journey.md) also contains details on the overall support for Atlas Streams Processing in Terraform.
