package streamconnection

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const streamConnectionName = "stream_connection"

// Connection type constants used to differentiate API field mapping by connection type.
const (
	ConnectionTypeAWSKinesisDataStreams = "AWSKinesisDataStreams"
	ConnectionTypeAWSLambda             = "AWSLambda"
	ConnectionTypeAzureBlobStorage      = "AzureBlobStorage"
	ConnectionTypeGCPPubSub             = "GCPPubSub"
	ConnectionTypeCluster               = "Cluster"
	ConnectionTypeHTTPS                 = "Https"
	ConnectionTypeKafka                 = "Kafka"
	ConnectionTypeS3                    = "S3"
	ConnectionTypeSample                = "Sample"
	ConnectionTypeSchemaRegistry        = "SchemaRegistry"
)

var _ resource.ResourceWithConfigure = &streamConnectionRS{}
var _ resource.ResourceWithImportState = &streamConnectionRS{}

func Resource() resource.Resource {
	return &streamConnectionRS{
		RSCommon: config.RSCommon{
			ResourceName: streamConnectionName,
		},
	}
}

type streamConnectionRS struct {
	config.RSCommon
}

// TFStreamConnectionCommonModel contains common fields shared between resource and data source models
type TFStreamConnectionCommonModel struct {
	ID               types.String `tfsdk:"id"`
	ProjectID        types.String `tfsdk:"project_id"`
	WorkspaceName    types.String `tfsdk:"workspace_name"`
	InstanceName     types.String `tfsdk:"instance_name"`
	ConnectionName   types.String `tfsdk:"connection_name"`
	Type             types.String `tfsdk:"type"`
	ClusterName      types.String `tfsdk:"cluster_name"`
	ClusterProjectID types.String `tfsdk:"cluster_project_id"`
	Authentication   types.Object `tfsdk:"authentication"`
	BootstrapServers types.String `tfsdk:"bootstrap_servers"`
	Config           types.Map    `tfsdk:"config"`
	Security         types.Object `tfsdk:"security"`
	DBRoleToExecute  types.Object `tfsdk:"db_role_to_execute"`
	Networking       types.Object `tfsdk:"networking"`
	AWS              types.Object `tfsdk:"aws"`
	GCP              types.Object `tfsdk:"gcp"`
	Azure            types.Object `tfsdk:"azure"`
	// https connection
	Headers types.Map    `tfsdk:"headers"`
	URL     types.String `tfsdk:"url"`
	// SchemaRegistry connection
	SchemaRegistryProvider       types.String `tfsdk:"schema_registry_provider"`
	SchemaRegistryURLs           types.List   `tfsdk:"schema_registry_urls"`
	SchemaRegistryAuthentication types.Object `tfsdk:"schema_registry_authentication"`
	FailoverConnections          types.List   `tfsdk:"failover_connections"`
}

