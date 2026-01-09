---
page_title: "Migration Guide: Programmatic API Keys (PAKs) to Service Accounts (SAs)"
---

# Migration Guide: Programmatic API Keys (PAKs) to Service Accounts (SAs)

## Overview

Service Accounts are the recommended method to manage authentication to the Atlas Administration API. Service Accounts provide improved security over API keys by using the industry standard OAuth 2.0 protocol with the Client Credentials flow. This guide covers migrating from Programmatic API Keys (PAKs) to Service Accounts (SAs) in MongoDB Atlas.

**Note:** Migration to Service Accounts is **not required**. If you are currently using API Key resources, you may continue to do so. This guide is for users who wish to adopt Service Accounts for greater security or best practices, but existing PAK configurations will continue to work and are supported.

## Before You Begin

- **Backup your Terraform state file** before making any changes.
- **Test the process in a non-production environment** if possible.
- **Secrets handling** - Managing Service Accounts with Terraform will expose sensitive organizational secrets in Terraform's state. We suggest following [Terraform's best practices](https://developer.hashicorp.com/terraform/language/state/sensitive-data).


---

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Organization-Level Migration</span></summary>

## Organization-Level API Keys to Service Accounts

**Objective**: Migrate from organization-level PAK resources (`mongodbatlas_api_key`, `mongodbatlas_api_key_project_assignment`, `mongodbatlas_access_list_api_key`) to organization-level Service Account resources (`mongodbatlas_service_account`, `mongodbatlas_service_account_project_assignment`, `mongodbatlas_service_account_access_list_entry`).

## Resource Mapping
The following table shows the mapping between organization-level PAK resources and their Service Account equivalents:

| PAK Resource | Service Account Resource | Notes |
|--------------|-------------------------|-------|
| `mongodbatlas_api_key` | `mongodbatlas_service_account` | Organization-level API key / Service Account |
| `mongodbatlas_api_key_project_assignment` | `mongodbatlas_service_account_project_assignment` | Project assignment |
| `mongodbatlas_access_list_api_key` | `mongodbatlas_service_account_access_list_entry` | IP access list entry |

---

## Migration Steps

For complete working examples, see the [organization-level migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_pak_to_service_account/org_level).

### Step 1: Initial State - PAK Resources Only

This is your starting configuration with organization-level PAK resources (org PAK + assignment to project + access list entry):

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
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}
```

### Step 2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources

Add the Service Account resources to your configuration while keeping the existing PAK resources. This allows both authentication methods to work simultaneously, enabling you to test Service Accounts before removing PAKs.

```terraform
# Service Account (new)
resource "mongodbatlas_service_account" "example" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account for project access"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

# Service Account Project Assignment (new)
resource "mongodbatlas_service_account_project_assignment" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.example.client_id
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

# Service Account Access List Entry (new)
resource "mongodbatlas_service_account_access_list_entry" "example" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.example.client_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}

# Output to capture the secret (add this before running apply)
output "service_account_first_secret" {
  description = "The secret value of the first secret created with the service account. Only available after initial creation."
  value       = try(mongodbatlas_service_account.example.secrets[0].secret, null)
  sensitive   = true
}
```

```terraform
# Keep existing PAK resources (for now)
resource "mongodbatlas_api_key" "example" {
  org_id      = var.org_id
  description = "Example API Key for project access"
  role_names  = ["ORG_READ_ONLY"]
}

resource "mongodbatlas_api_key_project_assignment" "example" {
  project_id = var.project_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_api_key.example.api_key_id
  cidr_block = "192.168.1.100/32"
}
```

**Apply and test:**

1. Run `terraform plan` to review the changes.
2. Run `terraform apply` to create the Service Account resources.
3. **Important**: Save the Service Account secret from the output. The secret value is only returned once at creation time.

   You can retrieve it using:

   ```bash
   terraform output -raw service_account_first_secret
   ```

4. Test your Service Account by updating your provider configuration or using it in your applications.
5. Verify that both PAK and SA authentication methods work correctly.
6. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

### Step 3: Final State - Remove PAK Resources, SA Resources Only

Once you've verified that the Service Account works correctly, remove the PAK resources from your Terraform configuration:

Once you've verified that the Service Account works correctly, remove the PAK resources from your Terraform configuration:

```terraform
# Service Account Resources (FINAL STATE)
resource "mongodbatlas_service_account" "example" {
  org_id                     = var.org_id
  name                       = "example-service-account"
  description                = "Example Service Account for project access"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_service_account_project_assignment" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_service_account.example.client_id
  roles      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
}

resource "mongodbatlas_service_account_access_list_entry" "example" {
  org_id     = var.org_id
  client_id  = mongodbatlas_service_account.example.client_id
  cidr_block = "192.168.1.100/32"
}

