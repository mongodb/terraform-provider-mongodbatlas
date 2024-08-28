package encryptionatrest

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20240805001/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// TODO: check for sensitive attr
// TODO: check about ID attr
// TODO: check if we can add 'valid' to resource & re-use models
func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"aws_kms_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"access_key_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).",
						MarkdownDescription: "Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).",
					},
					"customer_master_key_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.",
						MarkdownDescription: "Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.",
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
						MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
					},
					"region": schema.StringAttribute{
						Computed:            true,
						Description:         "Physical location where MongoDB Cloud deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. MongoDB Cloud assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.", //nolint:lll // reason: auto-generated from Open API spec.
						MarkdownDescription: "Physical location where MongoDB Cloud deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. MongoDB Cloud assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.", //nolint:lll // reason: auto-generated from Open API spec.
					},
					"role_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.",
						MarkdownDescription: "Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.",
					},
					"secret_access_key": schema.StringAttribute{
						Computed:            true,
						Description:         "Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.",
						MarkdownDescription: "Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.",
					},
					"valid": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.",
						MarkdownDescription: "Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.",
					},
				},
				Computed:            true,
				Description:         "Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project.",
				MarkdownDescription: "Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project.",
			},
			"azure_key_vault_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"azure_environment": schema.StringAttribute{
						Computed:            true,
						Description:         "Azure environment in which your account credentials reside.",
						MarkdownDescription: "Azure environment in which your account credentials reside.",
					},
					"client_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique 36-hexadecimal character string that identifies an Azure application associated with your Azure Active Directory tenant.",
						MarkdownDescription: "Unique 36-hexadecimal character string that identifies an Azure application associated with your Azure Active Directory tenant.",
					},
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
						MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
					},
					"key_identifier": schema.StringAttribute{
						Computed:            true,
						Description:         "Web address with a unique key that identifies for your Azure Key Vault.",
						MarkdownDescription: "Web address with a unique key that identifies for your Azure Key Vault.",
					},
					"key_vault_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique string that identifies the Azure Key Vault that contains your key.",
						MarkdownDescription: "Unique string that identifies the Azure Key Vault that contains your key.",
					},
					"require_private_networking": schema.BoolAttribute{
						Computed:            true,
						Description:         "Enable connection to your Azure Key Vault over private networking.",
						MarkdownDescription: "Enable connection to your Azure Key Vault over private networking.",
					},
					"resource_group_name": schema.StringAttribute{
						Computed:            true,
						Description:         "Name of the Azure resource group that contains your Azure Key Vault.",
						MarkdownDescription: "Name of the Azure resource group that contains your Azure Key Vault.",
					},
					"secret": schema.StringAttribute{
						Computed:            true,
						Description:         "Private data that you need secured and that belongs to the specified Azure Key Vault (AKV) tenant (**azureKeyVault.tenantID**). This data can include any type of sensitive data such as passwords, database connection strings, API keys, and the like. AKV stores this information as encrypted binary data.",
						MarkdownDescription: "Private data that you need secured and that belongs to the specified Azure Key Vault (AKV) tenant (**azureKeyVault.tenantID**). This data can include any type of sensitive data such as passwords, database connection strings, API keys, and the like. AKV stores this information as encrypted binary data.",
					},
					"subscription_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique 36-hexadecimal character string that identifies your Azure subscription.",
						MarkdownDescription: "Unique 36-hexadecimal character string that identifies your Azure subscription.",
					},
					"tenant_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Unique 36-hexadecimal character string that identifies the Azure Active Directory tenant within your Azure subscription.",
						MarkdownDescription: "Unique 36-hexadecimal character string that identifies the Azure Active Directory tenant within your Azure subscription.",
					},
					"valid": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether the Azure encryption key can encrypt and decrypt data.",
						MarkdownDescription: "Flag that indicates whether the Azure encryption key can encrypt and decrypt data.",
					},
				},
				Computed:            true,
				Description:         "Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV).",
				MarkdownDescription: "Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV).",
			},
			"google_cloud_kms_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
						MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
					},
					"key_version_resource_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Resource path that displays the key version resource ID for your Google Cloud KMS.",
						MarkdownDescription: "Resource path that displays the key version resource ID for your Google Cloud KMS.",
					},
					"service_account_key": schema.StringAttribute{
						Computed:            true,
						Description:         "JavaScript Object Notation (JSON) object that contains the Google Cloud Key Management Service (KMS). Format the JSON as a string and not as an object.",
						MarkdownDescription: "JavaScript Object Notation (JSON) object that contains the Google Cloud Key Management Service (KMS). Format the JSON as a string and not as an object.",
					},
					"valid": schema.BoolAttribute{
						Computed:            true,
						Description:         "Flag that indicates whether the Google Cloud Key Management Service (KMS) encryption key can encrypt and decrypt data.",
						MarkdownDescription: "Flag that indicates whether the Google Cloud Key Management Service (KMS) encryption key can encrypt and decrypt data.",
					},
				},
				Computed:            true,
				Description:         "Details that define the configuration of Encryption at Rest using Google Cloud Key Management Service (KMS).",
				MarkdownDescription: "Details that define the configuration of Encryption at Rest using Google Cloud Key Management Service (KMS).",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
		},
	}
}

