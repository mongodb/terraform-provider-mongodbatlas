package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	validators "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/validator"
)

const (
	encryptionAtRestResourceName = "encryption_at_rest"
	errorCreateEncryptionAtRest  = "error creating Encryption At Rest: %s"
	errorReadEncryptionAtRest    = "error getting Encryption At Rest: %s"
	errorDeleteEncryptionAtRest  = "error deleting Encryption At Rest: (%s): %s"
	errorUpdateEncryptionAtRest  = "error updating Encryption At Rest: %s"
)

var _ resource.Resource = &EncryptionAtRestRS{}
var _ resource.ResourceWithImportState = &EncryptionAtRestRS{}

func NewEncryptionAtRestRS() resource.Resource {
	return &EncryptionAtRestRS{}
}

type EncryptionAtRestRS struct {
	client *MongoDBClient
}

type tfEncryptionAtRestRSModel struct {
	ID                   types.String `tfsdk:"id"`
	ProjectID            types.String `tfsdk:"project_id"`
	AwsKmsConfig         types.List   `tfsdk:"aws_kms_config"`
	AzureKeyVaultConfig  types.List   `tfsdk:"azure_key_vault_config"`
	GoogleCloudKmsConfig types.List   `tfsdk:"google_cloud_kms_config"`
}

