---
subcategory: "Private Endpoint Services"
---

# Resource: mongodbatlas_private_endpoint_connection_string

`mongodbatlas_private_endpoint_connection_string` provides a resource for managing a customer-managed private endpoint connection string for one MongoDB Atlas cluster. The connection string is built from an explicit set of private endpoints, and MongoDB Atlas embeds a stable generated slug in its hostnames.

~> **NOTE:** The Atlas API for this resource is in preview and can change before it becomes generally available. Your project must be opted in to private endpoint connection strings through the project setting `isPrivateEndpointConnectionStringsEnabled` (see [Update Project Settings](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-updategroupsettings)). When a project opts in, Atlas copies each existing generated private endpoint connection string into a read-only `IMPORTED` record with the same hostname. This resource creates and manages `SELF_MANAGED` connection strings only; the data sources return both types, and Atlas rejects updates to `IMPORTED` records. To use this resource, the requesting Service Account or API Key must have the Project Owner role.

Create and update operations are asynchronous. Terraform waits until `status` is `AVAILABLE`, up to the `timeouts` values (default 2 hours). If Atlas reports `FAILED`, the apply fails and `error_message` explains the cause and the fix.

`private_endpoint_ids` accepts 1 to 50 private endpoints, at most one per cloud provider region. It can include private endpoints that the cluster does not use yet, which lets you pre-link a private endpoint for a future region. You can update `name` and `private_endpoint_ids` in place: the connection string ID and slug stay the same, and Atlas keeps returning the existing connection string until the updated one is available. `project_id` and `cluster_name` cannot be changed after creation.

`optimized_mode_enabled` requires an AWS sharded cluster and a private endpoint selection that meets the [optimized connection string requirements](https://www.mongodb.com/docs/atlas/reference/faq/connection-changes/#how-can-i-optimize-connection-performance-for-sharded-clusters-using-private-endpoints).

## Example Usages

```terraform
# Creates an AWS private endpoint for an Atlas project, links it to Atlas, and creates a customer-managed private endpoint connection string for an existing cluster.

resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = var.atlas_region_name
}

resource "aws_vpc_endpoint" "this" {
  vpc_id             = var.vpc_id
  service_name       = mongodbatlas_privatelink_endpoint.this.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = var.subnet_ids
  security_group_ids = var.security_group_ids
}

resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.this.private_link_id
  endpoint_service_id = aws_vpc_endpoint.this.id
  provider_name       = "AWS"
}

resource "mongodbatlas_private_endpoint_connection_string" "this" {
  project_id             = var.project_id
  cluster_name           = var.cluster_name
  name                   = var.connection_string_name
  optimized_mode_enabled = var.optimized_mode_enabled
  # Reference the Atlas registration, not aws_vpc_endpoint, so the string is created after Atlas accepts the endpoint and deleted before it.
  private_endpoint_ids = [mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id]
}

data "mongodbatlas_private_endpoint_connection_string" "this" {
  project_id           = mongodbatlas_private_endpoint_connection_string.this.project_id
  cluster_name         = mongodbatlas_private_endpoint_connection_string.this.cluster_name
  connection_string_id = mongodbatlas_private_endpoint_connection_string.this.connection_string_id
}

data "mongodbatlas_private_endpoint_connection_strings" "self_managed" {
  project_id      = mongodbatlas_private_endpoint_connection_string.this.project_id
  cluster_name    = mongodbatlas_private_endpoint_connection_string.this.cluster_name
  management_type = "SELF_MANAGED"
}

output "srv_connection_string" {
  description = "Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to the cluster through the selected private endpoints."
  value       = mongodbatlas_private_endpoint_connection_string.this.srv_connection_string
}

output "connection_string" {
  description = "Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to the cluster through the selected private endpoints."
  value       = mongodbatlas_private_endpoint_connection_string.this.connection_string
}

output "connection_string_status" {
  description = "Status of the connection string, read with the singular data source."
  value       = data.mongodbatlas_private_endpoint_connection_string.this.status
}

output "self_managed_connection_string_ids" {
  description = "IDs of all self-managed private endpoint connection strings for the cluster, read with the plural data source."
  value       = [for cs in data.mongodbatlas_private_endpoint_connection_strings.self_managed.results : cs.connection_string_id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_name` (String) Human-readable label that identifies the cluster.
- `name` (String) Human-readable label of at most 64 characters that identifies this private endpoint connection string. You can change this label without changing the connection string.
- `private_endpoint_ids` (Set of String) Set of private endpoints that this connection string covers. MongoDB Cloud stores this set as a snapshot of your selection and never adds newly created private endpoints to it. The set accepts at most one private endpoint per cloud provider region, and can include private endpoints that the cluster does not use yet. Order is not significant and MongoDB Cloud rejects duplicate values.
- `project_id` (String) Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.

### Optional

- `delete_on_create_timeout` (Boolean) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.
- `optimized_mode_enabled` (Boolean) Flag that indicates whether this private endpoint connection string uses the optimized load-balanced mode. This mode requires a sharded cluster and a private endpoint selection that satisfies the optimized connection string requirements. This value cannot be changed after creation. To use a different mode, create a new connection string, move your applications to it, and then delete the old one.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `connection_string` (String) Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to this cluster through the selected private endpoints. MongoDB Cloud returns this value when `status` is `AVAILABLE`. After you submit an update, MongoDB Cloud continues to return the existing connection string until the updated connection string becomes available.
- `connection_string_id` (String) Unique 24-hexadecimal digit string that identifies this private endpoint connection string. This value is the `connectionStringId` path parameter of this resource's endpoints.
- `created_at` (String) Date and time when MongoDB Cloud created this private endpoint connection string. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
- `error_message` (String) Error message that explains why this private endpoint connection string is in the `FAILED` state, including the remediation, for example when a selected endpoint was deleted mid-flight. Present only when `status` is `FAILED`.
- `management_type` (String) Party that manages this private endpoint connection string. MongoDB Cloud returns `SELF_MANAGED` for connection strings that you create and maintain through this resource, and `IMPORTED` for records that preserve the private endpoint connection strings that this cluster exposed before migration. You can delete an imported record but not update it.
- `slug` (String) Short generated pin of at most 64 characters, unique per project. Present only when `managementType` is `SELF_MANAGED`. A self-managed connection string embeds this value in its DNS names, and its replica set horizon name is `cs-` plus this value.
- `srv_connection_string` (String) Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to this cluster through the selected private endpoints. MongoDB Cloud returns this value when `status` is `AVAILABLE`. After you submit an update, MongoDB Cloud continues to return the existing connection string until the updated connection string becomes available.
- `status` (String) State of this private endpoint connection string. MongoDB Cloud returns `PENDING` while it reconfigures the cluster, `AVAILABLE` once the connection string is usable, and `FAILED` when creation cannot complete.
- `updated_at` (String) Date and time when MongoDB Cloud last updated this private endpoint connection string. This parameter expresses its value in the ISO 8601 timestamp format in UTC.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import a private endpoint connection string by using the project ID, cluster name, and connection string ID in the format `PROJECT_ID/CLUSTER_NAME/CONNECTION_STRING_ID`, for example:

```shell
terraform import mongodbatlas_private_endpoint_connection_string.this 650972848269185c55f40ca1/my-cluster/32b6e34b3d91647abb20e7b8
```

For more information see: [MongoDB Atlas API - Private Endpoint Connection String](https://mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Clusters/operation/createGroupClusterPrivateEndpointConnectionString) Documentation.
