# Data Source: private_endpoint_regional_mode

`private_endpoint_regional_mode` describes a Private Endpoint Regional Mode. This represents a Private Endpoint Regional Mode Connection that wants to retrieve settings of an Atlas project.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```terraform
resource "private_endpoint_regional_mode" "test" {
  project_id    = "<PROJECT-ID>"
}

data "private_endpoint_regional_mode" "test" {
	project_id = private_endpoint_regional_mode.test.project_id
}
```

## Argument Reference
* `project_id` - (Required) Unique identifier for the project.
* `enabled` - (Optional) Flag that indicates whether the regionalized private endpoitn setting is enabled for the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.

See detailed information for arguments and attributes: **Private Endpoints** [Get Regional Mode](https://www.mongodb.com/docs/atlas/reference/api/private-endpoints-get-regional-mode/) | [Update Regional Mode](https://www.mongodb.com/docs/atlas/reference/api/private-endpoints-update-regional-mode/)