type tfAwsKmsConfigModel struct {
	AccessKeyID         types.String `tfsdk:"access_key_id"`
	SecretAccessKey     types.String `tfsdk:"secret_access_key"`
	CustomerMasterKeyID types.String `tfsdk:"customer_master_key_id"`
	Region              types.String `tfsdk:"region"`
	RoleID              types.String `tfsdk:"role_id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}
type tfAzureKeyVaultConfigModel struct {
	ClientID          types.String `tfsdk:"client_id"`
	AzureEnvironment  types.String `tfsdk:"azure_environment"`
	SubscriptionID    types.String `tfsdk:"subscription_id"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	KeyVaultName      types.String `tfsdk:"key_vault_name"`
	KeyIdentifier     types.String `tfsdk:"key_identifier"`
	Secret            types.String `tfsdk:"secret"`
	TenantID          types.String `tfsdk:"tenant_id"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}
type tfGcpKmsConfigModel struct {
	ServiceAccountKey    types.String `tfsdk:"service_account_key"`
	KeyVersionResourceID types.String `tfsdk:"key_version_resource_id"`
	Enabled              types.Bool   `tfsdk:"enabled"`
}

var tfAwsKmsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"enabled":                types.BoolType,
	"access_key_id":          types.StringType,
	"secret_access_key":      types.StringType,
	"customer_master_key_id": types.StringType,
	"region":                 types.StringType,
	"role_id":                types.StringType,
}}
var tfAzureKeyVaultObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"enabled":             types.BoolType,
	"client_id":           types.StringType,
	"azure_environment":   types.StringType,
	"subscription_id":     types.StringType,
	"resource_group_name": types.StringType,
	"key_vault_name":      types.StringType,
	"key_identifier":      types.StringType,
	"secret":              types.StringType,
	"tenant_id":           types.StringType,
}}
var tfGcpKmsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"enabled":                 types.BoolType,
	"service_account_key":     types.StringType,
	"key_version_resource_id": types.StringType,
}}

func (r *EncryptionAtRestRS) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, encryptionAtRestResourceName)
}

func (r *EncryptionAtRestRS) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, err := ConfigureClientInResource(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(errorConfigureSummary, err.Error())
		return
	}
	r.client = client
}

func (r *EncryptionAtRestRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
			},
		},
		Blocks: map[string]schema.Block{
			"aws_kms_config": schema.ListNestedBlock{
				Validators: []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
						},
						"access_key_id": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"secret_access_key": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"customer_master_key_id": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"region": schema.StringAttribute{
							Optional: true,
						},
						"role_id": schema.StringAttribute{
							Optional: true,
						},
					},
					Validators: []validator.Object{validators.AwsKmsConfig()},
				},
			},
			"azure_key_vault_config": schema.ListNestedBlock{
				Validators: []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Required: true,
						},
						"client_id": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"azure_environment": schema.StringAttribute{
							Optional: true,
						},
						"subscription_id": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"resource_group_name": schema.StringAttribute{
							Optional: true,
						},
						"key_vault_name": schema.StringAttribute{
							Optional: true,
						},
						"key_identifier": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"secret": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"tenant_id": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			"google_cloud_kms_config": schema.ListNestedBlock{
				Validators: []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
						},
						"service_account_key": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
						"key_version_resource_id": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func (r *EncryptionAtRestRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var encryptionAtRestPlan *tfEncryptionAtRestRSModel
	var encryptionAtRestConfig *tfEncryptionAtRestRSModel
	conn := r.client.Atlas

	resp.Diagnostics.Append(req.Plan.Get(ctx, &encryptionAtRestPlan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &encryptionAtRestConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := encryptionAtRestPlan.ProjectID.ValueString()
	encryptionAtRestReq := &matlas.EncryptionAtRest{
		GroupID: projectID,
	}

	if !encryptionAtRestPlan.AwsKmsConfig.IsNull() {
		encryptionAtRestReq.AwsKms = *toAtlasAwsKms(ctx, encryptionAtRestPlan.AwsKmsConfig)
	}
	if !encryptionAtRestPlan.AzureKeyVaultConfig.IsNull() {
		encryptionAtRestReq.AzureKeyVault = *toAtlasAzureKeyVault(ctx, encryptionAtRestPlan.AzureKeyVaultConfig)
	}
	if !encryptionAtRestPlan.GoogleCloudKmsConfig.IsNull() {
		encryptionAtRestReq.GoogleCloudKms = *toAtlasGcpKms(ctx, encryptionAtRestPlan.GoogleCloudKmsConfig)
	}

	for i := 0; i < 5; i++ {
		_, _, err := conn.EncryptionsAtRest.Create(ctx, encryptionAtRestReq)
		if err != nil {
			if strings.Contains(err.Error(), "CANNOT_ASSUME_ROLE") || strings.Contains(err.Error(), "INVALID_AWS_CREDENTIALS") ||
				strings.Contains(err.Error(), "CLOUD_PROVIDER_ACCESS_ROLE_NOT_AUTHORIZED") {
				log.Printf("warning issue performing authorize EncryptionsAtRest not done try again: %s \n", err.Error())
				log.Println("retrying ")
				time.Sleep(10 * time.Second)
				encryptionAtRestReq.GroupID = projectID
				continue
			}
		}
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf(errorCreateEncryptionAtRest, projectID), err.Error())
			return
		}
		break
	}

	// read
	encryptionResp, response, err := conn.EncryptionsAtRest.Get(context.Background(), projectID)
	tflog.Debug(ctx, fmt.Sprintf("encryptionResp from api: %v", encryptionResp))
	if err != nil {
		if resp != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting encryption at rest resource after create", fmt.Sprintf(errorReadEncryptionAtRest, err.Error()))
		return
	}

	encryptionAtRestPlanNew := toTFEncryptionAtRestRSModel(ctx, projectID, encryptionResp)
	resetDefaults(ctx, encryptionAtRestPlan, encryptionAtRestPlanNew, encryptionAtRestConfig)

	// set state to fully populated data
	diags := resp.State.Set(ctx, encryptionAtRestPlanNew)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func resetDefaults(ctx context.Context, encryptionAtRestRS *tfEncryptionAtRestRSModel, encryptionAtRestRSNew *tfEncryptionAtRestRSModel, encryptionAtRestRSConfig *tfEncryptionAtRestRSModel) {
	if encryptionAtRestRS.AwsKmsConfig.IsNull() {
		tfAwsKmsConfigs := make([]tfAwsKmsConfigModel, 0)
		encryptionAtRestRSNew.AwsKmsConfig, _ = types.ListValueFrom(ctx, tfAwsKmsObjectType, tfAwsKmsConfigs)
	} else {

		var awsKmsConfigsNew []tfAwsKmsConfigModel
		encryptionAtRestRSNew.AwsKmsConfig.ElementsAs(ctx, &awsKmsConfigsNew, false)

		// user may set "region" equal to 'US_EAST_1' or 'US-EAST-1' in config, we ensure to update new plan/state with value that is used in the config
		if encryptionAtRestRSConfig != nil {
			var awsKmsConfigs []tfAwsKmsConfigModel
			encryptionAtRestRSConfig.AwsKmsConfig.ElementsAs(ctx, &awsKmsConfigs, false)

			awsKmsConfigsNew[0].Region = awsKmsConfigs[0].Region
		} else {
			var awsKmsConfigs []tfAwsKmsConfigModel
			encryptionAtRestRS.AwsKmsConfig.ElementsAs(ctx, &awsKmsConfigs, false)

			awsKmsConfigsNew[0].Region = awsKmsConfigs[0].Region
		}
		encryptionAtRestRSNew.AwsKmsConfig, _ = types.ListValueFrom(ctx, tfAwsKmsObjectType, awsKmsConfigsNew)
	}
	if encryptionAtRestRS.AzureKeyVaultConfig.IsNull() {
		// encryptionAtRestPlanNew.AzureKeyVaultConfig = types.ListNull(tfAzureKeyVaultObjectType)
		tfAzKeyVaultConfigs := make([]tfAzureKeyVaultConfigModel, 0)
		encryptionAtRestRSNew.AzureKeyVaultConfig, _ = types.ListValueFrom(ctx, tfAzureKeyVaultObjectType, tfAzKeyVaultConfigs)
	} else {

		var azureConfigsNew []tfAzureKeyVaultConfigModel
		encryptionAtRestRSNew.AzureKeyVaultConfig.ElementsAs(ctx, &azureConfigsNew, false)

		if encryptionAtRestRSConfig != nil {
			var azureConfigs []tfAzureKeyVaultConfigModel
			encryptionAtRestRSConfig.AzureKeyVaultConfig.ElementsAs(ctx, &azureConfigs, false)

			azureConfigsNew[0].Secret = azureConfigs[0].Secret
		} else {
			var azureConfigs []tfAzureKeyVaultConfigModel
			encryptionAtRestRS.AzureKeyVaultConfig.ElementsAs(ctx, &azureConfigs, false)

			azureConfigsNew[0].Secret = azureConfigs[0].Secret
		}
		encryptionAtRestRSNew.AzureKeyVaultConfig, _ = types.ListValueFrom(ctx, tfAzureKeyVaultObjectType, azureConfigsNew)
	}

	if encryptionAtRestRS.GoogleCloudKmsConfig.IsNull() {
		// encryptionAtRestPlanNew.GoogleCloudKmsConfig = types.ListNull(tfGcpKmsObjectType)
		tfGcpKmsConfigs := make([]tfGcpKmsConfigModel, 0)
		encryptionAtRestRSNew.GoogleCloudKmsConfig, _ = types.ListValueFrom(ctx, tfGcpKmsObjectType, tfGcpKmsConfigs)

	}
}

func (r *EncryptionAtRestRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var encryptionAtRestState tfEncryptionAtRestRSModel

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &encryptionAtRestState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get resource from API
	conn := r.client.Atlas
	projectID := encryptionAtRestState.ProjectID.ValueString()
	encryptionResp, _, err := conn.EncryptionsAtRest.Get(context.Background(), projectID)
	tflog.Debug(ctx, fmt.Sprintf("encryptionResp from api: %v", encryptionResp))
	if err != nil {
		resp.Diagnostics.AddError("error when getting encryption at rest resource during read", fmt.Sprintf(errorReadEncryptionAtRest, err.Error()))
		return
	}

	encryptionAtRestStateNew := toTFEncryptionAtRestRSModel(ctx, projectID, encryptionResp)
	resetDefaults(ctx, &encryptionAtRestState, encryptionAtRestStateNew, nil)
	// resetDefaultsForRead(&encryptionAtRestState, encryptionAtRestStateNew)

	// save read data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &encryptionAtRestStateNew)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EncryptionAtRestRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var encryptionAtRestState *tfEncryptionAtRestRSModel
	var encryptionAtRestConfig *tfEncryptionAtRestRSModel
	var encryptionAtRestPlan *tfEncryptionAtRestRSModel
	conn := r.client.Atlas

	// get current config
	resp.Diagnostics.Append(req.Config.Get(ctx, &encryptionAtRestConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &encryptionAtRestState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// get current plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &encryptionAtRestPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := encryptionAtRestState.ProjectID.ValueString()
	atlasEncryptionAtRest, atlasResp, err := conn.EncryptionsAtRest.Get(context.Background(), projectID)
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting encryption at rest resource during update", fmt.Sprintf(errorProjectRead, projectID, err.Error()))
		return
	}

	if hasAwsKmsConfigChanged(encryptionAtRestPlan.AwsKmsConfig, encryptionAtRestState.AwsKmsConfig) {
		atlasEncryptionAtRest.AwsKms = *toAtlasAwsKms(ctx, encryptionAtRestPlan.AwsKmsConfig)
	}
	if hasAzureKeyVaultConfigChanged(encryptionAtRestPlan.AzureKeyVaultConfig, encryptionAtRestState.AzureKeyVaultConfig) {
		atlasEncryptionAtRest.AzureKeyVault = *toAtlasAzureKeyVault(ctx, encryptionAtRestPlan.AzureKeyVaultConfig)
	}
	if hasGcpKmsConfigChanged(encryptionAtRestPlan.GoogleCloudKmsConfig, encryptionAtRestState.GoogleCloudKmsConfig) {
		atlasEncryptionAtRest.GoogleCloudKms = *toAtlasGcpKms(ctx, encryptionAtRestPlan.GoogleCloudKmsConfig)
	}

	atlasEncryptionAtRest.GroupID = projectID
	_, _, err = conn.EncryptionsAtRest.Create(ctx, atlasEncryptionAtRest)
	if err != nil {
		resp.Diagnostics.AddError("error updating encryption at rest", fmt.Sprintf(errorUpdateEncryptionAtRest, err.Error()))
		return
	}

	// read
	encryptionResp, response, err := conn.EncryptionsAtRest.Get(context.Background(), projectID)
	tflog.Debug(ctx, fmt.Sprintf("encryptionResp from api: %v", encryptionResp))
	if err != nil {
		if resp != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting encryption at rest resource after update", fmt.Sprintf(errorReadEncryptionAtRest, err.Error()))
		return
	}

	encryptionAtRestStateNew := toTFEncryptionAtRestRSModel(ctx, projectID, encryptionResp)
	resetDefaults(ctx, encryptionAtRestState, encryptionAtRestStateNew, encryptionAtRestConfig)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &encryptionAtRestStateNew)...)
}

func hasGcpKmsConfigChanged(gcpKmsConfigsPlan, gcpKmsConfigsState basetypes.ListValue) bool {
	return !reflect.DeepEqual(gcpKmsConfigsPlan, gcpKmsConfigsState)
}

func hasAzureKeyVaultConfigChanged(azureKeyVaultConfigPlan, azureKeyVaultConfigState basetypes.ListValue) bool {
	return !reflect.DeepEqual(azureKeyVaultConfigPlan, azureKeyVaultConfigState)
}

func hasAwsKmsConfigChanged(awsKmsConfigPlan, awsKmsConfigState basetypes.ListValue) bool {
	return !reflect.DeepEqual(awsKmsConfigPlan, awsKmsConfigState)
}

func (r *EncryptionAtRestRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var encryptionAtRestState *tfEncryptionAtRestRSModel

	// read prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &encryptionAtRestState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn := r.client.Atlas
	projectID := encryptionAtRestState.ProjectID.ValueString()
	_, err := conn.EncryptionsAtRest.Delete(ctx, projectID)

	if err != nil {
		resp.Diagnostics.AddError("error when destroying resource", fmt.Sprintf(errorDeleteEncryptionAtRest, projectID, err.Error()))
		return
	}
}

func (r *EncryptionAtRestRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func toTFEncryptionAtRestRSModel(ctx context.Context, projectID string, encryptionResp *matlas.EncryptionAtRest) *tfEncryptionAtRestRSModel {
	encryptionAtRest := tfEncryptionAtRestRSModel{
		ID:                   types.StringValue(projectID),
		ProjectID:            types.StringValue(projectID),
		AwsKmsConfig:         toTFAwsKmsConfig(ctx, &encryptionResp.AwsKms),
		AzureKeyVaultConfig:  toTFAzureKeyVaultConfig(ctx, &encryptionResp.AzureKeyVault),
		GoogleCloudKmsConfig: toTFGcpKmsConfig(ctx, &encryptionResp.GoogleCloudKms),
	}
	return &encryptionAtRest
}

func toTFAwsKmsConfig(ctx context.Context, awsKms *matlas.AwsKms) types.List {
	tfAwsKmsConfigs := make([]tfAwsKmsConfigModel, 1)

	if awsKms != nil {
		tfAwsKmsConfigs[0].Enabled = types.BoolPointerValue(awsKms.Enabled)
		tfAwsKmsConfigs[0].CustomerMasterKeyID = types.StringValue(awsKms.CustomerMasterKeyID)
		tfAwsKmsConfigs[0].Region = types.StringValue(awsKms.Region)

		if accessKeyID := awsKms.AccessKeyID; accessKeyID == "" {
			tfAwsKmsConfigs[0].AccessKeyID = types.StringNull()
		} else {
			tfAwsKmsConfigs[0].AccessKeyID = types.StringValue(accessKeyID)
		}
		// tfAwsKmsConfigs[0].AccessKeyID = types.StringValue(awsKms.AccessKeyID)

		if secretAccessKey := awsKms.SecretAccessKey; secretAccessKey == "" {
			tfAwsKmsConfigs[0].SecretAccessKey = types.StringNull()
		} else {
			tfAwsKmsConfigs[0].SecretAccessKey = types.StringValue(secretAccessKey)
		}
		// tfAwsKmsConfigs[0].SecretAccessKey = types.StringValue(awsKms.SecretAccessKey)

		if roleID := awsKms.RoleID; roleID == "" {
			tfAwsKmsConfigs[0].RoleID = types.StringNull()
		} else {
			tfAwsKmsConfigs[0].RoleID = types.StringValue(roleID)
		}
		// tfAwsKmsConfigs[0].RoleID = types.StringValue(awsKms.RoleID)
	}

	list, _ := types.ListValueFrom(ctx, tfAwsKmsObjectType, tfAwsKmsConfigs)
	return list
}

func toTFAzureKeyVaultConfig(ctx context.Context, az *matlas.AzureKeyVault) types.List {
	tfAzKeyVaultConfigs := make([]tfAzureKeyVaultConfigModel, 1)

	tfAzKeyVaultConfigs[0].Enabled = types.BoolPointerValue(az.Enabled)
	tfAzKeyVaultConfigs[0].ClientID = types.StringValue(az.ClientID)
	tfAzKeyVaultConfigs[0].AzureEnvironment = types.StringValue(az.AzureEnvironment)
	tfAzKeyVaultConfigs[0].SubscriptionID = types.StringValue(az.SubscriptionID)
	tfAzKeyVaultConfigs[0].ResourceGroupName = types.StringValue(az.ResourceGroupName)
	tfAzKeyVaultConfigs[0].KeyVaultName = types.StringValue(az.KeyVaultName)
	tfAzKeyVaultConfigs[0].KeyIdentifier = types.StringValue(az.KeyIdentifier)
	tfAzKeyVaultConfigs[0].TenantID = types.StringValue(az.TenantID)

	if secret := az.Secret; secret == "" {
		tfAzKeyVaultConfigs[0].Secret = types.StringNull()
	} else {
		tfAzKeyVaultConfigs[0].Secret = types.StringValue(secret)
	}
	// tfAzKeyVaultConfigs[0].Secret = types.StringValue(az.Secret)

	list, _ := types.ListValueFrom(ctx, tfAzureKeyVaultObjectType, tfAzKeyVaultConfigs)
	return list
}

func toTFGcpKmsConfig(ctx context.Context, gcpKms *matlas.GoogleCloudKms) types.List {
	tfGcpKmsConfigs := make([]tfGcpKmsConfigModel, 1)

	tfGcpKmsConfigs[0].Enabled = types.BoolPointerValue(gcpKms.Enabled)
	// tfGcpKmsConfigs[0].ServiceAccountKey = types.StringValue(gcpKms.ServiceAccountKey)
	tfGcpKmsConfigs[0].KeyVersionResourceID = types.StringValue(gcpKms.KeyVersionResourceID)

	if serviceAccountKey := gcpKms.ServiceAccountKey; serviceAccountKey == "" {
		tfGcpKmsConfigs[0].ServiceAccountKey = types.StringNull()
	} else {
		tfGcpKmsConfigs[0].ServiceAccountKey = types.StringValue(serviceAccountKey)
	}

	list, _ := types.ListValueFrom(ctx, tfGcpKmsObjectType, tfGcpKmsConfigs)
	return list
}

func toAtlasAwsKms(ctx context.Context, tfAwsKmsConfigList basetypes.ListValue) *matlas.AwsKms {
	if len(tfAwsKmsConfigList.Elements()) == 0 {
		return &matlas.AwsKms{}
	}
	var awsKmsConfigs []tfAwsKmsConfigModel
	tfAwsKmsConfigList.ElementsAs(ctx, &awsKmsConfigs, false)

	awsRegion, _ := valRegion(awsKmsConfigs[0].Region.ValueString())

	return &matlas.AwsKms{
		Enabled:             awsKmsConfigs[0].Enabled.ValueBoolPointer(),
		AccessKeyID:         awsKmsConfigs[0].AccessKeyID.ValueString(),
		SecretAccessKey:     awsKmsConfigs[0].SecretAccessKey.ValueString(),
		CustomerMasterKeyID: awsKmsConfigs[0].CustomerMasterKeyID.ValueString(),
		Region:              awsRegion,
		RoleID:              awsKmsConfigs[0].RoleID.ValueString(),
	}
}

func toAtlasGcpKms(ctx context.Context, tfGcpKmsConfigList basetypes.ListValue) *matlas.GoogleCloudKms {
	if len(tfGcpKmsConfigList.Elements()) == 0 {
		return &matlas.GoogleCloudKms{}
	}
	var gcpKmsConfigs []tfGcpKmsConfigModel
	tfGcpKmsConfigList.ElementsAs(ctx, &gcpKmsConfigs, false)

	return &matlas.GoogleCloudKms{
		Enabled:              gcpKmsConfigs[0].Enabled.ValueBoolPointer(),
		ServiceAccountKey:    gcpKmsConfigs[0].ServiceAccountKey.ValueString(),
		KeyVersionResourceID: gcpKmsConfigs[0].KeyVersionResourceID.ValueString(),
	}
}

func toAtlasAzureKeyVault(ctx context.Context, tfAzureKeyVaultList basetypes.ListValue) *matlas.AzureKeyVault {
	if len(tfAzureKeyVaultList.Elements()) == 0 {
		return &matlas.AzureKeyVault{}
	}
	var azureKeyVaultConfigs []tfAzureKeyVaultConfigModel
	tfAzureKeyVaultList.ElementsAs(ctx, &azureKeyVaultConfigs, false)

	return &matlas.AzureKeyVault{
		Enabled:           azureKeyVaultConfigs[0].Enabled.ValueBoolPointer(),
		ClientID:          azureKeyVaultConfigs[0].ClientID.ValueString(),
		AzureEnvironment:  azureKeyVaultConfigs[0].AzureEnvironment.ValueString(),
		SubscriptionID:    azureKeyVaultConfigs[0].SubscriptionID.ValueString(),
		ResourceGroupName: azureKeyVaultConfigs[0].ResourceGroupName.ValueString(),
		KeyVaultName:      azureKeyVaultConfigs[0].KeyVaultName.ValueString(),
		KeyIdentifier:     azureKeyVaultConfigs[0].KeyIdentifier.ValueString(),
		Secret:            azureKeyVaultConfigs[0].Secret.ValueString(),
		TenantID:          azureKeyVaultConfigs[0].TenantID.ValueString(),
	}
}
