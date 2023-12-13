# MongoDB Atlas Provider - Atlas Streams with Terraform

Atlas Stream Processing is composed of multiple components, and users can leverage Terraform to define a subset of these. To obtain more details on each of the components please refer to the [Atlas Stream Processing Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/overview/#atlas-stream-processing-overview).

### Resources supported by Terraform

- `mongodbatlas_stream_instance`: Enables creating, modifying, and deleteing Stream Instances. as part of this resource, a computed `hostnames` attribute is available for connecting to the created instance.
- `mongodbatlas_stream_connection`: Enables creating, modifying, and deleteing Stream Instance Connections, which serve as data sources and sinks for your instance.

**Note**: To leverage these resources you'll need to set the environment variable `MONGODB_ATLAS_ENABLE_BETA=true` as this functionality is currently in Preview.

### Managing Stream Processors

Once a stream instance and its connections have been defined, `Stream Processors` can be created to define how your data will be processed in your instance. There are currently no resources defined in Terraform to provide this configuration. To obtain information on how this can be configured refer to [Manage Stream Processors](https://www.mongodb.com/docs/atlas/atlas-sp/manage-stream-processor/#manage-stream-processors).

For connecting to your stream instance created in Terraform the following block can be used:
```
output "stream_instance_hostname" {
  value = mongodbatlas_stream_instance.test.hostnames
}
```

This value is then be used for connecting into the stream instance using `mongosh`, as described in the [Get Started Tutorial](https://www.mongodb.com/docs/atlas/atlas-sp/tutorial/). 

