# MongoDB Atlas Service Account Examples

Create a Service Account for an organization and decide how its first secret is created.

`mongodbatlas_service_account` supports two create flows:

- Atlas generates the first secret from the create request (`secret_expires_after_hours`).
- Atlas creates the Service Account without a secret (`without_initial_secret = true`), and you create secrets explicitly with `mongodbatlas_service_account_secret`.

## Sibling examples

- [`bootstrap_secret/`](bootstrap_secret/README.md): create the Service Account with `secret_expires_after_hours` and read the initial secret from the create response.
- [`without_initial_secret/`](without_initial_secret/README.md): create the Service Account with `without_initial_secret = true`, then create the first secret with `mongodbatlas_service_account_secret`.

For product details, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/service_account). For managing and rotating secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).
