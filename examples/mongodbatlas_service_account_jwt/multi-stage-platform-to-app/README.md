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
                                        Creates flex cluster
                                        No time limit
```

The hand-off is automated using `local-exec` provisioners that trigger downstream `terraform apply` runs, passing values as `TF_VAR_` environment variables. This mirrors how CI/CD pipelines inject masked runtime variables into downstream stages.

## Structure

| Directory | Phase | Description |
|---|---|---|
| `step-1-platform/` | 0 | Creates the Atlas project, generates a JWT, and triggers `step-2-app-bootstrap`. |
| `step-2-app-bootstrap/` | 1 | Uses the JWT to create a project-scoped SA, stores its credentials in AWS Secrets Manager, and triggers `step-3-app-ongoing`. |
| `step-3-app-ongoing/` | 2 | Reads SA credentials from Secrets Manager, configures the Atlas provider, and creates a flex cluster. |

## Prerequisites

- Terraform >= 1.11 (required for write-only attributes in `step-2-app-bootstrap`).
- An org-level MongoDB Atlas Service Account with permissions to create projects.
- AWS CLI configured with `secretsmanager:CreateSecret`, `secretsmanager:PutSecretValue`, and `secretsmanager:GetSecretValue` permissions.

## Usage

Set the required variables in `step-1-platform/terraform.tfvars`:

- `atlas_client_id`: Org-level Service Account Client ID.
- `atlas_client_secret`: Org-level Service Account Client Secret.
- `org_id`: Atlas Organization ID.

Then run a single apply from the `step-1-platform/` directory:

```bash
cd step-1-platform
terraform init
terraform apply
```

This triggers the entire chain automatically:
1. Creates the Atlas project and generates a JWT.
2. Triggers `step-2-app-bootstrap`, which creates a project-scoped SA and stores its credentials.
3. Triggers `step-3-app-ongoing`, which reads the SA credentials and creates a flex cluster.

## Mapping to CI/CD Pipelines

In a real deployment, each phase maps to a separate CI/CD pipeline stage:

- **Platform pipeline** creates the project and JWT, then triggers the app pipeline with the JWT as a masked runtime variable (e.g., GitHub Actions `workflow_dispatch` secret, GitLab CI trigger variable).
- **App bootstrap pipeline** receives the JWT, creates the project SA, and stores credentials in the app team's own secret store.
- **App ongoing pipeline** reads SA credentials from the secret store and runs independently on any schedule.

The `local-exec` provisioners in this example simulate that trigger injection pattern.

## Cleanup

Destroy resources in reverse order:

```bash
cd step-3-app-ongoing
terraform destroy

cd ../step-2-app-bootstrap
terraform destroy

cd ../step-1-platform
terraform destroy
```

Note: the flex cluster in `step-3-app-ongoing` may take several minutes to destroy.