type TfEncryptionAtRestDSModel struct {
	ID                   types.String               `tfsdk:"id"`
	ProjectID            types.String               `tfsdk:"project_id"`
	AzureKeyVaultConfig  TfAzureKeyVaultConfigModel `tfsdk:"azure_key_vault_config"`
	AwsKmsConfig         TfAwsKmsConfigModel        `tfsdk:"aws_kms_config"`
	GoogleCloudKmsConfig TfGcpKmsConfigModel        `tfsdk:"google_cloud_kms_config"`
}

func NewTfEncryptionAtRestDSModel(projectID string, encryptionResp *admin.EncryptionAtRest) *TfEncryptionAtRestDSModel {
	return &TfEncryptionAtRestDSModel{
		ID:                   types.StringValue(projectID),
		ProjectID:            types.StringValue(projectID),
		AwsKmsConfig:         *NewTFAwsKmsConfigItem(encryptionResp.AwsKms),
		AzureKeyVaultConfig:  *NewTFAzureKeyVaultConfigItem(encryptionResp.AzureKeyVault),
		GoogleCloudKmsConfig: *NewTFGcpKmsConfigItem(encryptionResp.GoogleCloudKms),
	}
}

func NewTFAwsKmsConfigItem(awsKms *admin.AWSKMSConfiguration) *TfAwsKmsConfigModel {
	if awsKms == nil {
		return nil
	}

	return &TfAwsKmsConfigModel{
		Enabled:             types.BoolPointerValue(awsKms.Enabled),
		CustomerMasterKeyID: types.StringValue(awsKms.GetCustomerMasterKeyID()),
		Region:              types.StringValue(awsKms.GetRegion()),
		AccessKeyID:         conversion.StringNullIfEmpty(awsKms.GetAccessKeyID()),
		SecretAccessKey:     conversion.StringNullIfEmpty(awsKms.GetSecretAccessKey()),
		RoleID:              conversion.StringNullIfEmpty(awsKms.GetRoleId()),
		Valid:               types.BoolPointerValue(awsKms.Valid),
	}
}

func NewTFAzureKeyVaultConfigItem(az *admin.AzureKeyVault) *TfAzureKeyVaultConfigModel {
	if az == nil {
		return nil
	}

	return &TfAzureKeyVaultConfigModel{
		Enabled:                  types.BoolPointerValue(az.Enabled),
		ClientID:                 types.StringValue(az.GetClientID()),
		AzureEnvironment:         types.StringValue(az.GetAzureEnvironment()),
		SubscriptionID:           types.StringValue(az.GetSubscriptionID()),
		ResourceGroupName:        types.StringValue(az.GetResourceGroupName()),
		KeyVaultName:             types.StringValue(az.GetKeyVaultName()),
		KeyIdentifier:            types.StringValue(az.GetKeyIdentifier()),
		TenantID:                 types.StringValue(az.GetTenantID()),
		Secret:                   conversion.StringNullIfEmpty(az.GetSecret()),
		RequirePrivateNetworking: types.BoolValue(az.GetRequirePrivateNetworking()),
		Valid:                    types.BoolPointerValue(az.Valid),
	}
}

func NewTFGcpKmsConfigItem(gcpKms *admin.GoogleCloudKMS) *TfGcpKmsConfigModel {
	if gcpKms == nil {
		return nil
	}

	return &TfGcpKmsConfigModel{
		Enabled:              types.BoolPointerValue(gcpKms.Enabled),
		KeyVersionResourceID: types.StringValue(gcpKms.GetKeyVersionResourceID()),
		ServiceAccountKey:    conversion.StringNullIfEmpty(gcpKms.GetServiceAccountKey()),
		Valid:                types.BoolPointerValue(gcpKms.Valid),
	}
}

// type TfAwsKmsConfigDSModel struct {
// 	AccessKeyID         types.String `tfsdk:"access_key_id"`
// 	SecretAccessKey     types.String `tfsdk:"secret_access_key"`
// 	CustomerMasterKeyID types.String `tfsdk:"customer_master_key_id"`
// 	Region              types.String `tfsdk:"region"`
// 	RoleID              types.String `tfsdk:"role_id"`
// 	Enabled             types.Bool   `tfsdk:"enabled"`
// 	Valid               types.Bool   `tfsdk:"valid"`
// }

// type TfAzureKeyVaultConfigDSModel struct {
// 	ClientID                 types.String `tfsdk:"client_id"`
// 	AzureEnvironment         types.String `tfsdk:"azure_environment"`
// 	SubscriptionID           types.String `tfsdk:"subscription_id"`
// 	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
// 	KeyVaultName             types.String `tfsdk:"key_vault_name"`
// 	KeyIdentifier            types.String `tfsdk:"key_identifier"`
// 	Secret                   types.String `tfsdk:"secret"`
// 	TenantID                 types.String `tfsdk:"tenant_id"`
// 	Enabled                  types.Bool   `tfsdk:"enabled"`
// 	RequirePrivateNetworking types.Bool   `tfsdk:"require_private_networking"`
// 	Valid                    types.Bool   `tfsdk:"valid"`
// }

// type TfGcpKmsConfigDSModel struct {
// 	ServiceAccountKey    types.String `tfsdk:"service_account_key"`
// 	KeyVersionResourceID types.String `tfsdk:"key_version_resource_id"`
// 	Enabled              types.Bool   `tfsdk:"enabled"`
// 	Valid                types.Bool   `tfsdk:"valid"`
// }
