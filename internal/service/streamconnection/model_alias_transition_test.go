package streamconnection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
