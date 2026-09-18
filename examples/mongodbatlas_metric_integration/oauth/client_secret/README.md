# OAuth 2.0 client-secret metric integration

Export Atlas metrics to an OTLP-compatible endpoint using OAuth 2.0 client-secret authentication.

For a real integration, point `token_endpoint` at your identity provider's token endpoint and use the `client_id` and `client_secret` registered there. The `client_secret` is write-only and never returned by the API.

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/metric_integration#limitations).

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- An OTLP-compatible endpoint URL.
- OAuth 2.0 client credentials (`client_id` and `client_secret`) registered with your identity provider's token endpoint.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

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
atlas_org_id     = "your-org-id"
otel_endpoint    = "https://otel-collector.example.com:4318/v1/metrics"
token_endpoint   = "https://idp.example.com/oauth2/token"
client_id        = "your-oauth-client-id"
client_secret    = "your-oauth-client-secret"

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
