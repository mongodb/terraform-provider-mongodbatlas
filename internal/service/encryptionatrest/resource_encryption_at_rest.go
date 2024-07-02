package encryptionatrest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

const (
	encryptionAtRestResourceName = "encryption_at_rest"
	errorCreateEncryptionAtRest  = "error creating Encryption At Rest: %s"
	errorReadEncryptionAtRest    = "error getting Encryption At Rest: %s"
	errorDeleteEncryptionAtRest  = "error deleting Encryption At Rest: (%s): %s"
	errorUpdateEncryptionAtRest  = "error updating Encryption At Rest: %s"
)

var _ resource.ResourceWithConfigure = &encryptionAtRestRS{}
var _ resource.ResourceWithImportState = &encryptionAtRestRS{}

func Resource() resource.Resource {
	return &encryptionAtRestRS{
		RSCommon: config.RSCommon{
			ResourceName: encryptionAtRestResourceName,
		},
	}
}

type encryptionAtRestRS struct {
	config.RSCommon
}

type TfEncryptionAtRestRSModel struct {
	ID                   types.String                 `tfsdk:"id"`
	ProjectID            types.String                 `tfsdk:"project_id"`
	AwsKmsConfig         []TfAwsKmsConfigModel        `tfsdk:"aws_kms_config"`
	AzureKeyVaultConfig  []TfAzureKeyVaultConfigModel `tfsdk:"azure_key_vault_config"`
	GoogleCloudKmsConfig []TfGcpKmsConfigModel        `tfsdk:"google_cloud_kms_config"`
}

