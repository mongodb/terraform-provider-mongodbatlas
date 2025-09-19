package unit

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	fwProvider "github.com/hashicorp/terraform-plugin-framework/provider"
)

type HTTPClientModifier interface {
	ModifyHTTPClient(*http.Client) error
	ResetHTTPClient(*http.Client)
}

type ProviderMocked struct {
	OriginalProvider *provider.MongodbtlasProvider
	ClientModifier   HTTPClientModifier
	t                *testing.T
}

func (p *ProviderMocked) Metadata(ctx context.Context, req fwProvider.MetadataRequest, resp *fwProvider.MetadataResponse) {
	p.OriginalProvider.Metadata(ctx, req, resp)
}
func (p *ProviderMocked) Schema(ctx context.Context, req fwProvider.SchemaRequest, resp *fwProvider.SchemaResponse) {
	p.OriginalProvider.Schema(ctx, req, resp)
}
func (p *ProviderMocked) Configure(ctx context.Context, req fwProvider.ConfigureRequest, resp *fwProvider.ConfigureResponse) {
	p.OriginalProvider.Configure(ctx, req, resp)
	rd := resp.ResourceData
	client, ok := rd.(*config.MongoDBClient)
	if !ok {
		p.t.Fatal("Failed to cast ResourceData to MongoDBClient")
	}

	// Create a copy of the HTTP client to avoid data races with OAuth2 background operations
	originalClient := client.AtlasV2.GetConfig().HTTPClient
	if originalClient == nil {
		p.t.Fatal("HTTPClient is nil, mocking will fail")
	}

	// Create a new HTTP client to avoid modifying the live one
	mockedClient := &http.Client{
		Transport: originalClient.Transport,
		Timeout:   originalClient.Timeout,
	}

	if p.ClientModifier != nil {
		err := p.ClientModifier.ModifyHTTPClient(mockedClient)
		if err != nil {
			p.t.Fatal(err)
		}
	}

	// Replace the HTTP client in the Atlas configuration
	client.AtlasV2.GetConfig().HTTPClient = mockedClient
}

func (p *ProviderMocked) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.OriginalProvider.DataSources(ctx)
}
func (p *ProviderMocked) Resources(ctx context.Context) []func() resource.Resource {
	return p.OriginalProvider.Resources(ctx)
}

// Similar to provider.go#muxProviderFactory
func muxProviderFactory(t *testing.T, clientModifier HTTPClientModifier) func() tfprotov6.ProviderServer {
	t.Helper()
	v2Provider := provider.NewSdkV2Provider()
	v2ProviderConfigureContextFunc := v2Provider.ConfigureContextFunc
	v2Provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		resp, diags := v2ProviderConfigureContextFunc(ctx, d)
		client, ok := resp.(*config.MongoDBClient)
		if !ok {
			t.Fatalf("Failed to cast response to MongoDBClient, Got type %T", resp)
		}

		// Create a copy of the HTTP client to avoid data races with OAuth2 background operations
		originalClient := client.AtlasV2.GetConfig().HTTPClient
		if originalClient == nil {
			t.Fatalf("HTTPClient is nil, mocking will fail")
		}

		// Create a new HTTP client to avoid modifying the live one
		mockedClient := &http.Client{
			Transport: originalClient.Transport,
			Timeout:   originalClient.Timeout,
		}

		err := clientModifier.ModifyHTTPClient(mockedClient)
		if err != nil {
			t.Fatalf("Failed to modify HTTPClient: %s", err)
		}

		// Replace the HTTP client in the Atlas configuration
		client.AtlasV2.GetConfig().HTTPClient = mockedClient
		return resp, diags
	}
	fwProviderInstance := provider.NewFrameworkProvider()
	fwProviderInstanceTyped, ok := fwProviderInstance.(*provider.MongodbtlasProvider)
	if !ok {
		log.Fatal("Failed to cast provider to MongodbtlasProvider")
	}
	mockedProvider := &ProviderMocked{
		OriginalProvider: fwProviderInstanceTyped,
		ClientModifier:   clientModifier,
		t:                t,
	}
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(t.Context(), v2Provider.GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}
	muxServer, err := tf6muxserver.NewMuxServer(t.Context(),
		func() tfprotov6.ProviderServer { return upgradedSdkProvider },
		providerserver.NewProtocol6(mockedProvider),
	)
	if err != nil {
		log.Fatal(err)
	}
	return muxServer.ProviderServer
}

func TestAccProviderV6FactoriesWithMock(t *testing.T, clientModifier HTTPClientModifier) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	return map[string]func() (tfprotov6.ProviderServer, error){
		acc.ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return muxProviderFactory(t, clientModifier)(), nil
		},
	}
}
