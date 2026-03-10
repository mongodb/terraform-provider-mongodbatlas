package config_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEphemeral struct {
	config.ESCommon
	capturedCtx                          context.Context
	openCalled, renewCalled, closeCalled bool
}

func (f *fakeEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (f *fakeEphemeral) Open(ctx context.Context, _ ephemeral.OpenRequest, _ *ephemeral.OpenResponse) {
	f.openCalled = true
	f.capturedCtx = ctx
}

func (f *fakeEphemeral) Renew(ctx context.Context, _ ephemeral.RenewRequest, _ *ephemeral.RenewResponse) {
	f.renewCalled = true
	f.capturedCtx = ctx
}

func (f *fakeEphemeral) Close(ctx context.Context, _ ephemeral.CloseRequest, _ *ephemeral.CloseResponse) {
	f.closeCalled = true
	f.capturedCtx = ctx
}

func newFakeEphemeral() *fakeEphemeral {
	return &fakeEphemeral{
		ESCommon: config.ESCommon{ResourceName: "test_ephemeral"},
	}
}

func TestAnalyticsEphemeralResourceFunc_WrapsCorrectly(t *testing.T) {
	fake := newFakeEphemeral()
	wrapped := config.AnalyticsEphemeralResourceFunc(fake)()
	require.NotNil(t, wrapped)
}

type nonCompliantEphemeral struct{}

func (n *nonCompliantEphemeral) Metadata(_ context.Context, _ ephemeral.MetadataRequest, _ *ephemeral.MetadataResponse) {
}
func (n *nonCompliantEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, _ *ephemeral.SchemaResponse) {
}
func (n *nonCompliantEphemeral) Open(_ context.Context, _ ephemeral.OpenRequest, _ *ephemeral.OpenResponse) {
}

func TestAnalyticsEphemeralResourceFunc_PanicsOnNonCompliant(t *testing.T) {
	assert.Panics(t, func() {
		config.AnalyticsEphemeralResourceFunc(&nonCompliantEphemeral{})
	})
}

func TestEphemeralWrapper_PreservesOptionalInterfaces(t *testing.T) {
	fake := newFakeEphemeral()
	wrapped := config.AnalyticsEphemeralResourceFunc(fake)()

	_, hasRenew := wrapped.(ephemeral.EphemeralResourceWithRenew)
	assert.True(t, hasRenew, "wrapped resource must preserve EphemeralResourceWithRenew")

	_, hasClose := wrapped.(ephemeral.EphemeralResourceWithClose)
	assert.True(t, hasClose, "wrapped resource must preserve EphemeralResourceWithClose")
}

func TestEphemeralWrapper_OpenSetsUserAgent(t *testing.T) {
	fake := newFakeEphemeral()
	wrapped := config.AnalyticsEphemeralResourceFunc(fake)()

	wrapped.Open(context.Background(), ephemeral.OpenRequest{}, &ephemeral.OpenResponse{})
	require.True(t, fake.openCalled)

	ua := config.ReadUserAgentExtra(fake.capturedCtx)
	require.NotNil(t, ua)
	assert.Equal(t, "test_ephemeral", ua.Name)
	assert.Equal(t, config.UserAgentOperationValueOpen, ua.Operation)
}

func TestEphemeralWrapper_RenewSetsUserAgent(t *testing.T) {
	fake := newFakeEphemeral()
	wrapped := config.AnalyticsEphemeralResourceFunc(fake)()

	wrapped.(ephemeral.EphemeralResourceWithRenew).Renew(context.Background(), ephemeral.RenewRequest{}, &ephemeral.RenewResponse{})
	require.True(t, fake.renewCalled)

	ua := config.ReadUserAgentExtra(fake.capturedCtx)
	require.NotNil(t, ua)
	assert.Equal(t, "test_ephemeral", ua.Name)
	assert.Equal(t, config.UserAgentOperationValueRenew, ua.Operation)
}

func TestEphemeralWrapper_CloseSetsUserAgent(t *testing.T) {
	fake := newFakeEphemeral()
	wrapped := config.AnalyticsEphemeralResourceFunc(fake)()

	wrapped.(ephemeral.EphemeralResourceWithClose).Close(context.Background(), ephemeral.CloseRequest{}, &ephemeral.CloseResponse{})
	require.True(t, fake.closeCalled)

	ua := config.ReadUserAgentExtra(fake.capturedCtx)
	require.NotNil(t, ua)
	assert.Equal(t, "test_ephemeral", ua.Name)
	assert.Equal(t, config.UserAgentOperationValueClose, ua.Operation)
}
