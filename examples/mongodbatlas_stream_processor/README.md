# MongoDB Atlas Provider - Atlas Stream Processor defined in a Project

This example shows how to use Atlas Stream Processors in Terraform. It also creates a project, which is a prerequisite.

You must set the following variables:

- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_id`: Unique 24-hexadecimal digit string that identifies the project where the stream instance will be created.
- `kafka_username`: Username used for connecting to your external Kafka Cluster. 
- `kafka_password`: Password used for connecting to your external Kafka Cluster.
- `cluster_name`: Name of Cluster that will be used for creating a connection.

To learn more, see the [Stream Processor Documentation](https://www.mongodb.com/docs/atlas/atlas-stream-processing/manage-stream-processor/).