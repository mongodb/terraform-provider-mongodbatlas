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

	fwProvider "github.com/hashicorp/terraform-plugin-framework/provider"
)

type ProviderMocked struct {
	OriginalProvider *provider.MongodbtlasProvider
	MockRoundTripper http.RoundTripper
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
	httpClient := client.AtlasV2.GetConfig().HTTPClient
	if httpClient == nil {
		p.t.Fatal("HTTPClient is nil, mocking will fail")
	}
	httpClient.Transport = p.MockRoundTripper
}

func (p *ProviderMocked) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.OriginalProvider.DataSources(ctx)
}
func (p *ProviderMocked) Resources(ctx context.Context) []func() resource.Resource {
	return p.OriginalProvider.Resources(ctx)
}

// Similar to provider.go#muxProviderFactory
func muxProviderFactory(t *testing.T, mockRoundTripper http.RoundTripper) func() tfprotov6.ProviderServer {
	t.Helper()
	v2Provider := provider.NewSdkV2Provider(nil)
	fwProviderInstance := provider.NewFrameworkProvider(nil)
	fwProviderInstanceTyped, ok := fwProviderInstance.(*provider.MongodbtlasProvider)
	if !ok {
		log.Fatal("Failed to cast provider to MongodbtlasProvider")
	}
	mockedProvider := &ProviderMocked{
		OriginalProvider: fwProviderInstanceTyped,
		MockRoundTripper: mockRoundTripper,
		t:                t,
	}
	ctx := context.Background()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, v2Provider.GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}
	muxServer, err := tf6muxserver.NewMuxServer(ctx,
		func() tfprotov6.ProviderServer { return upgradedSdkProvider },
		providerserver.NewProtocol6(mockedProvider),
	)
	if err != nil {
		log.Fatal(err)
	}
	return muxServer.ProviderServer
}

func TestAccProviderV6FactoriesWithMock(t *testing.T, mockRoundTripper http.RoundTripper) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	return map[string]func() (tfprotov6.ProviderServer, error){
		acc.ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return muxProviderFactory(t, mockRoundTripper)(), nil
		},
	}
}
