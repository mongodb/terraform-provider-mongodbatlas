---
page_title: "Migration Guide: Programmatic API Keys (PAKs) to Service Accounts (SAs)"
---

**Note:** Migration to Service Accounts is **not required**. If you are currently using `mongodbatlas_api_key`, `mongodbatlas_api_key_project_assignment`, and `mongodbatlas_access_list_api_key` resources, you may continue to do so. This guide is for users who wish to adopt Service Accounts for greater security, flexibility, or best practices, but existing PAK configurations will continue to work and are supported.

# Migration Guide: Programmatic API Keys (PAKs) to Service Accounts (SAs)

The goal of this guide is to help users transition from Programmatic API Keys (PAKs) to Service Accounts (SAs) in MongoDB Atlas. Service Accounts provide a more secure and flexible authentication method that eliminates the need to manage API key secrets in Terraform state.

### Jump to:
- [Stage 1: Initial State (PAK Resources Only)](#stage-1-initial-state-pak-resources-only)
- [Stage 2: Intermediate State (Both PAK and SA Resources)](#stage-2-intermediate-state-both-pak-and-sa-resources)
- [Stage 3: Final State (SA Resources Only)](#stage-3-final-state-sa-resources-only)

## Resource Mapping

The following table shows the mapping between PAK resources and their Service Account equivalents:

| PAK Resource | Service Account Resource | Notes |
|--------------|-------------------------|-------|
| `mongodbatlas_api_key` | `mongodbatlas_service_account` | Organization-level key/account |
| `mongodbatlas_api_key_project_assignment` | `mongodbatlas_service_account_project_assignment` | Project assignment |
| `mongodbatlas_access_list_api_key` | `mongodbatlas_service_account_access_list_entry` | IP access list entry |

## Stage 1: Initial State (PAK Resources Only)

Your current configuration using PAK resources:

```terraform
# Organization-level Programmatic API Key
resource "mongodbatlas_api_key" "example" {
  org_id      = var.org_id
  description = "Example API Key for project access"
  role_names  = ["ORG_READ_ONLY"]
}

# Project assignment for the API Key
resource "mongodbatlas_api_key_project_assignment" "example" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  roles      = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  cidr_block = "192.168.1.0/24"
  # Alternative: ip_address = "192.168.1.100"
}
```

## Stage 2: Intermediate State (Both PAK and SA Resources)

In this stage, you'll add Service Account resources alongside your existing PAK resources. This allows both authentication methods to work simultaneously, enabling you to test Service Accounts before removing PAKs.

The following steps explain how to add Service Account resources without affecting your existing PAK resources:

1. Add the Service Account resource to your configuration:

```terraform
# Service Account (new)
resource "mongodbatlas_service_account" "example" {
  org_id      = var.org_id
  name        = "example-service-account"
  description = "Example Service Account for project access"
  roles       = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

# Keep existing PAK resources (for now)
resource "mongodbatlas_api_key" "example" {
  org_id      = var.org_id
  description = "Example API Key for project access"
  role_names  = ["ORG_READ_ONLY"]
}
```

2. Add the project assignment for the Service Account:

```terraform
# Service Account Project Assignment (new)
resource "mongodbatlas_service_account_project_assignment" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.example.client_id
  roles      = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}

# Keep existing PAK project assignment (for now)
resource "mongodbatlas_api_key_project_assignment" "example" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  roles      = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}
```

3. Add the access list entry for the Service Account:

```terraform
# Service Account Access List Entry (new)
resource "mongodbatlas_service_account_access_list_entry" "example" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.example.client_id
  cidr_block = "192.168.1.0/24"
  # Alternative: ip_address = "192.168.1.100"
}

# Keep existing PAK access list entry (for now)
resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  cidr_block = "192.168.1.0/24"
}
```

4. Run `terraform plan` to review the changes.
5. Run `terraform apply` to create the Service Account resources.
6. **Important**: Save the Service Account secret from the output. When a Service Account is created, a secret is automatically generated. The secret value is only returned once at creation time.

   The example includes a sensitive output `service_account_first_secret` that captures this initial secret. 

   Add the data source and output blocks to capture the secret and service account information:
   ```terraform
   data "mongodbatlas_service_accounts" "this" {
     org_id = var.org_id
   }

   output "service_account_first_secret" {
     description = "The secret value of the first secret created with the service account. Only available after initial creation."
     value       = try(mongodbatlas_service_account.example.secrets[0].secret, null)
     sensitive   = true
   }

   output "service_accounts_results" {
     value = data.mongodbatlas_service_accounts.this.results
   }
   ```
   
   You can then retrieve it using:

   ```bash
   terraform output -raw service_account_first_secret
   ```
7. Test your Service Account by updating your provider configuration or using it in your applications.
8. Verify that both PAK and SA authentication methods work correctly.
9. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

## Stage 3: Final State (SA Resources Only)

Once you've verified that the Service Account works correctly, you can remove the PAK resources. The following steps explain how to remove PAK resources from your Terraform configuration:

1. Remove the PAK resource blocks from your Terraform configuration:

```terraform
# Service Account Resources (FINAL STATE)
resource "mongodbatlas_service_account" "example" {
  org_id      = var.org_id
  name        = "example-service-account"
  description = "Example Service Account for project access"
  roles       = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_service_account_project_assignment" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.example.client_id
  roles      = ["GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"]
}

resource "mongodbatlas_service_account_access_list_entry" "example" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.example.client_id
  cidr_block = "192.168.1.0/24"
}
```

2. Run `terraform plan` to verify:
   - PAK resources are planned for destruction
   - Only the Service Account resources remain
   - No unexpected changes

3. Run `terraform apply` to finalize the migration. This will delete the PAK resources from Atlas.

4. Verify that your applications and infrastructure continue to work with Service Accounts.

5. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`