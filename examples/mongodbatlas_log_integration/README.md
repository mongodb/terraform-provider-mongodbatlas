# MongoDB Atlas Log Integration Examples

This directory contains examples demonstrating how to configure log integrations to export MongoDB Atlas logs to various destinations.

## Available Examples

### [S3 Bucket](./s3bucket/)

A basic example that exports logs to a single S3 bucket. This is the simplest S3 setup and is suitable for most use cases.

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

### [GCS Bucket](./gcp/)

Exports logs to a Google Cloud Storage (GCS) bucket.

**Resources created:**
- GCS bucket
- MongoDB Atlas Log Integration

### [Azure Blob Storage](./azure/)

Exports logs to an Azure Blob Storage container.

**Resources created:**
- Azure resource group, storage account, and storage container
- MongoDB Atlas Log Integration

### [Datadog](./datadog/)

Exports logs to Datadog.

**Resources created:**
- MongoDB Atlas Log Integration

### [Splunk](./splunk/)

Exports logs to a Splunk HTTP Event Collector (HEC) endpoint.

**Resources created:**
- MongoDB Atlas Log Integration

### [OpenTelemetry (OTel)](./otel/)

Exports logs to an OpenTelemetry collector.

**Resources created:**
- MongoDB Atlas Log Integration

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- Cloud provider account and credentials appropriate for the chosen example.

## Log Types

The `log_types` attribute supports the following values:
- `MONGOD` - MongoDB server logs.
- `MONGOS` - MongoDB router logs.
- `MONGOD_AUDIT` - MongoDB server audit logs.
- `MONGOS_AUDIT` - MongoDB router audit logs.
