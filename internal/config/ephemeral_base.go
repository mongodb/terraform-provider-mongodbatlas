package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

const ephemeralErrorConfigureSummary = "Unexpected Ephemeral Resource Configure Type"
const ephemeralErrorConfigure = "expected *EphemeralResourceData, got: %T. Please report this issue to the provider developers"

type EphemeralResourceData struct {
	ClientID         string
	ClientSecret     string
	BaseURL          string
	TerraformVersion string
}

type ImplementedEphemeralResource interface {
	ephemeral.EphemeralResourceWithConfigure
	GetName() string
	SetClient(*EphemeralResourceData)
}

func AnalyticsEphemeralResourceFunc(iResource ephemeral.EphemeralResource) func() ephemeral.EphemeralResource {
	commonResource, ok := iResource.(ImplementedEphemeralResource)
	if !ok {
		panic(fmt.Sprintf("ephemeral resource %T didn't comply with the ImplementedEphemeralResource interface", iResource))
	}
	return func() ephemeral.EphemeralResource {
		return analyticsEphemeralResource(commonResource)
	}
}

// analyticsEphemeralResource wraps an ImplementedEphemeralResource with ESCommon to add analytics tracking.
// We cannot return iResource directly because we need to intercept the Open operation
// to inject user agent information into the context before calling the actual ephemeral resource method.
func analyticsEphemeralResource(iResource ImplementedEphemeralResource) ephemeral.EphemeralResource {
	return &ESCommon{
		ResourceName:                 iResource.GetName(),
		ImplementedEphemeralResource: iResource,
	}
}

// ESCommon is used as an embedded struct for framework ephemeral resources.
// It implements shared Metadata and Configure behavior.
//
// When used as a wrapper (ImplementedEphemeralResource is set), it intercepts Open to add analytics tracking.
// When embedded in an ephemeral resource struct, the resource's own Open method is used.
//
// EphemeralResourceData is left empty and populated by the framework when invoking Configure.
// ResourceName must be defined when creating an instance of an ephemeral resource.
type ESCommon struct {
	ImplementedEphemeralResource // Set when used as a wrapper, nil when embedded
	EphemeralResourceData        *EphemeralResourceData
	ResourceName                 string
}

func (e *ESCommon) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, e.ResourceName)
}

func (e *ESCommon) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	if e.ImplementedEphemeralResource != nil {
		e.ImplementedEphemeralResource.Schema(ctx, req, resp)
	}
}

func (e *ESCommon) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	data, err := configureEphemeralResourceData(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError(ephemeralErrorConfigureSummary, err.Error())
		return
	}
	e.EphemeralResourceData = data
	if e.ImplementedEphemeralResource != nil {
		e.ImplementedEphemeralResource.SetClient(data)
	}
}

func (e *ESCommon) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	if e.ImplementedEphemeralResource == nil {
		return
	}
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      userAgentNameValue(e.ResourceName),
		Operation: UserAgentOperationValueOpen,
	})
	e.ImplementedEphemeralResource.Open(ctx, req, resp)
}

// Renew delegates to the wrapped resource if it implements EphemeralResourceWithRenew.
func (e *ESCommon) Renew(ctx context.Context, req ephemeral.RenewRequest, resp *ephemeral.RenewResponse) {
	resourceWithRenew, ok := e.ImplementedEphemeralResource.(ephemeral.EphemeralResourceWithRenew)
	if !ok {
		return
	}
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      userAgentNameValue(e.ResourceName),
		Operation: UserAgentOperationValueRenew,
	})
	resourceWithRenew.Renew(ctx, req, resp)
}

// Close delegates to the wrapped resource if it implements EphemeralResourceWithClose.
func (e *ESCommon) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	resourceWithClose, ok := e.ImplementedEphemeralResource.(ephemeral.EphemeralResourceWithClose)
	if !ok {
		return
	}
	ctx = AddUserAgentExtra(ctx, UserAgentExtra{
		Name:      userAgentNameValue(e.ResourceName),
		Operation: UserAgentOperationValueClose,
	})
	resourceWithClose.Close(ctx, req, resp)
}

func (e *ESCommon) GetName() string {
	return e.ResourceName
}

func (e *ESCommon) SetClient(data *EphemeralResourceData) {
	e.EphemeralResourceData = data
}

func configureEphemeralResourceData(providerData any) (*EphemeralResourceData, error) {
	if providerData == nil {
		return nil, nil
	}

	if data, ok := providerData.(*EphemeralResourceData); ok {
		return data, nil
	}

	return nil, fmt.Errorf(ephemeralErrorConfigure, providerData)
}
