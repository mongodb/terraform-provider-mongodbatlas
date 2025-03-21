# MongoDB Atlas Provider - Atlas Cluster with dedicated Search Nodes Deployment

This example shows how you can configure maintenance windows for your Atlas project in Terraform.

Variables required to be set:

- `public_key`: Atlas public key
- `private_key`: Atlas private key
- `org_id`: Organization ID where the project and cluster will be created.

For additional information you can visit the [Maintenance Window Documentation](https://www.mongodb.com/docs/atlas/tutorial/cluster-maintenance-window/).