// TFFailoverConnectionModel represents one failover connection configuration within a stream connection.
type TFFailoverConnectionModel struct {
	ID                           types.String `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	Type                         types.String `tfsdk:"type"`
	Region                       types.String `tfsdk:"region"`
	ClusterName                  types.String `tfsdk:"cluster_name"`
	ClusterProjectID             types.String `tfsdk:"cluster_project_id"`
	DBRoleToExecute              types.Object `tfsdk:"db_role_to_execute"`
	BootstrapServers             types.String `tfsdk:"bootstrap_servers"`
	Authentication               types.Object `tfsdk:"authentication"`
	Config                       types.Map    `tfsdk:"config"`
	Security                     types.Object `tfsdk:"security"`
	Networking                   types.Object `tfsdk:"networking"`
	AWS                          types.Object `tfsdk:"aws"`
	Azure                        types.Object `tfsdk:"azure"`
	GCP                          types.Object `tfsdk:"gcp"`
	URL                          types.String `tfsdk:"url"`
	Headers                      types.Map    `tfsdk:"headers"`
	SchemaRegistryProvider       types.String `tfsdk:"schema_registry_provider"`
	SchemaRegistryURLs           types.List   `tfsdk:"schema_registry_urls"`
	SchemaRegistryAuthentication types.Object `tfsdk:"schema_registry_authentication"`
}

var FailoverConnectionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":                             types.StringType,
	"name":                           types.StringType,
	"type":                           types.StringType,
	"region":                         types.StringType,
	"cluster_name":                   types.StringType,
	"cluster_project_id":             types.StringType,
	"db_role_to_execute":             DBRoleToExecuteObjectType,
	"bootstrap_servers":              types.StringType,
	"authentication":                 ConnectionAuthenticationObjectType,
	"config":                         types.MapType{ElemType: types.StringType},
	"security":                       ConnectionSecurityObjectType,
	"networking":                     NetworkingObjectType,
	"aws":                            AWSObjectType,
	"azure":                          AzureObjectType,
	"gcp":                            GCPObjectType,
	"url":                            types.StringType,
	"headers":                        types.MapType{ElemType: types.StringType},
	"schema_registry_provider":       types.StringType,
	"schema_registry_urls":           types.ListType{ElemType: types.StringType},
	"schema_registry_authentication": SchemaRegistryAuthenticationObjectType,
}}

type TFStreamConnectionModel struct {
	TFStreamConnectionCommonModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// TFStreamConnectionDSModel is the data source model without timeouts (data sources don't support timeouts)
type TFStreamConnectionDSModel struct {
	TFStreamConnectionCommonModel
}

// ToDS converts the resource model to a data source model (without timeouts)
func (m *TFStreamConnectionModel) ToDS() *TFStreamConnectionDSModel {
	return &TFStreamConnectionDSModel{
		TFStreamConnectionCommonModel: m.TFStreamConnectionCommonModel,
	}
}

type TFConnectionAuthenticationModel struct {
	Mechanism                 types.String `tfsdk:"mechanism"`
	Method                    types.String `tfsdk:"method"`
	Password                  types.String `tfsdk:"password"`
	Username                  types.String `tfsdk:"username"`
	TokenEndpointURL          types.String `tfsdk:"token_endpoint_url"`
	ClientID                  types.String `tfsdk:"client_id"`
	ClientSecret              types.String `tfsdk:"client_secret"`
	Scope                     types.String `tfsdk:"scope"`
	SaslOauthbearerExtensions types.String `tfsdk:"sasl_oauthbearer_extensions"`
}

var ConnectionAuthenticationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"mechanism":                   types.StringType,
	"method":                      types.StringType,
	"password":                    types.StringType,
	"username":                    types.StringType,
	"token_endpoint_url":          types.StringType,
	"client_id":                   types.StringType,
	"client_secret":               types.StringType,
	"scope":                       types.StringType,
	"sasl_oauthbearer_extensions": types.StringType,
}}

type TFSchemaRegistryAuthenticationModel struct {
	Type     types.String `tfsdk:"type"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

var SchemaRegistryAuthenticationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":     types.StringType,
	"username": types.StringType,
	"password": types.StringType,
}}

type TFConnectionSecurityModel struct {
	BrokerPublicCertificate types.String `tfsdk:"broker_public_certificate"`
	Protocol                types.String `tfsdk:"protocol"`
}

var ConnectionSecurityObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"broker_public_certificate": types.StringType,
	"protocol":                  types.StringType,
}}

type TFDbRoleToExecuteModel struct {
	Role types.String `tfsdk:"role"`
	Type types.String `tfsdk:"type"`
}

var DBRoleToExecuteObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"role": types.StringType,
	"type": types.StringType,
}}

type TFNetworkingAccessModel struct {
	Type         types.String `tfsdk:"type"`
	ConnectionID types.String `tfsdk:"connection_id"`
}

var NetworkingAccessObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":          types.StringType,
	"connection_id": types.StringType,
}}

type TFNetworkingModel struct {
	Access types.Object `tfsdk:"access"`
}

var NetworkingObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"access": NetworkingAccessObjectType,
}}

type TFAWSModel struct {
	RoleArn types.String `tfsdk:"role_arn"`
}

var AWSObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"role_arn": types.StringType,
}}

type TFGCPModel struct {
	ServiceAccountID types.String `tfsdk:"service_account_id"`
}

var GCPObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"service_account_id": types.StringType,
}}

type TFAzureModel struct {
	ServicePrincipalID types.String `tfsdk:"service_principal_id"`
	StorageAccountName types.String `tfsdk:"storage_account_name"`
	Region             types.String `tfsdk:"region"`
}

var AzureObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"service_principal_id": types.StringType,
	"storage_account_name": types.StringType,
	"region":               types.StringType,
}}

func (r *streamConnectionRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

// getWorkspaceOrInstanceName returns the workspace name from workspace_name or instance_name field
func getWorkspaceOrInstanceName(model *TFStreamConnectionModel) string {
	if !model.WorkspaceName.IsNull() && !model.WorkspaceName.IsUnknown() {
		return model.WorkspaceName.ValueString()
	}
	if !model.InstanceName.IsNull() && !model.InstanceName.IsUnknown() {
		return model.InstanceName.ValueString()
	}
	return ""
}

func (r *streamConnectionRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var streamConnectionPlan TFStreamConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamConnectionPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionPlan.ProjectID.ValueString()
	workspaceOrInstanceName := getWorkspaceOrInstanceName(&streamConnectionPlan)
	if workspaceOrInstanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}

	streamConnectionReq, diags := NewStreamConnectionReq(ctx, &streamConnectionPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	apiResp, _, err := connV2.StreamsApi.CreateStreamConnection(ctx, projectID, workspaceOrInstanceName, streamConnectionReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error creating resource", err.Error())
		return
	}

	connectionName := conversion.SafeValue(apiResp.Name)

	createTimeout, diags := streamConnectionPlan.Timeouts.Create(ctx, DefaultConnectionTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// StateNotFound is a pending state for create - handles eventual consistency where
	// the resource may briefly return 404 after creation before becoming visible.
	apiResp, err = WaitStateTransition(ctx, projectID, workspaceOrInstanceName, connectionName, connV2.StreamsApi, createTimeout, []string{StatePending, StateNotFound}, []string{StateReady, StateFailed})
	if err != nil {
		resp.Diagnostics.AddError("error waiting for stream connection to be ready", err.Error())
		return
	}

	instanceName := streamConnectionPlan.InstanceName.ValueString()
	workspaceName := streamConnectionPlan.WorkspaceName.ValueString()

	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, workspaceName, &streamConnectionPlan.Authentication, &streamConnectionPlan.SchemaRegistryAuthentication, apiResp, &streamConnectionPlan.Timeouts)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Create failover connections if configured, building state directly from create responses
	// to avoid eventual-consistency issues with the list endpoint.
	newStreamConnectionModel.FailoverConnections = types.ListNull(FailoverConnectionObjectType)
	if !streamConnectionPlan.FailoverConnections.IsNull() && !streamConnectionPlan.FailoverConnections.IsUnknown() && len(streamConnectionPlan.FailoverConnections.Elements()) > 0 {
		var planFCs []TFFailoverConnectionModel
		resp.Diagnostics.Append(streamConnectionPlan.FailoverConnections.ElementsAs(ctx, &planFCs, false)...)
		if resp.Diagnostics.HasError() {
			// Save partial state so the primary connection is not orphaned.
			resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
			return
		}
		var fcModels []TFFailoverConnectionModel
		for i := range planFCs {
			fcItem := planFCs[i]
			conn, diags := newFailoverConnectionReq(ctx, &fcItem)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
				return
			}
			createdConn, _, err := connV2.StreamsApi.CreateFailoverConnection(ctx, projectID, workspaceOrInstanceName, connectionName, conn).Execute()
			if err != nil {
				resp.Diagnostics.AddError("error creating failover connection", err.Error())
				resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
				return
			}
			// Build state from the plan so sensitive fields (passwords) are always present.
			// Only inject the computed id from the API response.
			fcItem.ID = types.StringPointerValue(createdConn.Id)
			fcModels = append(fcModels, fcItem)
		}
		fcList, diags := types.ListValueFrom(ctx, FailoverConnectionObjectType, fcModels)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
			return
		}
		newStreamConnectionModel.FailoverConnections = fcList
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}

