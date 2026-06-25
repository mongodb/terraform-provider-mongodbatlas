# MongoDB Atlas Provider - Configure Org Maintenance Settings

This example demonstrates how to configure organization-level maintenance settings in Terraform using `ENV_TAG_MAPPING` mode, where Atlas automatically assigns project maintenance waves based on the project's `Environment` tag.

Required variables to set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Unique 24-hexadecimal digit string that identifies the organization.
- `dev_project_name`: Name of the development MongoDB Atlas project.
- `prod_project_name`: Name of the production MongoDB Atlas project.
