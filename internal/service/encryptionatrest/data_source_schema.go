package encryptionatrest

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312002/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"aws_kms_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"access_key_id": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).",
					},
					"customer_master_key_id": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.",
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
					},
					"region": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Physical location where MongoDB Atlas deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Atlas deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Atlas creates them as part of the deployment. MongoDB Atlas assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.", //nolint:lll // reason: auto-generated from Open API spec.
					},
					"role_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.",
					},
					"secret_access_key": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.",
					},
					"valid": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.",
					},
					"require_private_networking": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Enable connection to your Amazon Web Services (AWS) Key Management Service (KMS) over private networking.",
					},
				},
				Computed:            true,
				MarkdownDescription: "Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project.",
			},
			"azure_key_vault_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"azure_environment": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Azure environment in which your account credentials reside.",
					},
					"client_id": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Unique 36-hexadecimal character string that identifies an Azure application associated with your Azure Active Directory tenant.",
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
					},
					"key_identifier": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Web address with a unique key that identifies for your Azure Key Vault.",
					},
					"key_vault_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Unique string that identifies the Azure Key Vault that contains your key.",
					},
					"require_private_networking": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Enable connection to your Azure Key Vault over private networking.",
					},
					"resource_group_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Name of the Azure resource group that contains your Azure Key Vault.",
					},
					"secret": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Private data that you need secured and that belongs to the specified Azure Key Vault (AKV) tenant (**azureKeyVault.tenantID**). This data can include any type of sensitive data such as passwords, database connection strings, API keys, and the like. AKV stores this information as encrypted binary data.",
					},
					"subscription_id": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Unique 36-hexadecimal character string that identifies your Azure subscription.",
					},
					"tenant_id": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Unique 36-hexadecimal character string that identifies the Azure Active Directory tenant within your Azure subscription.",
					},
					"valid": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether the Azure encryption key can encrypt and decrypt data.",
					},
				},
				Computed:            true,
				MarkdownDescription: "Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV).",
			},
			"google_cloud_kms_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
					},
					"key_version_resource_id": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "Resource path that displays the key version resource ID for your Google Cloud KMS.",
					},
					"service_account_key": schema.StringAttribute{
						Computed:            true,
						Sensitive:           true,
						MarkdownDescription: "JavaScript Object Notation (JSON) object that contains the Google Cloud Key Management Service (KMS). Format the JSON as a string and not as an object.",
					},
					"valid": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Flag that indicates whether the Google Cloud Key Management Service (KMS) encryption key can encrypt and decrypt data.",
					},
				},
				Computed:            true,
				MarkdownDescription: "Details that define the configuration of Encryption at Rest using Google Cloud Key Management Service (KMS).",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"enabled_for_search_nodes": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Flag that indicates whether Encryption at Rest for Dedicated Search Nodes is enabled in the specified project.",
			},
		},
	}
}

type TFEncryptionAtRestDSModel struct {
	AzureKeyVaultConfig   *TFAzureKeyVaultConfigModel `tfsdk:"azure_key_vault_config"`
	AwsKmsConfig          *TFAwsKmsConfigModel        `tfsdk:"aws_kms_config"`
	GoogleCloudKmsConfig  *TFGcpKmsConfigModel        `tfsdk:"google_cloud_kms_config"`
	ID                    types.String                `tfsdk:"id"`
	ProjectID             types.String                `tfsdk:"project_id"`
	EnabledForSearchNodes types.Bool                  `tfsdk:"enabled_for_search_nodes"`
}

func NewTFEncryptionAtRestDSModel(projectID string, encryptionResp *admin.EncryptionAtRest) *TFEncryptionAtRestDSModel {
	return &TFEncryptionAtRestDSModel{
		ID:                    types.StringValue(projectID),
		ProjectID:             types.StringValue(projectID),
		AwsKmsConfig:          NewTFAwsKmsConfigItem(encryptionResp.AwsKms),
		AzureKeyVaultConfig:   NewTFAzureKeyVaultConfigItem(encryptionResp.AzureKeyVault),
		GoogleCloudKmsConfig:  NewTFGcpKmsConfigItem(encryptionResp.GoogleCloudKms),
		EnabledForSearchNodes: types.BoolPointerValue(encryptionResp.EnabledForSearchNodes),
	}
}
