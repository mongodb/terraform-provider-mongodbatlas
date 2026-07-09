# MongoDB Atlas Metric Integration Example

This directory contains an example demonstrating how to configure an OpenTelemetry metric integration to export MongoDB Atlas metrics to an OTLP-compatible endpoint such as Datadog, New Relic, or Dynatrace.

## Example

### [Metric Integration](.)

Configures a metric integration at the project level and demonstrates how to use the singular and plural data sources to read back the integration configuration.

**Resources created:**

- MongoDB Atlas Project
- MongoDB Atlas Metric Integration

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- An OTLP-compatible endpoint URL and authentication credentials.

## Variables

| Variable | Description |
| --- | --- |
| `atlas_org_id` | MongoDB Atlas Organization ID. |
| `atlas_project_name` | Name of the Atlas project to create. |
| `datadog_api_key` | Datadog API key used to authenticate the Datadog provider. |
| `datadog_app_key` | Datadog application key used to authenticate the Datadog provider. |
| `datadog_endpoint` | Datadog OTLP-compatible endpoint URL for metric ingestion. |
