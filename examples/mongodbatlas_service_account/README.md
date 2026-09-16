# MongoDB Atlas Provider -- Service Account

This example shows how to create a Service Account in MongoDB Atlas.

## Sibling examples

- [`without_initial_secret/`](without_initial_secret/README.md) — recommended flow. Omits `secret_expires_after_hours`, so Atlas creates the Service Account with no secret. Add a secret later with the `mongodbatlas_service_account_secret` resource.
- [`bootstrap_secret/`](bootstrap_secret/README.md) — sets `secret_expires_after_hours`, so Atlas generates a secret on create and the example captures it.

Pick `without_initial_secret/` when something else owns the secret, such as a rotation submodule or a secret store. Pick `bootstrap_secret/` when Atlas should generate the first secret and you read it once from the apply.

For secret rotation, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).
