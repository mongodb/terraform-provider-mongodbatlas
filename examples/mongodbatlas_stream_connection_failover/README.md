# MongoDB Atlas Provider — Stream Connection Failover

This example creates a stream workspace with a failover region, a primary `mongodbatlas_stream_connection`,
and a `mongodbatlas_stream_connection_failover` for the failover region. The failover connection shares the
primary connection's name and carries its own regional configuration.

## Dependencies

* Terraform MongoDB Atlas Provider
* A MongoDB Atlas account
