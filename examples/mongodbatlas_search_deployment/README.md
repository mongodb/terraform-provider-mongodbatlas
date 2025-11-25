# MongoDB Atlas Provider - Atlas Cluster with dedicated Search Nodes Deployment

This example shows how you can use Atlas Dedicated Search Nodes in Terraform. As part of it, a project and cluster resource are created as a prerequisite.

Variables Required to be set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Organization ID where the project and cluster will be created.

For additional information you can visit the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-nodes-for-workload-isolation).