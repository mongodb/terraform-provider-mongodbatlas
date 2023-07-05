package provider

// import (
// 	tfsdkv2 "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// type sdkv2 interface {
// 	GetProvider() *tfsdkv2.Provider
// }

// type Framework struct {
// 	SDKv2 sdkv2
// }

// func NewImpl(sdkv2 sdkv2) *Framework {
// 	return &Framework{
// 		SDKv2: sdkv2,
// 	}
// }

// func (p *Framework) TestFunc() {
// 	pframework := p.NewImpl(p)
// 	// pframework.HelloFromP2()
// }

// // func (p *Framework) HelloFromP2() {
// // 	fmt.Println("Hello from package p2")
// // }

// func (p *Framework) SDKv2Provider() {
// 	p.SDKv2.GetProvider()
// }

// func (p *Framework) GetTestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
// 	var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 		"mongodbatlas": func() (tfprotov6.ProviderServer, error) {
// 			ctx := context.Background()
// 			upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, p.SDKv2.GetProvider().GRPCProvider)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			providers := []func() tfprotov6.ProviderServer{
// 				func() tfprotov6.ProviderServer {
// 					return upgradedSdkProvider
// 				},
// 				providerserver.NewProtocol6(New()()),
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
