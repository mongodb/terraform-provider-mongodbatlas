# MongoDB Atlas Provider - Grant log access to MongoBD employees

This example shows how to use MongoDB Employee Access Grant in Terraform.

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Unique 24-hexadecimal digit string that identifies the project where the stream instance will be created.
- `cluster_name`: Name of Cluster that will be used for creating a connection.

To learn more, see the [MongoDB Employee Access Grant API doc](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-grantmongodbemployeeaccess).
