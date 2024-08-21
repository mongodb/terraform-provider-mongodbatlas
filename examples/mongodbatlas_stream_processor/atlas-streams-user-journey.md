# MongoDB Atlas Provider - Atlas Streams with Terraform

Atlas Stream Processing is composed of multiple components, and users can leverage Terraform to define a subset of these. To obtain more details on each of the components please refer to the [Atlas Stream Processing Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/overview/#atlas-stream-processing-overview).

### Resources supported by Terraform

- `mongodbatlas_stream_instance`: Enables creating, modifying, and deleting Stream Instances. as part of this resource, a computed `hostnames` attribute is available for connecting to the created instance.
- `mongodbatlas_stream_connection`: Enables creating, modifying, and deleting Stream Instance Connections, which serve as data sources and sinks for your instance.
- `mongodbatlas_stream_processor`: Enables creating, deleting, starting and stopping a Stream Processor, which define how your data will be processed in your instance.
