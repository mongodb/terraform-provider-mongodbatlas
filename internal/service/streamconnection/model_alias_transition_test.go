package streamconnection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
)

func TestIsAliasOnlyTransition(t *testing.T) {
	state := streamconnection.TFStreamConnectionModel{
		ID:               types.StringValue("workspace-project-connection"),
		InstanceName:     types.StringValue("workspace"),
		WorkspaceName:    types.StringNull(),
		ProjectID:        types.StringValue("project"),
		ConnectionName:   types.StringValue("connection"),
		Type:             types.StringValue("Kafka"),
		BootstrapServers: types.StringValue("broker:9092"),
	}
	plan := func(workspaceName, bootstrapServers string) streamconnection.TFStreamConnectionModel {
		return streamconnection.TFStreamConnectionModel{
			ID:               types.StringUnknown(),
			InstanceName:     types.StringNull(),
			WorkspaceName:    types.StringValue(workspaceName),
			ProjectID:        types.StringValue("project"),
			ConnectionName:   types.StringValue("connection"),
			Type:             types.StringValue("Kafka"),
			BootstrapServers: types.StringValue(bootstrapServers),
		}
	}

	testCases := map[string]struct {
		plan                streamconnection.TFStreamConnectionModel
		aliasOnlyTransition bool
	}{
		"alias_only":            {plan: plan("workspace", "broker:9092"), aliasOnlyTransition: true},
		"different_workspace":   {plan: plan("other-workspace", "broker:9092")},
		"configuration_changed": {plan: plan("workspace", "other-broker:9092")},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.aliasOnlyTransition, streamconnection.IsAliasOnlyTransition(&tc.plan, &state))
		})
	}
}
