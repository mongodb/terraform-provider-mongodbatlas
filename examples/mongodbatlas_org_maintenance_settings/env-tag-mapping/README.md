# MongoDB Atlas Provider - Org Maintenance Settings: ENV_TAG_MAPPING

This example demonstrates how to configure organization-level maintenance settings using `ENV_TAG_MAPPING` mode, where Atlas automatically assigns project maintenance waves based on the project's `environment` tag.

The tag key must be `environment`. The tag value determines the wave assignment: `development` or `test` map to Wave 1, `staging` maps to Wave 2, and `production` maps to Wave 3.

Required variables to set:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: MongoDB Atlas Organization ID.
- `dev_project_name`: Name of the development MongoDB Atlas project.
- `prod_project_name`: Name of the production MongoDB Atlas project.
