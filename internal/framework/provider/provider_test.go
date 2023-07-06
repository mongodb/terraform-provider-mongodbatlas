package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	mongodbatlasSDKv2 "github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas"
)

func TestMuxServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"mongodbatlas": func() (tfprotov6.ProviderServer, error) {
				ctx := context.Background()

				upgradedSdkServer, err := tf5to6server.UpgradeServer(
					ctx,
					mongodbatlasSDKv2.Provider().GRPCProvider,
				)

				if err != nil {
					return nil, err
				}

				providers := []func() tfprotov6.ProviderServer{
					func() tfprotov6.ProviderServer {
						return upgradedSdkServer
					},
					providerserver.NewProtocol6(New()()),
				}

				muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

				if err != nil {
					return nil, err
				}

				return muxServer.ProviderServer(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `resource mongodbatlas_example "test" {
					configurable_attribute = "config_attr_val"
				}`,
			},
		},
	})
}
