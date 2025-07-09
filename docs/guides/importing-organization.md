---
page_title: "Guide: Importing MongoDB Atlas Organizations"
---

# Importing MongoDB Atlas Organizations

**Objective**: Learn how to import existing MongoDB Atlas organizations into Terraform management, understand why certain attributes are optional for imports, and successfully manage your existing organizations using Infrastructure as Code.

## Overview

Starting with version 1.38.0 of the MongoDB Atlas Provider, you can now import existing MongoDB Atlas organizations into your Terraform state. This feature allows you to bring organizations that were created outside of Terraform under Infrastructure as Code (IaC) management.

## Why Import Organizations?

Before version 1.38.0, the `mongodbatlas_organization` resource did not support the import functionality. Organizations could only be created through Terraform, which meant:

- Existing organizations created through the Atlas UI or API couldn't be managed by Terraform.
- Teams had to maintain organizations outside of their IaC workflows.
- There was no migration path to bring existing organizations under Terraform management.

With import support, you can now:
- Bring existing organizations into Terraform management.
- Standardize organization configuration across your infrastructure.
- Apply consistent security settings and policies using IaC.

## Important Changes to Required Attributes

To enable import functionality, several attributes that were previously marked as required have been made optional:

- `org_owner_id`
- `description` 
- `role_names`
- `federation_settings_id` (was already optional but it's also creation-only)

### Why These Changes Were Necessary

These attributes were required attributes that:
1. Are only needed when creating a new organization.
2. Cannot be modified after the organization is created.
3. Are not returned by the Atlas API when reading existing organizations.

**Important:** While these attributes are now optional in the schema, they are still **required when creating a new organization**. The provider will validate that these attributes are present during the creation process.

## How to Import an Organization

### Prerequisites

Before importing an organization, ensure you have:
- MongoDB Atlas API credentials with Organization Owner permissions.
- The Organization ID of the organization you want to import.
- Terraform version 0.15 or later.
- MongoDB Atlas Provider version 1.38.0 or later.

### Finding Your Organization ID

You can find your organization ID in several ways:

1. **Atlas UI**: Navigate to your organization settings - the ID is visible in the URL.
2. **Atlas CLI**: Use `atlas organizations list` if you have the Atlas CLI installed.
3. **Atlas API**: Use the organizations endpoint to list all organizations.

### Import Process

#### Step 1: Create the Terraform Configuration

Create a Terraform configuration file that describes the organization you want to import. **Do not include** the creation-only attributes:

```hcl
resource "mongodbatlas_organization" "imported" {
  name = "My Existing Organization"
}
```

#### Step 2: Run the Import Command

Execute the import command with your organization ID:

```bash
terraform import mongodbatlas_organization.imported <ORG_ID>
```

Replace `<ORG_ID>` with your actual organization ID.

#### Step 3: Verify the Import

Run `terraform plan` to verify that the import was successful and see if any configuration adjustments are needed:

```bash
terraform plan
```

The plan should show minimal or no changes if your configuration matches the existing organization settings.

#### Step 4: Apply Any Updates

If you want to update any organization settings, modify your configuration and apply changs, e.g.:

```bash
resource "mongodbatlas_organization" "imported" {
  name = "My Existing Organization"

  # Add any settings that you want to modify
  api_access_list_required     = false
  multi_factor_auth_required   = true
  restrict_employee_access     = true
  gen_ai_features_enabled      = true
  security_contact             = "security@example.com"
  skip_default_alerts_settings = true
}
```

## Important Considerations

### API Credentials Usage  

- When importing an organization, the provider API credentials are used. Ensure that your API key has Organization Owner permissions for the organization being imported.
- When creating an organization, the process remains unchanged: new API keys are generated as part of the resource creation, and these keys will be used for subsequent `mongodbatlas_organization` resource operations.

**Important**: API credentials stored in the `mongodbatlas_organization` Terraform state will take precedence, regardless of their validity.

### Creation-Only Attributes

The newly-declared optional attributes **cannot be specified** when importing:
- `org_owner_id` - The organization owner is already set.
- `description` - API key description from organization creation.
- `role_names` - API key roles from organization creation.
- `federation_settings_id` - Federation settings from organization creation.

If you include these in your configuration when importing, Terraform will show them as changes that cannot be applied.

### State Management

After a successful import:
- The organization will be fully managed by Terraform.
- You can update any modifiable attributes through Terraform.
- The `private_key` and `public_key` attributes will be empty (as these are only generated during organization creation).

## Example: Complete Import Workflow

For a complete example of importing an organization, including all configuration files and detailed instructions, see the [organization-import example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_organization/organization-import) in the provider repository.

## Troubleshooting

### Common Issues

1. **"cannot be changed after creation" errors during import**: Ensure you haven't included creation-only attributes in your configuration
2. **Permission errors**: Verify your provider API key has Organization Owner role
3. **Resource not found**: Double-check the organization ID is correct

## See Also

- [mongodbatlas_organization Resource Documentation](../resources/organization)
- [Example: Creating a New Organization](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_organization) 
- [MongoDB Atlas Admin API Organization](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-organizations).
