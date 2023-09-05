package mongodbatlas

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/mongodb/terraform-provider-mongodbatlas/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/test/mock"
	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

const (
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

// frameworkTestProvider is a test version of the plugin-framework version of the provider
// that uses mock clients
type frameworkTestProvider struct {
	MongodbtlasProvider
	MockClient *config.MongoDBClient
}

// UnitTestMuxedProviderFactoryWithProvider creates mux provider using existing sdk v2 provider passed as parameter and creating new instance of framework provider.
// Used in testing where existing sdk v2 provider has to be used.
func UnitTestMuxedProviderFactoryWithProvider() func() tfprotov6.ProviderServer {
	ctx := context.Background()

	// Unit Tests
	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(NewFrameworkTestProvider()),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}
	return muxServer.ProviderServer
}

// Configure is here to overwrite the FrameworkProvider configure function for unit testing
func (p *frameworkTestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.DataSourceData = p.MockClient
	resp.ResourceData = p.MockClient
}

func (p *frameworkTestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return p.MongodbtlasProvider.DataSources(ctx)
}

func (p *frameworkTestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return p.MongodbtlasProvider.Resources(ctx)
}

func (p *frameworkTestProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mongodbatlas"
	resp.Version = version.ProviderVersion
}
func (p *frameworkTestProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	p.MongodbtlasProvider.Schema(ctx, req, resp)
}

func NewFrameworkTestProvider() provider.Provider {
	return &frameworkTestProvider{
		MongodbtlasProvider: MongodbtlasProvider{},
		MockClient:          mock.NewMockMongoDBClient(),
	}
}

func NewUnitTestProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return UnitTestMuxedProviderFactoryWithProvider()(), nil
		},
	}
}