output "service_account_first_secret" {
  description = "The secret value of the first secret created with the service account. Only available after initial creation."
  value       = try(mongodbatlas_service_account.example.secrets[0].secret, null)
  sensitive   = true
}
```

1. Run `terraform plan` to verify:
   - PAK resources are planned for destruction
   - Only the Service Account resources remain
   - No unexpected changes

2. Run `terraform apply` to finalize the migration. This will delete the PAK resources from Atlas.

3. Verify that your applications and infrastructure continue to work with Service Accounts.

4. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

---

- **Important:** The Service Account secret is only returned once at creation time. Make sure to save it securely before proceeding.

- After successful migration, ensure no references to PAK resources remain in your configuration.


</details>

---

<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Project-Level Migration</span></summary>

## Project-Level API Keys to Service Accounts

**Objective**: Migrate from project-level PAK resources (`mongodbatlas_project_api_key`, `mongodbatlas_access_list_api_key`) to project-level Service Account resources (`mongodbatlas_project_service_account`, `mongodbatlas_project_service_account_access_list_entry`).

**Important:** Organization-level resources (`mongodbatlas_service_account`) are the recommended approach. Project-level resources (`mongodbatlas_project_service_account`) should only be used if you do not have organization-level permissions to manage service accounts. Otherwise, use the [Organization-Level Migration](#organization-level-api-keys-to-service-accounts) approach.

## Resource Mapping

The following table shows the mapping between project-level PAK resources and their Service Account equivalents:

| PAK Resource | Service Account Resource | Notes |
|--------------|-------------------------|-------|
| `mongodbatlas_project_api_key` | `mongodbatlas_project_service_account` | Project-level API key / Service Account (use only if you don't have org-level permissions) |
| `mongodbatlas_access_list_api_key` | `mongodbatlas_project_service_account_access_list_entry` | IP access list entry for project-level service account |

**Important:** Organization-level resources (`mongodbatlas_service_account`) are the recommended approach. Project-level resources (`mongodbatlas_project_service_account`) should only be used if you do not have organization-level permissions to manage service accounts.

---

## Migration Steps

For complete working examples, see the [project-level migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_pak_to_service_account/project_level).

### Step 1: Initial State - PAK Resources Only

This is your starting configuration with project-level PAK resources (project PAK + access list entry):

```terraform
# Project-level Programmatic API Key
resource "mongodbatlas_project_api_key" "example" {
  description = "Example Project API Key"
  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  }
}

# IP Access List entry for the API Key
resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_project_api_key.example.api_key_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}
```

### Step 2: Intermediate State - Add Service Account Resources Alongside Existing PAK Resources

Add the Service Account resources to your configuration while keeping the existing PAK resources. This allows both authentication methods to work simultaneously, enabling you to test Service Accounts before removing PAKs.

```terraform
# Project Service Account (new)
resource "mongodbatlas_project_service_account" "example" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

# Project Service Account Access List Entry (new)
resource "mongodbatlas_project_service_account_access_list_entry" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.example.client_id
  cidr_block = "192.168.1.100/32"
  # Alternative: ip_address = "192.168.1.100"
}

# Output to capture the secret (add this before running apply)
output "project_service_account_first_secret" {
  description = "The secret value of the first secret created with the project service account. Only available after initial creation."
  value       = try(mongodbatlas_project_service_account.example.secrets[0].secret, null)
  sensitive   = true
}
```

```terraform
# Keep existing PAK resources (for now)
resource "mongodbatlas_project_api_key" "example" {
  description = "Example Project API Key"
  project_assignment {
    project_id = var.project_id
    role_names = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  }
}

resource "mongodbatlas_access_list_api_key" "example" {
  org_id     = var.org_id
  api_key_id = mongodbatlas_project_api_key.example.api_key_id
  cidr_block = "192.168.1.100/32"
}
```

**Apply and test:**

1. Run `terraform plan` to review the changes.
2. Run `terraform apply` to create the Service Account resource.
3. **Important**: Save the Service Account secret from the output. The secret value is only returned once at creation time.

   You can retrieve it using:

   ```bash
   terraform output -raw project_service_account_first_secret
   ```

4. Test your Service Account by updating your provider configuration or using it in your applications.
5. Verify that both PAK and SA authentication methods work correctly.
6. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

### Step 3: Final State - Remove PAK Resources, SA Resources Only

Once you've verified that the Service Account works correctly, remove the PAK resources from your Terraform configuration:

```terraform
# Project Service Account Resources (FINAL STATE)
resource "mongodbatlas_project_service_account" "example" {
  project_id                 = var.project_id
  name                       = "example-project-service-account"
  description                = "Example Project Service Account"
  roles                      = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
  secret_expires_after_hours = 2160 # 90 days
}

resource "mongodbatlas_project_service_account_access_list_entry" "example" {
  project_id = var.project_id
  client_id  = mongodbatlas_project_service_account.example.client_id
  cidr_block = "192.168.1.100/32"
}

output "project_service_account_first_secret" {
  description = "The secret value of the first secret created with the project service account. Only available after initial creation."
  value       = try(mongodbatlas_project_service_account.example.secrets[0].secret, null)
  sensitive   = true
}
```

1. Run `terraform plan` to verify:
   - PAK resources are planned for destruction
   - Only the Service Account resources remain
   - No unexpected changes

2. Run `terraform apply` to finalize the migration. This will delete the PAK resource from Atlas.

3. Verify that your applications and infrastructure continue to work with Service Accounts.

4. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

---

- **Important:** The Service Account secret is only returned once at creation time. Make sure to save it securely before proceeding.

- After successful migration, ensure no references to PAK resources remain in your configuration.
