package streamconnection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
)

func azureObject(t *testing.T, servicePrincipalID, storageAccountName, region types.String) types.Object {
	t.Helper()
	return types.ObjectValueMust(streamconnection.AzureObjectType.AttrTypes, map[string]attr.Value{
		"service_principal_id": servicePrincipalID,
		"storage_account_name": storageAccountName,
		"region":               region,
	})
}

func TestIsAliasOnlyTransition(t *testing.T) {
	legacyState := streamconnection.TFStreamConnectionModel{
		TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
			ID:             types.StringValue("workspace-project-connection"),
			InstanceName:   types.StringValue("workspace"),
			WorkspaceName:  types.StringNull(),
			ProjectID:      types.StringValue("project"),
			ConnectionName: types.StringValue("connection"),
			Type:           types.StringValue("Sample"),
		},
	}

	azureState := streamconnection.TFStreamConnectionModel{
		TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
			ID:             types.StringValue("workspace-project-connection"),
			InstanceName:   types.StringValue("workspace"),
			WorkspaceName:  types.StringNull(),
			ProjectID:      types.StringValue("project"),
			ConnectionName: types.StringValue("connection"),
			Type:           types.StringValue("AzureBlobStorage"),
			Azure:          azureObject(t, types.StringValue("principal"), types.StringValue("storage"), types.StringValue("eastus2")),
		},
	}

	testCases := map[string]struct {
		plan                streamconnection.TFStreamConnectionModel
		state               streamconnection.TFStreamConnectionModel
		aliasOnlyTransition bool
	}{
		"legacy_to_canonical": {
			state: legacyState,
			plan: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					ID:             types.StringUnknown(),
					InstanceName:   types.StringNull(),
					WorkspaceName:  types.StringValue("workspace"),
					ProjectID:      types.StringValue("project"),
					ConnectionName: types.StringValue("connection"),
					Type:           types.StringValue("Sample"),
				},
			},
			aliasOnlyTransition: true,
		},
		"azure_region_unknown_in_plan": {
			state: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					ID:             types.StringValue("workspace-project-connection"),
					InstanceName:   types.StringValue("workspace"),
					WorkspaceName:  types.StringNull(),
					ProjectID:      types.StringValue("project"),
					ConnectionName: types.StringValue("connection"),
					Type:           types.StringValue("AzureBlobStorage"),
					Azure:          azureObject(t, types.StringValue("principal"), types.StringValue("storage"), types.StringValue("eastus2")),
				},
			},
			plan: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					ID:             types.StringUnknown(),
					InstanceName:   types.StringNull(),
					WorkspaceName:  types.StringValue("workspace"),
					ProjectID:      types.StringValue("project"),
					ConnectionName: types.StringValue("connection"),
					Type:           types.StringValue("AzureBlobStorage"),
					Azure:          azureObject(t, types.StringValue("principal"), types.StringValue("storage"), types.StringUnknown()),
				},
			},
			aliasOnlyTransition: true,
		},
		"azure_storage_account_changed_with_unknown_region": {
			state: azureState,
			plan: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					ID:             types.StringUnknown(),
					InstanceName:   types.StringNull(),
					WorkspaceName:  types.StringValue("workspace"),
					ProjectID:      types.StringValue("project"),
					ConnectionName: types.StringValue("connection"),
					Type:           types.StringValue("AzureBlobStorage"),
					Azure:          azureObject(t, types.StringValue("principal"), types.StringValue("other-storage"), types.StringUnknown()),
				},
			},
		},
		"azure_service_principal_changed_with_unknown_region": {
			state: azureState,
			plan: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					ID:             types.StringUnknown(),
					InstanceName:   types.StringNull(),
					WorkspaceName:  types.StringValue("workspace"),
					ProjectID:      types.StringValue("project"),
					ConnectionName: types.StringValue("connection"),
					Type:           types.StringValue("AzureBlobStorage"),
					Azure:          azureObject(t, types.StringValue("other-principal"), types.StringValue("storage"), types.StringUnknown()),
				},
			},
		},
		"different_workspace": {
			state: legacyState,
			plan: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					InstanceName:   types.StringNull(),
					WorkspaceName:  types.StringValue("other-workspace"),
					ProjectID:      types.StringValue("project"),
					ConnectionName: types.StringValue("connection"),
					Type:           types.StringValue("Sample"),
				},
			},
		},
		"connection_configuration_changed": {
			state: legacyState,
			plan: streamconnection.TFStreamConnectionModel{
				TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
					InstanceName:     types.StringNull(),
					WorkspaceName:    types.StringValue("workspace"),
					ProjectID:        types.StringValue("project"),
					ConnectionName:   types.StringValue("connection"),
					Type:             types.StringValue("Kafka"),
					BootstrapServers: types.StringValue("broker:9092"),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.aliasOnlyTransition, streamconnection.IsAliasOnlyTransition(t.Context(), &tc.plan, &tc.state))
		})
	}
}

