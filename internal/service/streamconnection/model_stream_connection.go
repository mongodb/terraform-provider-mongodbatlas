package streamconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

func NewStreamConnectionReq(ctx context.Context, plan *TFStreamConnectionModel) (*admin.StreamsConnection, diag.Diagnostics) {
	streamConnection := admin.StreamsConnection{
		Name:             plan.ConnectionName.ValueStringPointer(),
		Type:             plan.Type.ValueStringPointer(),
		ClusterName:      plan.ClusterName.ValueStringPointer(),
		BootstrapServers: plan.BootstrapServers.ValueStringPointer(),
	}
	if !plan.Authentication.IsNull() {
		authenticationModel := &TFConnectionAuthenticationModel{}
		if diags := plan.Authentication.As(ctx, authenticationModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamConnection.Authentication = &admin.StreamsKafkaAuthentication{
			Mechanism: authenticationModel.Mechanism.ValueStringPointer(),
			Password:  authenticationModel.Password.ValueStringPointer(),
			Username:  authenticationModel.Username.ValueStringPointer(),
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

	return &streamConnection, nil
}

func NewTFStreamConnection(ctx context.Context, projID, instanceName string, currAuthConfig *types.Object, apiResp *admin.StreamsConnection) (*TFStreamConnectionModel, diag.Diagnostics) {
	rID := fmt.Sprintf("%s-%s-%s", instanceName, projID, conversion.SafeString(apiResp.Name))
	connectionModel := TFStreamConnectionModel{
		ID:               types.StringValue(rID),
		ProjectID:        types.StringValue(projID),
		InstanceName:     types.StringValue(instanceName),
		ConnectionName:   types.StringPointerValue(apiResp.Name),
		Type:             types.StringPointerValue(apiResp.Type),
		ClusterName:      types.StringPointerValue(apiResp.ClusterName),
		BootstrapServers: types.StringPointerValue(apiResp.BootstrapServers),
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

	return &connectionModel, nil
}

func newTFConnectionAuthenticationModel(ctx context.Context, currAuthConfig *types.Object, authResp *admin.StreamsKafkaAuthentication) (*types.Object, diag.Diagnostics) {
	if authResp != nil {
		resultAuthModel := TFConnectionAuthenticationModel{
			Mechanism: types.StringPointerValue(authResp.Mechanism),
			Username:  types.StringPointerValue(authResp.Username),
		}

		if currAuthConfig != nil && !currAuthConfig.IsNull() { // if config is available (create & update of resource) password value is set in new state
			configAuthModel := &TFConnectionAuthenticationModel{}
			if diags := currAuthConfig.As(ctx, configAuthModel, basetypes.ObjectAsOptions{}); diags.HasError() {
				return nil, diags
			}
			resultAuthModel.Password = configAuthModel.Password
		}

		resultObject, diags := types.ObjectValueFrom(ctx, ConnectionAuthenticationObjectType.AttrTypes, resultAuthModel)
		if diags.HasError() {
			return nil, diags
		}
		return &resultObject, nil
	}
	nullValue := types.ObjectNull(ConnectionAuthenticationObjectType.AttrTypes)
	return &nullValue, nil
}

func NewTFStreamConnections(ctx context.Context,
	streamConnectionsConfig *TFStreamConnectionsDSModel,
	paginatedResult *admin.PaginatedApiStreamsConnection) (*TFStreamConnectionsDSModel, diag.Diagnostics) {
	input := paginatedResult.GetResults()
	results := make([]TFStreamConnectionModel, len(input))
	for i := range input {
		projectID := streamConnectionsConfig.ProjectID.ValueString()
		instanceName := streamConnectionsConfig.InstanceName.ValueString()
		connectionModel, diags := NewTFStreamConnection(ctx, projectID, instanceName, nil, &input[i])
		if diags.HasError() {
			return nil, diags
		}
		results[i] = *connectionModel
	}
	return &TFStreamConnectionsDSModel{
		ID:           types.StringValue(id.UniqueId()),
		ProjectID:    streamConnectionsConfig.ProjectID,
		InstanceName: streamConnectionsConfig.InstanceName,
		Results:      results,
		PageNum:      streamConnectionsConfig.PageNum,
		ItemsPerPage: streamConnectionsConfig.ItemsPerPage,
		TotalCount:   types.Int64PointerValue(conversion.IntPtrToInt64Ptr(paginatedResult.TotalCount)),
	}, nil
}
