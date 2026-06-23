# MongoDB Atlas Provider - Configure Org Maintenance Settings

This example demonstrates how to configure organization-level maintenance settings in Terraform, including wave assignment mode for controlling how maintenance waves are assigned across projects.

Required variables to set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Unique 24-hexadecimal digit string that identifies the organization.
