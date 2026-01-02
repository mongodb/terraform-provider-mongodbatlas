# MongoDB Atlas Provider -- Service Account Project Assignment

This example shows how to create and assign a Service Account to a Project with specific roles.

## Prerequisites
- Service Account with Organization Owner permissions

## Variables Required to be set:
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Organization ID where the Service Account will be created
- `project_id`: Project ID to assign the Service Account to

## Outputs
- `service_account_project_roles`: The roles assigned to the Service Account in the Project
- `service_account_assigned_projects`: All Projects that the Service Account is assigned to
