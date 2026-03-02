# MongoDB Atlas Log Integration Examples

This directory contains examples demonstrating how to configure log integrations to export MongoDB Atlas logs to AWS S3, Microsoft Azure, Google Cloud Platform (GCP), OTel, Datadog and Splunk.

## Available Examples

### [S3 Bucket](./s3bucket/)

A basic example that exports logs to a single S3 bucket. This is the simplest setup and is suitable for most use cases.

**Resources created:**
- S3 bucket
- IAM role and policy
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

### [S3 Multi-Region Access Point (MRAP)](./s3bucket_mrap/)

An advanced example that exports logs to an S3 Multi-Region Access Point (MRAP) instead of a single bucket. This provides high availability and lower latency by automatically routing requests to the closest available S3 bucket.

**Resources created:**
- S3 buckets in multiple AWS regions
- S3 Multi-Region Access Point (MRAP)
- IAM role and policy with MRAP permissions
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

### [Azure Blob](./azure/)

A basic example that exports logs to a single Azure Blob. This is the simplest setup and is suitable for most use cases.

**Resources created:**
- Azure Blob
- IAM role and policy
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

### [GCP Container](./gcp/)

A basic example that exports logs to a single Google Cloud Platform Container. This is the simplest setup and is suitable for most use cases.

**Resources created:**
- Google Cloud Platform Container
- IAM role and policy
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

### [Datadog Datastores](./datadog/)

A basic example that configures MongoDB Atlas log integration to export logs to an existing Datadog Datastore. This is the simplest setup and is suitable for most use cases. This example does not create any Datadog resources; it assumes the Datadog Datastore already exists.

**Resources created:**
- Datadog Datastore
- IAM role and policy
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

### [OTel Collector](./otel/)

A basic example that configures MongoDB Atlas log integration to export logs to an existing OpenTelemetry Collector. This is the simplest setup and is suitable for most use cases. This example does not create any OpenTelemetry resources; it assumes the collector is already deployed.

**Resources created:**
- OpenTelemetry Collector
- IAM role and policy
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

### [Splunk Storage](./splunk/)

 basic example that configures MongoDB Atlas log integration to export logs to an existing Splunk storage destination. This is the simplest setup and is suitable for most use cases. This example does not create any Splunk resources; it assumes the Splunk environment and storage are already configured.

**Resources created:**
- Splunk Storage
- IAM role and policy
- MongoDB Atlas Cloud Provider Access
- MongoDB Atlas Log Integration

## Prerequisites

- MongoDB Atlas account with Organization Owner or Project Owner role.

## Log Types

The `log_types` attribute supports the following values:
- `MONGOD` - MongoDB server logs.
- `MONGOS` - MongoDB router logs.
- `MONGOD_AUDIT` - MongoDB server audit logs.
- `MONGOS_AUDIT` - MongoDB router audit logs.

## Notes

- The requesting Service Account or API Key must have the Organization Owner or Project Owner role.
- MongoDB Cloud will add sub-directories based on the log type under the specified `prefix_path`.
- Optional: Use `kms_key` to specify an AWS KMS key ID or ARN for server-side encryption.
