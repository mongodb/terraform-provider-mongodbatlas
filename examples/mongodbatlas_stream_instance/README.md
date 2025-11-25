# MongoDB Atlas Provider - Atlas Stream Instance defined in a Project

> **DEPRECATED:** This example uses the deprecated `mongodbatlas_stream_instance` resource. Please use the [`mongodbatlas_stream_workspace`](../mongodbatlas_stream_workspace/) example instead.

This example shows how to use Atlas Stream Instances in Terraform. It also creates a project, which is a prerequisite.

## Migration to stream_workspace

To migrate to the new `mongodbatlas_stream_workspace` resource, see the [stream_workspace example](../mongodbatlas_stream_workspace/) which demonstrates the updated syntax.

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Unique 24-hexadecimal digit string that identifies the Organization that must contain the project.

To learn more, see the [Stream Instance Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/manage-processing-instance/#configure-a-stream-processing-instance).