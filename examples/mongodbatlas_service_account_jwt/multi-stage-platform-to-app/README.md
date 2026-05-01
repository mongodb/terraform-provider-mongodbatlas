# MongoDB Atlas Provider -- Multi-Stage Project Vending with Ephemeral JWT

This example demonstrates a project vending pattern where a platform team creates an Atlas project, generates a short-lived JWT, and automatically hands it off to an application team. The app team uses the JWT to bootstrap their own project-scoped Service Account, then operates independently with long-lived SA credentials. This is useful for platform teams that vend Atlas projects to application teams without sharing long-lived org-level credentials.

## Flow

```
Platform (Phase 0)                    App Team
Org-level SA credentials              (nothing)
Create project + generate JWT ──────> App Bootstrap (Phase 1)
  (3600s lifetime)                      Uses JWT to create project-scoped SA
                                        Stores SA creds in Secrets Manager
JWT expires, nothing                    │
  to rotate or revoke                   v
                                      App Ongoing (Phase 2)
                                        Reads SA creds from Secrets Manager
                                        Creates dedicated cluster
                                        No time limit
```

The hand-off is automated using `local-exec` provisioners that trigger downstream `terraform apply` runs, passing values as `TF_VAR_` environment variables. This demonstrates how CI/CD pipelines inject masked runtime variables into downstream stages.

## Structure

| Directory | Description |
|---|---|
| `phase-0-platform/` | Creates the Atlas project, generates a JWT, and triggers `phase-1-app-bootstrap`. |
| `phase-1-app-bootstrap/` | Uses the JWT to create a project-scoped SA, stores its credentials in AWS Secrets Manager, and triggers `phase-2-app-ongoing`. |
| `phase-2-app-ongoing/` | Reads SA credentials from Secrets Manager, configures the Atlas provider, and creates a dedicated cluster. |

## Prerequisites

- Terraform >= 1.11 (required for write-only attributes in `phase-1-app-bootstrap`).
- An org-level MongoDB Atlas Service Account with permissions to create projects.
- AWS CLI configured with `secretsmanager:CreateSecret`, `secretsmanager:PutSecretValue`, and `secretsmanager:GetSecretValue` permissions.

## Usage

Set the required variables in `phase-0-platform/terraform.tfvars`:

- `atlas_client_id`: Org-level Service Account Client ID.
- `atlas_client_secret`: Org-level Service Account Client Secret.
- `org_id`: Atlas Organization ID.

Then run a single apply from the `phase-0-platform/` directory:

```bash
cd phase-0-platform
terraform init
terraform apply
```

This triggers the entire chain automatically:
1. Creates the Atlas project and generates a JWT.
2. Triggers `phase-1-app-bootstrap`, which creates a project-scoped SA and stores its credentials.
3. Triggers `phase-2-app-ongoing`, which reads the SA credentials and creates a dedicated cluster.

## Mapping to CI/CD Pipelines

In a real deployment, each phase can map to a separate CI/CD pipeline stage:

- **Platform pipeline** creates the project and JWT, then triggers the app pipeline with the JWT as a masked runtime variable (e.g., GitHub Actions `workflow_dispatch` secret, GitLab CI trigger variable).
- **App bootstrap pipeline** receives the JWT, creates the project SA, and stores credentials in the app team's own secret store.
- **App ongoing pipeline** reads SA credentials from the secret store and runs independently on any schedule.

The `local-exec` provisioners in this example simulate that trigger injection pattern.

