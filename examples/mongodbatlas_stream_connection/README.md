# MongoDB Atlas Provider - Atlas Stream Instance defined in a Project

This example shows how to create Atlas Stream Connections in Terraform. It also creates a stream instance, which is a prerequisite. The Kafka, HTTPs and Cluster connections types are defined to showcase their usage.

You must set the following variables:

- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `project_id`: Unique 24-hexadecimal digit string that identifies the project where the stream instance will be created.
- `kafka_username`: Username used for connecting to your external Kafka Cluster. 
- `kafka_password`: Password used for connecting to your external Kafka Cluster.
- `kafka_ssl_cert`: String value of public x509 certificate for connecting to Kafka over SSL.
- `cluster_name`: Name of Cluster that will be used for creating a connection.

To learn more, see the [Stream Instance Connection Registry Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/manage-processing-instance/#view-connections-in-the-connection-registry).