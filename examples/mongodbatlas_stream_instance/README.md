# MongoDB Atlas Provider - Atlas Stream Instance defined in a Project

This example shows how to use Atlas Stream Instances in Terraform. It also creates a project, which is a prerequisite.

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `org_id`: Unique 24-hexadecimal digit string that identifies the Organization that must contain the project.

To learn more, see the [Stream Instance Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/manage-processing-instance/#configure-a-stream-processing-instance).