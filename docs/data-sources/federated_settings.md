# Data Source: mongodbatlas_federated_settings

`mongodbatlas_federated_settings` provides a federated settings data source. Atlas Cloud federated settings provides federated settings outputs.


## Example Usage

```terraform
data "mongodbatlas_federated_settings" "settings" {
  org_id = "627a9683e7f7f7ff7fe306f14"
}
```

## Argument Reference
* `org_id` - Unique 24-hexadecimal digit string that identifies the organization that contains your projects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


### FederatedSettings
          
* `federated_domains` - List that contains the domains associated with the organization's identity provider.
* `has_role_mappings` - Flag that indicates whether this organization has role mappings configured.
* `id` - Unique 24-hexadecimal digit string that identifies this federation.
* `identity_provider_id` - Unique 20-hexadecimal digit string that identifies the identity provider connected to this organization.
* `identity_provider_status` - Value that indicates whether the identity provider is active. Atlas returns ACTIVE if the identity provider is active and INACTIVE if the identity provider is inactive.


For more information see: [MongoDB Atlas API Reference.](https://www.mongodb.com/docs/atlas/reference/api/federation-configuration/)
