# MongoDB Atlas Provider — Advanced Cluster (Replica Set)

This example creates a project and a Replica Set cluster. 

### Migrating from v1.x to v2.0.0 or later
If you are migrating from v1.x of our provider to v2.0.0 or later, the `v1.x.x/` sub-directory shows how your current configuration might look like (with added inline comments to demonstrate what has changed in v2.0.0+ for migration reference).

Variables Required:
- `atlas_org_id`: ID of the Atlas organization
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (Free: M0, Dedicated: M10+.)
