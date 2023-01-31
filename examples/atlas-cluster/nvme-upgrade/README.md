# MongoDB Atlas Provider -- Cluster NVME (Non-Volatile Memory Express) Upgrade
This example creates a project and cluster. It is intended to show how to upgrade from Standard, to PROVISIONED storage tier.

Variables Required:
- `atlas_org_id`: ID of the Atlas organization
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `backing_provider_name`: If provider_name is tenant, the backing provider (AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (Shared: M0, M2, M5, Dedicated: M10+.)
- `provider_volume_type`: Provider storage type STANDARD vs PROVISIONED (NVME)
- `provider_disk_iops`: The maximum input/output operations per second (IOPS) the system can perform. The possible values depend on the selected `provider_instance_size_name` and `disk_size_gb`.  This setting requires that `provider_instance_size_name` to be M30 or greater and cannot be used with clusters with local NVMe SSDs.  The default value for `provider_disk_iops` is the same as the cluster tier's Standard IOPS value, as viewable in the Atlas console.  It is used in cases where a higher number of IOPS is needed and possible.  If a value is submitted that is lower or equal to the default IOPS value for the cluster tier Atlas ignores the requested value and uses the default.  More details available under the providerSettings.diskIOPS parameter: [MongoDB API Clusters](https://docs.atlas.mongodb.com/reference/api/clusters-create-one/)
  * You do not need to configure IOPS for a STANDARD disk configuration but only for a PROVISIONED configuration.

For this example, first we'll start out on the standard tier, then upgrade to a NVME storage tier.


Utilize the following to execute a working example, replacing the org id, public and private key with your values:

Apply with the following `terraform.tfvars` to first create a shared tier cluster:
```
atlas_org_id                = "627a9687f7f7f7f774de306f14"
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "AWS"
provider_instance_size_name = "M40"
provider_volume_type        = "STANDARD"
provider_disk_iops          = 3000
```

Apply with the following `terraform.tfvars` to upgrade the standard storage tier cluster you just created to provisioned storage NVME tier:
```
atlas_org_id                = "627a9687f7f7f7f774de306f14"
public_key                  = <REDACTED>
private_key                 = <REDACTED>
provider_name               = "AWS"
provider_instance_size_name = "M40_NVME"
provider_volume_type        = "PROVISIONED"
provider_disk_iops          = 135125
```