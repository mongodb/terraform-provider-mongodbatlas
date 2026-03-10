package config_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

type mockEphemeralResource struct {
	config.ESCommon
}

func (f *mockEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (f *mockEphemeralResource) Open(_ context.Context, _ ephemeral.OpenRequest, _ *ephemeral.OpenResponse) {
}
func (f *mockEphemeralResource) Renew(_ context.Context, _ ephemeral.RenewRequest, _ *ephemeral.RenewResponse) {
}
func (f *mockEphemeralResource) Close(_ context.Context, _ ephemeral.CloseRequest, _ *ephemeral.CloseResponse) {
}

func TestNoEphemeralInterfaceLoss(t *testing.T) {
	mock := &mockEphemeralResource{ESCommon: config.ESCommon{ResourceName: "test_ephemeral"}}
	wrapped := config.AnalyticsEphemeralResourceFunc(mock)()
	_, ok := wrapped.(ephemeral.EphemeralResourceWithRenew)
	assert.True(t, ok)
	_, ok = wrapped.(ephemeral.EphemeralResourceWithClose)
	assert.True(t, ok)
}
