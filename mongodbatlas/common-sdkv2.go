package mongodbatlas

// import (
// 	// "github.com/hashicorp/terraform-plugin-framework/provider"
// 	tfsdkv2 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	// frameworkProvider "github.com/mongodb/terraform-provider-mongodbatlas/internal/framework/provider"
// )

// type SDKv2 struct{}

// func NewImpl() *SDKv2 {
// 	return &SDKv2{}
// }

// func (v *SDKv2) GetProvider() *tfsdkv2.Provider {
// 	return Provider()
// }

// func (p *SDKv2) TestFunc() {
// 	// pframework := framework.NewImpl(p)
// 	// pframework.HelloFromP2()
// }

// // func (v *SDKv2) GetFrameworkProvider() func() provider.Provider {
// // 	return frameworkProvider.New()
// // }

// func (v *SDKv2) GetTestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
// 	var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 		"mongodbatlas": func() (tfprotov6.ProviderServer, error) {
// 			ctx := context.Background()
// 			upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, v.GetProvider().GRPCProvider)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			providers := []func() tfprotov6.ProviderServer{
// 				func() tfprotov6.ProviderServer {
// 					return upgradedSdkProvider
// 				},
// 				providerserver.NewProtocol6(frameworkProvider.New()()),
// 			}

// 			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

// 			if err != nil {
// 				return nil, err
// 			}

// 			return muxServer.ProviderServer(), nil
// 		},
// 	}
// 	return TestAccProtoV6ProviderFactories
// }

// var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 	"mongodbatlas": func() (tfprotov6.ProviderServer, error) {
// 		sdkv2 := SDKv2{}

// 		upgradedSdkProvider, err := tf5to6server.UpgradeServer(context.Background(), Provider().GRPCProvider)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		providers := []func() tfprotov6.ProviderServer{
// 			func() tfprotov6.ProviderServer {
// 				return upgradedSdkProvider
// 			},
// 			providerserver.NewProtocol6(sdkv2.GetFrameworkProvider()()),
// 		}

// 		muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)

// 		if err != nil {
// 			return nil, err
// 		}

// 		return muxServer.ProviderServer(), nil
// 	},
// }

// func GetTestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
// 	sdkv2 := SDKv2{}
// 	return sdkv2.GetTestAccProtoV6ProviderFactories()

// }
