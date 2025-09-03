# MongoDB Atlas Cloud Provider Access with Google Cloud Platform (GCP)

This example demonstrates how to set up MongoDB Atlas Cloud Provider Access with Google Cloud Platform (GCP) to enable Customer-Managed Encryption Keys (CMEK) using Google Cloud KMS.

## What This Example Does

This Terraform configuration:

1. **Creates Cloud Provider Access Setup**: Establishes the initial setup for MongoDB Atlas to access GCP resources
2. **Authorizes the Access Role**: Completes the authorization process, creating a GCP service account that MongoDB Atlas can use
3. **Creates GCP KMS Resources**: Sets up a Key Ring and Crypto Key in Google Cloud KMS for encryption
4. **Configures IAM Permissions**: Grants the MongoDB Atlas service account the necessary permissions to use the KMS key
5. **Enables Encryption at Rest**: Configures MongoDB Atlas to use the GCP KMS key for encrypting data at rest

## Architecture Overview

```
┌─────────────────────┐    ┌─────────────────────────┐
│   MongoDB Atlas     │    │     Google Cloud        │
│                     │    │                         │
│  ┌───────────────┐  │    │  ┌─────────────────┐    │
│  │   Project     │  │    │  │   KMS Key Ring  │    │
│  │               │  │◄───┤  │                 │    │
│  │ Cloud Provider│  │    │  │ ┌─────────────┐ │    │
│  │    Access     │  │    │  │ │ Crypto Key  │ │    │
│  │               │  │    │  │ └─────────────┘ │    │
│  └───────────────┘  │    │  └─────────────────┘    │
│                     │    │                         │
│  Service Account ───┼────┤► IAM Permissions        │
│  (Created by Atlas) │    │  - cryptoKeyEncrypter   │
│                     │    │    Decrypter            │
│                     │    │  - viewer               │
└─────────────────────┘    └─────────────────────────┘
```

## Prerequisites

Before running this example, you need:

1. **MongoDB Atlas Account**: With API keys that have project owner permissions
2. **Google Cloud Platform Account**: With a project and appropriate permissions to:
   - Create KMS resources (Key Rings and Crypto Keys)
   - Manage IAM bindings on KMS resources
3. **Terraform**: Version 0.13 or later
4. **Google Cloud Terraform Provider**: The Google provider must be configured in your Terraform environment. See the [Google provider documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs) for authentication methods 

### Required GCP Permissions

Your user or service account needs the following IAM roles:
- `roles/cloudkms.admin` - To create and manage KMS resources
- `roles/resourcemanager.projectIamAdmin` - To manage IAM bindings

## Usage

### 1. Set Up Variables

Create a `terraform.tfvars` file:

```hcl
atlas_public_key  = <ATLAS_PUBLIC_KEY>
atlas_private_key = <ATLAS_PRIVATE_KEY>
atlas_project_id  = <ATLAS_PROJECT_ID>
gcp_project_id    = <GCP_PROJECT_ID>

# Optional: Customize KMS resources
key_ring_name    = <KEY_RING_NAME>
crypto_key_name  = <CRYPTO_KEY_NAME>
location         = <GCP_LOCATION>
```

### 2. Deploy the Infrastructure

```bash
terraform init
terraform plan
terraform apply
```

### 3. Verify the Setup

After successful deployment, you should see outputs including:
- Atlas role ID
- GCP service account email created by Atlas
- KMS key ring and crypto key IDs
- Key version resource ID used for encryption

## Important Notes

### GCP-Specific Behavior

Unlike AWS and Azure, GCP Cloud Provider Access:
- **No Configuration Updates**: GCP authorization only requires a role ID and has no additional configuration parameters
- **Immutable After Creation**: Once authorized, you cannot "update" a GCP cloud provider access role
- **New Authorization = New Resource**: If you need to change GCP settings, create a new `mongodbatlas_cloud_provider_access_setup` and `mongodbatlas_cloud_provider_access_authorization` resource and then delete the old one

### Resource Dependencies

The configuration manages dependencies between resources:
- KMS resources are created first
- Atlas authorization completes before encryption configuration
- IAM bindings are established before Atlas tries to use the key

### Security Considerations

- The Atlas service account gets minimal required permissions (`cryptoKeyEncrypterDecrypter` and `viewer`)
- KMS keys are created with `ENCRYPT_DECRYPT` purpose only

## Variables

| Variable | Description | Type | Default | Required |
|----------|-------------|------|---------|----------|
| `atlas_public_key` | MongoDB Atlas public API key | string | - | Yes |
| `atlas_private_key` | MongoDB Atlas private API key | string | - | Yes |
| `atlas_project_id` | MongoDB Atlas project ID | string | - | Yes |
| `gcp_project_id` | GCP project ID | string | - | Yes |
| `key_ring_name` | Name of the KMS key ring | string | `"atlas-key-ring"` | No |
| `crypto_key_name` | Name of the crypto key | string | `"atlas-crypto-key"` | No |
| `location` | GCP region for KMS resources | string | `"us-central1"` | No |

## Outputs

| Output | Description |
|--------|-------------|
| `atlas_role_id` | The MongoDB Atlas cloud provider access role ID |
| `gcp_service_account_email` | GCP service account email created by Atlas |
| `kms_key_ring_id` | Full ID of the created KMS key ring |
| `kms_crypto_key_id` | Full ID of the created crypto key |
| `kms_key_version_resource_id` | Resource ID of the primary key version |

## Related Documentation

- [MongoDB Atlas Cloud Provider Access](https://www.mongodb.com/docs/atlas/security/customer-key-management/)
- [Google Cloud KMS Documentation](https://cloud.google.com/kms/docs)
- [Terraform MongoDB Atlas Provider](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs)
