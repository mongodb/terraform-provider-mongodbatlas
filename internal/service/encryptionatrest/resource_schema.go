package encryptionatrest

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
		},
		Blocks: map[string]schema.Block{
			"aws_kms_config": schema.ListNestedBlock{
				MarkdownDescription: "Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project.",
				Validators:          []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
						},
						"access_key_id": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).",
						},
						"secret_access_key": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.",
						},
						"customer_master_key_id": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.",
						},
						"region": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Physical location where MongoDB Atlas deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Atlas creates them as part of the deployment. MongoDB Atlas assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.", //nolint:lll // reason: auto-generated from Open API spec.
						},
						"role_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.",
						},
						"valid": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.",
						},
					},
					Validators: []validator.Object{validate.AwsKmsConfig()},
				},
			},
			"azure_key_vault_config": schema.ListNestedBlock{
				MarkdownDescription: "Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV).",
				Validators:          []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
						},
						"client_id": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Unique 36-hexadecimal character string that identifies an Azure application associated with your Azure Active Directory tenant.",
						},
						"azure_environment": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Azure environment in which your account credentials reside.",
						},
						"subscription_id": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Unique 36-hexadecimal character string that identifies your Azure subscription.",
						},
						"resource_group_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Name of the Azure resource group that contains your Azure Key Vault.",
						},
						"key_vault_name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Unique string that identifies the Azure Key Vault that contains your key.",
						},
						"key_identifier": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Web address with a unique key that identifies for your Azure Key Vault.",
						},
						"secret": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Private data that you need secured and that belongs to the specified Azure Key Vault (AKV) tenant (**azureKeyVault.tenantID**). This data can include any type of sensitive data such as passwords, database connection strings, API keys, and the like. AKV stores this information as encrypted binary data.",
						},
						"tenant_id": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Unique 36-hexadecimal character string that identifies the Azure Active Directory tenant within your Azure subscription.",
						},
						"require_private_networking": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: "Enable connection to your Azure Key Vault over private networking.",
						},
						"valid": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag that indicates whether the Azure encryption key can encrypt and decrypt data.",
						},
					},
				},
			},
			"google_cloud_kms_config": schema.ListNestedBlock{
				MarkdownDescription: "Details that define the configuration of Encryption at Rest using Google Cloud Key Management Service (KMS).",
				Validators:          []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
							MarkdownDescription: "Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.",
						},
						"service_account_key": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "JavaScript Object Notation (JSON) object that contains the Google Cloud Key Management Service (KMS). Format the JSON as a string and not as an object.",
						},
						"key_version_resource_id": schema.StringAttribute{
							Optional:            true,
							Sensitive:           true,
							MarkdownDescription: "Resource path that displays the key version resource ID for your Google Cloud KMS.",
						},
						"valid": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag that indicates whether the Google Cloud Key Management Service (KMS) encryption key can encrypt and decrypt data.",
						},
					},
				},
			},
		},
	}
}

type TfEncryptionAtRestRSModel struct {
	ID                   types.String                 `tfsdk:"id"`
	ProjectID            types.String                 `tfsdk:"project_id"`
	AwsKmsConfig         []TFAwsKmsConfigModel        `tfsdk:"aws_kms_config"`
	AzureKeyVaultConfig  []TFAzureKeyVaultConfigModel `tfsdk:"azure_key_vault_config"`
	GoogleCloudKmsConfig []TFGcpKmsConfigModel        `tfsdk:"google_cloud_kms_config"`
}

type TFAwsKmsConfigModel struct {
	AccessKeyID         types.String `tfsdk:"access_key_id"`
	SecretAccessKey     types.String `tfsdk:"secret_access_key"`
	CustomerMasterKeyID types.String `tfsdk:"customer_master_key_id"`
	Region              types.String `tfsdk:"region"`
	RoleID              types.String `tfsdk:"role_id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	Valid               types.Bool   `tfsdk:"valid"`
}
type TFAzureKeyVaultConfigModel struct {
	ClientID                 types.String `tfsdk:"client_id"`
	AzureEnvironment         types.String `tfsdk:"azure_environment"`
	SubscriptionID           types.String `tfsdk:"subscription_id"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	KeyVaultName             types.String `tfsdk:"key_vault_name"`
	KeyIdentifier            types.String `tfsdk:"key_identifier"`
	Secret                   types.String `tfsdk:"secret"`
	TenantID                 types.String `tfsdk:"tenant_id"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	RequirePrivateNetworking types.Bool   `tfsdk:"require_private_networking"`
	Valid                    types.Bool   `tfsdk:"valid"`
}
type TFGcpKmsConfigModel struct {
	ServiceAccountKey    types.String `tfsdk:"service_account_key"`
	KeyVersionResourceID types.String `tfsdk:"key_version_resource_id"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	Valid                types.Bool   `tfsdk:"valid"`
}
