package mongodbatlas

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	framework "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/provider"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"mongodbatlas": func() (tfprotov6.ProviderServer, error) {
		upgradedSdkProvider, err := tf5to6server.UpgradeServer(context.Background(), Provider().GRPCProvider)
		if err != nil {
			log.Fatal(err)
		}
		providers := []func() tfprotov6.ProviderServer{
			func() tfprotov6.ProviderServer {
				return upgradedSdkProvider
			},
			providerserver.NewProtocol6(framework.New()()),
		}

		return tf6muxserver.NewMuxServer(context.Background(), providers...)
	},
}
