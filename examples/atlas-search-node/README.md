# MongoDB Atlas Provider - Atlas Cluster with dedicated Search Nodes

This example shows how Atlas Dedicated Search Nodes can used in Terraform. As a prerequisite, a project and cluster resource are created.

Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `org_id`: Organization ID where the project and cluster will be created.

For additional information you can visit the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-nodes-for-workload-isolation).