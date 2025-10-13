# MongoDB Atlas Provider — Advanced Cluster (Replica Set)

This example creates a project and a Replica Set cluster. 

_NOTE: If you are migrating from v1.x of our provider, the `v1.x.x/` sub-directory shows how your current configuration might look like (with added inline comments to demonstrate what has changed in v2.0.0+ for migration reference)._ 

Variables Required:
- `atlas_org_id`: ID of the Atlas organization
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `provider_name`: Name of provider to use for cluster (TENANT, AWS, GCP)
- `provider_instance_size_name`: Size of the cluster (Free: M0, Dedicated: M10+.)
