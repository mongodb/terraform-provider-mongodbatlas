//nolint:gocritic
package encryptionatrestprivateendpoint_test

// TODO: if acceptance test will be run in an existing CI group of resources, the name should include the group in the prefix followed by the name of the resource e.i. TestAccStreamRSStreamInstance_basic
// In addition, if acceptance test contains testing of both resource and data sources, the RS/DS can be omitted.
// func TestAccEncryptionatrestprivateendpointRS_basic(t *testing.T) {
// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:                 func() { acc.PreCheckBasic(t) },
// 		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
// 		//		CheckDestroy:             checkDestroyEncryptionatrestprivateendpoint,
// 		Steps: []resource.TestStep{ // TODO: verify updates and import in case of resources
// 			//			{
// 			//				Config: encryptionatrestprivateendpointConfig(),
// 			//				Check:  encryptionatrestprivateendpointAttributeChecks(),
// 			//			},
// 			//          {
// 			//				Config: encryptionatrestprivateendpointConfig(),
// 			//				Check:  encryptionatrestprivateendpointAttributeChecks(),
// 			//			},
// 			//			{
// 			//				Config:            encryptionatrestprivateendpointConfig(),
// 			//				ResourceName:      resourceName,
// 			//				ImportStateIdFunc: checkEncryptionatrestprivateendpointImportStateIDFunc,
// 			//				ImportState:       true,
// 			//				ImportStateVerify: true,
// 		},
// 	},
// 	)
// }
