# MongoDB Atlas Provider - Configure Maintenance Window

This example demonstrates how to configure maintenance windows for your Atlas project in Terraform.

Required variables to set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Unique 24-hexadecimal digit string that identifies the organization that contains the project and cluster.

For additional information you can visit the [Maintenance Window Documentation](https://www.mongodb.com/docs/atlas/tutorial/cluster-maintenance-window/).