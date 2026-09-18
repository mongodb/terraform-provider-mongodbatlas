# MongoDB Atlas Metric Integration Examples

Configure a metric integration to export Atlas metrics to an OTLP-compatible endpoint, using header-based or OAuth 2.0 authentication.

## Sibling examples

Header-based:

- [`header/`](header/README.md) — export to Datadog using header-based authentication (creates the Datadog API key with the `datadog` provider).

OAuth 2.0:

- [`oauth/client_secret/`](oauth/client_secret/README.md) — export using a shared client secret.
- [`oauth/private_key_jwt/`](oauth/private_key_jwt/README.md) — export using an Atlas-managed private-key JWT signing assertion; register the returned `jwks_uri` with your identity provider.

Each example creates a MongoDB Atlas project and the metric integration resource, and reads it back with the singular and plural data sources.

For product limits and supported providers, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/metric_integration) and [MongoDB Atlas - OTel Integration](https://www.mongodb.com/docs/atlas/tutorial/otel-integration/) Documentation.
