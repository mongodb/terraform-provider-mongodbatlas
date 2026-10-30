# OAuth 2.0 client-secret metric integration

Export Atlas metrics to an OTLP-compatible endpoint using OAuth 2.0 client-secret authentication.

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/metric_integration#limitations).

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- An OTLP-compatible endpoint URL.
- OAuth 2.0 client credentials (`client_id` and `client_secret`) registered with your identity provider's token endpoint.

## Defaults

Omit optional variables to use these values from `variables.tf`:

- `metric_selection`: `["ATLAS_STREAM_PROCESSING"]`
- `oauth_scopes`: `[]` — no scopes requested
- `atlas_project_name`: `tf-metric-integration-oauth-client-secret`

## Usage

**1. Set credentials.**

```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

**2. Create `terraform.tfvars`.**

Required inputs:

```hcl
atlas_org_id        = "your-org-id"
atlas_client_id     = "your-service-account-client-id"
atlas_client_secret = "your-service-account-client-secret"
otel_endpoint       = "https://otel-collector.example.com:4318/v1/metrics"
token_endpoint      = "https://idp.example.com/oauth2/token"
client_id           = "your-oauth-client-id"
client_secret       = "your-oauth-client-secret"

# metric_selection = ["ATLAS_STREAM_PROCESSING"]   # default
# oauth_scopes     = ["metrics.write"]              # optional
```

**3. Plan and apply.**

```bash
terraform plan
terraform apply
```

**4. Destroy.**

```bash
terraform destroy
```
