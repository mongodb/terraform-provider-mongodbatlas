package streamconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312021/admin"
)

func NewStreamConnectionReq(ctx context.Context, plan *TFStreamConnectionModel) (*admin.StreamsConnection, diag.Diagnostics) {
	streamConnection := admin.StreamsConnection{
		Name:             plan.ConnectionName.ValueStringPointer(),
		Type:             plan.Type.ValueStringPointer(),
		ClusterName:      plan.ClusterName.ValueStringPointer(),
		ClusterGroupId:   plan.ClusterProjectID.ValueStringPointer(),
		BootstrapServers: plan.BootstrapServers.ValueStringPointer(),
		Url:              plan.URL.ValueStringPointer(),
		Provider:         plan.SchemaRegistryProvider.ValueStringPointer(),
	}
	if !plan.Authentication.IsNull() {
		authenticationModel := &TFConnectionAuthenticationModel{}
		if diags := plan.Authentication.As(ctx, authenticationModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.Authentication = &admin.StreamsKafkaAuthentication{
			Mechanism:                 authenticationModel.Mechanism.ValueStringPointer(),
			Method:                    authenticationModel.Method.ValueStringPointer(),
			Password:                  authenticationModel.Password.ValueStringPointer(),
			Username:                  authenticationModel.Username.ValueStringPointer(),
			TokenEndpointUrl:          authenticationModel.TokenEndpointURL.ValueStringPointer(),
			ClientId:                  authenticationModel.ClientID.ValueStringPointer(),
			ClientSecret:              authenticationModel.ClientSecret.ValueStringPointer(),
			Scope:                     authenticationModel.Scope.ValueStringPointer(),
			SaslOauthbearerExtensions: authenticationModel.SaslOauthbearerExtensions.ValueStringPointer(),
		}
	}
	if !plan.Security.IsNull() {
		securityModel := &TFConnectionSecurityModel{}
		if diags := plan.Security.As(ctx, securityModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.Security = &admin.StreamsKafkaSecurity{
			BrokerPublicCertificate: securityModel.BrokerPublicCertificate.ValueStringPointer(),
			Protocol:                securityModel.Protocol.ValueStringPointer(),
		}
	}

	if !plan.Config.IsNull() {
		configMap := &map[string]string{}
		if diags := plan.Config.ElementsAs(ctx, configMap, true); diags.HasError() {
			return nil, diags
		}
		streamConnection.Config = configMap
	}

	if !plan.DBRoleToExecute.IsNull() {
		dbRoleToExecuteModel := &TFDbRoleToExecuteModel{}
		if diags := plan.DBRoleToExecute.As(ctx, dbRoleToExecuteModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.DbRoleToExecute = &admin.DBRoleToExecute{
			Role: dbRoleToExecuteModel.Role.ValueStringPointer(),
			Type: dbRoleToExecuteModel.Type.ValueStringPointer(),
		}
	}

	if !plan.Networking.IsNull() && !plan.Networking.IsUnknown() {
		networkingModel := &TFNetworkingModel{}
		if diags := plan.Networking.As(ctx, networkingModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		networkingAccessModel := &TFNetworkingAccessModel{}
		if diags := networkingModel.Access.As(ctx, networkingAccessModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		switch plan.Type.ValueString() {
		case ConnectionTypeAzureBlobStorage, ConnectionTypeGCPPubSub:
			streamConnection.PublicPrivateNetworking = &admin.StreamsPublicPrivateLinkNetworking{
				Access: &admin.StreamsPublicPrivateLinkNetworkingAccess{
					Type:         networkingAccessModel.Type.ValueStringPointer(),
					ConnectionId: networkingAccessModel.ConnectionID.ValueStringPointer(),
				},
			}
		case ConnectionTypeKafka, ConnectionTypeAWSKinesisDataStreams, ConnectionTypeS3:
			streamConnection.Networking = &admin.StreamsKafkaNetworking{
				Access: &admin.StreamsKafkaNetworkingAccess{
					Type:         networkingAccessModel.Type.ValueStringPointer(),
					ConnectionId: networkingAccessModel.ConnectionID.ValueStringPointer(),
				},
			}
		default:
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				"invalid connection type with networking",
				fmt.Sprintf("connection type %q does not support networking configuration", plan.Type.ValueString()),
			)}
		}
	}

	if !plan.AWS.IsNull() {
		awsModel := &TFAWSModel{}
		if diags := plan.AWS.As(ctx, awsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.Aws = &admin.StreamsAWSConnectionConfig{
			RoleArn: awsModel.RoleArn.ValueStringPointer(),
		}
	}

	if !plan.GCP.IsNull() {
		gcpModel := &TFGCPModel{}
		if diags := plan.GCP.As(ctx, gcpModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.Gcp = &admin.StreamsGCPConnectionConfig{
			ServiceAccountId: gcpModel.ServiceAccountID.ValueStringPointer(),
		}
	}

	if !plan.Azure.IsNull() {
		azureModel := &TFAzureModel{}
		if diags := plan.Azure.As(ctx, azureModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.Azure = &admin.AzureConnection{
			ServicePrincipalId: azureModel.ServicePrincipalID.ValueStringPointer(),
			StorageAccountName: azureModel.StorageAccountName.ValueStringPointer(),
			Region:             azureModel.Region.ValueStringPointer(),
		}
	}

	if !plan.Headers.IsNull() {
		headersMap := make(map[string]string)
		if diags := plan.Headers.ElementsAs(ctx, &headersMap, true); diags.HasError() {
			return nil, diags
		}
		streamConnection.Headers = &headersMap
	}

	// SchemaRegistry
	if !plan.SchemaRegistryURLs.IsNull() {
		var schemaRegistryURLs []string
		diags := plan.SchemaRegistryURLs.ElementsAs(ctx, &schemaRegistryURLs, false)
		if diags.HasError() {
			return nil, diags
		}
		streamConnection.SchemaRegistryUrls = &schemaRegistryURLs
	}

	if !plan.SchemaRegistryAuthentication.IsNull() {
		schemaRegistryAuthenticationModel := &TFSchemaRegistryAuthenticationModel{}
		if diags := plan.SchemaRegistryAuthentication.As(ctx, schemaRegistryAuthenticationModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.SchemaRegistryAuthentication = &admin.SchemaRegistryAuthentication{
			Type:     schemaRegistryAuthenticationModel.Type.ValueString(),
			Username: schemaRegistryAuthenticationModel.Username.ValueStringPointer(),
			Password: schemaRegistryAuthenticationModel.Password.ValueStringPointer(),
		}
	}

	return &streamConnection, nil
}

func NewStreamConnectionUpdateReq(ctx context.Context, plan *TFStreamConnectionModel) (*admin.StreamsConnection, diag.Diagnostics) {
	streamConnection, diags := NewStreamConnectionReq(ctx, plan)
	if diags.HasError() {
		return nil, diags
	}

	headersMap := make(map[string]string)
	// only set headers if the plan is not empty, otherwise the headers will be removed by sending an empty headers map to the PATCH endpoint
	if !plan.Headers.IsNull() && !plan.Headers.IsUnknown() {
		if diags := plan.Headers.ElementsAs(ctx, &headersMap, true); diags.HasError() {
			return nil, diags
		}
	}
	streamConnection.Headers = &headersMap
	return streamConnection, nil
}

// NewTFStreamConnection determines if the original model was created with instance_name or workspace_name and sets the appropriate field.
// The planTimeouts parameter is optional and used to preserve user-configured timeouts across create/update operations.
func NewTFStreamConnection(ctx context.Context, projID, instanceName, workspaceName string, currAuthConfig, currSchemaRegistryAuthConfig *types.Object, apiResp *admin.StreamsConnection, planTimeouts *timeouts.Value) (*TFStreamConnectionModel, diag.Diagnostics) {
	streamWorkspaceName := workspaceName
	if instanceName != "" {
		streamWorkspaceName = instanceName
	}

	rID := fmt.Sprintf("%s-%s-%s", streamWorkspaceName, projID, conversion.SafeValue(apiResp.Name))

	connectionModel := TFStreamConnectionModel{
		TFStreamConnectionCommonModel: TFStreamConnectionCommonModel{
			ID:                     types.StringValue(rID),
			ProjectID:              types.StringValue(projID),
			ConnectionName:         types.StringPointerValue(apiResp.Name),
			Type:                   types.StringPointerValue(apiResp.Type),
			ClusterName:            types.StringPointerValue(apiResp.ClusterName),
			ClusterProjectID:       types.StringPointerValue(apiResp.ClusterGroupId),
			BootstrapServers:       types.StringPointerValue(apiResp.BootstrapServers),
			URL:                    types.StringPointerValue(apiResp.Url),
			SchemaRegistryURLs:     types.ListNull(types.StringType),
			SchemaRegistryProvider: types.StringPointerValue(apiResp.Provider),
		},
	}

	// Preserve user-configured timeouts
	if planTimeouts != nil {
		connectionModel.Timeouts = *planTimeouts
	}

	// Set the appropriate field based on the original model
	if workspaceName != "" {
		connectionModel.WorkspaceName = types.StringValue(workspaceName)
		connectionModel.InstanceName = types.StringNull()
	} else {
		// Default to instance_name for backward compatibility
		connectionModel.InstanceName = types.StringValue(instanceName)
		connectionModel.WorkspaceName = types.StringNull()
	}

	authModel, diags := newTFConnectionAuthenticationModel(ctx, currAuthConfig, apiResp.Authentication)
	if diags.HasError() {
		return nil, diags
	}
	connectionModel.Authentication = *authModel

	connectionModel.Config = types.MapNull(types.StringType)
	if apiResp.Config != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, apiResp.Config)
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.Config = mapValue
	}

	if apiResp.SchemaRegistryUrls != nil {
		schemaRegistryURLs, diags := types.ListValueFrom(ctx, types.StringType, apiResp.SchemaRegistryUrls)
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.SchemaRegistryURLs = schemaRegistryURLs
	}

	schemaRegistryAuthModel, diags := newTFSchemaRegistryAuthentication(ctx, currSchemaRegistryAuthConfig, apiResp.SchemaRegistryAuthentication)
	if diags.HasError() {
		return nil, diags
	}
	connectionModel.SchemaRegistryAuthentication = *schemaRegistryAuthModel

	connectionModel.Security = types.ObjectNull(ConnectionSecurityObjectType.AttrTypes)
	if apiResp.Security != nil {
		securityModel, diags := types.ObjectValueFrom(ctx, ConnectionSecurityObjectType.AttrTypes, TFConnectionSecurityModel{
			BrokerPublicCertificate: types.StringPointerValue(apiResp.Security.BrokerPublicCertificate),
			Protocol:                types.StringPointerValue(apiResp.Security.Protocol),
		})
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.Security = securityModel
	}

	connectionModel.DBRoleToExecute = types.ObjectNull(DBRoleToExecuteObjectType.AttrTypes)
	if apiResp.DbRoleToExecute != nil {
		dbRoleToExecuteModel, diags := types.ObjectValueFrom(ctx, DBRoleToExecuteObjectType.AttrTypes, TFDbRoleToExecuteModel{
			Role: types.StringPointerValue(apiResp.DbRoleToExecute.Role),
			Type: types.StringPointerValue(apiResp.DbRoleToExecute.Type),
		})
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.DBRoleToExecute = dbRoleToExecuteModel
	}

	// The API returns networking in either Networking (Kafka, S3) or PublicPrivateNetworking (GCPPubSub, Azure) depending on connection type.
	connectionModel.Networking = types.ObjectNull(NetworkingObjectType.AttrTypes)
	var networkingAccessType, networkingConnectionID *string
	if apiResp.Networking != nil && apiResp.Networking.Access != nil {
		networkingAccessType = apiResp.Networking.Access.Type
		networkingConnectionID = apiResp.Networking.Access.ConnectionId
	} else if apiResp.PublicPrivateNetworking != nil && apiResp.PublicPrivateNetworking.Access != nil {
		networkingAccessType = apiResp.PublicPrivateNetworking.Access.Type
		networkingConnectionID = apiResp.PublicPrivateNetworking.Access.ConnectionId
	}
	if networkingAccessType != nil {
		networkingAccessModel, diags := types.ObjectValueFrom(ctx, NetworkingAccessObjectType.AttrTypes, TFNetworkingAccessModel{
			Type:         types.StringPointerValue(networkingAccessType),
			ConnectionID: types.StringPointerValue(networkingConnectionID),
		})
		if diags.HasError() {
			return nil, diags
		}
		networkingModel, diags := types.ObjectValueFrom(ctx, NetworkingObjectType.AttrTypes, TFNetworkingModel{
			Access: networkingAccessModel,
		})
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.Networking = networkingModel
	}

	connectionModel.Azure = types.ObjectNull(AzureObjectType.AttrTypes)
	if apiResp.Azure != nil {
		azure, diags := types.ObjectValueFrom(ctx, AzureObjectType.AttrTypes, TFAzureModel{
			ServicePrincipalID: types.StringPointerValue(apiResp.Azure.ServicePrincipalId),
			StorageAccountName: types.StringPointerValue(apiResp.Azure.StorageAccountName),
			Region:             types.StringPointerValue(apiResp.Azure.Region),
		})
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.Azure = azure
	}

	connectionModel.AWS = types.ObjectNull(AWSObjectType.AttrTypes)
	if apiResp.Aws != nil {
		aws, diags := types.ObjectValueFrom(ctx, AWSObjectType.AttrTypes, TFAWSModel{
			RoleArn: types.StringPointerValue(apiResp.Aws.RoleArn),
		})
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.AWS = aws
	}

	connectionModel.GCP = types.ObjectNull(GCPObjectType.AttrTypes)
	if apiResp.Gcp != nil {
		gcp, diags := types.ObjectValueFrom(ctx, GCPObjectType.AttrTypes, TFGCPModel{
			ServiceAccountID: types.StringPointerValue(apiResp.Gcp.ServiceAccountId),
		})
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.GCP = gcp
	}

	connectionModel.Headers = types.MapNull(types.StringType)
	// this is to handle the case where empty headers are returned as an empty map from the API, which is equivalent to a null value
	if apiResp.Headers != nil && len(*apiResp.Headers) > 0 {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, apiResp.Headers)
		if diags.HasError() {
			return nil, diags
		}
		connectionModel.Headers = mapValue
	}

	connectionModel.FailoverConnections = types.ListNull(FailoverConnectionObjectType)

	return &connectionModel, nil
}

// NewTFStreamConnectionDS creates a data source model by reusing NewTFStreamConnection and converting.
func NewTFStreamConnectionDS(ctx context.Context, projID, instanceName, workspaceName string, currAuthConfig, currSchemaRegistryAuthConfig *types.Object, apiResp *admin.StreamsConnection) (*TFStreamConnectionDSModel, diag.Diagnostics) {
	model, diags := NewTFStreamConnection(ctx, projID, instanceName, workspaceName, currAuthConfig, currSchemaRegistryAuthConfig, apiResp, nil)
	if diags.HasError() {
		return nil, diags
	}
	return model.ToDS(), nil
}

func newTFConnectionAuthenticationModel(ctx context.Context, currAuthConfig *types.Object, authResp *admin.StreamsKafkaAuthentication) (*types.Object, diag.Diagnostics) {
	if authResp != nil {
		resultAuthModel := TFConnectionAuthenticationModel{
			Mechanism:                 types.StringPointerValue(authResp.Mechanism),
			Method:                    types.StringPointerValue(authResp.Method),
			Username:                  types.StringPointerValue(authResp.Username),
			TokenEndpointURL:          types.StringPointerValue(authResp.TokenEndpointUrl),
			ClientID:                  types.StringPointerValue(authResp.ClientId),
			Scope:                     types.StringPointerValue(authResp.Scope),
			SaslOauthbearerExtensions: types.StringPointerValue(authResp.SaslOauthbearerExtensions),
		}

		if currAuthConfig != nil && !currAuthConfig.IsNull() { // if config is available (create & update of resource) password value is set in new state
			configAuthModel := &TFConnectionAuthenticationModel{}
			if diags := currAuthConfig.As(ctx, configAuthModel, basetypes.ObjectAsOptions{}); diags.HasError() {
				return nil, diags
			}
			resultAuthModel.Password = configAuthModel.Password
			resultAuthModel.ClientSecret = configAuthModel.ClientSecret
		}

		resultObject, diags := types.ObjectValueFrom(ctx, ConnectionAuthenticationObjectType.AttrTypes, resultAuthModel)
		if diags.HasError() {
			return nil, diags
		}
		return &resultObject, nil
	}
	return new(types.ObjectNull(ConnectionAuthenticationObjectType.AttrTypes)), nil
}

func newTFSchemaRegistryAuthentication(ctx context.Context, currAuthConfig *types.Object, authResp *admin.SchemaRegistryAuthentication) (*types.Object, diag.Diagnostics) {
	if authResp != nil {
		resultAuthModel := TFSchemaRegistryAuthenticationModel{
			Type:     types.StringValue(authResp.Type),
			Username: types.StringPointerValue(authResp.Username),
		}

		if currAuthConfig != nil && !currAuthConfig.IsNull() { // if config is available (create & update of resource) password value is set in new state
			configAuthModel := &TFSchemaRegistryAuthenticationModel{}
			if diags := currAuthConfig.As(ctx, configAuthModel, basetypes.ObjectAsOptions{}); diags.HasError() {
				return nil, diags
			}
			resultAuthModel.Password = configAuthModel.Password
		}

		resultObject, diags := types.ObjectValueFrom(ctx, SchemaRegistryAuthenticationObjectType.AttrTypes, resultAuthModel)
		if diags.HasError() {
			return nil, diags
		}
		return &resultObject, nil
	}
	return new(types.ObjectNull(SchemaRegistryAuthenticationObjectType.AttrTypes)), nil
}

// newFailoverConnectionReq converts one TFFailoverConnectionModel to a StreamsConnection API request.
// A failover connection has the same shape as a primary connection (plus a region), so it maps the
// failover model onto a TFStreamConnectionModel and reuses NewStreamConnectionReq.
func newFailoverConnectionReq(ctx context.Context, fc *TFFailoverConnectionModel) (*admin.StreamsConnection, diag.Diagnostics) {
	plan := &TFStreamConnectionModel{
		TFStreamConnectionCommonModel: TFStreamConnectionCommonModel{
			ConnectionName:               fc.Name,
			Type:                         fc.Type,
			ClusterName:                  fc.ClusterName,
			ClusterProjectID:             fc.ClusterProjectID,
			DBRoleToExecute:              fc.DBRoleToExecute,
			BootstrapServers:             fc.BootstrapServers,
			Authentication:               fc.Authentication,
			Config:                       fc.Config,
			Security:                     fc.Security,
			Networking:                   fc.Networking,
			AWS:                          fc.AWS,
			Azure:                        fc.Azure,
			GCP:                          fc.GCP,
			URL:                          fc.URL,
			Headers:                      fc.Headers,
			SchemaRegistryProvider:       fc.SchemaRegistryProvider,
			SchemaRegistryURLs:           fc.SchemaRegistryURLs,
			SchemaRegistryAuthentication: fc.SchemaRegistryAuthentication,
		},
	}
	conn, diags := NewStreamConnectionReq(ctx, plan)
	if diags.HasError() {
		return nil, diags
	}
	// Region is specific to failover connections and not part of NewStreamConnectionReq.
	conn.Region = fc.Region.ValueStringPointer()
	return conn, nil
}

// newTFFailoverConnectionModel converts an API StreamsConnection response into a TFFailoverConnectionModel.
// newTFFailoverConnectionsList returns the failover connections as a types.List.
// knownFailoverList carries the desired connections (plan on Update, prior state on Read), which are
// the source of truth for the connection fields: the list endpoint returns sparse/primary-shaped data,
// so we build state from the known values and only refresh the computed id (and drop connections that
// no longer exist server-side).
func newTFFailoverConnectionsList(ctx context.Context, streamsAPI admin.StreamsApi, projectID, workspaceName, primaryConnectionName string, knownFailoverList types.List) (types.List, diag.Diagnostics) {
	var allDiags diag.Diagnostics

	if knownFailoverList.IsNull() || knownFailoverList.IsUnknown() {
		return types.ListNull(FailoverConnectionObjectType), nil
	}
	var known []TFFailoverConnectionModel
	if diags := knownFailoverList.ElementsAs(ctx, &known, false); diags.HasError() {
		allDiags.Append(diags...)
		return types.ListNull(FailoverConnectionObjectType), allDiags
	}
	if len(known) == 0 {
		return types.ListNull(FailoverConnectionObjectType), nil
	}

	// Fetch current failover connections to refresh ids and detect out-of-band deletions.
	apiConns, _, err := streamsAPI.ListFailoverConnections(ctx, projectID, workspaceName, primaryConnectionName).Execute()
	if err != nil {
		allDiags.AddError("error listing failover connections", err.Error())
		return types.ListNull(FailoverConnectionObjectType), allDiags
	}
	apiResults := apiConns.GetResults()
	idByName := make(map[string]*string, len(apiResults))
	for i := range apiResults {
		if apiResults[i].Name != nil {
			idByName[*apiResults[i].Name] = apiResults[i].Id
		}
	}

	var fcModels []TFFailoverConnectionModel
	for i := range known {
		fcID, ok := idByName[known[i].Name.ValueString()]
		if !ok {
			continue // connection was removed out-of-band; drop it so it is recreated on next apply
		}
		item := known[i]
		item.ID = types.StringPointerValue(fcID)
		fcModels = append(fcModels, item)
	}

	if len(fcModels) == 0 {
		return types.ListNull(FailoverConnectionObjectType), nil
	}

	listVal, diags := types.ListValueFrom(ctx, FailoverConnectionObjectType, fcModels)
	allDiags.Append(diags...)
	return listVal, allDiags
}

// syncFailoverConnections computes the diff between state and plan failover connections and calls
// create/update/delete API endpoints as needed.
func syncFailoverConnections(ctx context.Context, streamsAPI admin.StreamsApi, projectID, workspaceName, primaryConnectionName string, planList, stateList types.List) diag.Diagnostics {
	var allDiags diag.Diagnostics

	var planFCs, stateFCs []TFFailoverConnectionModel
	if !planList.IsNull() && !planList.IsUnknown() {
		if diags := planList.ElementsAs(ctx, &planFCs, false); diags.HasError() {
			return diags
		}
	}
	if !stateList.IsNull() && !stateList.IsUnknown() {
		if diags := stateList.ElementsAs(ctx, &stateFCs, false); diags.HasError() {
			return diags
		}
	}

	// Index state by name to look up existing IDs.
	stateByName := make(map[string]TFFailoverConnectionModel, len(stateFCs))
	for i := range stateFCs {
		stateByName[stateFCs[i].Name.ValueString()] = stateFCs[i]
	}
	planByName := make(map[string]TFFailoverConnectionModel, len(planFCs))
	for i := range planFCs {
		planByName[planFCs[i].Name.ValueString()] = planFCs[i]
	}

	// Connections to create (in plan but not in state).
	for i := range planFCs {
		name := planFCs[i].Name.ValueString()
		if _, exists := stateByName[name]; !exists {
			conn, diags := newFailoverConnectionReq(ctx, &planFCs[i])
			if diags.HasError() {
				allDiags.Append(diags...)
				return allDiags
			}
			if _, _, err := streamsAPI.CreateFailoverConnection(ctx, projectID, workspaceName, primaryConnectionName, conn).Execute(); err != nil {
				allDiags.AddError("error creating failover connection", fmt.Sprintf("connection %q: %s", name, err.Error()))
				return allDiags
			}
		}
	}

	// Connections to update (in both plan and state).
	for i := range planFCs {
		name := planFCs[i].Name.ValueString()
		if stateFC, exists := stateByName[name]; exists {
			existingID := stateFC.ID.ValueString()
			if existingID == "" {
				continue
			}
			conn, diags := newFailoverConnectionReq(ctx, &planFCs[i])
			if diags.HasError() {
				allDiags.Append(diags...)
				return allDiags
			}
			if _, _, err := streamsAPI.UpdateStreamFailoverConnection(ctx, projectID, workspaceName, primaryConnectionName, existingID, conn).Execute(); err != nil {
				allDiags.AddError("error updating failover connection", fmt.Sprintf("connection %q: %s", name, err.Error()))
				return allDiags
			}
		}
	}

	// Connections to delete (in state but not in plan).
	for i := range stateFCs {
		name := stateFCs[i].Name.ValueString()
		if _, exists := planByName[name]; !exists {
			if fcID := stateFCs[i].ID.ValueString(); fcID != "" {
				if _, err := streamsAPI.DeleteStreamFailoverConnection(ctx, projectID, workspaceName, primaryConnectionName, fcID).Execute(); err != nil {
					allDiags.AddError("error deleting failover connection", fmt.Sprintf("connection %q: %s", name, err.Error()))
					return allDiags
				}
			}
		}
	}

	return allDiags
}

// deleteAllFailoverConnections removes all failover connections tracked in stateList.
func deleteAllFailoverConnections(ctx context.Context, streamsAPI admin.StreamsApi, projectID, workspaceName, primaryConnectionName string, stateList types.List) diag.Diagnostics {
	if stateList.IsNull() || stateList.IsUnknown() || len(stateList.Elements()) == 0 {
		return nil
	}
	var stateFCs []TFFailoverConnectionModel
	if diags := stateList.ElementsAs(ctx, &stateFCs, false); diags.HasError() {
		return diags
	}
	for i := range stateFCs {
		fcID := stateFCs[i].ID.ValueString()
		if fcID == "" {
			continue
		}
		if _, err := streamsAPI.DeleteStreamFailoverConnection(ctx, projectID, workspaceName, primaryConnectionName, fcID).Execute(); err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("error deleting failover connection", fmt.Sprintf("connection %q: %s", stateFCs[i].Name.ValueString(), err.Error()))}
		}
	}
	return nil
}

func NewTFStreamConnectionsDS(ctx context.Context,
	streamConnectionsConfig *TFStreamConnectionsDSModel,
	paginatedResult *admin.PaginatedApiStreamsConnection) (*TFStreamConnectionsDSModel, diag.Diagnostics) {
	input := paginatedResult.GetResults()
	results := make([]TFStreamConnectionDSModel, len(input))

	workspaceName := streamConnectionsConfig.WorkspaceName.ValueString()
	instanceName := streamConnectionsConfig.InstanceName.ValueString()

	for i := range input {
		projectID := streamConnectionsConfig.ProjectID.ValueString()
		connectionModel, diags := NewTFStreamConnectionDS(ctx, projectID, instanceName, workspaceName, nil, nil, &input[i])
		if diags.HasError() {
			return nil, diags
		}
		results[i] = *connectionModel
	}

	return &TFStreamConnectionsDSModel{
		ID:            types.StringValue(id.UniqueId()),
		ProjectID:     streamConnectionsConfig.ProjectID,
		InstanceName:  streamConnectionsConfig.InstanceName,
		WorkspaceName: streamConnectionsConfig.WorkspaceName,
		Results:       results,
		PageNum:       streamConnectionsConfig.PageNum,
		ItemsPerPage:  streamConnectionsConfig.ItemsPerPage,
		TotalCount:    types.Int64PointerValue(conversion.IntPtrToInt64Ptr(paginatedResult.TotalCount)),
	}, nil
}
