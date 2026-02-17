---
page_title: "Migration Guide: Encryption at Rest (Azure) Client Credentials to Role-based Auth"
---

# Migration Guide: Encryption at Rest (Azure) Client Credentials to Role-based Auth

**Objective**: Migrate from using Azure client credentials (`client_id`, `tenant_id`, and `secret`) to role-based authentication using an Atlas-managed role via `mongodbatlas_encryption_at_rest.azure_key_vault_config.role_id`.

## Best Practices Before Migrating

- Test the migration in a non-production environment if possible.

## Migration Steps

### Current (using Client Credentials)

This configuration serves as the starting point for the migration.

```hcl
resource "mongodbatlas_encryption_at_rest" "this" {
  project_id = var.atlas_project_id

  azure_key_vault_config {
    enabled           = true
    azure_environment = "AZURE"

    tenant_id       = var.azure_tenant_id
    subscription_id = var.azure_subscription_id
    client_id       = var.azure_client_id
    secret          = var.azure_client_secret

    resource_group_name = var.azure_resource_group_name
    key_vault_name      = var.azure_key_vault_name
    key_identifier      = var.azure_key_identifier
  }
}
```

### 1) Obtain the Atlas-managed Azure role

Add the following resources to enable and authorize Atlas Cloud Provider Access for Azure for your project:

```hcl
data "azuread_service_principal" "atlas_sp" {
  client_id = var.atlas_azure_app_id
}

resource "mongodbatlas_cloud_provider_access_setup" "this" {
    project_id    = var.atlas_project_id
    provider_name = "AZURE"

    azure_config {
        atlas_azure_app_id   = var.atlas_azure_app_id
        service_principal_id = data.azuread_service_principal.atlas_sp.object_id
        tenant_id            = var.azure_tenant_id
    }
}

resource "mongodbatlas_cloud_provider_access_authorization" "this" {
    project_id = var.atlas_project_id
    role_id    = mongodbatlas_cloud_provider_access_setup.this.role_id

    azure {
        atlas_azure_app_id   = var.atlas_azure_app_id
        service_principal_id = data.azuread_service_principal.atlas_sp.object_id
        tenant_id            = var.azure_tenant_id
    }
}
```

The value `mongodbatlas_cloud_provider_access_authorization.this.role_id` is the Azure role identifier you use as `role_id` in `azure_key_vault_config`.

### 2) Grant Key Vault permissions to the Atlas role

Grant your Azure Service Principal the required permissions on the Key Vault used by Atlas. You can use either access policies or RBAC role assignments. Below are examples with the AzureRM provider.

First, ensure you have a reference to your existing Key Vault:

```hcl
data "azurerm_key_vault" "existing_kv" {
  name                = var.azure_key_vault_name
  resource_group_name = var.azure_resource_group_name
}
```

Then, grant access policy permissions for key operations (Get/Encrypt/Decrypt):

```hcl
resource "azurerm_key_vault_access_policy" "kv_crypto_perms" {
  key_vault_id = data.azurerm_key_vault.existing_kv.id

  tenant_id = var.azure_tenant_id
  object_id = data.azuread_service_principal.atlas_sp.object_id

  key_permissions = [
    "Get",
    "Encrypt",
    "Decrypt"
  ]
}
```

### 3) Update the Encryption at Rest resource to use role-based auth

Replace the `client_id`, `secret`, and `tenant_id` with `role_id` using the value from the authorization resource:

```hcl
resource "mongodbatlas_encryption_at_rest" "this" {
  project_id = var.atlas_project_id

  azure_key_vault_config {
    enabled           = true
    azure_environment = "AZURE"

    subscription_id     = var.azure_subscription_id
    resource_group_name = var.azure_resource_group_name
    key_vault_name      = var.azure_key_vault_name
    key_identifier      = var.azure_key_identifier

    role_id = mongodbatlas_cloud_provider_access_authorization.this.role_id
  }

  depends_on = [
    azurerm_key_vault_access_policy.kv_crypto_perms
  ]
}
```

**Note:** The `depends_on` block ensures that Terraform configures the Key Vault permissions before it configures the encryption at rest resource. This is necessary when you create both resources in the same `terraform apply` execution.

### 4) Review your plan

Run `terraform plan` and review the plan output carefully before confirming. The plan should be similar to the following example:

```shell
# mongodbatlas_encryption_at_rest.this will be updated in-place
  ~ resource "mongodbatlas_encryption_at_rest" "this" {
        id = "<project_encryption_at_rest_id>"
        # (2 unchanged attributes hidden)

      ~ azure_key_vault_config {
          + role_id       = "<YOUR_ROLE_ID>"
          - secret              = (sensitive value) -> null
          ~ valid               = true -> (known after apply)
            # (other unchanged attributes hidden)
        }
    }
```

### 5) Apply the changes

```shell
terraform apply
```

Once applied, your encryption at rest configuration starts using role-based authentication instead of client credentials.

## Additional Resources

- [Azure Encryption at Rest Example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_encryption_at_rest/azure)
