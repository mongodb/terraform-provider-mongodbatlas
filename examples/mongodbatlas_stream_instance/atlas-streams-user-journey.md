# MongoDB Atlas Provider - Atlas Streams with Terraform

Atlas Stream Processing is composed of multiple components, and users can leverage Terraform to define a subset of these. To obtain more details on each of the components please refer to the [Atlas Stream Processing Documentation](https://www.mongodb.com/docs/atlas/atlas-sp/overview/#atlas-stream-processing-overview).

### Resources supported by Terraform

- `mongodbatlas_stream_instance`: Enables creating, modifying, and deleting Stream Instances. as part of this resource, a computed `hostnames` attribute is available for connecting to the created instance.
- `mongodbatlas_stream_connection`: Enables creating, modifying, and deleting Stream Instance Connections, which serve as data sources and sinks for your instance.

**NOTE**: To leverage these resources you'll need to set the environment variable `MONGODB_ATLAS_ENABLE_PREVIEW=true` as this functionality is currently in preview.  Also see [Limitations](https://www.mongodb.com/docs/atlas/atlas-sp/limitations/#std-label-atlas-sp-limitations) of Atlas Streams Processing during this preview period.

### Managing Stream Processors

Once a stream instance and its connections have been defined, `Stream Processors` can be created to define how your data will be processed in your instance. There are currently no resources defined in Terraform to provide this configuration. To obtain information on how this can be configured refer to [Manage Stream Processors](https://www.mongodb.com/docs/atlas/atlas-sp/manage-stream-processor/#manage-stream-processors).

Connect to your stream instance defined in terraform using the following code block:

```
output "stream_instance_hostname" {
  value = mongodbatlas_stream_instance.test.hostnames
}
```

This value can then be used to connect to the stream instance using `mongosh`, as described in the [Get Started Tutorial](https://www.mongodb.com/docs/atlas/atlas-sp/tutorial/).
