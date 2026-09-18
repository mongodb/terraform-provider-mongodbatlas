# OAuth 2.0 private-key JWT metric integration

Export Atlas metrics to an OTLP-compatible endpoint using OAuth 2.0 private-key JWT authentication.

Atlas generates and manages the signing key server-side. After apply, register the `jwks_uri` from the `signing_key_info` output with your identity provider so it can verify Atlas-signed client assertions. No client secret is sent.

For product limits, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/metric_integration#limitations).

## Prerequisites

- MongoDB Atlas Service Account with Organization Owner or Project Owner role.
- An OTLP-compatible endpoint URL.
- An identity provider that accepts private-key JWT (RFC 7523) client assertions, with a `client_id` registered for Atlas.

Credentials come from the environment (`MONGODB_ATLAS_CLIENT_ID` and `MONGODB_ATLAS_CLIENT_SECRET`).

## Defaults

Omit optional variables to use these values from `variables.tf`:

- `metric_selection`: `["ATLAS_STREAM_PROCESSING"]`
- `token_request_params`: `{}` — no extra token request parameters
- `atlas_project_name`: `tf-metric-integration-oauth-private-key-jwt`

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

# metric_selection = ["ATLAS_STREAM_PROCESSING"]   # default
# token_request_params = { resource = "your-resource-uri" }   # optional
```

**3. Plan and apply.**

```bash
terraform plan
terraform apply
```

After apply, read `signing_key_info` and register the returned `jwks_uri` with your identity provider.

**4. Destroy.**

```bash
terraform destroy
```
