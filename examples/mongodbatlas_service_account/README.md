# MongoDB Atlas Service Account Examples

Create a Service Account for an organization and decide how its first secret is created.

The recommended flow creates the Service Account without a secret and creates secrets explicitly with `mongodbatlas_service_account_secret`:

- [`without_initial_secret/`](without_initial_secret/README.md): create the Service Account with `without_initial_secret = true`, then create the first secret with `mongodbatlas_service_account_secret`.

Alternatively, omit `without_initial_secret` and set `secret_expires_after_hours` to let Atlas generate a bootstrap secret from the create request. The secret value is returned only once, in the create response, so the configuration cannot manage it afterwards. Use this only when you do not need to own the first secret.

`without_initial_secret` and `secret_expires_after_hours` are mutually exclusive: set one and omit the other.

For product details, see the [resource documentation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/service_account). For managing and rotating secrets, see [Guide: Service Account Secret Rotation](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/service-account-secret-rotation).