type TfAwsKmsConfigModel struct {
	AccessKeyID         types.String `tfsdk:"access_key_id"`
	SecretAccessKey     types.String `tfsdk:"secret_access_key"`
	CustomerMasterKeyID types.String `tfsdk:"customer_master_key_id"`
	Region              types.String `tfsdk:"region"`
	RoleID              types.String `tfsdk:"role_id"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}
type TfAzureKeyVaultConfigModel struct {
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
type TfGcpKmsConfigModel struct {
	ServiceAccountKey    types.String `tfsdk:"service_account_key"`
	KeyVersionResourceID types.String `tfsdk:"key_version_resource_id"`
	Enabled              types.Bool   `tfsdk:"enabled"`
}

func (r *encryptionAtRestRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
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
					Validators: []validator.Object{validate.AwsKmsConfig()},
				},
			},
			"azure_key_vault_config": schema.ListNestedBlock{
				Validators: []validator.List{listvalidator.SizeAtMost(1)},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
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
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
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

func (r *encryptionAtRestRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var encryptionAtRestPlan *TfEncryptionAtRestRSModel
	var encryptionAtRestConfig *TfEncryptionAtRestRSModel
	connV2 := r.Client.AtlasV2

	resp.Diagnostics.Append(req.Plan.Get(ctx, &encryptionAtRestPlan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &encryptionAtRestConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := encryptionAtRestPlan.ProjectID.ValueString()
	encryptionAtRestReq := &admin.EncryptionAtRest{}
	if encryptionAtRestPlan.AwsKmsConfig != nil {
		encryptionAtRestReq.AwsKms = NewAtlasAwsKms(encryptionAtRestPlan.AwsKmsConfig)
	}
	if encryptionAtRestPlan.AzureKeyVaultConfig != nil {
		encryptionAtRestReq.AzureKeyVault = NewAtlasAzureKeyVault(encryptionAtRestPlan.AzureKeyVaultConfig)
	}
	if encryptionAtRestPlan.GoogleCloudKmsConfig != nil {
		encryptionAtRestReq.GoogleCloudKms = NewAtlasGcpKms(encryptionAtRestPlan.GoogleCloudKmsConfig)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{retrystrategy.RetryStrategyPendingState},
		Target:     []string{retrystrategy.RetryStrategyCompletedState, retrystrategy.RetryStrategyErrorState},
		Refresh:    ResourceMongoDBAtlasEncryptionAtRestCreateRefreshFunc(ctx, projectID, connV2.EncryptionAtRestUsingCustomerKeyManagementApi, encryptionAtRestReq),
		Timeout:    1 * time.Minute,
		MinTimeout: 1 * time.Second,
		Delay:      0,
	}

	var encryptionResp any
	var err error
	if encryptionResp, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorCreateEncryptionAtRest, projectID), err.Error())
		return
	}

	encryptionAtRestPlanNew := NewTfEncryptionAtRestRSModel(ctx, projectID, encryptionResp.(*admin.EncryptionAtRest))
	resetDefaultsFromConfigOrState(ctx, encryptionAtRestPlan, encryptionAtRestPlanNew, encryptionAtRestConfig)

	// set state to fully populated data
	diags := resp.State.Set(ctx, encryptionAtRestPlanNew)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ResourceMongoDBAtlasEncryptionAtRestCreateRefreshFunc(ctx context.Context, projectID string, client admin.EncryptionAtRestUsingCustomerKeyManagementApi, encryptionAtRestReq *admin.EncryptionAtRest) retry.StateRefreshFunc {
	return func() (any, string, error) {
		encryptionResp, _, err := client.UpdateEncryptionAtRest(ctx, projectID, encryptionAtRestReq).Execute()
		if err != nil {
			if errors.Is(err, errors.New("CANNOT_ASSUME_ROLE")) ||
				errors.Is(err, errors.New("INVALID_AWS_CREDENTIALS")) ||
				errors.Is(err, errors.New("CLOUD_PROVIDER_ACCESS_ROLE_NOT_AUTHORIZED")) {
				log.Printf("warning issue performing authorize EncryptionsAtRest not done try again: %s \n", err.Error())
				log.Println("retrying ")

				return encryptionResp, retrystrategy.RetryStrategyPendingState, nil
			}
			return encryptionResp, retrystrategy.RetryStrategyErrorState, err
		}
		return encryptionResp, retrystrategy.RetryStrategyCompletedState, nil
	}
}

func (r *encryptionAtRestRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var encryptionAtRestState TfEncryptionAtRestRSModel
	var isImport bool

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &encryptionAtRestState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := encryptionAtRestState.ProjectID.ValueString()

	// Use the ID only with the IMPORT operation
	if encryptionAtRestState.ID.ValueString() != "" && (projectID == "") {
		projectID = encryptionAtRestState.ID.ValueString()
		isImport = true
	}

	connV2 := r.Client.AtlasV2

	encryptionResp, getResp, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), projectID).Execute()
	if err != nil {
		if getResp != nil && getResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting encryption at rest resource during read", fmt.Sprintf(errorReadEncryptionAtRest, err.Error()))
		return
	}

	encryptionAtRestStateNew := NewTfEncryptionAtRestRSModel(ctx, projectID, encryptionResp)
	if isImport {
		setEmptyArrayForEmptyBlocksReturnedFromImport(encryptionAtRestStateNew)
	} else {
		resetDefaultsFromConfigOrState(ctx, &encryptionAtRestState, encryptionAtRestStateNew, nil)
	}

	// save read data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &encryptionAtRestStateNew)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *encryptionAtRestRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var encryptionAtRestState *TfEncryptionAtRestRSModel
	var encryptionAtRestConfig *TfEncryptionAtRestRSModel
	var encryptionAtRestPlan *TfEncryptionAtRestRSModel
	connV2 := r.Client.AtlasV2

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
	atlasEncryptionAtRest, atlasResp, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), projectID).Execute()
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting encryption at rest resource during update", fmt.Sprintf(project.ErrorProjectRead, projectID, err.Error()))
		return
	}

	if hasAwsKmsConfigChanged(encryptionAtRestPlan.AwsKmsConfig, encryptionAtRestState.AwsKmsConfig) {
		atlasEncryptionAtRest.AwsKms = NewAtlasAwsKms(encryptionAtRestPlan.AwsKmsConfig)
	}
	if hasAzureKeyVaultConfigChanged(encryptionAtRestPlan.AzureKeyVaultConfig, encryptionAtRestState.AzureKeyVaultConfig) {
		atlasEncryptionAtRest.AzureKeyVault = NewAtlasAzureKeyVault(encryptionAtRestPlan.AzureKeyVaultConfig)
	}
	if hasGcpKmsConfigChanged(encryptionAtRestPlan.GoogleCloudKmsConfig, encryptionAtRestState.GoogleCloudKmsConfig) {
		atlasEncryptionAtRest.GoogleCloudKms = NewAtlasGcpKms(encryptionAtRestPlan.GoogleCloudKmsConfig)
	}

	encryptionResp, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.UpdateEncryptionAtRest(ctx, projectID, atlasEncryptionAtRest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating encryption at rest", fmt.Sprintf(errorUpdateEncryptionAtRest, err.Error()))
		return
	}

	encryptionAtRestStateNew := NewTfEncryptionAtRestRSModel(ctx, projectID, encryptionResp)
	resetDefaultsFromConfigOrState(ctx, encryptionAtRestState, encryptionAtRestStateNew, encryptionAtRestConfig)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &encryptionAtRestStateNew)...)
}

func (r *encryptionAtRestRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var encryptionAtRestState *TfEncryptionAtRestRSModel

	// read prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &encryptionAtRestState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabled := false
	connV2 := r.Client.AtlasV2
	projectID := encryptionAtRestState.ProjectID.ValueString()

	_, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), projectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error when destroying resource", fmt.Sprintf(errorDeleteEncryptionAtRest, projectID, err.Error()))
		return
	}

	softDelete := admin.EncryptionAtRest{
		AwsKms:         &admin.AWSKMSConfiguration{Enabled: &enabled},
		AzureKeyVault:  &admin.AzureKeyVault{Enabled: &enabled},
		GoogleCloudKms: &admin.GoogleCloudKMS{Enabled: &enabled},
	}
	_, _, err = connV2.EncryptionAtRestUsingCustomerKeyManagementApi.UpdateEncryptionAtRest(ctx, projectID, &softDelete).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error when destroying resource", fmt.Sprintf(errorDeleteEncryptionAtRest, projectID, err.Error()))
		return
	}
}

func (r *encryptionAtRestRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func hasGcpKmsConfigChanged(gcpKmsConfigsPlan, gcpKmsConfigsState []TfGcpKmsConfigModel) bool {
	return !reflect.DeepEqual(gcpKmsConfigsPlan, gcpKmsConfigsState)
}

func hasAzureKeyVaultConfigChanged(azureKeyVaultConfigPlan, azureKeyVaultConfigState []TfAzureKeyVaultConfigModel) bool {
	return !reflect.DeepEqual(azureKeyVaultConfigPlan, azureKeyVaultConfigState)
}

func hasAwsKmsConfigChanged(awsKmsConfigPlan, awsKmsConfigState []TfAwsKmsConfigModel) bool {
	return !reflect.DeepEqual(awsKmsConfigPlan, awsKmsConfigState)
}

// resetDefaultsFromConfigOrState resets certain values that are not returned by the Atlas APIs from the Config
// However, during Read() and ImportState() since there is no access to the Config object, we use the State/Plan
// to achieve the same and encryptionAtRestRSConfig in that case is passed as nil in the calling method.
//
// encryptionAtRestRSCurrent - current State/Plan for this resource
// encryptionAtRestRSNew - final object that will be written in the State once the CRUD operation succeeds
// encryptionAtRestRSConfig - Config object for this resource
func resetDefaultsFromConfigOrState(ctx context.Context, encryptionAtRestRSCurrent, encryptionAtRestRSNew, encryptionAtRestRSConfig *TfEncryptionAtRestRSModel) {
	HandleAwsKmsConfigDefaults(ctx, encryptionAtRestRSCurrent, encryptionAtRestRSNew, encryptionAtRestRSConfig)
	HandleAzureKeyVaultConfigDefaults(ctx, encryptionAtRestRSCurrent, encryptionAtRestRSNew, encryptionAtRestRSConfig)
	HandleGcpKmsConfig(ctx, encryptionAtRestRSCurrent, encryptionAtRestRSNew, encryptionAtRestRSConfig)
}

func HandleGcpKmsConfig(ctx context.Context, earRSCurrent, earRSNew, earRSConfig *TfEncryptionAtRestRSModel) {
	// this is required to avoid unnecessary change detection during plan after migration to Plugin Framework if user didn't set this block
	if earRSCurrent.GoogleCloudKmsConfig == nil {
		earRSNew.GoogleCloudKmsConfig = []TfGcpKmsConfigModel{}
		return
	}

	// handling sensitive values that are not returned in the API response, so we sync them from the config
	// that user provided. encryptionAtRestRSConfig is nil during Read(), so we use the current plan
	if earRSConfig != nil && len(earRSConfig.GoogleCloudKmsConfig) > 0 {
		earRSNew.GoogleCloudKmsConfig[0].ServiceAccountKey = earRSConfig.GoogleCloudKmsConfig[0].ServiceAccountKey
	} else {
		earRSNew.GoogleCloudKmsConfig[0].ServiceAccountKey = earRSCurrent.GoogleCloudKmsConfig[0].ServiceAccountKey
	}
}

func HandleAwsKmsConfigDefaults(ctx context.Context, currentStateFile, newStateFile, earRSConfig *TfEncryptionAtRestRSModel) {
	// this is required to avoid unnecessary change detection during plan after migration to Plugin Framework if user didn't set this block
	if currentStateFile.AwsKmsConfig == nil {
		newStateFile.AwsKmsConfig = []TfAwsKmsConfigModel{}
		return
	}

	// handling sensitive values that are not returned in the API response, so we sync them from the config
	// that user provided. encryptionAtRestRSConfig is nil during Read(), so we use the current plan
	if earRSConfig != nil && len(earRSConfig.AwsKmsConfig) > 0 {
		newStateFile.AwsKmsConfig[0].Region = earRSConfig.AwsKmsConfig[0].Region
	} else {
		newStateFile.AwsKmsConfig[0].Region = currentStateFile.AwsKmsConfig[0].Region
	}

	// Secret access key is not returned by the API response
	if len(currentStateFile.AwsKmsConfig) == 1 && conversion.IsStringPresent(currentStateFile.AwsKmsConfig[0].SecretAccessKey.ValueStringPointer()) {
		newStateFile.AwsKmsConfig[0].SecretAccessKey = currentStateFile.AwsKmsConfig[0].SecretAccessKey
	}
}

func HandleAzureKeyVaultConfigDefaults(ctx context.Context, earRSCurrent, earRSNew, earRSConfig *TfEncryptionAtRestRSModel) {
	// this is required to avoid unnecessary change detection during plan after migration to Plugin Framework if user didn't set this block
	if earRSCurrent.AzureKeyVaultConfig == nil {
		earRSNew.AzureKeyVaultConfig = []TfAzureKeyVaultConfigModel{}
		return
	}

	// handling sensitive values that are not returned in the API response, so we sync them from the config
	// that user provided. encryptionAtRestRSConfig is nil during Read(), so we use the current plan
	if earRSConfig != nil && len(earRSConfig.AzureKeyVaultConfig) > 0 {
		earRSNew.AzureKeyVaultConfig[0].Secret = earRSConfig.AzureKeyVaultConfig[0].Secret
	} else {
		earRSNew.AzureKeyVaultConfig[0].Secret = earRSCurrent.AzureKeyVaultConfig[0].Secret
	}
}

// setEmptyArrayForEmptyBlocksReturnedFromImport sets the blocks AwsKmsConfig, GoogleCloudKmsConfig, TfAzureKeyVaultConfigModel
// to an empty array to avoid unnecessary change detection during plan after migration to Plugin Framework.
// Why:
// - the API returns the block AwsKmsConfig{enable=false} when the user does not provide the AWS KMS.
// - the API returns the block GoogleCloudKmsConfig{enable=false} if the user does not provider Google KMS
// - the API returns the block TfAzureKeyVaultConfigModel{enable=false} if the user does not provider AZURE KMS
func setEmptyArrayForEmptyBlocksReturnedFromImport(newStateFromImport *TfEncryptionAtRestRSModel) {
	if len(newStateFromImport.AwsKmsConfig) == 1 && !newStateFromImport.AwsKmsConfig[0].Enabled.ValueBool() {
		newStateFromImport.AwsKmsConfig = []TfAwsKmsConfigModel{}
	}

	if len(newStateFromImport.GoogleCloudKmsConfig) == 1 && !newStateFromImport.GoogleCloudKmsConfig[0].Enabled.ValueBool() {
		newStateFromImport.GoogleCloudKmsConfig = []TfGcpKmsConfigModel{}
	}

	if len(newStateFromImport.AzureKeyVaultConfig) == 1 && !newStateFromImport.AzureKeyVaultConfig[0].Enabled.ValueBool() {
		newStateFromImport.AzureKeyVaultConfig = []TfAzureKeyVaultConfigModel{}
	}
}
