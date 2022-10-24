# MongoDB Atlas Provider -- Cluster Tenant Upgrade
This example creates a project and cluster. It is intended to show how to upgrade from shared to dedicated tier.

Variables Required:
- `atlas_org_id`: ID of atlas organization
- `public_key`: Atlas provider public_key
- `private_key`: Atlas provider private key
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `backing_provider_name`: If provider_name is tenant, the backing provider (AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (M0, M2, M5, M10, etc...)

For this example, first we'll start out on the shared tier, then upgrade to a dedicated tier.

As such, utilize the following example `terraform.tfvars` and pseudo-code to execute a working example:

Apply with the following `terraform.tfvars` for shared tier cluster:
```
atlas_org_id                = "627a9687f7f7f7f774de306f14"
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "TENANT"
backing_provider_name       = "AWS"
provider_instance_size_name = "M2"
```

Apply with the following `terraform.tfvars` to upgrade to dedicated tier:
```
atlas_org_id                = "627a9687f7f7f7f774de306f14"
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "GCP"
provider_instance_size_name = "M10"