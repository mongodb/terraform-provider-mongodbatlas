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
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
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
func newFailoverConnectionReq(ctx context.Context, fc *TFFailoverConnectionModel) (*admin.StreamsConnection, diag.Diagnostics) {
	conn := admin.StreamsConnection{
		Name:             fc.Name.ValueStringPointer(),
		Type:             fc.Type.ValueStringPointer(),
		ClusterName:      fc.ClusterName.ValueStringPointer(),
		ClusterGroupId:   fc.ClusterProjectID.ValueStringPointer(),
		BootstrapServers: fc.BootstrapServers.ValueStringPointer(),
		Url:              fc.URL.ValueStringPointer(),
		Provider:         fc.SchemaRegistryProvider.ValueStringPointer(),
	}

	if !fc.Authentication.IsNull() && !fc.Authentication.IsUnknown() {
		authModel := &TFConnectionAuthenticationModel{}
		if diags := fc.Authentication.As(ctx, authModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.Authentication = &admin.StreamsKafkaAuthentication{
			Mechanism:                 authModel.Mechanism.ValueStringPointer(),
			Method:                    authModel.Method.ValueStringPointer(),
			Password:                  authModel.Password.ValueStringPointer(),
			Username:                  authModel.Username.ValueStringPointer(),
			TokenEndpointUrl:          authModel.TokenEndpointURL.ValueStringPointer(),
			ClientId:                  authModel.ClientID.ValueStringPointer(),
			ClientSecret:              authModel.ClientSecret.ValueStringPointer(),
			Scope:                     authModel.Scope.ValueStringPointer(),
			SaslOauthbearerExtensions: authModel.SaslOauthbearerExtensions.ValueStringPointer(),
		}
	}

	if !fc.Security.IsNull() && !fc.Security.IsUnknown() {
		secModel := &TFConnectionSecurityModel{}
		if diags := fc.Security.As(ctx, secModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.Security = &admin.StreamsKafkaSecurity{
			BrokerPublicCertificate: secModel.BrokerPublicCertificate.ValueStringPointer(),
			Protocol:                secModel.Protocol.ValueStringPointer(),
		}
	}

	if !fc.Config.IsNull() && !fc.Config.IsUnknown() {
		configMap := &map[string]string{}
		if diags := fc.Config.ElementsAs(ctx, configMap, true); diags.HasError() {
			return nil, diags
		}
		conn.Config = configMap
	}

	if !fc.DBRoleToExecute.IsNull() && !fc.DBRoleToExecute.IsUnknown() {
		dbRoleModel := &TFDbRoleToExecuteModel{}
		if diags := fc.DBRoleToExecute.As(ctx, dbRoleModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.DbRoleToExecute = &admin.DBRoleToExecute{
			Role: dbRoleModel.Role.ValueStringPointer(),
			Type: dbRoleModel.Type.ValueStringPointer(),
		}
	}

	if !fc.Networking.IsNull() && !fc.Networking.IsUnknown() {
		networkingModel := &TFNetworkingModel{}
		if diags := fc.Networking.As(ctx, networkingModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		networkingAccessModel := &TFNetworkingAccessModel{}
		if diags := networkingModel.Access.As(ctx, networkingAccessModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		switch fc.Type.ValueString() {
		case ConnectionTypeAzureBlobStorage, ConnectionTypeGCPPubSub:
			conn.PublicPrivateNetworking = &admin.StreamsPublicPrivateLinkNetworking{
				Access: &admin.StreamsPublicPrivateLinkNetworkingAccess{
					Type:         networkingAccessModel.Type.ValueStringPointer(),
					ConnectionId: networkingAccessModel.ConnectionID.ValueStringPointer(),
				},
			}
		default:
			conn.Networking = &admin.StreamsKafkaNetworking{
				Access: &admin.StreamsKafkaNetworkingAccess{
					Type:         networkingAccessModel.Type.ValueStringPointer(),
					ConnectionId: networkingAccessModel.ConnectionID.ValueStringPointer(),
				},
			}
		}
	}

	if !fc.AWS.IsNull() && !fc.AWS.IsUnknown() {
		awsModel := &TFAWSModel{}
		if diags := fc.AWS.As(ctx, awsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.Aws = &admin.StreamsAWSConnectionConfig{RoleArn: awsModel.RoleArn.ValueStringPointer()}
	}

	if !fc.GCP.IsNull() && !fc.GCP.IsUnknown() {
		gcpModel := &TFGCPModel{}
		if diags := fc.GCP.As(ctx, gcpModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.Gcp = &admin.StreamsGCPConnectionConfig{ServiceAccountId: gcpModel.ServiceAccountID.ValueStringPointer()}
	}

	if !fc.Azure.IsNull() && !fc.Azure.IsUnknown() {
		azureModel := &TFAzureModel{}
		if diags := fc.Azure.As(ctx, azureModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.Azure = &admin.AzureConnection{
			ServicePrincipalId: azureModel.ServicePrincipalID.ValueStringPointer(),
			StorageAccountName: azureModel.StorageAccountName.ValueStringPointer(),
			Region:             azureModel.Region.ValueStringPointer(),
		}
	}

	if !fc.Headers.IsNull() && !fc.Headers.IsUnknown() {
		headersMap := make(map[string]string)
		if diags := fc.Headers.ElementsAs(ctx, &headersMap, true); diags.HasError() {
			return nil, diags
		}
		conn.Headers = &headersMap
	}

	if !fc.SchemaRegistryURLs.IsNull() && !fc.SchemaRegistryURLs.IsUnknown() {
		var urls []string
		if diags := fc.SchemaRegistryURLs.ElementsAs(ctx, &urls, false); diags.HasError() {
			return nil, diags
		}
		conn.SchemaRegistryUrls = &urls
	}

	if !fc.SchemaRegistryAuthentication.IsNull() && !fc.SchemaRegistryAuthentication.IsUnknown() {
		srAuthModel := &TFSchemaRegistryAuthenticationModel{}
		if diags := fc.SchemaRegistryAuthentication.As(ctx, srAuthModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		conn.SchemaRegistryAuthentication = &admin.SchemaRegistryAuthentication{
			Type:     srAuthModel.Type.ValueString(),
			Username: srAuthModel.Username.ValueStringPointer(),
			Password: srAuthModel.Password.ValueStringPointer(),
		}
	}

	return &conn, nil
}

// newTFFailoverConnectionModel converts an API StreamsConnection response into a TFFailoverConnectionModel.
// currAuthConfig and currSRAuthConfig preserve sensitive fields (passwords) from prior state.
func newTFFailoverConnectionModel(ctx context.Context, apiResp *admin.StreamsConnection, currAuthConfig, currSRAuthConfig *types.Object) (*TFFailoverConnectionModel, diag.Diagnostics) {
	fc := &TFFailoverConnectionModel{
		ID:                     types.StringPointerValue(apiResp.Id),
		Name:                   types.StringPointerValue(apiResp.Name),
		Type:                   types.StringPointerValue(apiResp.Type),
		ClusterName:            types.StringPointerValue(apiResp.ClusterName),
		ClusterProjectID:       types.StringPointerValue(apiResp.ClusterGroupId),
		BootstrapServers:       types.StringPointerValue(apiResp.BootstrapServers),
		URL:                    types.StringPointerValue(apiResp.Url),
		SchemaRegistryProvider: types.StringPointerValue(apiResp.Provider),
		SchemaRegistryURLs:     types.ListNull(types.StringType),
	}

	authModel, diags := newTFConnectionAuthenticationModel(ctx, currAuthConfig, apiResp.Authentication)
	if diags.HasError() {
		return nil, diags
	}
	fc.Authentication = *authModel

	fc.Config = types.MapNull(types.StringType)
	if apiResp.Config != nil {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, apiResp.Config)
		if diags.HasError() {
			return nil, diags
		}
		fc.Config = mapValue
	}

	if apiResp.SchemaRegistryUrls != nil {
		urls, diags := types.ListValueFrom(ctx, types.StringType, apiResp.SchemaRegistryUrls)
		if diags.HasError() {
			return nil, diags
		}
		fc.SchemaRegistryURLs = urls
	}

	srAuthModel, diags := newTFSchemaRegistryAuthentication(ctx, currSRAuthConfig, apiResp.SchemaRegistryAuthentication)
	if diags.HasError() {
		return nil, diags
	}
	fc.SchemaRegistryAuthentication = *srAuthModel

	fc.Security = types.ObjectNull(ConnectionSecurityObjectType.AttrTypes)
	if apiResp.Security != nil {
		secObj, diags := types.ObjectValueFrom(ctx, ConnectionSecurityObjectType.AttrTypes, TFConnectionSecurityModel{
			BrokerPublicCertificate: types.StringPointerValue(apiResp.Security.BrokerPublicCertificate),
			Protocol:                types.StringPointerValue(apiResp.Security.Protocol),
		})
		if diags.HasError() {
			return nil, diags
		}
		fc.Security = secObj
	}

	fc.DBRoleToExecute = types.ObjectNull(DBRoleToExecuteObjectType.AttrTypes)
	if apiResp.DbRoleToExecute != nil {
		dbObj, diags := types.ObjectValueFrom(ctx, DBRoleToExecuteObjectType.AttrTypes, TFDbRoleToExecuteModel{
			Role: types.StringPointerValue(apiResp.DbRoleToExecute.Role),
			Type: types.StringPointerValue(apiResp.DbRoleToExecute.Type),
		})
		if diags.HasError() {
			return nil, diags
		}
		fc.DBRoleToExecute = dbObj
	}

	fc.Networking = types.ObjectNull(NetworkingObjectType.AttrTypes)
	var networkingAccessType, networkingConnectionID *string
	if apiResp.Networking != nil && apiResp.Networking.Access != nil {
		networkingAccessType = apiResp.Networking.Access.Type
		networkingConnectionID = apiResp.Networking.Access.ConnectionId
	} else if apiResp.PublicPrivateNetworking != nil && apiResp.PublicPrivateNetworking.Access != nil {
		networkingAccessType = apiResp.PublicPrivateNetworking.Access.Type
		networkingConnectionID = apiResp.PublicPrivateNetworking.Access.ConnectionId
	}
	if networkingAccessType != nil {
		accessObj, diags := types.ObjectValueFrom(ctx, NetworkingAccessObjectType.AttrTypes, TFNetworkingAccessModel{
			Type:         types.StringPointerValue(networkingAccessType),
			ConnectionID: types.StringPointerValue(networkingConnectionID),
		})
		if diags.HasError() {
			return nil, diags
		}
		netObj, diags := types.ObjectValueFrom(ctx, NetworkingObjectType.AttrTypes, TFNetworkingModel{Access: accessObj})
		if diags.HasError() {
			return nil, diags
		}
		fc.Networking = netObj
	}

	fc.AWS = types.ObjectNull(AWSObjectType.AttrTypes)
	if apiResp.Aws != nil {
		awsObj, diags := types.ObjectValueFrom(ctx, AWSObjectType.AttrTypes, TFAWSModel{RoleArn: types.StringPointerValue(apiResp.Aws.RoleArn)})
		if diags.HasError() {
			return nil, diags
		}
		fc.AWS = awsObj
	}

	fc.GCP = types.ObjectNull(GCPObjectType.AttrTypes)
	if apiResp.Gcp != nil {
		gcpObj, diags := types.ObjectValueFrom(ctx, GCPObjectType.AttrTypes, TFGCPModel{ServiceAccountID: types.StringPointerValue(apiResp.Gcp.ServiceAccountId)})
		if diags.HasError() {
			return nil, diags
		}
		fc.GCP = gcpObj
	}

	fc.Azure = types.ObjectNull(AzureObjectType.AttrTypes)
	if apiResp.Azure != nil {
		azureObj, diags := types.ObjectValueFrom(ctx, AzureObjectType.AttrTypes, TFAzureModel{
			ServicePrincipalID: types.StringPointerValue(apiResp.Azure.ServicePrincipalId),
			StorageAccountName: types.StringPointerValue(apiResp.Azure.StorageAccountName),
			Region:             types.StringPointerValue(apiResp.Azure.Region),
		})
		if diags.HasError() {
			return nil, diags
		}
		fc.Azure = azureObj
	}

	fc.Headers = types.MapNull(types.StringType)
	if apiResp.Headers != nil && len(*apiResp.Headers) > 0 {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, apiResp.Headers)
		if diags.HasError() {
			return nil, diags
		}
		fc.Headers = mapValue
	}

	return fc, nil
}

// newTFFailoverConnectionsList reads the current failover connections from the API and returns a types.List.
// stateFailoverConnections is used to retrieve the names of failover connections already tracked in state.
func newTFFailoverConnectionsList(ctx context.Context, streamsAPI admin.StreamsApi, projectID, workspaceName string, planFailoverList types.List, stateFailoverList types.List) (types.List, diag.Diagnostics) {
	var allDiags diag.Diagnostics

	// Collect names in deterministic order: plan order first, then state-only entries.
	seen := make(map[string]bool)
	var orderedNames []string
	planAuthByName := make(map[string]types.Object)
	planSRAuthByName := make(map[string]types.Object)
	for listIdx, list := range []types.List{planFailoverList, stateFailoverList} {
		isPlanList := listIdx == 0
		if list.IsNull() || list.IsUnknown() {
			continue
		}
		var items []TFFailoverConnectionModel
		if diags := list.ElementsAs(ctx, &items, false); diags.HasError() {
			allDiags.Append(diags...)
			return types.ListNull(FailoverConnectionObjectType), allDiags
		}
		for _, item := range items {
			name := item.Name.ValueString()
			if name == "" {
				continue
			}
			if !seen[name] {
				seen[name] = true
				orderedNames = append(orderedNames, name)
			}
			// Only plan items carry auth config for password preservation.
			if isPlanList {
				planAuthByName[name] = item.Authentication
				planSRAuthByName[name] = item.SchemaRegistryAuthentication
			}
		}
	}

	if len(orderedNames) == 0 {
		return types.ListNull(FailoverConnectionObjectType), nil
	}

	var fcModels []TFFailoverConnectionModel
	for _, name := range orderedNames {
		apiResp, getResp, err := streamsAPI.GetStreamConnection(ctx, projectID, workspaceName, name).Execute()
		if err != nil {
			if validate.StatusNotFound(getResp) {
				continue // connection was removed out-of-band; skip it
			}
			allDiags.AddError("error reading failover connection", fmt.Sprintf("connection %q: %s", name, err.Error()))
			return types.ListNull(FailoverConnectionObjectType), allDiags
		}
		auth := planAuthByName[name]
		srAuth := planSRAuthByName[name]
		fc, diags := newTFFailoverConnectionModel(ctx, apiResp, &auth, &srAuth)
		if diags.HasError() {
			allDiags.Append(diags...)
			return types.ListNull(FailoverConnectionObjectType), allDiags
		}
		fcModels = append(fcModels, *fc)
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
	for _, fc := range stateFCs {
		stateByName[fc.Name.ValueString()] = fc
	}
	planByName := make(map[string]TFFailoverConnectionModel, len(planFCs))
	for _, fc := range planFCs {
		planByName[fc.Name.ValueString()] = fc
	}

	// Connections to create (in plan but not in state).
	var toCreate []admin.StreamsConnection
	for _, fc := range planFCs {
		name := fc.Name.ValueString()
		if _, exists := stateByName[name]; !exists {
			fcItem := fc
			conn, diags := newFailoverConnectionReq(ctx, &fcItem)
			if diags.HasError() {
				allDiags.Append(diags...)
				return allDiags
			}
			toCreate = append(toCreate, *conn)
		}
	}
	if len(toCreate) > 0 {
		if _, _, err := streamsAPI.CreateFailoverConnections(ctx, projectID, workspaceName, primaryConnectionName, &admin.StreamsCreateFailoverConnectionsRequest{Connections: toCreate}).Execute(); err != nil {
			allDiags.AddError("error creating failover connections", err.Error())
			return allDiags
		}
	}

	// Connections to update (in both plan and state).
	var toUpdate []admin.StreamsFailoverConnectionUpdate
	for _, fc := range planFCs {
		name := fc.Name.ValueString()
		if stateFC, exists := stateByName[name]; exists {
			existingID := stateFC.ID.ValueString()
			if existingID == "" {
				continue
			}
			fcItem := fc
			conn, diags := newFailoverConnectionReq(ctx, &fcItem)
			if diags.HasError() {
				allDiags.Append(diags...)
				return allDiags
			}
			toUpdate = append(toUpdate, admin.StreamsFailoverConnectionUpdate{Id: existingID, Connection: *conn})
		}
	}
	if len(toUpdate) > 0 {
		if _, _, err := streamsAPI.UpdateFailoverConnections(ctx, projectID, workspaceName, primaryConnectionName, &admin.StreamsUpdateFailoverConnectionsRequest{Connections: toUpdate}).Execute(); err != nil {
			allDiags.AddError("error updating failover connections", err.Error())
			return allDiags
		}
	}

	// Connections to delete (in state but not in plan).
	var toDeleteIDs []string
	for _, fc := range stateFCs {
		name := fc.Name.ValueString()
		if _, exists := planByName[name]; !exists {
			if id := fc.ID.ValueString(); id != "" {
				toDeleteIDs = append(toDeleteIDs, id)
			}
		}
	}
	if len(toDeleteIDs) > 0 {
		if _, err := streamsAPI.DeleteFailoverConnections(ctx, projectID, workspaceName, primaryConnectionName, &admin.StreamsDeleteFailoverConnectionsRequest{ConnectionIds: toDeleteIDs}).Execute(); err != nil {
			allDiags.AddError("error deleting failover connections", err.Error())
			return allDiags
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
	var ids []string
	for _, fc := range stateFCs {
		if id := fc.ID.ValueString(); id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	if _, err := streamsAPI.DeleteFailoverConnections(ctx, projectID, workspaceName, primaryConnectionName, &admin.StreamsDeleteFailoverConnectionsRequest{ConnectionIds: ids}).Execute(); err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("error deleting failover connections", err.Error())}
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
