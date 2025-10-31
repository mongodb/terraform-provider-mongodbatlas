# MongoDB Atlas Provider - Atlas Stream Workspace defined in a Project

This example shows how to use Atlas Stream Workspaces in Terraform. It also creates a project, which is a prerequisite.

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Unique 24-hexadecimal digit string that identifies the Organization that must contain the project.

To learn more, see the [Stream Workspace Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/manage-processing-instance/#configure-a-stream-processing-instance).

## Migration from stream_instance

This example replaces the deprecated `mongodbatlas_stream_instance` resource. The main changes are:
- Resource name: `mongodbatlas_stream_instance` → `mongodbatlas_stream_workspace`
- Field name: `instance_name` → `workspace_name`