func (r *streamConnectionRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var streamConnectionState TFStreamConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamConnectionState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionState.ProjectID.ValueString()
	workspaceOrInstanceName := getWorkspaceOrInstanceName(&streamConnectionState)
	if workspaceOrInstanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}
	connectionName := streamConnectionState.ConnectionName.ValueString()
	apiResp, getResp, err := connV2.StreamsApi.GetStreamConnection(ctx, projectID, workspaceOrInstanceName, connectionName).Execute()
	if err != nil {
		if validate.StatusNotFound(getResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	instanceName := streamConnectionState.InstanceName.ValueString()
	workspaceName := streamConnectionState.WorkspaceName.ValueString()
	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, workspaceName, &streamConnectionState.Authentication, &streamConnectionState.SchemaRegistryAuthentication, apiResp, &streamConnectionState.Timeouts)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	fcList, diags := newTFFailoverConnectionsList(ctx, connV2.StreamsApi, projectID, workspaceOrInstanceName, connectionName, streamConnectionState.FailoverConnections, streamConnectionState.FailoverConnections)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newStreamConnectionModel.FailoverConnections = fcList

	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}

func (r *streamConnectionRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var streamConnectionPlan TFStreamConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &streamConnectionPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionPlan.ProjectID.ValueString()
	workspaceOrInstanceName := getWorkspaceOrInstanceName(&streamConnectionPlan)
	if workspaceOrInstanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}
	connectionName := streamConnectionPlan.ConnectionName.ValueString()
	streamConnectionReq, diags := NewStreamConnectionUpdateReq(ctx, &streamConnectionPlan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	_, _, err := connV2.StreamsApi.UpdateStreamConnection(ctx, projectID, workspaceOrInstanceName, connectionName, streamConnectionReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error updating resource", err.Error())
		return
	}

	updateTimeout, diags := streamConnectionPlan.Timeouts.Update(ctx, DefaultConnectionTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := WaitStateTransition(ctx, projectID, workspaceOrInstanceName, connectionName, connV2.StreamsApi, updateTimeout, []string{StatePending}, []string{StateReady, StateFailed})
	if err != nil {
		resp.Diagnostics.AddError("error waiting for stream connection to be ready", err.Error())
		return
	}

	instanceName := streamConnectionPlan.InstanceName.ValueString()
	workspaceName := streamConnectionPlan.WorkspaceName.ValueString()
	newStreamConnectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, workspaceName, &streamConnectionPlan.Authentication, &streamConnectionPlan.SchemaRegistryAuthentication, apiResp, &streamConnectionPlan.Timeouts)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Sync failover connections: create new ones, update changed ones, delete removed ones.
	var streamConnectionState TFStreamConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamConnectionState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if diags := syncFailoverConnections(ctx, connV2.StreamsApi, projectID, workspaceOrInstanceName, connectionName, streamConnectionPlan.FailoverConnections, streamConnectionState.FailoverConnections); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	fcList, diags := newTFFailoverConnectionsList(ctx, connV2.StreamsApi, projectID, workspaceOrInstanceName, connectionName, streamConnectionPlan.FailoverConnections, streamConnectionState.FailoverConnections)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	newStreamConnectionModel.FailoverConnections = fcList

	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}

func (r *streamConnectionRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var streamConnectionState *TFStreamConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &streamConnectionState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := r.Client.AtlasV2
	projectID := streamConnectionState.ProjectID.ValueString()
	instanceName := getWorkspaceOrInstanceName(streamConnectionState)
	if instanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}
	connectionName := streamConnectionState.ConnectionName.ValueString()

	// Delete failover connections first before removing the primary connection.
	if diags := deleteAllFailoverConnections(ctx, connV2.StreamsApi, projectID, instanceName, connectionName, streamConnectionState.FailoverConnections); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	deleteTimeout, diags := streamConnectionState.Timeouts.Delete(ctx, DefaultConnectionTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := DeleteStreamConnection(ctx, connV2.StreamsApi, projectID, instanceName, connectionName, deleteTimeout); err != nil {
		resp.Diagnostics.AddError("error deleting resource", err.Error())
		return
	}
}

func (r *streamConnectionRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	workspaceName, projectID, connectionName, err := splitStreamConnectionImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error splitting stream connection import ID", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_name"), workspaceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workspace_name"), workspaceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_name"), connectionName)...)
}

func splitStreamConnectionImportID(id string) (workspaceName, projectID, connectionName string, err error) {
	var re = regexp.MustCompile(`^(.*)-([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 4 {
		err = errors.New("use the format {workspace_name}-{project_id}-{connection_name}")
		return
	}

	workspaceName = parts[1]
	projectID = parts[2]
	connectionName = parts[3]
	return
}