// The resource has no Atlas client configured, so any Atlas request made by Update would panic on the nil client.
func TestUpdateAliasOnlyTransitionMakesNoAtlasRequests(t *testing.T) {
	ctx := t.Context()
	rs := streamconnection.Resource()
	var schemaResp resource.SchemaResponse
	rs.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	newModel := func(id, instanceName, workspaceName, region types.String) streamconnection.TFStreamConnectionModel {
		return streamconnection.TFStreamConnectionModel{
			TFStreamConnectionCommonModel: streamconnection.TFStreamConnectionCommonModel{
				ID:             id,
				InstanceName:   instanceName,
				WorkspaceName:  workspaceName,
				ProjectID:      types.StringValue("project"),
				ConnectionName: types.StringValue("connection"),
				Type:           types.StringValue("AzureBlobStorage"),
				Azure:          azureObject(t, types.StringValue("principal"), types.StringValue("storage"), region),

				ClusterName:                  types.StringNull(),
				ClusterProjectID:             types.StringNull(),
				Authentication:               types.ObjectNull(streamconnection.ConnectionAuthenticationObjectType.AttrTypes),
				BootstrapServers:             types.StringNull(),
				Config:                       types.MapNull(types.StringType),
				Security:                     types.ObjectNull(streamconnection.ConnectionSecurityObjectType.AttrTypes),
				DBRoleToExecute:              types.ObjectNull(streamconnection.DBRoleToExecuteObjectType.AttrTypes),
				Networking:                   types.ObjectNull(streamconnection.NetworkingObjectType.AttrTypes),
				AWS:                          types.ObjectNull(streamconnection.AWSObjectType.AttrTypes),
				GCP:                          types.ObjectNull(streamconnection.GCPObjectType.AttrTypes),
				Headers:                      types.MapNull(types.StringType),
				URL:                          types.StringNull(),
				SchemaRegistryProvider:       types.StringNull(),
				SchemaRegistryURLs:           types.ListNull(types.StringType),
				SchemaRegistryAuthentication: types.ObjectNull(streamconnection.SchemaRegistryAuthenticationObjectType.AttrTypes),
			},
			Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType, "update": types.StringType, "delete": types.StringType,
			})},
		}
	}
	state := newModel(types.StringValue("workspace-project-connection"), types.StringValue("workspace"), types.StringNull(), types.StringValue("eastus2"))
	plan := newModel(types.StringUnknown(), types.StringNull(), types.StringValue("workspace"), types.StringUnknown())

	req := resource.UpdateRequest{
		Plan:  tfsdk.Plan{Schema: schemaResp.Schema},
		State: tfsdk.State{Schema: schemaResp.Schema},
	}
	assert.False(t, req.Plan.Set(ctx, &plan).HasError())
	assert.False(t, req.State.Set(ctx, &state).HasError())
	resp := resource.UpdateResponse{State: tfsdk.State{Schema: schemaResp.Schema}}

	rs.Update(ctx, req, &resp)

	assert.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	var got streamconnection.TFStreamConnectionModel
	assert.False(t, resp.State.Get(ctx, &got).HasError())
	assert.Equal(t, state.ID, got.ID)
	assert.True(t, got.InstanceName.IsNull())
	assert.Equal(t, "workspace", got.WorkspaceName.ValueString())
	assert.Equal(t, state.ProjectID, got.ProjectID)
	assert.Equal(t, state.ConnectionName, got.ConnectionName)
	assert.Equal(t, state.Type, got.Type)
	assert.True(t, state.Azure.Equal(got.Azure), "azure config, including the known region, must be preserved")
}
