# MongoDB Atlas Provider — Stream Connection Failover

This example creates a `mongodbatlas_stream_connection` (the primary connection) and a
`mongodbatlas_stream_connection_failover` for one of the workspace's failover regions.

A failover connection shares the primary connection's `connection_name` and carries its own
regional configuration (for example, a different `bootstrap_servers` for the failover region).
The stream workspace must have failover regions enabled and the failover connection's `region`
must be one of them.

## Dependencies

* Terraform MongoDB Atlas Provider
* A MongoDB Atlas account
* A stream workspace (`mongodbatlas_stream_workspace`) with failover regions enabled
